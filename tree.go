package bpt

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

/*
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
*/

func (t *tree) Del(key []byte) {
	t.root = delete(t.root, key)
}

func (t *tree) Close() {
	destroy_tree(t.root)
}
