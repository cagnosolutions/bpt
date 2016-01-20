package main

import "testing"

func BenchmarkLoop1(b *testing.B) {
	var items = make([]int, 250)
	for j := 0; j < b.N; j++ {
		for i := 0; i < len(items); i++ {
			if items[i] != 0 {
				b.Logf("Wasn't expecting %d\n", items[i])
				break
			}
		}
	}
}

func BenchmarkLoop2(b *testing.B) {
	var items = make([]int, 250)
	var count = len(items)
	var i int
	for j := 0; j < b.N; j++ {
		for i = 0; i < count && items[i] == 0; i++ {
			// stuff...
		}
		if i < count && items[i] != 0 {
			b.Logf("Wasn't expecting %d\n", items[i])
		}
	}
}

/*
func BenchmarkForLoop1(b *testing.B) {
	var items = make([]byte, 25000)
	var n = len(items)
	var j int
	for i := 0; i < b.N; i++ {
		for j = 0; j < n; j++ {
			if items[j] != 0 {
				log.Printf("items[%d] != 0, got %d instead\n", j, items[j])
			}
		}
	}
}

func BenchmarkForLoop2(b *testing.B) {
	var items = make([]byte, 25000)
	var n = len(items)
	var j int
	for i := 0; i < b.N; i++ {
		j = 0
		for j < n {
			if items[j] != 0 {
				log.Printf("items[%d] != 0, got %d instead\n", j, items[j])
			}
			j++
		}
	}
}

func BenchmarkForLoop3(b *testing.B) {
	var items = make([]byte, 25000)
	var n = len(items)
	for i := 0; i < b.N; i++ {
		var j int
		for j = 0; j < n && items[j] != 0; j++ {
			// nothing to do...
		}
		if j < n && items[j] != 0 {
			log.Printf("items[%d] != 0, got %d instead\n", j, items[j])
		}
	}
}
*/

/*
func BenchmarkFindLeafNode(b *testing.B) {
	b.StopTimer()
	var root *node
	uuid := UUID()
	root = insert(root, uuid, uuid)
	for i := 0; i < 100000; i++ {
		kv := UUID()
		root = insert(root, kv, kv)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n := find_leaf(root, uuid); n == nil {
			log.Println("found nil node...")
		}
	}
	b.StopTimer()
}

func BenchmarkFindLeafNodeUpdated(b *testing.B) {
	b.StopTimer()
	var root *node
	uuid := UUID()
	root = insert(root, uuid, uuid)
	for i := 0; i < 100000; i++ {
		kv := UUID()
		root = insert(root, kv, kv)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		if n := find_leaf_updated(root, uuid); n == nil {
			log.Println("found nil node...")
		}
	}
	b.StopTimer()
}

func find_leaf_updated(n *node, key []byte) *node {
	if n == nil {
		return n
	}
	for i := 0; i < n.num_keys && !n.is_leaf && bytes.Compare(key, n.keys[i]) >= 0; i++ {
		n = n.ptrs[i].(*node)
	}
	return n
}

func BenchmarkSwitchOnItems(b *testing.B) {
	for i := 0; i < b.N; i++ {
		switchOnItems(i)
	}
}

func switchOnItems(x int) {
	switch {
	case x > 5:
		break
	case x > 25:
		break
	case x > 50:
		break
	case x > 100:
		break
	case x > 200:
		break
	case x > 400:
		break
	case x > 800:
		break
	case x > 1600:
		break
	case x > 3200:
		break
	case x > 6400:
		break
	case x > 12800:
		break
	default:
		break
	}
	return
}

func BenchmarkIfStmtOnItems(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ifStmtOnItems(i)
	}
}

func ifStmtOnItems(x int) {
	if x > 5 {
		return
	}
	if x > 25 {
		return
	}
	if x > 50 {
		return
	}
	if x > 100 {
		return
	}
	if x > 200 {
		return
	}
	if x > 400 {
		return
	}
	if x > 800 {
		return
	}
	if x > 1600 {
		return
	}
	if x > 3200 {
		return
	}
	if x > 6400 {
		return
	}
	if x > 12800 {
		return
	}
	return
}
*/
