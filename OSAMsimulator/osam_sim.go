package osamsimulator

// ------ OSAM ------- //
type addr = string

type OSAM struct {
	counter int
	stash   []block
}

func CreateOSAM(stashSize int) *OSAM {
	o := &OSAM{}
	o.counter = 0
	o.stash = make([]block, stashSize)
	return o
}

// func (osam *OSAM) Alloc() addr {

// }

// func (osam *OSAM) Read(i addr) block {

// }

// func (osam *OSAM) Write(i addr) block {

// }
