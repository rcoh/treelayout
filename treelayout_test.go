package treelayout

import (
	"math/rand"
	"os/exec"
	"testing"
)

func TestTreePreorderStress(t *testing.T) {
	n := 1 << 10
	rand.Seed(0)
	values := make([]float64, n)
	for i := range values {
		values[i] = rand.Float64()
	}

	treeArray := createTreePreorder(values)
	for _, v := range values {
		found := findValue(v, treeArray, 0)
		if found != v {
			t.Errorf("Inserted value v not found")
		}
	}
}

func TestTreeLeveledStress(t *testing.T) {
	n := 1<<5 - 1
	rand.Seed(0)
	values := make([]float64, n)
	for i := range values {
		values[i] = float64(i) //rand.Float64()
	}

	treeArray := createTreeLeveled(values)
	for _, v := range values {
		found := findValue(v, treeArray, 0)
		if found != v {
			t.Errorf("Inserted value %v not found", v)
		}
	}
	findValue(float64(n+1), treeArray, 0)
}

func TestTreeCBLayout(t *testing.T) {
	rand.Seed(0)
	values := []float64{4, 2, 6, 1, .5, 1.5, 3, 2.5, 3.5, 5, 4.5, 5.5, 7, 6.5, 7.5}

	treeArray := createTreeCacheOblivious(values)
	for _, v := range values {
		found := findValue(v, treeArray, 0)
		if found != v {
			t.Errorf("Inserted value v not found")
		}
	}

	//          4
	//      2             6
	//  1       3       5       7
	//.5 1.5 2.5 3.5  4.5 5.5 6.5  7.5

	expected := []float64{4, 2, 6, 1, .5, 1.5, 3, 2.5, 3.5, 5, 4.5, 5.5, 7, 6.5, 7.5}
	for pos := range expected {
		if treeArray[pos].Value != expected[pos] {
			t.Errorf("Unexpected value at position %v. Expected %v got %v", pos, expected[pos], treeArray[pos].Value)
		}
	}
}

func TestTreeLeveledLayout(t *testing.T) {
	values := []float64{1, 2, 3, 4, 5, 6, 7}
	//       4
	//    2      6
	//  1   3  5   7
	expected := []float64{4, 2, 6, 1, 3, 5, 7}
	treeArray := createTreeLeveled(values)
	for pos := range expected {
		if treeArray[pos].Value != expected[pos] {
			t.Errorf("Unexpected value at position %v. Expected %v got %v", pos, expected[pos], treeArray[pos].Value)
		}
	}
}

func TestSplitTree(t *testing.T) {
	rand.Seed(0)
	n := 31
	values := make([]float64, n)
	for i := range values {
		values[i] = rand.Float64()
	}

	treeArray := createTree(values)
	root, leaves := splitTree(treeArray, treeArray.Height/2)
	if root.Height != 2 {
		t.Errorf("Root has wrong height")
	}
	if len(leaves) != 4 {
		t.Errorf("Wrong number of leaves expected 2 got %v", len(leaves))
	}
}

func treeData(n int) []float64 {
	rand.Seed(0)
	v := make([]float64, n)
	pos := rand.Perm(n)
	for i := range v {
		v[pos[i]] = rand.Float64() // + float64(i)
	}
	return v
}

var n = 1<<24 - 1
var values []float64 = treeData(n)

var leveledTree []TreeNodeInline

func BenchmarkLeveledTree(b *testing.B) {
	if leveledTree == nil {
		leveledTree = createTreeLeveled(values)
	}

	treeArray := leveledTree
	b.ResetTimer()
	benchmarkTree(treeArray, b)
}

var preorderTree []TreeNodeInline

func BenchmarkPreorderTree(b *testing.B) {
	if preorderTree == nil {
		preorderTree = createTreePreorder(values)
	}
	treeArray := preorderTree
	b.ResetTimer()
	benchmarkTree(treeArray, b)
}

var cbTree []TreeNodeInline

func BenchmarkCBTree(b *testing.B) {
	if cbTree == nil {
		cbTree = createTreeCacheOblivious(values)
	}
	treeArray := cbTree
	b.ResetTimer()
	benchmarkTree(treeArray, b)
}

func dropCaches() {
	cmd := "free && sync && echo 3 > /proc/sys/vm/drop_caches && free"
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		panic(err)
	}
}

/*
func stressTestTree(values, []float64, treeArray []TreeNodeInline, t *testing.T) {
	treeValues := []float64{}
	for _, node := range treeArray {
		if findValue(v, treeArray, 0) != v {
			t.Errorf("Looking for %v, couldn't find", v)
		}
	}
	for v, _ := range treeValues {
	}
}*/

func benchmarkTree(treeArray []TreeNodeInline, b *testing.B) {
	dropCaches()
	res := float64(0)
	for i := 0; i < b.N; i++ {
		v := rand.Float64()
		found := findValue(v, treeArray, 0)
		res += found
	}
}
