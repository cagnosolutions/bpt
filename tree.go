package main

//package bpt

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type tree struct {
	root *node
}

func NewTree() *tree {
	return &tree{}
}

func (t *tree) Set(key []byte, val []byte) {
	t.root = insert(t.root, key, val)
}

func (t *tree) Get(key []byte) []byte {
	r := find(t.root, key)
	return r.value
}

func (t *tree) GetAll() [][]byte {
	r := find_all_records(t.root)
	var data [][]byte
	if r != nil {
		for _, v := range r {
			data = append(data, v.value)
		}
	}
	return data
}

func (t *tree) Del(key []byte) {
	t.root = delete(t.root, key)
}

func (t *tree) Close() {
	destroy_tree(t.root)
}

var COUNT = 50000

func main() {

	t := NewTree()

	log.Printf("Inserting %d key/value pairs into tree...\n", COUNT)
	for i := 0; i < COUNT; i++ {
		uuid := UUID()
		t.Set(uuid, uuid)
	}

	log.Printf("Enumerating all leaf nodes; finding all records...\n")
	d := t.GetAll()
	for i, e := range d {
		fmt.Printf("Record #%d: [% x]\n", i, e)
	}
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
