package fastio

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func newTestReader(s string) *FastReader {
	return NewReader(strings.NewReader(s))
}

func TestNextIntSimple(t *testing.T) {
	r := newTestReader("123")
	v, err := r.NextInt()
	if err != nil {
		t.Fatalf("NextInt error: %v", err)
	}
	if v != 123 {
		t.Fatalf("NextInt = %d; want 123", v)
	}
}

func TestNextIntWithSpacesAndNewLines(t *testing.T) {
	r := newTestReader("   10\n  20\t30  \n40")
	want := []int{10, 20, 30, 40}

	for i, w := range want {
		v, err := r.NextInt()
		if err != nil {
			t.Fatalf("NextInt error at index %d: %v", i, err)
		}
		if v != w {
			t.Fatalf("NextInt at index %d = %d; want %d", i, v, w)
		}
	}

	_, err := r.NextInt()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Expected EOF error, got: %v", err)
	}
}

func TestNextIntWithSign(t *testing.T) {
	r := newTestReader(" -15 +25 0 -0 ")
	want := []int{-15, 25, 0, 0}

	for i, w := range want {
		v, err := r.NextInt()
		if err != nil {
			t.Fatalf("NextInt error at index %d: %v", i, err)
		}
		if v != w {
			t.Fatalf("NextInt at index %d = %d; want %d", i, v, w)
		}
	}

	_, err := r.NextInt()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Expected EOF error, got: %v", err)
	}
}

func TestNextIntNoDigits(t *testing.T) {
	r := newTestReader("   abc")
	_, err := r.NextInt()
	if err == nil {
		t.Fatalf("Expected error for no digits, got nil")
	}
}

func TestNextWordBasic(t *testing.T) {
	r := newTestReader("hello world\tthis is a test\n")
	want := []string{"hello", "world", "this", "is", "a", "test"}

	for i, w := range want {
		v, err := r.NextWord()
		if err != nil {
			t.Fatalf("NextWord error at index %d: %v", i, err)
		}
		if v != w {
			t.Fatalf("NextWord at index %d = %q; want %q", i, v, w)
		}
	}

	_, err := r.NextWord()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Expected EOF error, got: %v", err)
	}
}

func TestNextLineUnix(t *testing.T) {
	r := newTestReader("first line\nsecond line\nthird line\n")
	want := []string{"first line", "second line", "third line"}

	for i, w := range want {
		v, err := r.NextLine()
		if err != nil {
			t.Fatalf("NextLine error at index %d: %v", i, err)
		}
		if v != w {
			t.Fatalf("NextLine at index %d = %q; want %q", i, v, w)
		}
	}

	_, err := r.NextLine()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Expected EOF error, got: %v", err)
	}
}

func TestNextLineWindowsCRLF(t *testing.T) {
	r := newTestReader("line one\r\nline two\r\nline three\r\n")
	want := []string{"line one", "line two", "line three"}

	for i, w := range want {
		v, err := r.NextLine()
		if err != nil {
			t.Fatalf("NextLine error at index %d: %v", i, err)
		}
		if v != w {
			t.Fatalf("NextLine at index %d = %q; want %q", i, v, w)
		}
	}

	_, err := r.NextLine()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Expected EOF error, got: %v", err)
	}
}

func TestPeekByteAndReadByte(t *testing.T) {
	r := newTestReader("abc")

	b1, err := r.PeekByte()
	if err != nil {
		t.Fatalf("PeekByte error: %v", err)
	}
	if b1 != 'a' {
		t.Fatalf("PeekByte = %q; want 'a'", b1)
	}

	b2, err := r.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte error: %v", err)
	}
	if b2 != 'a' {
		t.Fatalf("ReadByte = %q; want after PeekByte 'a'", b2)
	}

	b3, err := r.ReadByte()
	if err != nil {
		t.Fatalf("ReadByte #2 error: %v", err)
	}
	if b3 != 'b' {
		t.Fatalf("ReadByte #2 = %q; want 'b'", b3)
	}

	b4, err := r.PeekByte()
	if err != nil {
		t.Fatalf("PeekByte #3 error: %v", err)
	}
	if b4 != 'c' {
		t.Fatalf("PeekByte #3 = %q; want 'c'", b4)
	}

	_, err = r.ReadByte()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Expected EOF error after reading all bytes, got: %v", err)
	}
}

func TestSkipSpaceAtEOF(t *testing.T) {
	r := newTestReader("   \n\t  ")
	err := r.SkipSpaces()
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("Expected EOF or nil error after skipping spaces, got: %v", err)
	}
}
