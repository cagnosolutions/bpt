package main

import (
	"log"

	"github.com/cagnosolutions/bpt"
)

var N = 64

func main() {
	t := bpt.NewTree()
	for i := 0; i < N; i++ {
		t.Set([]byte{byte(i)}, []byte{byte(i)})
	}
	log.Printf("GETTING BYTE: %d, %d\n", byte(N/2), t.Get([]byte{byte(N / 2)}))
}
