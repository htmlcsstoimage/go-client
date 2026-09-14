package hcti

import (
	"fmt"
	"net/url"
	"strings"
	"testing"
)

// Includes the C# encoder's ASCII, é, 漢, emoji, mixed, and long-input cases.
// Go strings can also hold arbitrary bytes, unlike C#'s UTF-16 strings.
var encodingCases = []struct{ name, value string }{
	{"empty", ""},
	{"unreserved", "abc-_.~XYZ123"},
	{"latin", "é"},
	{"cjk", "漢"},
	{"emoji", "😀👀"},
	{"mixed", "hello é 漢 😀 world مرحبا"},
	{"emoji_sequences", "👩🏽‍💻👨‍👩‍👧‍👦🇺🇸❤️"},
	{"combining_marks", "é e\u0301"},
	{"three_byte_boundary", "\u0800"},
	{"unicode_limits", "\u007f\u0080\u07ff\u0800\ud7ff\ue000\uffff\U00010000\U0010ffff"},
	{"reserved", " +%&=?/#:;[]!$'()*,-._~"},
	{"already_escaped", "%F0%9F%98%80%20%2B"},
	{"controls", "\x00\t\r\n\x1f\x7f"},
	{"invalid_leading", "\xffabc"},
	{"isolated_continuation", "a\x80b"},
	{"truncated", "\xf0\x9f\x98"},
	{"overlong", "\xc0\xaf"},
	{"encoded_surrogate", "\xed\xa0\xbdabc"},
	{"out_of_range", "\xf4\x90\x80\x80"},
	{"mixed_invalid", "é\xff漢\xc0\xaf😀"},
	{"long_unicode", strings.Repeat("é漢😀abc123", 1000)},
}

func TestQueryEncodingCorpus(t *testing.T) {
	for _, tc := range encodingCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, capacity := range []int{0, 1, 4, 8, 1024} {
				t.Run(fmt.Sprint(capacity), func(t *testing.T) {
					prefix := []byte("prefix=")
					buf := make([]byte, len(prefix), len(prefix)+capacity)
					copy(buf, prefix)
					got := string(appendQueryEscaped(buf, tc.value))
					want := "prefix=" + url.QueryEscape(tc.value)
					if got != want {
						t.Fatalf("encoding differs: got %q; want %q", got, want)
					}
					decoded, err := url.QueryUnescape(got[len(prefix):])
					if err != nil || decoded != tc.value {
						t.Fatalf("byte round trip failed: %v", err)
					}
				})
			}
		})
	}
}

func TestDynamicKeysAndSafeValuesCorpus(t *testing.T) {
	for _, tc := range encodingCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := newSignedURLBuilder(make([]byte, 0, 4), "https://example.test", "/test", PNG)
			if err != nil {
				t.Fatal(err)
			}
			b.stringValue(tc.value, tc.value)
			b.safeStringValue("css", tc.value)
			want := url.QueryEscape(tc.value) + "=" + url.QueryEscape(tc.value) + "&css=" + url.QueryEscape(tc.value)
			if got := string(b.buf[b.queryStart:]); got != want {
				t.Fatal("key/value encoding or query boundary differs")
			}
			assertSignedQuery(t, b.finish("secret"), "secret", PNG, want)
		})
	}
}

func FuzzAppendQueryEscaped(f *testing.F) {
	for _, tc := range encodingCases {
		f.Add(tc.value)
	}
	f.Fuzz(func(t *testing.T, value string) {
		got := string(appendQueryEscaped([]byte("prefix="), value))
		if want := "prefix=" + url.QueryEscape(value); got != want {
			t.Fatalf("got %q; want %q", got, want)
		}
		decoded, err := url.QueryUnescape(got[len("prefix="):])
		if err != nil || decoded != value {
			t.Fatalf("round trip failed: %v", err)
		}
	})
}
