package main

//package bpt

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

type tree struct {
	root *node
	sync.RWMutex
}

func NewTree() *tree {
	return &tree{}
}

func (t *tree) Set(key []byte, val []byte) {
	t.Lock()
	t.root = insert(t.root, key, val)
	t.Unlock()
}

func (t *tree) Get(key []byte) []byte {
	t.RLock()
	r := find(t.root, key)
	t.RUnlock()
	return r.value
}

func (t *tree) GetAll() [][]byte {
	t.RLock()
	r := find_all_records(t.root)
	t.RUnlock()
	var data [][]byte
	if r != nil {
		for _, v := range r {
			data = append(data, v.value)
		}
	}
	return data
}

func (t *tree) Del(key []byte) {
	t.Lock()
	t.root = delete(t.root, key)
	t.Unlock()
}

func (t *tree) Close() {
	t.Lock()
	destroy_tree(t.root)
	t.Unlock()
}

var COUNT = 50000 / 100

func main() {

	t := NewTree()
	for j := 0; j < 100; j++ {
		go func() {
			log.Printf("Inserting %d key/value pairs into tree...\n", COUNT)
			for i := 0; i < COUNT; i++ {
				uuid := UUID()
				t.Set(uuid, uuid)
			}
		}()
	}

	//log.Printf("Enumerating all leaf nodes; finding all records...\n")
	//d := t.GetAll()
	//for i, e := range d {
	//	fmt.Printf("Record #%d: [% x]\n", i, e)
	//}
	//log.Printf("Found %d records, the last 5 records are...\n[% x]\n[% x]\n[% x]\n[% x]\n[% x]\n", len(d), d[len(d)-5], d[len(d)-4], d[len(d)-3], d[len(d)-2], d[len(d)-1])

	/*
		log.Printf("Attempting to get all key/value paris from tree...\n")
		var r int
		for i := COUNT - 1; i >= 0; i-- {
			k := fmt.Sprintf("%d", i)
			v := t.Get([]byte(k))
			if v != nil && bytes.Equal(v, []byte(k)) {
				r++
			}
		}
		log.Printf("Successfully found %d out of %d keys\n", r, COUNT)

		print_tree(t.root)

		log.Printf("Attempting to delete %d keys from tree\n", COUNT/2)
		for i := 0; i < COUNT/2; i++ {
			k := fmt.Sprintf("%d", i)
			t.Del([]byte(k))
		}
		log.Printf("Successfully removed %d out of %d keys\n", COUNT/2, COUNT)

		print_tree(t.root)
	*/

	fmt.Println("Press any key to exit...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
