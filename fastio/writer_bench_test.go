package fastio

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"
)

const benchNumCountWriter = 10000

func BenchmarkFastWriter_WriteInt(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		fw := NewWriter(io.Discard)

		for j := 0; j < benchNumCountWriter; j++ {
			if err := fw.WriteInt(j); err != nil {
				b.Fatalf("WriteInt error: %v", err)
			}
			if err := fw.WriteByte('\n'); err != nil {
				b.Fatalf("WriteByte error: %v", err)
			}
		}
		if err := fw.Flush(); err != nil {
			b.Fatalf("Flush error: %v", err)
		}
	}
}

func BenchmarkFmtFprintWithBufio(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		bw := bufio.NewWriter(&buf)

		for j := 0; j < benchNumCountWriter; j++ {
			if _, err := fmt.Fprintln(bw, j); err != nil {
				b.Fatalf("Fprintln error: %v", err)
			}
		}
		if err := bw.Flush(); err != nil {
			b.Fatalf("Flush error: %v", err)
		}
	}
}
