package main

import (
	"sync"
	"testing"
)

func BenchmarkInsert(b *testing.B) {
	t := NewTree()
	for i := 0; i < b.N; i++ {
		d := UUID()
		t.Set(d, d)
	}
}

func BenchmarkMapInsert(b *testing.B) {
	m := struct {
		mm map[string][]byte
		sync.RWMutex
	}{
		mm: map[string][]byte{},
	}
	for i := 0; i < b.N; i++ {
		d := UUID()
		m.mm[string(d)] = d
	}
}

func BenchmarkInsertParallel(b *testing.B) {
	t := NewTree()
	for j := 0; j < b.N; j++ {
		go func() {
			for i := 0; i < 200; i++ {
				d := UUID()
				t.Set(d, d)
			}
		}()
	}
}

func BenchmarkMapInsertParallel(b *testing.B) {
	m := struct {
		mm map[string][]byte
		sync.RWMutex
	}{
		mm: map[string][]byte{},
	}
	for j := 0; j < b.N; j++ {
		go func() {
			for i := 0; i < 200; i++ {
				d := UUID()
				m.Lock()
				m.mm[string(d)] = d
				m.Unlock()
			}
		}()
	}
}
