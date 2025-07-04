package osam_simulator

import "fmt"

type SmartPointer struct {
	osam      *OSAM
	print     bool
	nodeId    int
	printPath bool
}

func (sp *SmartPointer) log(str string, newline bool) {
	if newline && !suppressPrint {
		fmt.Println()
	}
	if sp.print && !suppressPrint {
		fmt.Println("[SP] " + str)
	}
}

func CreateSP(osam *OSAM, print bool, printPath bool) *SmartPointer {
	return &SmartPointer{osam, print, 0, printPath}
}

// defaults to all NIL values & intermediate node parameters otherwise
func (sp *SmartPointer) newNode() *Node {
	sp.nodeId++
	return &Node{id: sp.nodeId, tailL: NIL, tailR: NIL, headP: NIL, isRoot: false, content: Block{Data: NONE, IsNone: true}}
}

// ------------ SmartPointer helper functions ------------ //

func (sp *SmartPointer) chase(head addr) *Node {
	target := NIL
	latest := NIL
	tail := NIL
	for head != NIL {
		latest = target
		tail = head
		target, head = sp.osam.dequeue(head)
	}
	nd := sp.osam.Read(latest).Data.(*Node)
	if nd.tailL == tail {
		nd.tailL = NIL
	} else {
		nd.tailR = NIL
	}
	return nd
}

func (sp *SmartPointer) saveNode(nd *Node) {
	a := sp.osam.Alloc(fmt.Sprintf("saveNode %v", nd))
	if nd.tailL != NIL {
		nd.tailL = sp.osam.enqueue(nd.tailL, a)
	}
	if nd.tailR != NIL {
		nd.tailR = sp.osam.enqueue(nd.tailR, a)
	}
	sp.osam.writeN(a, nd)
}

func (sp *SmartPointer) addTail(nd *Node) addr {
	head, tail := sp.osam.initQueue()
	if nd.tailL == NIL {
		nd.tailL = tail
	} else {
		nd.tailR = tail
	}
	return head
}

// Helper function that is the main body of Get and Put
func (sp *SmartPointer) retrieve(p *Ptr, printPath bool) *Node {
	nd := sp.chase(p.head)
	p.head = sp.addTail(nd)
	for !nd.isRoot {
		if printPath {
			fmt.Printf("Fetched SP-node: %v \n", nd.id)
		}
		parent := sp.chase(nd.headP)
		nd.headP = sp.addTail(parent)
		sp.saveNode(nd)
		nd = parent
	}
	if printPath {
		fmt.Printf("Fetched SP-node: %v \n", nd.id)
	}
	assert(nd.isRoot, "Node returned from [retrieve] is not root node")
	return nd
}

// ------------ SmartPointer: MAIN API ------------ //
//  Get(p: Ptr) -> Block
//  Put(p: Ptr, c: Block)
//  IsNull(p: Ptr)
//  Copy(p1: Ptr) -> Ptr
//  New(c: Block) -> Ptr
//  Delete(p: Ptr)

func (sp *SmartPointer) Get(p *Ptr) Block {
	sp.log(fmt.Sprintf("GET: %v", p.head), true)
	nd := sp.retrieve(p, sp.printPath)
	// invariant after [retrieve]: nd.isRoot should be true
	out := nd.content
	sp.saveNode(nd)
	return out
}

func (sp *SmartPointer) Put(p *Ptr, c Block) {
	sp.log(fmt.Sprintf("PUT: content '%v' @ %v", c.Data, p.head), true)
	nd := sp.retrieve(p, sp.printPath)
	nd.content = c
	sp.saveNode(nd)
}

func (sp *SmartPointer) IsNull(p Ptr) bool {
	return p.head == NIL
}

func (sp *SmartPointer) Copy(p1 *Ptr) Ptr {
	sp.log(fmt.Sprintf("COPY: starting to copy pointer %v", p1.head), true)
	nd := sp.chase(p1.head)
	if nd.tailL != NIL || nd.tailR != NIL {
		ndNew := sp.newNode()
		if sp.printPath {
			if nd.tailL == NIL {
				fmt.Printf("Created node %v as L child of node %v \n", ndNew.id, nd.id)
			} else {
				fmt.Printf("Created node %v as R child of node %v \n", ndNew.id, nd.id)
			}
		}
		ndNew.headP = sp.addTail(nd)
		// NOTE: paper says to do sp.saveNode(ndNew), but I think this is a typo and should be saveNode(nd)
		sp.saveNode(nd)
		nd = ndNew
	}
	p0 := Ptr{head: sp.addTail(nd)}
	p1.head = sp.addTail(nd)
	sp.saveNode(nd)
	return p0
}

func (sp *SmartPointer) New(c Block) Ptr {
	sp.log(fmt.Sprintf("NEW: starting to create pointer to content %v", c.Data), true)
	nd := sp.newNode()
	nd.isRoot = true
	nd.content = c
	p := Ptr{head: sp.addTail(nd)}
	sp.saveNode(nd)
	return p
}

func (sp *SmartPointer) Delete(p *Ptr) {
	sp.log(fmt.Sprintf("DELETE: %v", p.head), true)
	if p.head != NIL {
		nd := sp.chase(p.head)
		if nd.isRoot {
			if nd.tailL == NIL && nd.tailR == NIL {
				// note: [chase] will have recently nulled-out one
				sp.log(fmt.Sprintf("All pointers to %v deleted; should delete its content", nd), false)
			} else {
				sp.saveNode(nd)
			}
		} else {
			var tail addr
			if nd.tailL != NIL {
				tail = nd.tailL
			} else {
				tail = nd.tailR
			}
			nd = sp.chase(nd.headP)
			if nd.tailL == NIL {
				nd.tailL = tail
			} else {
				nd.tailR = tail
			}
			sp.saveNode(nd)
		}
	}
}
