package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// adapted from https://github.com/andipro123/merkle-tree-Python/blob/master/merkleTree.py

type Node struct {
	data  string
	left  *Node
	right *Node
}

func NewNode(data string) *Node {
	return &Node{
		data: data,
	}
}

func (n *Node) isFull() bool {
	return n.left != nil && n.right != nil
}

func (n *Node) isLeaf() bool {
	return n.left == nil && n.right == nil
}

type MerkleTree struct {
	merkleRootS string
	root        *Node
}

func NewMerkleTree() *MerkleTree {
	return &MerkleTree{}
}

func (mt *MerkleTree) makeTreeFromTxs(txs []string) {
	nodesNeeded := numberOfNodesFromSlice(txs)
	nodes := make([]string, nodesNeeded)
	for i := range nodesNeeded {
		nodes[i] = ""
	}
	mt.root = mt.buildTree(nodes, nil, 0, len(nodes))
	mt.addLeaves(txs, mt.root)
}

func (mt *MerkleTree) buildTree(arr []string, root *Node, i, n int) *Node {
	if i < n {
		temp := NewNode(arr[i])
		root = temp

		root.left = mt.buildTree(arr, root.left, 2*i+1, n)
		root.right = mt.buildTree(arr, root.right, 2*i+2, n)
	}

	return root
}

func (mt *MerkleTree) getHash(x string) string {
	hash := sha256.Sum256([]byte(x))
	return hex.EncodeToString(hash[:])
}

func (mt *MerkleTree) addLeaves(arr []string, node *Node) {
	if node == nil {
		return
	}

	mt.addLeaves(arr, node.left)
	if node.isLeaf() {
		node.data = mt.getHash(arr[0])
		arr = arr[1:]
	} else {
		node.data = ""
	}
	mt.addLeaves(arr, node.right)
}

func numberOfNodesFromSlice(arr []string) int {
	l := len(arr)
	return (2*l - 1)
}

func travInorder(n *Node) {
	if n == nil {
		return
	}
	travInorder(n.left)
	fmt.Println(n.data)
	travInorder(n.right)
}

func (mt *MerkleTree) merkleHash(node *Node) *Node {
	if node.isLeaf() {
		return node
	}
	left := mt.merkleHash(node.left).data
	right := mt.merkleHash(node.right).data
	node.data = mt.getHash(left + right)
	return node
}

func (mt *MerkleTree) computeMerkleHash() string {
	merkleRoot := mt.merkleHash(mt.root)
	mt.merkleRootS = merkleRoot.data

	return mt.merkleRootS
}

func (mt *MerkleTree) verify(txs []string) bool {
	hash1 := mt.merkleRootS
	newTree := NewMerkleTree()
	newTree.makeTreeFromTxs(txs)
	newTree.computeMerkleHash()
	hash2 := newTree.merkleRootS
	return hash1 == hash2
}

func main() {
	txs := []string{
		"0x0000000000000000000000000",
		"0x0000000000000000000000011",
		"0x0000000000000000000000012",
		"0x0000000000000000000000014",
		"0x0000000000000000000000015",
		"0x0000000000000000000000016",
		"0x0000000000000000000000017",
	}
	tree := NewMerkleTree()
	tree.makeTreeFromTxs(txs)
	tree.computeMerkleHash()

	if tree.verify(txs) {
		fmt.Println("trees are equal, no tampering with txs has been done")
	} else {
		fmt.Println("trees are not equal, somebody tampered with txs")
	}
}
