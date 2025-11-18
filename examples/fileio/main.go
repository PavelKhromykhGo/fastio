package main

import (
	"log"
	"os"

	"github.com/PavelKhromykhGo/fastio/fastio"
)

func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("open input.txt: %v", err)
	}
	defer in.Close()

	out, err := os.Create("output.txt")
	if err != nil {
		log.Fatalf("create output.txt: %v", err)
	}
	defer out.Close()

	r := fastio.NewReader(in)
	w := fastio.NewWriter(out)
	defer w.Flush()

	n, err := r.NextInt()
	if err != nil {
		log.Fatalf("read N: %v", err)
	}

	sum := 0
	for i := 0; i < n; i++ {
		x, err := r.NextInt()
		if err != nil {
			log.Fatalf("read int #%d: %v", i, err)
		}
		sum += x
	}

	if err := w.WriteInt(sum); err != nil {
		log.Fatalf("write sum: %v", err)
	}
	if err := w.WriteByte('\n'); err != nil {
		log.Fatalf("write newline: %v", err)
	}
}
