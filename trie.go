package main

import "sort"

type TrieNode struct {
	char     rune
	children map[rune]*TrieNode
	num	 int
}

type dict struct {
	str string
	num int
}

func newTrieNode() *TrieNode {
	return &TrieNode{char: 0, children: make(map[rune]*TrieNode), num: 0}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: newTrieNode()}
}

func (trie *Trie) Insert(word string) {
	node := trie.root
	for _, v := range word{
		_, ok := node.children[v]
		if !ok {
			node.children[v] = newTrieNode()
			(node.children[v]).char = v
		}
		node = node.children[v]
	}
	node.num ++
}

func dfs(rt *TrieNode, tmp string, arr *[]dict) {		//搜索字典树中的字符串及其出现次数
	if rt == nil {
		return
	}
	if rt.num != 0 {
		*arr = append(*arr, dict{tmp+string(rt.char), rt.num})
	}
	for _, v := range rt.children {
		dfs(v, tmp+string(rt.char), arr)
	}
}

type Dicts []dict

func (d Dicts) Len() int { return len(d) }

func (d Dicts) Less(i, j int) bool {
	return d[i].num > d[j].num
}

func (d Dicts) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (trie *Trie) GetDictArr() []dict {
	var res []dict
	rt := trie.root
	dfs(rt, "", &res)
	sort.Sort(Dicts(res))
	return res
}

