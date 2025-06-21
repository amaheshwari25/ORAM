package osam_simulator

import (
	"fmt"
	"math"
)

// Implementation of balanced-tree version of Smart Pointer tree implementation
// Three main differences from the paper pseudocode that seem necessary:
// (1) definition of "rightmost" in [descend] function
// (2) "else if" in [addTail] instead of a second "if"
// (3) changes to [chase] (not shown/mentioned in paper)

type BSP struct {
	osam      *OSAM
	print     bool
	nodeId    int
	printPath bool
}

func (bsp *BSP) log(str string, newline bool) {
	if newline && !suppressPrint {
		fmt.Println()
	}
	if bsp.print && !suppressPrint {
		fmt.Println("[BSP] " + str)
	}
}

func CreateBSP(osam *OSAM, print bool, printPath bool) *BSP {
	return &BSP{osam, print, 0, printPath}
}

func (bsp *BSP) newNode() *BNode {
	bsp.nodeId++
	return &BNode{id: bsp.nodeId, tailL: NIL, headL: NIL, tailR: NIL, headR: NIL,
		isRoot: false, count: NONE, headP: NIL, tailP: NIL, content: Block{Data: NONE, IsNone: true}}
}

// ------------ BSP helper functions ------------ //

// copy of SP.chase
func (bsp *BSP) chase(head addr) *BNode {
	target := NIL
	latest := NIL
	tail := NIL
	for head != NIL {
		latest = target
		tail = head
		target, head = bsp.osam.dequeue(head)
	}
	nd := bsp.osam.Read(latest).Data.(*BNode)
	if nd.tailL == tail {
		nd.tailL = NIL
	} else if nd.tailP == tail {
		nd.tailP = NIL // NEW: adding this change to [chase] too (not in paper)?
	} else {
		nd.tailR = NIL
	}
	return nd
}

func (bsp *BSP) saveNode(nd *BNode) {
	a := bsp.osam.Alloc(fmt.Sprintf("saveNode %v", nd.id))
	if nd.tailL != NIL {
		nd.tailL = bsp.osam.enqueue(nd.tailL, a)
	}
	if nd.tailR != NIL {
		nd.tailR = bsp.osam.enqueue(nd.tailR, a)
	}
	if !nd.isRoot && nd.tailP != NIL {
		nd.tailP = bsp.osam.enqueue(nd.tailP, a)
	}
	bsp.osam.writeBN(a, nd)
}

// TBD: check what's correct here? against paper pseudocode
func (bsp *BSP) addTail(nd *BNode) addr {
	head, tail := bsp.osam.initQueue()
	if !nd.isRoot && nd.tailP == NIL {
		nd.tailP = tail
		// NOTE: NEW: "else if" instead of "if" in paper
	} else if nd.tailL == NIL {
		nd.tailL = tail
	} else {
		nd.tailR = tail
	}
	return head
}

// Note: equivalent code to sp.retrieve in smartpointers.go
func (bsp *BSP) ascend(p *Ptr, printPath bool) *BNode {
	nd := bsp.chase(p.head)
	p.head = bsp.addTail(nd)
	for !nd.isRoot {
		if printPath {
			fmt.Printf("Fetched BSP-node: %v \n", nd.id)
		}
		parent := bsp.chase(nd.headP)
		nd.headP = bsp.addTail(parent)
		bsp.saveNode(nd)
		nd = parent
	}
	if printPath {
		fmt.Printf("Fetched BSP-node: %v \n", nd.id)
	}
	assert(nd.isRoot, "Node returned from [ascend] is not root node")
	return nd
}

func getBits(n, len int) []int {
	bits := make([]int, len)
	for i := 0; i < len; i++ {
		bits[len-1-i] = n % 2
		n /= 2
	}
	return bits
}

func (bsp *BSP) descend(root *BNode) *BNode {
	assert(root.isRoot, "Node passed to [descend] is not root node")
	if root.count <= 1 { // short-circuit: should not create / find new node; stay at root
		return root
	}
	pow := int(math.Floor(math.Log2(float64(root.count))))

	// TBD: check what the right formula is?
	rmost := root.count - (1 << pow) // MY SOLUTION
	// rmost := int(math.Floor(float64(root.count -(1<<pow) - 1)/(2.0)))
	nd := root
	rmostbits := getBits(rmost, pow)
	for _, b := range rmostbits {
		var nextNd *BNode
		if b == 0 {
			if nd.headL != NIL {
				nextNd = bsp.chase(nd.headL)
			} else {
				nextNd = bsp.newNode()
				if bsp.printPath {
					fmt.Printf("Created node %v as L child of node %v \n", nextNd.id, nd.id)
				}
				nextNd.tailL = nd.tailL
				nd.tailL = NIL
				nextNd.headP = bsp.addTail(nd)
			}
			nd.headL = bsp.addTail(nextNd)
		} else {
			if nd.headR != NIL {
				nextNd = bsp.chase(nd.headR)
			} else {
				nextNd = bsp.newNode()
				if bsp.printPath {
					fmt.Printf("Created node %v as R child of node %v \n", nextNd.id, nd.id)
				}
				nextNd.tailR = nd.tailR
				nd.tailR = NIL
				nextNd.headP = bsp.addTail(nd)
			}
			nd.headR = bsp.addTail(nextNd)
		}
		bsp.saveNode(nd)
		nd = nextNd
	}
	return nd
}

// ------------ BalancedSmartPointer: MAIN API ------------
//  Get(p: Ptr) -> Block
//  Put(p: Ptr, c: Block)
//  IsNull(p: Ptr)
//  Copy(p1: Ptr) -> Ptr
//  New(c: Block) -> Ptr
//  Delete(p: Ptr)

func (bsp *BSP) Copy(p1 *Ptr) Ptr {
	bsp.log(fmt.Sprintf("COPY: copy pointer %v", p1.head), true)
	root := bsp.ascend(p1, false)
	root.count++
	nd := bsp.descend(root)
	p0 := Ptr{head: bsp.addTail(nd)}
	bsp.saveNode(nd)
	return p0
}

// Same as SP.Get (with different saveNode implementation)
func (bsp *BSP) Get(p *Ptr) Block {
	bsp.log(fmt.Sprintf("GET: %v", p.head), true)
	nd := bsp.ascend(p, bsp.printPath)
	out := nd.content
	bsp.saveNode(nd)
	return out
}

// Same as SP.Put (with different saveNode implementation)
func (bsp *BSP) Put(p *Ptr, c Block) {
	bsp.log(fmt.Sprintf("PUT: content '%v' @ %v", c.Data, p.head), true)
	nd := bsp.ascend(p, bsp.printPath)
	nd.content = c
	bsp.saveNode(nd)
}

func (bsp *BSP) IsNull(p *Ptr) bool {
	return p.head == NIL
}

func (bsp *BSP) New(c Block) Ptr {
	bsp.log(fmt.Sprintf("NEW: create pointer to content %v", c.Data), true)
	nd := bsp.newNode()
	// set root node properties
	nd.content = c
	nd.isRoot = true
	nd.count = 0
	p := Ptr{head: bsp.addTail(nd)}
	bsp.saveNode(nd)
	return p
}

func (bsp *BSP) Delete(p *Ptr) {
	bsp.log(fmt.Sprintf("DELETE: %v", p.head), true)
	root := bsp.ascend(p, false)
	nd := bsp.descend(root)
	root.count -= 1  // NOTE: NEW -- doing this AFTER the descend
	bsp.saveNode(nd) // NEW

	if nd.isRoot {
		bsp.chase(p.head) // to destroy the AQ between the root and p
		if nd.tailL == NIL && nd.tailR == NIL {
			fmt.Printf("All pointers to Node %v deleted; should delete its content \n", nd.id)
		} else {
			bsp.saveNode(nd)
		}
		return
	}

	tailLatest := nd.tailR
	ndPrime := bsp.chase(p.head) // important: ndPrime *could be the same* as nd
	if ndPrime.tailR == NIL {
		ndPrime.tailR = tailLatest
	} else {
		ndPrime.tailL = tailLatest
	}
	bsp.saveNode(ndPrime)
	parent := bsp.chase(nd.headP)
	if parent.tailL == NIL {
		parent.tailL = nd.tailL
		parent.headL = NIL // NEW
	} else {
		parent.tailR = nd.tailL
		parent.headR = NIL // NEW
	}
	bsp.saveNode(parent)

	// OLD CODE: based on paper pseudodcode
	// root := bsp.ascend(p, false)
	// nd := bsp.descend(root)
	// root.count -= 1
	// fmt.Println(nd.id)
	// // NOTE: NEW -- doing this AFTER the descend?
	// tailLatest := nd.tailR
	// parent := bsp.chase(nd.headP)
	// if parent.tailL == NIL {
	// 	parent.tailL = nd.tailL
	// } else {
	// 	parent.tailR = nd.tailL
	// }
	// bsp.saveNode(parent)
	// ndPrime := bsp.chase(p.head)
	// if ndPrime.tailR == NIL {
	// 	ndPrime.tailR = tailLatest
	// } else {
	// 	ndPrime.tailL = tailLatest
	// }
	// bsp.saveNode(ndPrime)
}
