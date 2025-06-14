package osam_simulator

import "fmt"

type SmartPointer struct {
	osam    *OSAM
	printSP bool
}

func CreateSP(osam *OSAM, print bool) *SmartPointer {
	return &SmartPointer{osam, print}
}

type Ptr struct {
	head addr
}

func (sp *SmartPointer) log(str string, newline bool) {
	var nlStr string
	if newline {
		nlStr = "\n"
	} else {
		nlStr = ""
	}
	if sp.printSP {
		fmt.Println(nlStr + "[SP] " + str)
	}
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
func (sp *SmartPointer) retrieve(p *Ptr) *Node {
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
func (sp *SmartPointer) Get(p *Ptr) Block {
	sp.log(fmt.Sprintf("GET: %v", p), true)
	nd := sp.retrieve(p)
	// invariant: nd.isRoot should be true
	out := nd.content
	sp.saveNode(nd)
	return out
}

func (sp *SmartPointer) Put(p *Ptr, c Block) {
	sp.log(fmt.Sprintf("PUT: %v, content %v", p, c.Data), true)
	nd := sp.retrieve(p)
	nd.content = c
	sp.saveNode(nd)
}

func (sp *SmartPointer) IsNull(p Ptr) bool {
	return p.head == NIL
}

func (sp *SmartPointer) Copy(p1 *Ptr) Ptr {
	sp.log(fmt.Sprintf("COPY: starting to copy pointer %v", p1), true)
	nd := sp.chase(p1.head)
	if nd.tailL != NIL || nd.tailR != NIL {
		ndNew := &Node{headP: sp.addTail(nd), isRoot: false, content: Block{Data: NONE, IsNone: true}}
		sp.saveNode(ndNew)
		nd = ndNew
	}
	p0 := Ptr{head: sp.addTail(nd)}
	p1.head = sp.addTail(nd)
	sp.saveNode(nd)
	sp.log(fmt.Sprintf("COPY: finished copying pointer %v to create pointer %v", p1, p0), false)
	return p0
}

func (sp *SmartPointer) New(c Block) Ptr {
	sp.log(fmt.Sprintf("NEW: starting to create pointer to content %v", c.Data), true)
	nd := &Node{tailL: NIL, tailR: NIL, content: c, isRoot: true, headP: NIL}
	p := Ptr{head: sp.addTail(nd)}
	sp.saveNode(nd)
	sp.log(fmt.Sprintf("NEW: finished creating pointer %v to content %v", p, c.Data), false)
	return p
}

func (sp *SmartPointer) Delete(p *Ptr) {
	if p.head != NIL {
		nd := sp.chase(p.head)
		if nd.isRoot {
			if nd.tailL == NIL && nd.tailR == NIL { // note: chase will have recently nulled-out one
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
