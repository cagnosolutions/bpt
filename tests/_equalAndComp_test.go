package main

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkBuiltinEquals(bench *testing.B) {
	var a, b [][]byte
	for i := 0; i < bench.N; i++ {
		a, b = append(a, []byte(fmt.Sprintf("%d", i))), append(b, []byte(fmt.Sprintf("%d", i)))
		if ok := bytes.Equal(a[i], b[i]); !ok {
			bench.Logf("Equal(a, b): Expected [%s] and [%s] to be true, got %v\n", a[i], b[i], ok)
		}
	}
}

func BenchmarkBuiltinCompare(bench *testing.B) {
	var c, d [][]byte
	for i := 0; i < bench.N; i++ {
		c, d = append(c, []byte(fmt.Sprintf("%d%d", i, i))), append(d, []byte(fmt.Sprintf("%d", i)))
		if n := bytes.Compare(c[i], d[i]); n < 0 {
			bench.Logf("Compare(c, d): Expected [%s] to be greater or equal to [%s], got %d instead\n", c[i], d[i], n)
		}
	}
}

func BenchmarkCustomEquals(bench *testing.B) {
	var a, b [][]byte
	for i := 0; i < bench.N; i++ {
		a, b = append(a, []byte(fmt.Sprintf("%d", i))), append(b, []byte(fmt.Sprintf("%d", i)))
		if ok := customEqual(a[i], b[i]); !ok {
			bench.Logf("customEqual(a, b): Expected [%s] and [%s] to be true, got %v\n", a[i], b[i], ok)
		}
	}
}

/*
func BenchmarkCustomCompare(bench *testing.B) {
	var c, d [][]byte
	for i := 0; i < bench.N; i++ {
		c, d = append(c, []byte(fmt.Sprintf("%d%d", i, i))), append(d, []byte(fmt.Sprintf("%d", i)))
		if n := customCompare(c[i], d[i]); n < 0 {
			bench.Logf("customCompare(c, d): Expected [%s] to be greater or equal to [%s], got %d instead\n", c[i], d[i], n)
		}
	}
}
*/

func customEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, j := 0, 1; j < len(a); i, j = i+1, j+1 {
		if a[i] != b[i] && a[j] != b[j] {
			return false
		}
	}
	return true
}

/*
func customCompare(a, b []byte) int {

}
*/
