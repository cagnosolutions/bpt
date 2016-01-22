package bptx

type tree struct {
	root *node
}

func NewTree() *tree {
	return &tree{}
}

func (t *tree) Set(k, v []byte) {
	t.root = insert(t.root, k, v)
}

func (t *tree) Get(k []byte) []byte {
	r := find(t.root, k)
	return r.data
}

func (t *tree) GetAll() [][]byte {
	r := find_all(t.root)
	var data [][]byte
	if r != nil {
		for _, v := range r {
			data = append(data, v.data)
		}
	}
	return data
}

func (t *tree) Del(k []byte) {
	t.root = delete(t.root, k)
}

func (t *tree) Close() {
	destroy_tree(t.root)
}
