package hashmap

type color bool

const black color = false
const red color = true

type matchPosition int8

const greater matchPosition = 1
const same matchPosition = 0
const lower matchPosition = -1

type rbTreeNode struct {
	color color

	keyHash int64
	key     interface{}
	value   interface{}

	parent *rbTreeNode

	left  *rbTreeNode
	right *rbTreeNode

	collisions map[interface{}]interface{}
}

type rbTree struct {
	root     *rbTreeNode
	hashFunc func(interface{}) int64
}

// New creates a new hash map with supplied hashing function
func New(hashFunc func(i interface{}) int64) *rbTree {
	return &rbTree{hashFunc: hashFunc}
}

func (rb *rbTree) Insert(key, value interface{}) {
	keyHash := rb.hashFunc(key)

	child := &rbTreeNode{
		keyHash: keyHash,
		key:     key,
		value:   value,
		left:    &rbTreeNode{},
		right:   &rbTreeNode{},
		color:   red,
	}
	child.collisions = map[interface{}]interface{}{}

	if rb.root != nil {
		//find insertion parent and position where we should place child
		parent, position := findInsertionParent(rb.root, keyHash)

		//insert the child node
		switch position {
		case greater:
			parent.right = child
			child.parent = parent
		case lower:
			parent.left = child
			child.parent = parent
		case same:
			if key == parent.key {
				parent.value = value
			} else {
				if parent.collisions == nil {
					parent.collisions = map[interface{}]interface{}{}
				}
				parent.collisions[key] = value
			}

			return
		}
	}

	insertCase1(child)

	//crawl to root and assign it
	for {
		if child.parent == nil {
			rb.root = child
			break
		}

		child = child.parent
	}
}

func insertCase1(node *rbTreeNode) {
	if node.parent == nil {
		node.color = black
		return
	}

	insertCase2(node)
}

func insertCase2(node *rbTreeNode) {
	if node.parent.color == black {
		return
	}

	insertCase3(node)
}

func insertCase3(node *rbTreeNode) {
	uncle := getUncle(node)
	if uncle != nil && uncle.color == red {
		node.parent.color = black
		uncle.color = black
		grandparent := getGrandparent(node)
		grandparent.color = red
		insertCase1(grandparent)
		return
	}

	insertCase4(node)
}

func insertCase4(node *rbTreeNode) {
	grandparent := getGrandparent(node)

	if node == node.parent.right && node.parent == grandparent.left {
		rotateLeft(node.parent)
		node = node.left
	} else if node == node.parent.left && node.parent == grandparent.right {
		rotateRight(node.parent)
		node = node.right
	}

	insertCase5(node)
}

func insertCase5(node *rbTreeNode) {
	grandparent := getGrandparent(node)
	node.parent.color = black
	grandparent.color = red
	if node == node.parent.left {
		rotateRight(grandparent)
	} else {
		rotateLeft(grandparent)
	}
}

func (rb *rbTree) Get(key interface{}) (value interface{}, found bool) {
	if rb.root == nil {
		return nil, false
	}

	keyHash := rb.hashFunc(key)

	node, found := findByKeyHash(rb.root, key, keyHash)
	if !found {
		return nil, false
	}

	if node.key == key {
		return node.value, true
	}

	value, found = node.collisions[key]
	return
}

func (rb *rbTree) Remove(key interface{}) (found bool) {
	keyHash := rb.hashFunc(key)
	node, found := findByKeyHash(rb.root, key, keyHash)
	if !found {
		return true
	}

	if len(node.collisions) > 0 {
		if key == node.key {
			for k, v := range node.collisions {
				node.key = k
				node.value = v
				break
			}
			key = node.key
		}

		delete(node.collisions, key)
		return true
	}

	//return a node with at most one non leaf sibling that should be used to replace node
	replacementNode := getReplacementNode(node)
	//copy the replacement value into original
	copyNodeValue(replacementNode, node)

	//select the replacement node's child
	replacementNodeChild := replacementNode.right
	if isLeaf(replacementNodeChild) {
		replacementNodeChild = replacementNode.left
	}

	//replace the replacementNode with its child
	replacementNodeChild.parent = replacementNode.parent
	if replacementNode.parent != nil {
		if replacementNode == replacementNode.parent.left {
			replacementNode.parent.left = replacementNodeChild
		} else {
			replacementNode.parent.right = replacementNodeChild
		}
	}

	//if it was red we don't care
	if replacementNode.color == red {
		return true
	}

	//if it was black and the new node is red, repaint the new node to black, preserves black depth
	if replacementNodeChild.color == red {
		replacementNodeChild.color = black
		return true
	}

	deleteCase1(replacementNodeChild)

	//crawl to root and assign it
	for {
		if replacementNodeChild.parent == nil {
			rb.root = replacementNodeChild

			if isLeaf(rb.root) {
				rb.root = nil
			}

			break
		}

		replacementNodeChild = replacementNodeChild.parent
	}

	return true
}

func isLeaf(node *rbTreeNode) bool {
	return node.left == nil && node.right == nil && node.color == black
}

//if node is the new root, finish
func deleteCase1(node *rbTreeNode) {
	if node.parent != nil {
		deleteCase2(node)
	}
}

//if sibling is red, we can switch sibling and parent colours and rotate
func deleteCase2(node *rbTreeNode) {
	sibling := getSibling(node)

	if sibling.color == red {
		node.parent.color = red
		sibling.color = black

		if node == node.parent.left {
			rotateLeft(node.parent)
		} else {
			rotateRight(node.parent)
		}
	}

	deleteCase3(node)
}

func deleteCase3(node *rbTreeNode) {
	sibling := getSibling(node)
	if node.parent.color == black && sibling.color == black && sibling.left.color == black && sibling.right.color == black {
		sibling.color = red
		deleteCase1(node.parent)
	} else {
		deleteCase4(node)
	}
}

func deleteCase4(node *rbTreeNode) {
	sibling := getSibling(node)
	if node.parent.color == red && sibling.color == black && sibling.left.color == black && sibling.right.color == black {
		sibling.color = red
		node.parent.color = black
	} else {
		deleteCase5(node)
	}
}

func deleteCase5(node *rbTreeNode) {
	sibling := getSibling(node)
	if sibling.color == black {
		if node.parent.left == node && sibling.right.color == black && sibling.left.color == red {
			sibling.color = red
			sibling.left.color = black
			rotateRight(sibling)
		} else if node.parent.right == node && sibling.right.color == red && sibling.left.color == black {
			sibling.color = red
			sibling.right.color = black
			rotateLeft(sibling)
		}
	}

	deleteCase6(node)
}

func deleteCase6(node *rbTreeNode) {
	sibling := getSibling(node)

	sibling.color = node.parent.color
	node.parent.color = black

	if node == node.parent.left {
		sibling.right.color = black
		rotateLeft(node.parent)
	} else {
		sibling.left.color = black
		rotateRight(node.parent)
	}
}

func copyNodeValue(fromNode *rbTreeNode, toNode *rbTreeNode) {
	//todo: optimize this to just do pointer magic, instead of copying values
	toNode.key = fromNode.key
	toNode.keyHash = fromNode.keyHash
	toNode.value = fromNode.value
	toNode.collisions = fromNode.collisions
}

func getReplacementNode(node *rbTreeNode) *rbTreeNode {
	if !isLeaf(node.right) {
		return getLeftmostNode(node.right)
	} else if !isLeaf(node.left) {
		return getRightmostNode(node.left)
	}

	return node
}

func getLeftmostNode(node *rbTreeNode) *rbTreeNode {
	if isLeaf(node.left) {
		return node
	}

	return getLeftmostNode(node.left)
}

func getRightmostNode(node *rbTreeNode) *rbTreeNode {
	if isLeaf(node.right) {
		return node
	}

	return getRightmostNode(node.right)
}

func findByKeyHash(node *rbTreeNode, key interface{}, keyHash int64) (res *rbTreeNode, found bool) {
	if node == nil {
		return
	} else if keyHash > node.keyHash && !isLeaf(node.right) {
		return findByKeyHash(node.right, key, keyHash)
	} else if keyHash < node.keyHash && !isLeaf(node.left) {
		return findByKeyHash(node.left, key, keyHash)
	} else if keyHash == node.keyHash {
		return node, true
	}

	return
}

func rotateLeft(root *rbTreeNode) {
	pivot := root.right
	if isLeaf(pivot) {
		return
	}

	rootParent := root.parent
	if rootParent != nil && rootParent.left == root {
		rootParent.left = pivot
	} else if rootParent != nil && rootParent.right == root {
		rootParent.right = pivot
	}

	pivotLeftChild := root.right.left
	pivot.parent = rootParent
	root.parent = pivot
	pivot.left = root
	root.right = pivotLeftChild
	pivotLeftChild.parent = root
}

func rotateRight(root *rbTreeNode) {
	pivot := root.left
	if isLeaf(root.left) {
		return
	}

	rootParent := root.parent
	if rootParent != nil && rootParent.right == root {
		rootParent.right = pivot
	} else if rootParent != nil && rootParent.left == root {
		rootParent.left = pivot
	}

	pivotRightChild := root.left.right
	pivot.parent = root.parent
	root.parent = pivot
	pivot.right = root
	root.left = pivotRightChild
	pivotRightChild.parent = root
}

func findInsertionParent(n *rbTreeNode, keyHash int64) (*rbTreeNode, matchPosition) {
	if keyHash > n.keyHash {
		if isLeaf(n.right) {
			return n, greater
		}

		return findInsertionParent(n.right, keyHash)

	} else if keyHash < n.keyHash {
		if isLeaf(n.left) {
			return n, lower
		}

		return findInsertionParent(n.left, keyHash)
	}

	return n, same
}

func getGrandparent(n *rbTreeNode) (g *rbTreeNode) {
	if n.parent != nil && n.parent.parent != nil {
		g = n.parent.parent
	}

	return
}

func getUncle(n *rbTreeNode) (u *rbTreeNode) {
	g := getGrandparent(n)
	if g == nil {
		return
	} else if n.parent == g.left {
		return g.right
	}

	return g.left
}

func getSibling(n *rbTreeNode) (u *rbTreeNode) {
	if n.parent == nil {
		return nil
	}

	if n.parent.left == n {
		return n.parent.right
	}

	return n.parent.left
}
