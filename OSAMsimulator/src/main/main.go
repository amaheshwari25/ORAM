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

func testBSP() {
	or := osam.CreateORAM(50, printORAM)
	os := osam.CreateOSAM(or, printOSAM)
	bsp := osam.CreateBSP(os, printSP)

	osam.Suppress()

	fmt.Println("\n[main] Create pointer A to Node with data='MYDATA'")
	A := bsp.New(Block{Data: "MYDATA", IsNone: false})

	fmt.Println("\n[main] Copy pointer A to create pointer B")
	B := bsp.Copy(&A)

	fmt.Println("\n[main] Copy pointer A to create pointer C")
	C := bsp.Copy(&A)

	fmt.Println("\n[main] Copy pointer A to create pointer D")
	D := bsp.Copy(&A)

	fmt.Println("\n[main] Copy pointer A to create pointers E, F, G, H")
	E := bsp.Copy(&A)
	F := bsp.Copy(&A)
	G := bsp.Copy(&A)
	H := bsp.Copy(&A)

	// osam.Unsupress()

	fmt.Println("\n[main] GET on pointer A")
	_ = bsp.Get(&A, true).Data
	fmt.Println("\n[main] GET on pointer B")
	_ = bsp.Get(&B, true).Data
	fmt.Println("\n[main] GET on pointer C")
	_ = bsp.Get(&C, true).Data
	fmt.Println("\n[main] GET on pointer D")
	_ = bsp.Get(&D, true).Data
	fmt.Println("\n[main] GET on pointer E")
	_ = bsp.Get(&E, true).Data
	fmt.Println("\n[main] GET on pointer F")
	_ = bsp.Get(&F, true).Data
	fmt.Println("\n[main] GET on pointer G")
	_ = bsp.Get(&G, true).Data
	fmt.Println("\n[main] GET on pointer H")
	val := bsp.Get(&H, true).Data
	fmt.Printf("[main] (IGNORE) RESULT: GET on pointer H = %v \n", val)
}

func testSP() {
	// Efficient (balanced) SP program
	or := osam.CreateORAM(50, printORAM)
	os := osam.CreateOSAM(or, printOSAM)
	sp := osam.CreateSP(os, printSP)

	osam.Suppress()

	fmt.Println("\n[main] Create pointer A to Node with data='MYDATA'")
	A := sp.New(Block{Data: "MYDATA", IsNone: false})

	fmt.Println("\n[main] Copy pointer A to create pointer B")
	B := sp.Copy(&A)

	fmt.Println("\n[main] Copy pointer A to create pointer C")
	C := sp.Copy(&A)

	fmt.Println("\n[main] Copy pointer A to create pointer D")
	D := sp.Copy(&A)

	fmt.Println("\n[main] Copy pointer A to create pointer E")
	E := sp.Copy(&A)

	// osam.Unsupress()
	fmt.Println("\n[main] GET on pointer A")
	_ = sp.Get(&A, true).Data
	fmt.Println("\n[main] GET on pointer B")
	_ = sp.Get(&B, true).Data
	fmt.Println("\n[main] GET on pointer C")
	_ = sp.Get(&C, true).Data
	fmt.Println("\n[main] GET on pointer D")
	_ = sp.Get(&D, true).Data
	fmt.Println("\n[main] GET on pointer E")
	_ = sp.Get(&E, true).Data
}

func main() {

	testBSP()

	// Basic SP program
	// or := osam.CreateORAM(12, printORAM)
	// os := osam.CreateOSAM(or, printOSAM)
	// sp := osam.CreateSP(os, printSP)

	// fmt.Println("\n[main] Create pointer A to Node with data='MYDATA'")
	// A := sp.New(Block{Data: "MYDATA", IsNone: false})

	// fmt.Println("\n[main] Copy pointer A to create pointer B")
	// B := sp.Copy(&A)

	// fmt.Println("\n[main] PUT 'newC' via pointer A")
	// sp.Put(&A, Block{Data: "newC", IsNone: false})

	// fmt.Println("\n[main] GET on pointer B")
	// valB := sp.Get(&B).Data
	// fmt.Printf("[main] RESULT: GET on pointer B = %v \n", valB)

	// OLD: OSAM-level program
	// a1 := os.Alloc()
	// fmt.Printf("[main] Read: %v \n", osam.GetData(os.Read(a1)))
	// a2 := os.Alloc()
	// os.Write(a2, "data_a2")
	// d2 := os.Read(a2)
	// fmt.Printf("[main] Read: %v \n", osam.GetData(d2))
}
