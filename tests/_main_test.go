package main

import (
	"testing"

	"github.com/cagnosolutions/bpt"
	"github.com/cagnosolutions/bpt/bptx"
)

var size int = 1000
var data = make([][]byte, size)

func init() {
	for i := 0; i < size; i++ {
		data[i] = bpt.UUID()
	}
}

func Benchmark_bpt_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := bpt.NewTree()
		for i := 0; i < size; i++ {
			t.Set(data[i], data[i])
		}
	}
}

func Benchmark_bptx_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t := bptx.NewTree()
		for i := 0; i < size; i++ {
			t.Set(data[i], data[i])
		}
	}
}

/*
func BenchmarkIterate(b *testing.B) {
	tree := btree.TreeNew(cznicCmp)
	for i := 0; i < len(fixture.TestData); i++ {
		tree.Set(fixture.TestData[i].Key, fixture.TestData[i].Value)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// only errors on empty trees; meh api
		iter, err := tree.SeekFirst()
		if err != nil {
			b.Fatalf("tree.SeekFirst: %v", err)
		}
		for {
			k, v, err := iter.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatalf("iter.Next: %v", err)
			}

			_ = k.(fixture.Key)
			_ = v.(fixture.Value)
		}
		iter.Close()
	}
}
*/
