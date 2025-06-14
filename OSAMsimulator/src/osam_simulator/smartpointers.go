package osam_simulator

import "fmt"

type SmartPointer struct {
	osam *OSAM
	// debug bool
}

func CreateSP(osam *OSAM) *SmartPointer {
	return &SmartPointer{osam}
}

type Ptr struct {
	head addr
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
	sp.osam.WriteN(a, nd)
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
func (sp *SmartPointer) retrieve(p Ptr) *Node {
	nd := sp.chase(p.head)
	p.head = sp.addTail(nd)
	for !nd.isRoot {
		parent := sp.chase(nd.headP)
		nd.headP = sp.addTail(parent)
		sp.saveNode(nd)
		nd = parent
	}
	return nd
}

// ------------ SmartPointer: MAIN API ------------ //
func (sp *SmartPointer) Get(p Ptr) Block {
	fmt.Printf("GET: %v \n", p)
	nd := sp.retrieve(p)
	// invariant: nd.isRoot should be true
	out := nd.content
	sp.saveNode(nd)
	return out
}

func (sp *SmartPointer) Put(p Ptr, c Block) {
	fmt.Printf("PUT: %v, content %v \n", p, c.Data)
	nd := sp.retrieve(p)
	nd.content = c
	sp.saveNode(nd)
}

func (sp *SmartPointer) IsNull(p Ptr) bool {
	return p.head == NIL
}

func (sp *SmartPointer) Copy(p1 Ptr) Ptr {
	fmt.Printf("COPY: starting to copy pointer %v \n", p1)
	nd := sp.chase(p1.head)
	if nd.tailL != NIL || nd.tailR != NIL {
		ndNew := &Node{headP: sp.addTail(nd), isRoot: false, content: Block{Data: NONE, IsNone: true}}
		sp.saveNode(ndNew)
		nd = ndNew
	}
	p0 := Ptr{head: sp.addTail(nd)}
	p1.head = sp.addTail(nd)
	fmt.Printf("COPY: copied pointer %v to create pointer %v \n", p1, p0)

	return p0
}

func (sp *SmartPointer) New(c Block) Ptr {
	fmt.Printf("NEW: starting to create pointer to content %v \n", c.Data)
	nd := &Node{tailL: NIL, tailR: NIL, content: c, isRoot: true, headP: NIL}
	p := Ptr{head: sp.addTail(nd)}
	sp.saveNode(nd)
	fmt.Printf("NEW: created pointer %v to content %v \n", p, c.Data)
	return p
}

func (sp *SmartPointer) Delete(p Ptr) {
	if p.head != NIL {
		nd := sp.chase(p.head)
		if nd.isRoot {
			if nd.tailL == NIL && nd.tailR == NIL { // note: chase will have recently nulled-out one
				fmt.Printf("All pointers to %v deleted; should delete its content", nd)
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
