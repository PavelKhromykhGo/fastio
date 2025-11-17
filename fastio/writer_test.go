package fastio

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriteIntAndLine(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.WriteInt(123); err != nil {
		t.Fatalf("WriteInt failed: %v", err)
	}
	if err := w.WriteByte(' '); err != nil {
		t.Fatalf("WriteByte failed: %v", err)
	}
	if err := w.WriteInt(-456); err != nil {
		t.Fatalf("WriteInt failed: %v", err)
	}
	if err := w.WriteByte('\n'); err != nil {
		t.Fatalf("WriteByte failed: %v", err)
	}
	if err := w.WriteLine("test"); err != nil {
		t.Fatalf("WriteLine failed: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	got := buf.String()
	want := "123 -456\ntest\n"
	if got != want {
		t.Errorf("Output mismatch: got %q, want %q", got, want)
	}
}

func TestWriteStringAndBytes(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.WriteString("hello"); err != nil {
		t.Fatalf("WriteString failed: %v", err)
	}
	if err := w.WriteBytes([]byte(" world")); err != nil {
		t.Fatalf("WriteBytes failed: %v", err)
	}
	if err := w.WriteByte('\n'); err != nil {
		t.Fatalf("WriteByte failed: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	got := buf.String()
	want := "hello world\n"
	if got != want {
		t.Errorf("Output mismatch: got %q, want %q", got, want)
	}
}

func TestWriteFloat64(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	if err := w.WriteFloat64(3.14159, 2); err != nil {
		t.Fatalf("WriteFloat64 failed: %v", err)
	}
	if err := w.WriteByte('\n'); err != nil {
		t.Fatalf("WriteByte failed: %v", err)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	got := buf.String()
	want := "3.14\n"
	if got != want {
		t.Errorf("Output mismatch: got %q, want %q", got, want)
	}
}

type errorWriter struct{}

func (e *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}

func TestWriteErrorPropagation(t *testing.T) {
	w := NewWriter(&errorWriter{})

	err := w.WriteString("test")
	if !errors.Is(err, errors.New("write error")) {
		t.Fatalf("Expected error, got: %v", err)
	}
	if w.Err() == nil {
		t.Fatalf("Expected writer to have error state")
	}

	err2 := w.WriteByte('x')
	if !errors.Is(err2, errors.New("write error")) {
		t.Fatalf("Expected error on subsequent write, got: %v", err2)
	}
}

func TestWriterImplementsIoWriter(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	data := []byte("hello io.Writer")
	n, err := w.Write(data)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	if n != len(data) {
		t.Fatalf("Expected to write %d bytes, wrote %d", len(data), n)
	}
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	got := buf.String()
	want := "hello io.Writer"
	if got != want {
		t.Errorf("Output mismatch: got %q, want %q", got, want)
	}
}

type writerFunc func(p []byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

func TestFlushEmptyDoesNotWrite(t *testing.T) {
	var called bool
	w := NewWriter(writerFunc(func(p []byte) (int, error) {
		called = true
		return len(p), nil
	}))
	if err := w.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}
	if called {
		t.Errorf("Expected Flush on empty buffer to not call underlying Write")
	}
}
