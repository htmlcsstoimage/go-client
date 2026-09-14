package hcti

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
)

// signedURLBuilder reserves space for the signature, then writes the query into
// the same buffer. It signs the encoded query bytes and fills the reserved space.
type signedURLBuilder struct {
	buf        []byte
	tokenStart int
	queryStart int
	hasQuery   bool
}

func newSignedURLBuilder(buf []byte, base, path string, format ImageFormat) (signedURLBuilder, error) {
	if format != "" && format != PNG && format != JPG && format != WebP && format != PDF {
		return signedURLBuilder{}, fmt.Errorf("hcti: unsupported image format %q", format)
	}
	buf = append(buf, base...)
	buf = append(buf, path...)
	buf = append(buf, '/')
	tokenStart := len(buf)
	buf = append(buf, make([]byte, sha256.Size*2)...)
	if format != "" {
		buf = append(buf, '/')
		buf = append(buf, format...)
	}
	return signedURLBuilder{buf: buf, tokenStart: tokenStart, queryStart: len(buf)}, nil
}

// safeQueryKey is an internal, unescaped API parameter name. Only fixed ASCII
// names containing letters, digits, '-', '_', '.', or '~' may use this path.
// Dynamic keys must use key/stringValue instead.
type safeQueryKey string

func (b *signedURLBuilder) separator() {
	if b.hasQuery {
		b.buf = append(b.buf, '&')
	} else {
		b.buf = append(b.buf, '?')
		b.queryStart = len(b.buf)
		b.hasQuery = true
	}
}

func (b *signedURLBuilder) key(key string) {
	b.separator()
	b.buf = appendQueryEscaped(b.buf, key)
	b.buf = append(b.buf, '=')
}

func (b *signedURLBuilder) safeKey(key safeQueryKey) {
	b.separator()
	b.buf = append(b.buf, key...)
	b.buf = append(b.buf, '=')
}

func (b *signedURLBuilder) stringValue(key, value string) {
	b.key(key)
	b.buf = appendQueryEscaped(b.buf, value)
}

func (b *signedURLBuilder) safeStringValue(key safeQueryKey, value string) {
	b.safeKey(key)
	b.buf = appendQueryEscaped(b.buf, value)
}

func (b *signedURLBuilder) optionalSafeString(key safeQueryKey, value string) {
	if value != "" {
		b.safeStringValue(key, value)
	}
}

func (b *signedURLBuilder) optionalSafeInt(key safeQueryKey, value *int) {
	if value != nil {
		b.safeKey(key)
		b.buf = strconv.AppendInt(b.buf, int64(*value), 10)
	}
}

func (b *signedURLBuilder) optionalSafeBool(key safeQueryKey, value *bool) {
	if value != nil {
		b.safeKey(key)
		b.buf = strconv.AppendBool(b.buf, *value)
	}
}

func (b *signedURLBuilder) finish(secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(b.buf[b.queryStart:])
	var digest [sha256.Size]byte
	hex.Encode(b.buf[b.tokenStart:b.tokenStart+sha256.Size*2], mac.Sum(digest[:0]))
	return string(b.buf)
}

// appendQueryEscaped matches url.QueryEscape without allocating an intermediate string.
func appendQueryEscaped(dst []byte, value string) []byte {
	const hexDigits = "0123456789ABCDEF"
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '-', c == '_', c == '.', c == '~':
			dst = append(dst, c)
		case c == ' ':
			dst = append(dst, '+')
		default:
			dst = append(dst, '%', hexDigits[c>>4], hexDigits[c&15])
		}
	}
	return dst
}

func appendOptionalString[T ~string](b *signedURLBuilder, key safeQueryKey, value *T) {
	if value != nil {
		b.safeStringValue(key, string(*value))
	}
}

func (b *signedURLBuilder) appendJSONFloat(value float64, bits int) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return fmt.Errorf("hcti: numeric query values must be finite")
	}
	format := byte('f')
	abs := math.Abs(value)
	if bits == 32 {
		abs = float64(float32(abs))
	}
	if abs != 0 && ((bits == 64 && (abs < 1e-6 || abs >= 1e21)) || (bits == 32 && (float32(abs) < 1e-6 || float32(abs) >= 1e21))) {
		format = 'e'
	}
	var scratch [32]byte
	encoded := strconv.AppendFloat(scratch[:0], value, format, -1, bits)
	if n := len(encoded); format == 'e' && n >= 4 && string(encoded[n-4:n-1]) == "e-0" {
		encoded[n-2] = encoded[n-1]
		encoded = encoded[:n-1]
	}
	b.buf = appendQueryEscaped(b.buf, string(encoded))
	return nil
}
