// PathTrie
package Router

import (
	"strings"
)

/*
	All Paths Examples:
	/ids/:id -> /ids/:
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
	Key    byte
	IsVar  bool   //true -> :id,false -> 5
	Var    string //:id,id
	Flag   bool   //word
	Childs []*PathTrieNode
	//此模式路径相应的处理函数
	Handler HandlerFunc
}

func NewPathTrieNode(key byte) *PathTrieNode {
	node := &PathTrieNode{
		Key:     key,
		IsVar:   false,
		Var:     "",
		Flag:    false,
		Childs:  []*PathTrieNode{},
		Handler: nil,
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

func (trie *PathTrie) AddTmpPath(path string) *PathTrieNode {
	key := []byte(path)
	node, lastI := trie.FindLastMatchNode(key)
	if node == nil {
		node = trie.Root
		lastI = 0
	}
	for i := lastI; i < len(key); i++ {
		newNode := NewPathTrieNode(key[i])
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
	return node
}

func (trie *PathTrie) AddPath(path string, handler HandlerFunc) {
	paths := strings.Split(path, "/")
	tPath := ""
	pLen := len(paths)
	if pLen <= 0 {
		return
	}
	var node *PathTrieNode = nil
	for i := 0; i < pLen-1; i++ {
		tPath += paths[i]
		node = trie.AddTmpPath(tPath)
	}
	index := strings.Index(paths[pLen-1], ":")
	if index == 0 {
		node.IsVar = true
		node.Var = paths[pLen-1]
	} else {
		tPath += paths[pLen-1]
		node = trie.AddTmpPath(tPath)
	}
	node.Handler = handler
}

func (trie *PathTrie) MatchPrefixPath(path string) *PathTrieNode {
	key := []byte(path)
	node := trie.Root
	if node == nil {
		return nil
	}
	kLen := len(key)
	var i int = 0
	for i = 0; i < kLen; i++ {
		childs := node.Childs
		pos, found := trie.FindInChilds(childs, key[i])
		if !found {
			break
		}
		node = childs[pos]
	}
	if i != kLen {
		return nil
	}
	return node
}

func (trie *PathTrie) MatchPath(path string) (bool, *Processor) {
	paths := strings.Split(path, "/")
	pLen := len(paths)
	if pLen <= 0 {
		return false, nil
	}
	tPath := ""
	var params *Params = nil
	for i := 0; i < pLen-1; i++ {
		tPath += paths[i]
	}
	node := trie.MatchPrefixPath(tPath)
	if node == nil {
		return false, nil
	}
	if node.IsVar {
		params = &Params{
			Value: paths[pLen-1],
		}
		p := &Processor{
			params:  params,
			Handler: node.Handler,
		}
		return true, p
	} else {
		node = trie.MatchPrefixPath(tPath + paths[pLen-1])
		if node != nil {
			p := &Processor{
				params:  nil,
				Handler: node.Handler,
			}
			return true, p
		}
	}
	return false, nil
}
