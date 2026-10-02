# urlshortener

**[Read this in English / Leia em inglês](README.md)**

Um encurtador de URL pequeno em Go. Ele transforma uma URL longa em um código
curto de tamanho fixo e faz o caminho de volta, em quatro camadas dentro de um
único módulo sem dependências.

```
main → operator → { url, storage }
                     url → hash
```

## Design

A fronteira entre pacotes é o design. `operator` é uma fachada fina sobre as
duas coisas de que o encurtamento precisa — um código (`url`) e um lugar para
guardá-lo (`storage`) — e cada camada depende só das que estão abaixo dela,
nunca para o lado ou para cima.

| Pacote | Responsabilidade |
|---|---|
| `operator` | Fachada: encurtar, guardar, recuperar. O único ponto de entrada de que quem chama precisa. |
| `url` | Deriva o código curto e é dono do formato público do link. |
| `hash` | Hash do conteúdo (interface `Operator`, `MD5Hasher` como padrão). |
| `storage` | As referências código → URL original (interface `Store`, em memória como padrão). |

`url` e `storage` não sabem nada de `net/http`; `hash` não sabe nada de URLs.
Essa separação é garantida, não convencional: o store usa só o código curto como
chave, então o formato do link pode mudar sem tocar nos dados guardados, e a
estratégia de hash pode ser trocada por trás de uma interface.

## Códigos curtos

Um código tem **7 caracteres base62** — `62^7 ≈ 3,5 trilhões` de valores. A
primeira colisão esperada fica por volta de ~2,6 milhões de links (limite do
aniversário), contra ~1,2 mil de um código hexadecimal de 5 caracteres. Em uma
colisão, o store **se recusa** a sobrescrever um código ocupado, e o `operator`
testa alternativas determinísticas antes de falhar com `ErrCollisionExhausted`.
Encurtar a mesma URL duas vezes é idempotente e devolve a mesma URL curta.

O MD5 sustenta a derivação do código pela velocidade e pela distribuição
uniforme, **não** por alguma propriedade de segurança — os digests endereçam
códigos curtos, nunca protegem segredos.

## Uso

```go
short, err := operator.CreateNewShortURL("https://www.example.com/some/long/path")
// short == "https://me.li/MuAuhzJ"

original, err := operator.RecoverOriginalURL(short)
// original == "https://www.example.com/some/long/path"
```

Os erros são sentinelas, verificáveis com `errors.Is`:

- `url.ErrInvalidURL` — a entrada não é uma URL válida.
- `url.ErrInvalidShortURL` — a string não foi emitida por este serviço.
- `storage.ErrNotFound` — não há referência para esse código.
- `operator.ErrCollisionExhausted` — nenhum código livre pôde ser alocado.

Rode a demo:

```sh
go run .
```

## Trocando uma camada

O store e o hasher podem ser injetados. Para usar um store durável, implemente
`storage.Store` e atribua:

```go
storage.Default = myRedisStore{}
```

O mesmo padrão vale para `hash.Default` (qualquer `hash.Operator`).

## Testes

```sh
go test ./... -race      # testes unitários, com detector de race
go test ./operator/ -bench . -count=10
```

O store é seguro para uso concorrente (`sync.RWMutex`); um benchmark o exercita
com 100 goroutines. Toda camada devolve erros em vez de entrar em panic, e nenhum
pacote faz log no caminho quente.

## Nota de desempenho

`CreateNewShortURL` foi perfilado com o benchmark em `operator/`. A validação
passava originalmente por um validador de URL baseado em regexp, responsável por
~72% das amostras de CPU; trocá-lo pelo parsing de `net/url` reduziu o
encurtamento de ~12,3µs para ~1,1µs (−91%) e deixou o módulo **sem nenhuma
dependência externa**.
