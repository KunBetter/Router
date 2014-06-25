// PathTrie
package Router

/*
	All Paths Examples:
	/ids/:id
	/users/:name
	/city/beijing
	/city
	/a/b
	/a/b/c
	/a/b/d
	/a/b/c/d
	/aaaa/bbbb/cccc/dddd
*/

type PathTrieNode struct {
	Key    byte   //哈希值
	isVar  bool   //true -> :id,false -> 5
	Var    string //:id,id
	Flag   bool   //word
	Childs []*PathTrieNode
}

func NewPathTrieNode(key byte) *PathTrieNode {
	node := &PathTrieNode{
		Key:    key,
		isVar:  false,
		Var:    "",
		Flag:   false,
		Childs: []*PathTrieNode{},
	}
	return node
}

type PathTrie struct {
	Root *PathTrieNode
}

func NewPathTrie() *PathTrie {
	trie := &PathTrie{
		Root: NewPathTrieNode(0),
	}
	return trie
}

func (trie *PathTrie) FindInChilds(childs []*PathTrieNode, key byte) (int, bool) {
	clen := len(childs)
	if clen <= 0 {
		return 0, false
	}
	l := 0
	h := clen - 1
	for l <= h {
		m := (l + h) / 2
		if childs[m].Key < key {
			l = m + 1
		} else if childs[m].Key > key {
			h = m - 1
		} else {
			return m, true
		}
	}
	return l, false
}

func (trie *PathTrie) FindLastMatchNode(key []byte) (*PathTrieNode, int) {
	node := trie.Root
	var last *PathTrieNode = nil
	if node == nil {
		return nil, -1
	}
	lastI := 0
	i := 0
	for i = 0; i < len(key); i++ {
		childs := node.Childs
		pos, found := trie.FindInChilds(childs, key[i])
		if !found {
			break
		}
		node = childs[pos]
		last = node
	}
	lastI = i
	return last, lastI
}

func (trie *PathTrie) AddPath(path string) {
	paths := strings.Split(path, "/")
	for i := 0; i < len(paths); i++ {
	}
	key := []byte(token.Text)
	node, lastI := trie.FindLastMatchNode(key)
	if node == nil {
		node = trie.Root
		lastI = 0
	}
	for i := lastI; i < len(key); i++ {
		newNode := NewSliceTrieNode(key[i])
		clen := len(node.Childs)
		if clen <= 0 {
			node.Childs = append(node.Childs, newNode)
		} else if len(node.Childs) == 1 {
			node.Childs = append(node.Childs, newNode)
			if node.Childs[0].Key > key[i] {
				node.Childs[0], node.Childs[1] = node.Childs[1], node.Childs[0]
			}
		} else {
			pos, _ := trie.FindInChilds(node.Childs, key[i])
			node.Childs = append(node.Childs, newNode)
			copy(node.Childs[pos+1:], node.Childs[pos:clen])
			node.Childs[pos] = newNode
		}
		node = newNode
	}
	node.Flag = true
	node.Token = token
}
