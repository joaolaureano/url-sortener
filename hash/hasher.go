// Package hash provides the content hashing primitives used to derive short
// codes. It knows nothing about URLs or storage.
package hash

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
)

// Size is the length in bytes of the digest returned by Digest.
const Size = md5.Size

// Operator hashes content and checks content against a hash.
//
// It exists so callers can swap the hashing strategy — a different algorithm,
// a salted variant, a stub in tests — without touching the layers above.
type Operator interface {
	// Encode returns the hex-encoded digest of content.
	Encode(content string) string
	// Compare reports whether content hashes to the given hex-encoded hash.
	Compare(content string, hash string) bool
}

// MD5Hasher is the default Operator.
//
// MD5 is used for its speed and uniform distribution, not for its security
// properties: the digests here address short codes, they never protect
// secrets. Do not reuse this type for authentication.
type MD5Hasher struct{}

// Default is the Operator backing the package-level functions. Replace it
// to change the hashing strategy process-wide.
var Default Operator = MD5Hasher{}

// Encode implements Operator.
func (MD5Hasher) Encode(content string) string {
	return hex.EncodeToString(digest(content))
}

// Compare implements Operator. The comparison is constant-time.
func (MD5Hasher) Compare(content string, hash string) bool {
	return subtle.ConstantTimeCompare([]byte(hex.EncodeToString(digest(content))), []byte(hash)) == 1
}

// Encode returns the hex-encoded digest of content, using Default.
func Encode(content string) string {
	return Default.Encode(content)
}

// Compare reports whether content hashes to hash, using Default.
func Compare(content string, hash string) bool {
	return Default.Compare(content, hash)
}

// Digest returns the raw digest of content, for callers that need the bytes
// rather than their hex form — the base62 encoder in package url, for one.
// The returned slice is always Size bytes long and is owned by the caller.
func Digest(content string) []byte {
	return digest(content)
}

func digest(content string) []byte {
	sum := md5.Sum([]byte(content))
	return sum[:]
}
