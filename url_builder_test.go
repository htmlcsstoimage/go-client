package hcti

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"
)

func TestAppendQueryEscaped(t *testing.T) {
	var allBytes [256]byte
	for i := range allBytes {
		allBytes[i] = byte(i)
	}
	for _, value := range []string{"", "hello world", "a+b&c=d?e/#f", "你好 ✓ café", string(allBytes[:])} {
		got := string(appendQueryEscaped([]byte("prefix="), value))
		if want := "prefix=" + url.QueryEscape(value); got != want {
			t.Fatalf("got %q; want %q", got, want)
		}
	}
}

func TestSafeKeysStillEscapeValues(t *testing.T) {
	b, err := newSignedURLBuilder(nil, "https://hcti.io", "/test", PNG)
	if err != nil {
		t.Fatal(err)
	}
	b.safeStringValue("url", "https://example.com/?a=1&b=hello world")
	b.optionalSafeString("css", "a::before { content: '✓'; }")
	b.optionalSafeBool("full_screen", Ptr(false))
	b.optionalSafeInt("ms_delay", Ptr(0))
	b.optionalSafeString("selector", "")
	b.optionalSafeBool("viewport_mobile", nil)
	b.optionalSafeInt("max_wait_ms", nil)
	// Caller-supplied keys take the escaped path and cannot add query parameters.
	b.stringValue("title&admin=true ✓", "a+b&c=d")
	want := "url=" + url.QueryEscape("https://example.com/?a=1&b=hello world") +
		"&css=" + url.QueryEscape("a::before { content: '✓'; }") +
		"&full_screen=false&ms_delay=0&" + url.QueryEscape("title&admin=true ✓") +
		"=" + url.QueryEscape("a+b&c=d")
	if got := string(b.buf[b.queryStart:]); got != want {
		t.Fatalf("got %q; want %q", got, want)
	}
}

func TestSignedURLBuilder(t *testing.T) {
	for _, value := range []string{"", "a+b & ✓", strings.Repeat("<&>你好", 4096)} {
		var buffer [128]byte // Force growth while preserving the reserved token position.
		b, err := newSignedURLBuilder(buffer[:0], "https://hcti.io", "/test", WebP)
		if err != nil {
			t.Fatal(err)
		}
		b.stringValue("first", value)
		b.stringValue("first", "second")
		b.optionalSafeBool("enabled", Ptr(false))
		b.optionalSafeInt("delay", Ptr(0))
		query := "first=" + url.QueryEscape(value) + "&first=second&enabled=false&delay=0"
		mac := hmac.New(sha256.New, []byte("test-secret"))
		_, _ = mac.Write([]byte(query))
		want := "https://hcti.io/test/" + hex.EncodeToString(mac.Sum(nil)) + "/webp?" + query
		got := b.finish("test-secret")
		if got != want {
			t.Fatal("signature or query differs after buffer growth")
		}
		// Returned strings must not alias the mutable buffer.
		b.buf[len(b.buf)-1] = 'x'
		if got != want {
			t.Fatal("returned URL aliases builder storage")
		}
	}
	b, err := newSignedURLBuilder(nil, "https://hcti.io", "/test", "")
	if err != nil {
		t.Fatal(err)
	}
	mac := hmac.New(sha256.New, []byte("test-secret"))
	want := "https://hcti.io/test/" + hex.EncodeToString(mac.Sum(nil))
	if got := b.finish("test-secret"); got != want {
		t.Fatalf("empty query: %s", got)
	}
}
