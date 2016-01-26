package bpt

import (
	"bytes"
	"log"
	"sort"
	"sync"
)

const order = 64

type record struct {
	value []byte
}

type node struct {
	num_keys int
	keys     [order - 1][]byte
	ptrs     [order]interface{}
	parent   *node
	is_leaf  bool
}

func find_all_records(root *node) []*record {
	// node is empty, can't do much
	if root == nil {
		return nil
	}
	var c *node = root
	// traverse down left side of tree
	// until we find the first leaf node
	for !c.is_leaf {
		c = c.ptrs[0].(*node)
	}
	var i int
	var r []*record
	for {
		// found first leaf. enumerate leaf node's
		// ptrs sending each data record onto channel.
		for i = 0; i < c.num_keys; i++ {
			if c.ptrs[i] != nil {
				r = append(r, c.ptrs[i].(*record))
			}
		}
		// finally, utilize last ptr of leaf node
		// to jump to the next leaf. continue until
		// all leaves have been enumerated
		if c.ptrs[order-1] != nil {
			c = c.ptrs[order-1].(*node)
		} else {
			break // last leaf, stop enumerating
		}
	}
	return r
}

func find_leaf(n *node, key []byte) *node {
	if n == nil {
		return n
	}
	for !n.is_leaf {
		i := sort.Search(n.num_keys, func(i int) bool {
			return bytes.Compare(key, n.keys[i]) >= 0
		})
		n = n.ptrs[i].(*node)
	}
	return n
}

// find leaf type node for a given key
func find_leafX(root *node, key []byte) *node {
	var c *node = root
	if c == nil {
		return c
	}
	var i int
	for !c.is_leaf {
		i = 0
		for i < c.num_keys {
			if bytes.Compare(key, c.keys[i]) >= 0 { // TODO: SLOW...
				i++
			} else {
				break
			}
		}
		c = c.ptrs[i].(*node) // TODO: SLOW...
	}
	return c
}

// find first leaf
func find_first_leaf(root *node) *node {
	if root == nil {
		return root
	}
	var c *node = root
	for !c.is_leaf {
		c = c.ptrs[0].(*node)
	}
	return c
}

func find(n *node, key []byte) *record {
	l := find_leaf(n, key)
	if l == nil {
		return nil
	}
	i := sort.Search(l.num_keys, func(i int) bool {
		return bytes.Equal(l.keys[i], key)
	})
	if i == l.num_keys {
		return nil
	}
	return l.ptrs[i].(*record)
}

// find record for a given key
func findX(root *node, key []byte) *record {
	var c *node = find_leaf(root, key)
	if c == nil {
		return nil
	}
	var i int
	for i = 0; i < c.num_keys; i++ {
		if len(c.keys[i]) == len(key) { // ADDED
			if bytes.Equal(c.keys[i], key) { // TODO: SLOW...
				break
			}
		}
	}
	if i == c.num_keys {
		return nil
	}
	return c.ptrs[i].(*record)
}

// find split point of full node
func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}
	return length/2 + 1
}

// create a record to hold a value for a given key
func make_record(val []byte) *record {
	return &record{
		value: val,
	}
}

/*
// create a new general node... can be adapted
// to serve as either an internal or leaf node
func make_node() *node {
	return &node{
		keys:     [order-1][]byte,
		ptrs:     [order]interface{},
	}
}
*/

/*
// create a new leaf node by addapting a general node
func make_leaf() *node {
	var leaf *node
	leaf = make_node()
	leaf.is_leaf = true
	return leaf
}
*/

// helper->insert_into_parent
// used to find index of the parent's ptr to the
// node to the left of the key to be inserted
func get_left_index(parent, left *node) int {
	left_index := 0
	for left_index <= parent.num_keys && parent.ptrs[left_index] != left {
		left_index++
	}
	return left_index
}

// inserts a new key and *record into a leaf, then returns leaf
func insert_into_leaf(leaf *node, key []byte, ptr *record) *node {
	var insertion_point int
	for insertion_point < leaf.num_keys && bytes.Compare(leaf.keys[insertion_point], key) == -1 {
		insertion_point++
	}
	for i := leaf.num_keys; i > insertion_point; i-- { //TODO: SLOW...
		leaf.keys[i] = leaf.keys[i-1] // TODO: SLOW...
		leaf.ptrs[i] = leaf.ptrs[i-1] // TODO: SLOW...
	}
	leaf.keys[insertion_point] = key
	leaf.ptrs[insertion_point] = ptr
	leaf.num_keys++
	return leaf
}

// inserts a new key and *record into a leaf, so as
// to exceed the order, causing the leaf to be split
func insert_into_leaf_after_splitting(root, leaf *node, key []byte, ptr *record) *node {

	insertion_index := (0 + leaf.num_keys + 1) >> 1

	var tmp_keys [order][]byte
	var tmp_ptrs [order]interface{}

	for i, j := 0, 0; i < leaf.num_keys; i, j = i+1, j+1 {
		if j == insertion_index {
			j++
		}
		tmp_keys[j] = leaf.keys[i]
		tmp_ptrs[j] = leaf.ptrs[i]
	}

	//var tmp_keys [order][]byte
	//var tmp_ptrs [order]interface{}
	//var insertion_index, split, i, j int

	//var new_key []byte

	/*
		for insertion_index < order-1 && bytes.Compare(leaf.keys[insertion_index], key) == -1 {
			insertion_index++
		}
		for i < leaf.num_keys {
			if j == insertion_index {
				j++
			}
			tmp_keys[j] = leaf.keys[i]
			tmp_ptrs[j] = leaf.ptrs[i]
			i++
			j++
		}
	*/

	tmp_keys[insertion_index] = key
	tmp_ptrs[insertion_index] = ptr
	leaf.num_keys = 0

	split := cut(order - 1)

	for i := 0; i < split; i++ {
		leaf.ptrs[i] = tmp_ptrs[i]
		leaf.keys[i] = tmp_keys[i]
		leaf.num_keys++
	}

	new_leaf := &node{is_leaf: true} //make_leaf() // TODO: LOTS OF RAM USED...

	for i, j := split, 0; i < order; i, j = i+1, j+1 {
		new_leaf.ptrs[j] = tmp_ptrs[i]
		new_leaf.keys[j] = tmp_keys[i]
		new_leaf.num_keys++
	}
	// freeing tmps...
	for i := 0; i < order; i++ {
		tmp_ptrs[i] = nil
		tmp_keys[i] = nil
	}
	new_leaf.ptrs[order-1] = leaf.ptrs[order-1]
	leaf.ptrs[order-1] = new_leaf
	for i := leaf.num_keys; i < order-1; i++ {
		leaf.ptrs[i] = nil
	}
	for i := new_leaf.num_keys; i < order-1; i++ {
		new_leaf.ptrs[i] = nil
	}
	new_leaf.parent = leaf.parent
	return insert_into_parent(root, leaf, new_leaf.keys[0], new_leaf)
	//new_key := new_leaf.keys[0]
	//return insert_into_parent(root, leaf, new_key, new_leaf)
}

// insert a new key, ptr to a node
func insert_into_node(root, n *node, left_index int, key []byte, right *node) *node {
	var i int
	for i = n.num_keys; i > left_index; i-- {
		n.ptrs[i+1] = n.ptrs[i]
		n.keys[i] = n.keys[i-1]
	}
	n.ptrs[left_index+1] = right
	n.keys[left_index] = key
	n.num_keys++
	return root
}

// insert a new key, ptr to a node causing node to split
func insert_into_node_after_splitting(root, old_node *node, left_index int, key []byte, right *node) *node {
	var i, j, split int
	var new_node, child *node
	var tmp_keys [order][]byte
	var tmp_ptrs [order + 1]interface{}
	var k_prime []byte

	for i < old_node.num_keys+1 {
		if j == left_index+1 {
			j++
		}
		tmp_ptrs[j] = old_node.ptrs[i]
		i++
		j++
	}

	i = 0
	j = 0

	for i < old_node.num_keys {
		if j == left_index {
			j++
		}
		tmp_keys[j] = old_node.keys[i]
		i++
		j++
	}

	tmp_ptrs[left_index+1] = right
	tmp_keys[left_index] = key

	split = cut(order)
	//new_node = make_node()
	new_node = &node{}
	old_node.num_keys = 0

	for i = 0; i < split-1; i++ {
		old_node.ptrs[i] = tmp_ptrs[i]
		old_node.keys[i] = tmp_keys[i]
		old_node.num_keys++
	}

	old_node.ptrs[i] = tmp_ptrs[i]
	k_prime = tmp_keys[split-1]

	j = 0
	for i += 1; i < order; i++ {
		new_node.ptrs[j] = tmp_ptrs[i]
		new_node.keys[j] = tmp_keys[i]
		new_node.num_keys++
		j++
	}

	new_node.ptrs[j] = tmp_ptrs[i]

	// free tmps...
	for i = 0; i < order; i++ {
		tmp_keys[i] = nil
		tmp_ptrs[i] = nil
	}
	tmp_ptrs[order] = nil

	new_node.parent = old_node.parent

	for i = 0; i <= new_node.num_keys; i++ {
		child = new_node.ptrs[i].(*node)
		child.parent = new_node
	}
	return insert_into_parent(root, old_node, k_prime, new_node)
}

// insert a new node (leaf or internal) into tree, return root of tree
func insert_into_parent(root, left *node, key []byte, right *node) *node {
	var left_index int
	var parent *node
	parent = left.parent
	if parent == nil {
		return insert_into_new_root(left, key, right)
	}
	left_index = get_left_index(parent, left)
	if parent.num_keys < order-1 {
		return insert_into_node(root, parent, left_index, key, right)
	}
	return insert_into_node_after_splitting(root, parent, left_index, key, right)
}

// creates a new root for two sub-trees and inserts the key into the new root
func insert_into_new_root(left *node, key []byte, right *node) *node {
	//var root *node = make_node()
	root := &node{}
	root.keys[0] = key
	root.ptrs[0] = left
	root.ptrs[1] = right
	root.num_keys++
	//root.parent = nil
	left.parent = root
	right.parent = root
	return root
}

// first insertion, start a new tree
func start_new_tree(key []byte, ptr *record) *node {
	//var root *node = make_leaf()
	root := &node{is_leaf: true}
	root.keys[0] = key
	root.ptrs[0] = ptr
	root.ptrs[order-1] = nil
	//root.parent = nil
	root.num_keys++
	return root
}

var recordPool = sync.Pool{New: func() interface{} { return &record{} }}

// master insert function
func insert(root *node, key []byte, val []byte) *node {

	if find(root, key) != nil {
		return root
	}

	//ptr = make_record(val) // TODO: LOTS OF RAM USED...

	//ptr := &record{value: val}

	if root == nil {
		return start_new_tree(key, &record{val})
	}

	leaf := find_leaf(root, key)
	if leaf.num_keys < order-1 {
		insert_into_leaf(leaf, key, &record{val})
		return root
	}

	return insert_into_leaf_after_splitting(root, leaf, key, &record{val})
}

// helper for delete methods... returns index of
// a nodes nearest sibling to the left if one exists
func get_neighbor_index(n *node) int {
	var i int
	for i = 0; i <= n.parent.num_keys; i++ {
		if n.parent.ptrs[i] == n {
			return i - 1
		}
	}
	log.Fatalf("Search for nonexistent ptr to node in parent.\nNode: %p\n", n)
	return 1
}

func remove_entry_from_node(n *node, key []byte, ptr interface{}) *node {
	var i, num_ptrs int
	// remove key and shift over keys accordingly
	for !bytes.Equal(n.keys[i], key) {
		i++
	}
	for i += 1; i < n.num_keys; i++ {
		n.keys[i-1] = n.keys[i]
	}
	// remove ptr and shift other ptrs accordingly
	// first determine the number of ptrs
	if n.is_leaf {
		num_ptrs = n.num_keys
	} else {
		num_ptrs = n.num_keys + 1
	}
	i = 0
	for n.ptrs[i] != ptr {
		i++
	}

	//for n.ptrs[i].(*node) != ptr {
	//	i++
	//}
	for i += 1; i < num_ptrs; i++ {
		n.ptrs[i-1] = n.ptrs[i]
	}
	// one key has been removed
	n.num_keys--
	// set other ptrs to nil for tidiness; remember leaf
	// nodes use the last ptr to point to the next leaf
	if n.is_leaf {
		for i = n.num_keys; i < order-1; i++ {
			n.ptrs[i] = nil
		}
	} else {
		for i = n.num_keys + 1; i < order; i++ {
			n.ptrs[i] = nil
		}
	}
	return n
}

func adjust_root(root *node) *node {
	// if non-empty root key and ptr
	// have already been deleted, so
	// nothing to be done here
	if root.num_keys > 0 {
		return root
	}
	var new_root *node
	// if root is empty and has a child
	// promote first (only) child as the
	// new root node. If it's a leaf then
	// the while tree is empty...
	if !root.is_leaf {
		new_root = root.ptrs[0].(*node)
		new_root.parent = nil
	} else {
		new_root = nil
	}
	root = nil // free root
	return new_root
}

// merge (underflow)
func coalesce_nodes(root, n, neighbor *node, neighbor_index int, k_prime []byte) *node {
	var i, j, neighbor_insertion_index, n_end int
	var tmp *node
	// swap neight with node if nod eis on the
	// extreme left and neighbor is to its right
	if neighbor_index == -1 {
		tmp = n
		n = neighbor
		neighbor = tmp
	}
	// starting index for merged pointers
	neighbor_insertion_index = neighbor.num_keys
	// case nonleaf node, append k_prime and the following ptr.
	// append all ptrs and keys for the neighbors.
	if !n.is_leaf {
		// append k_prime (key)
		neighbor.keys[neighbor_insertion_index] = k_prime
		neighbor.num_keys++
		n_end = n.num_keys
		i = neighbor_insertion_index + 1
		j = 0
		for j < n_end {
			neighbor.keys[i] = n.keys[j]
			neighbor.ptrs[i] = n.ptrs[j]
			neighbor.num_keys++
			n.num_keys--
			i++
			j++
		}
		neighbor.ptrs[i] = n.ptrs[j]
		for i = 0; i < neighbor.num_keys+1; i++ {
			tmp = neighbor.ptrs[i].(*node)
			tmp.parent = neighbor
		}
	} else {
		// in a leaf; append the keys and ptrs.
		i = neighbor_insertion_index
		j = 0
		for j < n.num_keys {
			neighbor.keys[i] = n.keys[j]
			neighbor.ptrs[i] = n.ptrs[j]
			neighbor.num_keys++
			i++
			j++
		}
		neighbor.ptrs[order-1] = n.ptrs[order-1]
	}
	root = delete_entry(root, n.parent, k_prime, n)
	//n.keys = nil
	//n.ptrs = nil
	n = nil // free n
	return root
}

// merge / redistribute
func redistribute_nodes(root, n, neighbor *node, neighbor_index, k_prime_index int, k_prime []byte) *node {
	var i int
	var tmp *node

	// case: node n has a neighnor to the left
	if neighbor_index != -1 {
		if !n.is_leaf {
			n.ptrs[n.num_keys+1] = n.ptrs[n.num_keys]
		}
		for i = n.num_keys; i > 0; i-- {
			n.keys[i] = n.keys[i-1]
			n.ptrs[i] = n.ptrs[i-1]
		}
		if !n.is_leaf {
			n.ptrs[0] = neighbor.ptrs[neighbor.num_keys]
			tmp = n.ptrs[0].(*node)
			tmp.parent = n
			neighbor.ptrs[neighbor.num_keys] = nil
			n.keys[0] = k_prime
			n.parent.keys[k_prime_index] = neighbor.keys[neighbor.num_keys-1]
		} else {
			n.ptrs[0] = neighbor.ptrs[neighbor.num_keys-1]
			neighbor.ptrs[neighbor.num_keys-1] = nil
			n.keys[0] = neighbor.keys[neighbor.num_keys-1]
			n.parent.keys[k_prime_index] = n.keys[0]
		}
	} else {
		// case: n is left most child (n has no left neighbor)
		if n.is_leaf {
			n.keys[n.num_keys] = neighbor.keys[0]
			n.ptrs[n.num_keys] = neighbor.ptrs[0]
			n.parent.keys[k_prime_index] = neighbor.keys[1]
		} else {
			n.keys[n.num_keys] = k_prime
			n.ptrs[n.num_keys+1] = neighbor.ptrs[0]
			tmp = n.ptrs[n.num_keys+1].(*node)
			tmp.parent = n
			n.parent.keys[k_prime_index] = neighbor.keys[0]
		}
		for i = 0; i < neighbor.num_keys-1; i++ {
			neighbor.keys[i] = neighbor.keys[i+1]
			neighbor.ptrs[i] = neighbor.ptrs[i+1]
		}
		if !n.is_leaf {
			neighbor.ptrs[i] = neighbor.ptrs[i+1]
		}
	}
	n.num_keys++
	neighbor.num_keys--
	return root
}

// deletes an entry from the tree; removes record, key, and ptr from leaf and rebalances tree
func delete_entry(root, n *node, key []byte, ptr interface{}) *node {
	var min_keys, neighbor_index, k_prime_index, capacity int
	var neighbor *node
	var k_prime []byte

	// remove key, ptr from node
	n = remove_entry_from_node(n, key, ptr)
	//switch ptr.(type) {
	//case *node:
	//	n = remove_entry_from_node(n, key, ptr.(*node))
	//case *record:
	//	n = remove_entry_from_node(n, key, ptr.(*record))
	//}
	if n == root {
		return adjust_root(root)
	}

	// case: delete from inner node
	if n.is_leaf {
		min_keys = cut(order - 1)
	} else {
		min_keys = cut(order) - 1
	}
	// case: node stays at or above min order
	if n.num_keys >= min_keys {
		return root
	}

	// case: node is bellow min order; coalescence or redistribute
	neighbor_index = get_neighbor_index(n)
	if neighbor_index == -1 {
		k_prime_index = 0
	} else {
		k_prime_index = neighbor_index
	}
	k_prime = n.parent.keys[k_prime_index]
	if neighbor_index == -1 {
		neighbor = n.parent.ptrs[1].(*node)
	} else {
		neighbor = n.parent.ptrs[neighbor_index].(*node)
	}
	if n.is_leaf {
		capacity = order
	} else {
		capacity = order - 1
	}

	// coalescence
	if neighbor.num_keys+n.num_keys < capacity {
		return coalesce_nodes(root, n, neighbor, neighbor_index, k_prime)
	}
	return redistribute_nodes(root, n, neighbor, neighbor_index, k_prime_index, k_prime)
}

// master delete
func delete(root *node, key []byte) *node {
	var key_leaf *node
	var key_record *record

	key_record = find(root, key)
	key_leaf = find_leaf(root, key)
	if key_record != nil && key_leaf != nil {
		root = delete_entry(root, key_leaf, key, key_record)
		key_record = nil // free key_record
	}
	return root
}

func destroy_tree_nodes(root *node) {
	var i int
	if root.is_leaf {
		for i = 0; i < root.num_keys; i++ {
			root.ptrs[i] = nil
		}
	} else {
		for i = 0; i < root.num_keys+1; i++ {
			destroy_tree_nodes(root.ptrs[i].(*node))
		}
	}
	//root.ptrs = nil // free
	//root.keys = nil // free
	root = nil // free
}

func destroy_tree(root *node) {
	destroy_tree_nodes(root)
}
