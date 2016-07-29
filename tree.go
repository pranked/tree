package tree

import (
	"strings"
	"sync"
)

// TreeNode is a tree structure
// TreeNode was setup to facilitate the structure of topic.sub
type TreeNode struct {
	Topic    string
	Root     *TreeNode
	Parent   *TreeNode
	Children []*TreeNode
	Object   interface{}
	mutex    *sync.RWMutex
}

func New() *TreeNode {
	tn := &TreeNode{
		mutex: &sync.RWMutex{},
	}
	tn.Parent = tn
	tn.Root = tn
	return tn
}

func (n *TreeNode) Create(topic string) *TreeNode {
	n.Root.mutex.Lock()
	defer n.Root.mutex.Unlock()
	return n.put(strings.Split(topic, "."))
}

func (n *TreeNode) put(topics []string) *TreeNode {
	newNode := &TreeNode{Topic:topics[0], Root: n.Root, Parent: n}
	if len(topics) == 1 {
		// We only have one more topic so we keep it if it doesn't exist..
		for _, node := range n.Children {
			if node.Topic == newNode.Topic {
				return nil
			}
		}
		n.Children = append(n.Children, newNode)
		return newNode
	}
	// We have to make some children
	// See if that child exists and if it does tell that node to create the new nodes.
	for _, node := range n.Children {
		if node.Topic == topics[0] {
			return node.put(topics[1:])
		}
	}
	// Child doesn't exist we have to make it
	n.Children = append(n.Children, newNode)
	// New node now needs to make the remaining nodes
	return newNode.put(topics[1:])
}

func (n *TreeNode) Search(topic string) *TreeNode {
	n.Root.mutex.RLock()
	defer n.Root.mutex.RUnlock()
	return n.find(strings.Split(topic, "."))
}

func (n *TreeNode) find(topics []string) *TreeNode {
	if len(topics) == 1 {
		// Last topic, we have to have it or it doesn't exist.
		for _, node := range n.Children {
			if node.Topic == topics[0] {
				return node
			}
		}
		return nil
	}
	// See if our children have the right topic name
	for _, node := range n.Children {
		if node.Topic == topics[0] {
			return node.find(topics[1:])
		}
	}
	// We don't have it and our children didn't either
	return nil
}

func (n *TreeNode) Remove(topic string) bool {
	n.Root.mutex.Lock()
	defer n.Root.mutex.Unlock()
	return n.pop(strings.Split(topic, "."))
}

func (n *TreeNode) pop(topics []string) bool {
	if len(topics) == 1 {
		// Last topic, we have to have it or it doesn't exist.
		for i, node := range n.Children {
			node.Children = append(node.Children[:i], node.Children[i + 1:]...)
			return true
		}
		return false
	}
	// See if our children have the right topic name
	for _, node := range n.Children {
		if node.Topic == topics[0] {
			return node.pop(topics[1:])
		}
	}
	// We don't have it and our children didn't either
	return false
}