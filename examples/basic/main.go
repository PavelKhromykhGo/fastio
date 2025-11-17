package main

import (
	"os"

	"github.com/PavelKhromykhGo/fastio/fastio"
)

func main() {
	fr := fastio.NewReader(os.Stdin)
	fw := fastio.NewWriter(os.Stdout)
	defer fw.Flush()

	n, err := fr.NextInt()
	if err != nil {
		return
	}

	sum := 0
	for i := 0; i < n; i++ {
		x, err := fr.NextInt()
		if err != nil {
			return
		}
		sum += x
	}

	_ = fw.WriteInt(sum)
	_ = fw.WriteByte('\n')
}
