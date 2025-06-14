package main

import (
	"fmt"
	osam "src/osam_simulator"
)

type Block = osam.Block

// ----------- toggle printing -----------
const printORAM = true // ORAM calls (oram_sim.go)
const printOSAM = true // OSAM calls (osam.go)
const printSP = true   // SmartPointer interface calls (smartpointers.go)

func main() {
	or := osam.CreateORAM(12, printORAM)
	os := osam.CreateOSAM(or, printOSAM)
	sp := osam.CreateSP(os, printSP)

	fmt.Println("\n[main] Create pointer A to Node with data='C'")
	A := sp.New(Block{Data: "C", IsNone: false})

	fmt.Println("\n[main] Copy pointer A to create pointer B")
	B := sp.Copy(&A)

	fmt.Println("\n[main] GET on pointer A ")
	valA := sp.Get(&A).Data
	fmt.Printf("[main] GET on pointer A result = %v \n", valA)

	fmt.Println("\n[main] PUT 'newC' via pointer A")
	sp.Put(&A, Block{Data: "newC", IsNone: false})

	fmt.Println("\n[main] GET on pointer B")
	valB := sp.Get(&B).Data
	fmt.Printf("[main] GET on pointer B result = %v \n", valB)

	// OLD: OSAM-level program
	// a1 := os.Alloc()
	// fmt.Printf("[main] Read: %v \n", osam.GetData(os.Read(a1)))
	// a2 := os.Alloc()
	// os.Write(a2, "data_a2")
	// d2 := os.Read(a2)
	// fmt.Printf("[main] Read: %v \n", osam.GetData(d2))
}
