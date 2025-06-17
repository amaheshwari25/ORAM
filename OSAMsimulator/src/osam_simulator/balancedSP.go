package osam_simulator

import (
	"fmt"
	"math"
)

// implementation of balanced-tree version of Smart Pointer tree implementation
type BSP struct {
	osam   *OSAM
	print  bool
	nodeId int
}

func (bsp *BSP) log(str string, newline bool) {
	if newline && !suppressPrint {
		fmt.Println()
	}
	if bsp.print && !suppressPrint {
		fmt.Println("[BSP] " + str)
	}
}

func CreateBSP(osam *OSAM, print bool) *BSP {
	return &BSP{osam, print, 0}
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
		// NOTE: NEW -- not in paper: doesn't [chase] also need this change? (TBD)
		nd.tailP = NIL
	} else {
		nd.tailR = NIL
	}
	return nd
}

func (bsp *BSP) saveNode(nd *BNode) {
	a := bsp.osam.Alloc(fmt.Sprintf("saveNode %v", nd))
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

// TBD: check what's correct here?
func (bsp *BSP) addTail(nd *BNode) addr {
	head, tail := bsp.osam.initQueue()
	if !nd.isRoot && nd.tailP == NIL {
		nd.tailP = tail
		// TBD: is this correct?
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

func (bsp *BSP) Copy(p1 *Ptr) Ptr {
	bsp.log(fmt.Sprintf("COPY: starting to copy pointer %v", p1.head), true)
	root := bsp.ascend(p1, false)
	root.count++
	nd := bsp.descend(root)
	p0 := Ptr{head: bsp.addTail(nd)}
	bsp.saveNode(nd)
	return p0
}

func (bsp *BSP) Get(p *Ptr, printPath bool) Block {
	bsp.log(fmt.Sprintf("GET: %v", p.head), true)
	nd := bsp.ascend(p, printPath)
	out := nd.content
	bsp.saveNode(nd)
	return out
}

func (bsp *BSP) Put(p *Ptr, c Block, printPath bool) {
	bsp.log(fmt.Sprintf("PUT: content '%v' @ %v", c.Data, p.head), true)
	nd := bsp.ascend(p, printPath)
	nd.content = c
	bsp.saveNode(nd)
}

func (bsp *BSP) IsNull(p *Ptr) bool {
	return p.head == NIL
}

// TBD CHECK: initialization of BNode?
func (bsp *BSP) New(c Block) Ptr {
	bsp.log(fmt.Sprintf("NEW: starting to create pointer to content %v", c.Data), true)
	nd := bsp.newNode()
	// set root node properties
	nd.content = c
	nd.isRoot = true
	nd.count = 0
	p := Ptr{head: bsp.addTail(nd)}
	bsp.saveNode(nd)
	return p
}
