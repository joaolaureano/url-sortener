# urlshortener

**[Leia em português / Read this in Portuguese](README.pt-BR.md)**

A small URL shortener in Go. It maps a long URL to a short, fixed-width code
and back again, in four layers with a single dependency-free module.

```
main → operator → { url, storage }
                     url → hash
```

## Design

The package boundary is the design. `operator` is a thin facade over the two
things shortening needs — a code (`url`) and somewhere to keep it (`storage`) —
and each layer depends only on the ones below it, never sideways or up.

| Package | Responsibility |
|---|---|
| `operator` | Facade: shorten, store, recover. The only entry point a caller needs. |
| `url` | Derives the short code and owns the public link format. |
| `hash` | Content hashing (`Operator` interface, `MD5Hasher` default). |
| `storage` | The code → original-URL references (`Store` interface, in-memory default). |

`url` and `storage` know nothing about `net/http`; `hash` knows nothing about
URLs. That separation is enforced, not conventional: the store is keyed by short
code alone, so the link format can change without touching stored data, and the
hashing strategy is swappable behind an interface.

## Short codes

A code is **7 base62 characters** — `62^7 ≈ 3.5 trillion` values. The first
expected collision sits around ~2.6M links (birthday bound), versus ~1.2K for a
5-character hex code. On a collision the store **refuses** to overwrite an
occupied code, and `operator` probes deterministic alternates before failing
with `ErrCollisionExhausted`. Shortening the same URL twice is idempotent and
returns the same short URL.

MD5 backs the code derivation for speed and uniform distribution, **not** for
any security property — the digests address short codes, they never protect
secrets.

## Usage

```go
short, err := operator.CreateNewShortURL("https://www.example.com/some/long/path")
// short == "https://me.li/MuAuhzJ"

original, err := operator.RecoverOriginalURL(short)
// original == "https://www.example.com/some/long/path"
```

Errors are sentinels, checkable with `errors.Is`:

- `url.ErrInvalidURL` — the input is not a valid URL.
- `url.ErrInvalidShortURL` — the string was not issued by this service.
- `storage.ErrNotFound` — no reference for that code.
- `operator.ErrCollisionExhausted` — no free code could be allocated.

Run the demo:

```sh
go run .
```

## Swapping a layer

Both the store and the hasher are injectable. To back the service with a durable
store, implement `storage.Store` and assign it:

```go
storage.Default = myRedisStore{}
```

The same pattern applies to `hash.Default` (any `hash.Operator`).

## Testing

```sh
go test ./... -race      # unit tests, race detector
go test ./operator/ -bench . -count=10
```

The store is safe for concurrent use (`sync.RWMutex`); a benchmark exercises it
under 100 goroutines. Every layer returns errors rather than panicking, and no
package logs from its hot path.

## Performance note

`CreateNewShortURL` was profiled with the benchmark under `operator/`. Validation
originally ran through a regexp-based URL validator, which accounted for ~72% of
CPU samples; replacing it with `net/url` parsing cut shortening from ~12.3µs to
~1.1µs (−91%) and left the module with **no external dependencies**.
