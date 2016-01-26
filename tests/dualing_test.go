package main

import (
	"bytes"
	"testing"

	"github.com/cagnosolutions/bplus"
	"github.com/cagnosolutions/bpt"
)

func BenchmarkBPlusTree(b *testing.B) {
	t := bplus.NewTree(bytes.Compare)
	for i := 0; i < b.N; i++ {
		x := bpt.UUID()
		t.Set(x, x)
	}
}

func BenchmarkBPT(b *testing.B) {
	t := bpt.NewTree()
	for i := 0; i < b.N; i++ {
		x := bpt.UUID()
		t.Set(x, x)
	}
}

/* NOTE: NO IMPROVEMENTS WITH DIYComp
func BenchmarkBytesCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := bpt.UUID()
		if n := bytes.Compare(x, x); n != 0 {
			b.Logf("Expected n to be 0, got %d instead...\n", n)
		}
	}
}

func BenchmarkDIYComp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := bpt.UUID()
		if n := comp(x, x); !n {
			b.Logf("Expected n to be true, got %v instead...\n", n)
		}
	}
}

func comp(a, b []byte) bool {
	la, lb := len(a), len(b)
	x := lb + ((la - lb) & ((la - lb) >> 0xf7)) - 1
	for l, m, h := 0, x>>1, x; l <= m && m <= h; l, h = l+1, h-1 {
		if (a[l]&0xff) < (b[l]&0xff) || (a[h]&0xff) < (b[h]&0xff) {
			return false
		}
	}
	return true
}
*/
