package treelayout

import (
	"sort"
)

// layout in memory, using pointers

type TreeNodeInline struct {
	Value      float64
	LeftChild  int
	RightChild int
}

type TreeNode struct {
	Value      float64
	LeftChild  *TreeNode
	RightChild *TreeNode
	Height     int
}

func findRoot(values []float64) float64 {
	sort.Float64s(values)
	return values[len(values)/2]
}

func createTree(values []float64) *TreeNode {
	if len(values) == 0 {
		return nil
	}

	value := findRoot(values)
	leftValues := make([]float64, 0, 0)
	rightValues := make([]float64, 0, 0)
	for _, v := range values {
		if v < value {
			leftValues = append(leftValues, v)
		} else if v > value {
			rightValues = append(rightValues, v)
		}

	}
	leftSubtree := createTree(leftValues)
	rightSubtree := createTree(rightValues)
	height := 0
	if leftSubtree != nil {
		height = max(height, leftSubtree.Height)
	}

	if rightSubtree != nil {
		height = max(height, rightSubtree.Height)
	}

	return &TreeNode{value, leftSubtree, rightSubtree, height + 1}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func createTreePreorder(values []float64) []TreeNodeInline {
	tree := createTree(values)
	treeArray := make([]TreeNodeInline, len(values), len(values))
	fillTreeArray(treeArray, 0, tree)
	return treeArray
}

func createTreeCacheOblivious(values []float64) []TreeNodeInline {
	treeArray := make([]TreeNodeInline, len(values), len(values))
	tree := createTree(values)
	fillTreeArraySplitting(treeArray, 0, tree, func(t *TreeNode) int { return t.Height / 2 })
	return treeArray
}

type splitLevel func(*TreeNode) int

func createTreeLeveled(values []float64) []TreeNodeInline {
	treeArray := make([]TreeNodeInline, len(values))
	tree := createTree(values)
	currentPos := 0
	for level := uint(0); level < uint(tree.Height); level++ {
		nodesAtLevel := getLevel(tree, int(level))
		if len(nodesAtLevel) != 1<<level {
			panic("Unexpected nodes at level")
		}
		nextLevelStart := currentPos + len(nodesAtLevel)
		for upperLevelNodePosition, node := range nodesAtLevel {
			lChild := -1
			rChild := -1
			if uint(tree.Height)-1 > level {
				lChild = nextLevelStart + upperLevelNodePosition*2
				rChild = nextLevelStart + upperLevelNodePosition*2 + 1
			}
			treeArray[currentPos] = TreeNodeInline{node.Value, lChild, rChild}
			currentPos++
		}
	}
	return treeArray
}

func getLevel(tree *TreeNode, level int) []*TreeNode {
	if level == 0 {
		return []*TreeNode{tree}
	}
	leftLeaves := getLevel(tree.LeftChild, level-1)
	rightLeaves := getLevel(tree.RightChild, level-1)
	return append(leftLeaves, rightLeaves...)
}

func fillTreeArraySplitting(treeArray []TreeNodeInline, nextOpenSpot int, tree *TreeNode, splitFunc splitLevel) int {
	if tree.Height == 1 {
		treeArray[nextOpenSpot] = TreeNodeInline{tree.Value, -1, -1}
		return nextOpenSpot + 1
	}
	top, leaves := splitTree(tree, splitFunc(tree))
	firstLeafSpot := fillTreeArraySplitting(treeArray, nextOpenSpot, top, splitFunc)
	leafLocations := make([]int, 0)
	nextLeafSpot := firstLeafSpot
	for _, leaf := range leaves {
		leafLocations = append(leafLocations, nextLeafSpot)
		nextLeafSpot = fillTreeArraySplitting(treeArray, nextLeafSpot, leaf, splitFunc)
	}

	unboundLeaves := leafIndexes(treeArray, nextOpenSpot)
	if len(leafLocations) != len(unboundLeaves)*2 {
		panic("Number of leaves is wrong")
	}

	nextLeaf := 0
	for _, leafIndex := range unboundLeaves {
		treeArray[leafIndex].LeftChild = leafLocations[nextLeaf]
		nextLeaf++
		treeArray[leafIndex].RightChild = leafLocations[nextLeaf]
		nextLeaf++
	}
	return nextLeafSpot
}

func leafIndexes(treeArray []TreeNodeInline, rootIndex int) []int {
	node := treeArray[rootIndex]
	if node.LeftChild == -1 && node.RightChild == -1 {
		return []int{rootIndex}
	}
	return append(leafIndexes(treeArray, node.LeftChild), leafIndexes(treeArray, node.RightChild)...)
}

func wireLeaves(treeArray []TreeNodeInline, rootIndex int, leafPositions []int) {

}

func splitTree(tree *TreeNode, splitHeight int) (*TreeNode, []*TreeNode) {
	if tree == nil {
		panic("Tree can't be nil")
	}
	root, leaves := splitTreeBrokenHeight(splitHeight, tree)
	rootHeight := fixHeight(root)
	if rootHeight != splitHeight {
		println("Root height is wrong %s vs %s", rootHeight, splitHeight)
		panic("see above")
	}
	for _, leaf := range leaves {
		lHeight := fixHeight(leaf)
		if lHeight != tree.Height-splitHeight {
			panic("Leaf height is wrong")
		}

	}
	return root, leaves
}

func splitTreeBrokenHeight(height int, tree *TreeNode) (*TreeNode, []*TreeNode) {
	if height == 1 {
		leaves := make([]*TreeNode, 0)
		if tree.LeftChild != nil {
			leaves = append(leaves, tree.LeftChild)
		}
		if tree.RightChild != nil {
			leaves = append(leaves, tree.RightChild)
		}
		return &TreeNode{tree.Value, nil, nil, -1}, leaves
	}
	leftChild, lLeaves := splitTreeBrokenHeight(height-1, tree.LeftChild)
	rightChild, rLeaves := splitTreeBrokenHeight(height-1, tree.RightChild)
	return &TreeNode{tree.Value, leftChild, rightChild, -1}, append(lLeaves, rLeaves...)
}

func fixHeight(t *TreeNode) int {
	if t == nil {
		return 0
	}
	t.Height = 1 + max(fixHeight(t.LeftChild), fixHeight(t.RightChild))
	return t.Height
}

func fillTreeArray(treeArray []TreeNodeInline, nextOpenSlot int, tree *TreeNode) int {
	if tree == nil {
		return nextOpenSlot
	}
	if nextOpenSlot >= len(treeArray) {
		panic("No spot for current node")
	}
	finalSlot := fillTreeArray(treeArray, nextOpenSlot+1, tree.LeftChild)
	lChild := -1
	rChild := -1
	if finalSlot != nextOpenSlot+1 {
		lChild = nextOpenSlot + 1
	}

	nextAvailable := fillTreeArray(treeArray, finalSlot, tree.RightChild)
	if nextAvailable != finalSlot {
		rChild = finalSlot
	}
	if lChild > len(treeArray) || rChild > len(treeArray) {
		panic("oob pointer")
	}
	node := TreeNodeInline{tree.Value, lChild, rChild}
	treeArray[nextOpenSlot] = node
	return nextAvailable
}

func findValue(v float64, tree []TreeNodeInline, index int) float64 {
	if index == -1 {
		return -1
	}
	node := tree[index]
	if node.Value == v {
		return v
	} else if node.Value > v && node.LeftChild != -1 {
		return findValue(v, tree, node.LeftChild)
	} else if node.Value < v && node.RightChild != -1 {
		return findValue(v, tree, node.RightChild)
	} else {
		return v
	}
}
