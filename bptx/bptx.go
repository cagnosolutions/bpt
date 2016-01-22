package bptx

import (
	"bytes"
	"log"
)

const order = 64

func comp(k1, k2 []byte) bool {
	return bytes.Compare(k1, k2) >= 0
}

func equal(k1, k2 []byte) bool {
	return bytes.Equal(k1, k2)
}

// find split point of full node
func cut(length int) int {
	if length%2 == 0 {
		return length / 2
	}
	return length/2 + 1
}

// data record
type rec struct {
	data []byte
}

// create a rec to hold a data for a given key
func make_rec(v []byte) *rec {
	return &rec{
		data: v,
	}
}

// tree node
type node struct {
	num_keys int
	keys     [][]byte
	ptrs     []interface{}
	parent   *node
	next     *node
	is_leaf  bool
}

// create a new general node... can be adapted
// to serve as either an internal or leaf node
func make_node() *node {
	return &node{
		keys: make([][]byte, order-1),
		ptrs: make([]interface{}, order),
	}
}

func find_all(n *node) []*rec {
	if n == nil {
		return nil
	}
	for !n.is_leaf {
		n = n.ptrs[0].(*node)
	}
	var r []*rec
	for {
		for i := 0; i < n.num_keys && n.ptrs[i] != nil; i++ {
			r = append(r, n.ptrs[i].(*rec))
		}
		if n.ptrs[order-1] != nil {
			n = n.ptrs[order-1].(*node)
		} else {
			break
		}
	}
	return r
}

func find_leaf(n *node, k []byte) *node {
	if n == nil {
		return n
	}
	var i int
	for !n.is_leaf {
		i = 0
		for i < n.num_keys {
			if comp(k, n.keys[i]) {
				i++
			} else {
				break
			}
		}
		n = n.ptrs[i].(*node)
	}
	return n
}

func find(n *node, k []byte) *rec {
	var l *node
	l = find_leaf(n, k)
	if l == nil {
		return nil
	}
	var i int
	for i < l.num_keys && !equal(l.keys[i], k) {
		i++
	}
	if i < l.num_keys {
		return l.ptrs[i].(*rec)
	}
	return nil
}

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

// inserts a new key and *rec into a leaf, then returns leaf
func insert_into_leaf(leaf *node, k []byte, ptr *rec) *node {
	var i, insertion_point int
	for insertion_point < leaf.num_keys && !comp(leaf.keys[insertion_point], k) {
		insertion_point++
	}
	for i = leaf.num_keys; i > insertion_point; i-- {
		leaf.keys[i] = leaf.keys[i-1]
		leaf.ptrs[i] = leaf.ptrs[i-1]
	}
	leaf.keys[insertion_point] = k
	leaf.ptrs[insertion_point] = ptr
	leaf.num_keys++
	return leaf
}

// inserts a new key and *rec into a leaf, so as
// to exceed the order, causing the leaf to be split
func insert_into_leaf_after_splitting(root, leaf *node, k []byte, ptr *rec) *node {

	var new_leaf *node
	var tmp_keys [order][]byte
	var tmp_ptrs [order]interface{}
	var insertion_index, split, i, j int
	var new_key []byte

	new_leaf = make_node()
	new_leaf.is_leaf = true

	for insertion_index < order-1 && !comp(leaf.keys[insertion_index], k) {
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
	tmp_keys[insertion_index] = k
	tmp_ptrs[insertion_index] = ptr

	leaf.num_keys = 0

	split = cut(order - 1)

	for i = 0; i < split; i++ {
		leaf.ptrs[i] = tmp_ptrs[i]
		leaf.keys[i] = tmp_keys[i]
		leaf.num_keys++
	}

	j = 0
	for i = split; i < order; i++ {
		new_leaf.ptrs[j] = tmp_ptrs[i]
		new_leaf.keys[j] = tmp_keys[i]
		new_leaf.num_keys++
		j++
	}

	// freeing tmps...
	for i = 0; i < order; i++ {
		tmp_ptrs[i] = nil
		tmp_keys[i] = nil
	}

	new_leaf.ptrs[order-1] = leaf.ptrs[order-1]
	leaf.ptrs[order-1] = new_leaf

	for i = leaf.num_keys; i < order-1; i++ {
		leaf.ptrs[i] = nil
	}

	for i = new_leaf.num_keys; i < order-1; i++ {
		new_leaf.ptrs[i] = nil
	}

	new_leaf.parent = leaf.parent
	new_key = new_leaf.keys[0]

	return insert_into_parent(root, leaf, new_key, new_leaf)
}

// insert a new key, ptr to a node
func insert_into_node(root, n *node, left_index int, k []byte, right *node) *node {
	var i int
	for i = n.num_keys; i > left_index; i-- {
		n.ptrs[i+1] = n.ptrs[i]
		n.keys[i] = n.keys[i-1]
	}
	n.ptrs[left_index+1] = right
	n.keys[left_index] = k
	n.num_keys++
	return root
}

// insert a new key, ptr to a node causing node to split
func insert_into_node_after_splitting(root, old_node *node, left_index int, k []byte, right *node) *node {
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
	tmp_keys[left_index] = k

	split = cut(order)
	new_node = make_node()
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
func insert_into_parent(root, left *node, k []byte, right *node) *node {
	var left_index int
	var parent *node
	parent = left.parent
	if parent == nil {
		return insert_into_new_root(left, k, right)
	}
	left_index = get_left_index(parent, left)
	if parent.num_keys < order-1 {
		return insert_into_node(root, parent, left_index, k, right)
	}
	return insert_into_node_after_splitting(root, parent, left_index, k, right)
}

// creates a new root for two sub-trees and inserts the key into the new root
func insert_into_new_root(left *node, k []byte, right *node) *node {
	n := make_node()
	n.keys[0], n.ptrs[0], n.ptrs[1] = k, left, right
	n.num_keys++
	left.parent = n
	right.parent = n
	return n
}

// first insertion, start a new tree
func start_new_tree(k []byte, ptr *rec) *node {
	r := make_node()
	r.is_leaf = true
	r.keys[0], r.ptrs[0] = k, ptr
	r.ptrs[order-1] = nil
	r.num_keys++
	return r
}

// master insert function
func insert(n *node, k, v []byte) *node {
	if find(n, k) != nil {
		return n
	}
	ptr := make_rec(v)
	if n == nil {
		return start_new_tree(k, ptr)
	}
	leaf := find_leaf(n, k)
	if leaf.num_keys < order-1 {
		leaf = insert_into_leaf(leaf, k, ptr)
		return n
	}
	return insert_into_leaf_after_splitting(n, leaf, k, ptr)
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

func remove_entry_from_node(n *node, k []byte, ptr interface{}) *node {
	var i, num_ptrs int
	// remove key and shift over keys accordingly
	for !equal(n.keys[i], k) {
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
	n.keys = nil
	n.ptrs = nil
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

// deletes an entry from the tree; removes rec, key, and ptr from leaf and rebalances tree
func delete_entry(root, n *node, k []byte, ptr interface{}) *node {
	var min_keys, neighbor_index, k_prime_index, capacity int
	var neighbor *node
	var k_prime []byte

	// remove key, ptr from node
	n = remove_entry_from_node(n, k, ptr)
	//switch ptr.(type) {
	//case *node:
	//	n = remove_entry_from_node(n, key, ptr.(*node))
	//case *rec:
	//	n = remove_entry_from_node(n, key, ptr.(*rec))
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
func delete(root *node, k []byte) *node {
	var key_leaf *node
	var key_rec *rec

	key_rec = find(root, k)
	key_leaf = find_leaf(root, k)
	if key_rec != nil && key_leaf != nil {
		root = delete_entry(root, key_leaf, k, key_rec)
		key_rec = nil // free key_rec
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
	root.ptrs = nil // free
	root.keys = nil // free
	root = nil      // free
}

func destroy_tree(root *node) {
	destroy_tree_nodes(root)
}
