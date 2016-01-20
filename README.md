# bpt
B+Tree re-write / transpile from C to Go


# RETURN FUNCTIONS
------------------------------------------------------------------------------------
- find(root *node) *record
	- find_leaf(root, key)

- find_leaf(root *node, key []byte)

- find_all_records(root *node) []*record


# INSERTION FUNCTIONS
------------------------------------------------------------------------------------
- insert(root *node, key []byte, val []byte) *node
	- find(root, key)
	- make_record(val) 											*[if find was false]
	- start_new_tree(key, ptr) 									*[if root == nil]
	- find_leaf(root, key) 										*[if root != nil]
	- insert_into_leaf(leaf, key, ptr) 							*[if leaf is not full]
	- insert_into_leaf_after_splitting(root, leaf, key ptr) 	*[if leaf is full]

- find(root *node, key []byte) *record
	- find_leaf(root, key)

- start_new_tree(key []byte, ptr *record)
	- make_leaf()
		- make_node()

- insert_into_leaf_after_splitting(root, leaf *node, key []byte, ptr *record) *node
	- make_leaf()
	- cut(ORDER - 1)
	- insert_into_parent(root, leaf, new_key, new_leaf)

- make_leaf()
	- make_node()

- insert_into_parent(root, left *node, key []byte, right *node) *node
	- insert_into_new_root(left, key, right) 									*[if parent == nil]
	- get_left_index(parent, left)												*[if parent != nil]
	- insert_into_node(root, parent, left_index, key, right) 					*[if parent is not full]
	- insert_into_node_after_splitting(root, parent, left_index, key, right) 	*[if parent is full]

- insert_into_new_root(left *node, key []byte, right *node)
	- make_node()

- insert_into_node_after_splitting(root, old_node *node, left_index int, key []byte, right *node) *node
	- cut(ORDER - 1)
	- make_node()
	- insert_into_parent(root, old_node, k_prime, new_node)


# DELETION FUNCTIONS
-------------------------------------------------------------------------------------

