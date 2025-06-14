package main

import (
	"fmt"
	osam "src/osam_simulator"
)

type Block = osam.Block

func main() {
	or := osam.CreateORAM(12, true)
	os := osam.CreateOSAM(or, true)
	sp := osam.CreateSP(os, true)

	A := sp.New(Block{Data: "C", IsNone: false})
	B := sp.Copy(&A)
	valA := sp.Get(&A).Data
	fmt.Printf("\n [main] GET on pointer A = %v \n", valA)
	valB := sp.Get(&B).Data
	fmt.Printf("\n [main] GET on pointer B = %v \n", valB)
	// a1 := os.Alloc()
	// fmt.Printf("[main] Read: %v \n", osam.GetData(os.Read(a1)))
	// a2 := os.Alloc()
	// os.Write(a2, "data_a2")
	// d2 := os.Read(a2)
	// fmt.Printf("[main] Read: %v \n", osam.GetData(d2))

}
