package fastio

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const benchNumCount = 10000

func makeIntInput(count int) []byte {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(i))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func BenchmarkFastReader_NextInt(b *testing.B) {
	data := makeIntInput(benchNumCount)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r := NewReader(bytes.NewReader(data))

		sum := 0
		for j := 0; j < benchNumCount; j++ {
			x, err := r.NextInt()
			if err != nil {
				b.Fatalf("NextInt error: %v", err)
			}
			sum += x
		}
		_ = sum
	}
}

func BenchmarkFmtFscan(b *testing.B) {
	data := makeIntInput(benchNumCount)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		br := bufio.NewReader(bytes.NewReader(data))

		sum := 0
		for j := 0; j < benchNumCount; j++ {
			var x int
			if _, err := fmt.Fscan(br, &x); err != nil {
				b.Fatalf("fmt.Fscan error: %v", err)
			}
			sum += x
		}
		_ = sum
	}
}

func BenchmarkBufioScanner(b *testing.B) {
	data := makeIntInput(benchNumCount)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sc := bufio.NewScanner(bytes.NewReader(data))
		sc.Split(bufio.ScanWords)

		sum := 0
		read := 0
		for sc.Scan() {
			txt := sc.Text()
			x, err := strconv.Atoi(txt)
			if err != nil {
				b.Fatalf("Atoi error: %v", err)
			}
			sum += x
			read++
			if read >= benchNumCount {
				break
			}
		}
		if err := sc.Err(); err != nil {
			b.Fatalf("scanner error: %v", err)
		}
		if read != benchNumCount {
			b.Fatalf("scanner read %d numbers, expected %d", read, benchNumCount)
		}
		_ = sum
	}
}
