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
const printPath = true

// ------ OSAM: Smart Pointer frameworks ------
func testBSPBaseCase() {
	or := osam.CreateORAM(50, printORAM)
	os := osam.CreateOSAM(or, printOSAM)
	bsp := osam.CreateBSP(os, printSP, printPath)

	osam.Suppress()

	fmt.Println("\n[main] Create pointer A to Node with data='MYDATA'")
	A := bsp.New(Block{Data: "MYDATA", IsNone: false})

	fmt.Println("\n[main] Copy pointer A to create pointer B")
	B := bsp.Copy(&A)

	fmt.Println("\n[main] Copy pointer B to create pointer C")
	C := bsp.Copy(&B)

	fmt.Println("\n[main] Copy pointer A to create pointer D")
	D := bsp.Copy(&A)

	fmt.Println("\n[main] GET A")
	_ = bsp.Get(&A)

	fmt.Println("\n[main] GET B")
	_ = bsp.Get(&B)

	fmt.Println("\n[main] GET C")
	_ = bsp.Get(&C)

	fmt.Println("\n[main] GET D")
	_ = bsp.Get(&D)

	fmt.Println("\n[main] DELETE C")
	bsp.Delete(&C)

	fmt.Println("\n[main] GET A")
	_ = bsp.Get(&A)

	fmt.Println("\n[main] GET B")
	_ = bsp.Get(&B)

	fmt.Println("\n[main] GET D")
	_ = bsp.Get(&D)

	fmt.Println("\n[main] DELETE B")
	bsp.Delete(&B)

	fmt.Println("\n[main] GET A")
	_ = bsp.Get(&A)

	fmt.Println("\n[main] GET D")
	_ = bsp.Get(&D)

	fmt.Println("\n[main] DELETE D")
	bsp.Delete(&D)

	fmt.Println("\n[main] GET A")
	_ = bsp.Get(&A)

	fmt.Println("\n[main] DELETE A")
	bsp.Delete(&A)

}

func testBSP() {
	or := osam.CreateORAM(50, printORAM)
	os := osam.CreateOSAM(or, printOSAM)
	bsp := osam.CreateBSP(os, printSP, printPath)

	osam.Suppress()

	fmt.Println("\n[main] Create pointer A to Node with data='MYDATA'")
	A := bsp.New(Block{Data: "MYDATA", IsNone: false})

	fmt.Println("\n[main] Copy pointer A to create pointer B, C, D, E, F, G, H")
	B := bsp.Copy(&A)

	// fmt.Println("\n[main] Copy pointer A to create pointer C")
	C := bsp.Copy(&A)

	// fmt.Println("\n[main] Copy pointer A to create pointer D")
	D := bsp.Copy(&A)

	// fmt.Println("\n[main] Copy pointer A to create pointers E, F, G, H")
	E := bsp.Copy(&A)
	F := bsp.Copy(&A)
	G := bsp.Copy(&A)
	H := bsp.Copy(&A)

	// osam.Unsupress()

	fmt.Println("\n[main] GET on pointer A")
	_ = bsp.Get(&A).Data

	fmt.Println("\n[main] PUT 'MYDATA_B' with pointer B")
	bsp.Put(&B, Block{Data: "MYDATA_B", IsNone: false})

	fmt.Println("\n[main] GET on pointer C")
	valC := bsp.Get(&C).Data
	fmt.Printf("[main] RESULT: GET on pointer C = %v \n", valC)

	fmt.Println("\n[main] GET on pointer D")
	valD := bsp.Get(&D).Data
	fmt.Printf("[main] RESULT: GET on pointer D = %v \n", valD)

	fmt.Println("\n[main] PUT 'MYDATA_E' with pointer E")
	bsp.Put(&E, Block{Data: "MYDATA_E", IsNone: false})

	fmt.Println("\n[main] PUT 'MYDATA_F' with pointer F")
	bsp.Put(&F, Block{Data: "MYDATA_F", IsNone: false})

	fmt.Println("\n[main] GET on pointer G")
	valG := bsp.Get(&G).Data
	fmt.Printf("[main] RESULT: GET on pointer G = %v \n", valG)

	fmt.Println("\n[main] GET on pointer H")
	valH := bsp.Get(&H).Data
	fmt.Printf("[main] RESULT: GET on pointer H = %v \n", valH)

	fmt.Println("\n[main] Deleting pointers (A, B, G)")
	bsp.Delete(&A)
	bsp.Delete(&G)
	bsp.Delete(&B)
	// bsp.Delete(&H)

	fmt.Println("\n[main] Recreating B, G pointers via copy(D), copy(C)")
	B = bsp.Copy(&D)
	G = bsp.Copy(&C)

	fmt.Println("\n[main] Deleting C pointer")
	bsp.Delete(&C)

	fmt.Println("\n[main] GET on pointer G")
	_ = bsp.Get(&G)

	fmt.Println("\n[main] GET on pointer B")
	_ = bsp.Get(&B)

	fmt.Println("\n[main] GET on pointer D")
	_ = bsp.Get(&D)

	fmt.Println("\n[main] GET on pointer E")
	_ = bsp.Get(&E)

	fmt.Println("\n[main] GET on pointer F")
	_ = bsp.Get(&F)

	fmt.Println("\n[main] GET on pointer H")
	_ = bsp.Get(&H)

}

func testBasicSP() {
	or := osam.CreateORAM(12, printORAM)
	os := osam.CreateOSAM(or, printOSAM)
	sp := osam.CreateSP(os, printSP, printPath)

	fmt.Println("\n[main] Create pointer A to Node with data='DATA'")
	A := sp.New(Block{Data: "DATA", IsNone: false})

	fmt.Println("\n[main] Copy pointer A to create pointer B")
	B := sp.Copy(&A)

	fmt.Println("\n[main] PUT 'NEW_DATA' via pointer A")
	sp.Put(&A, Block{Data: "NEW_DATA", IsNone: false})

	fmt.Println("\n[main] GET on pointer B")
	valB := sp.Get(&B).Data
	fmt.Printf("[main] RESULT: GET on pointer B = %v \n", valB)
}

func testSP() {
	// Efficient (balanced) SP program
	or := osam.CreateORAM(50, printORAM)
	os := osam.CreateOSAM(or, printOSAM)
	sp := osam.CreateSP(os, printSP, printPath)

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

	fmt.Println("\n[main] Delete pointer C")
	sp.Delete(&C)

	// osam.Unsupress()
	fmt.Println("\n[main] GET on pointer A")
	_ = sp.Get(&A).Data
	fmt.Println("\n[main] GET on pointer B")
	_ = sp.Get(&B).Data
	// fmt.Println("\n[main] GET on pointer C")
	// _ = sp.Get(&C).Data
	fmt.Println("\n[main] GET on pointer D")
	_ = sp.Get(&D).Data
	fmt.Println("\n[main] GET on pointer E")
	_ = sp.Get(&E).Data
}

// ----------------------------------------------
// ----- OSAM: Emulated Graph construction ------
const printGr = true

func testGraph() {
	// NOTE: needs to include self-vertices
	// us := []int{1, 2, 3, 4, 5, 6, 7}
	// vs := []int{1, 1, 1, 1, 1, 1, 1}
	// ws := []int{osam.NONE, 13, 13, 13, 13, 13, 13}
	us := []int{1, 1, 1, 1, 1, 2, 2, 3, 4, 4, 4, 5, 5, 6, 6}
	vs := []int{1, 2, 3, 4, 6, 2, 3, 3, 4, 3, 5, 5, 3, 6, 4}
	ws := []int{osam.NONE, 13, 13, 13, 13, osam.NONE, 13, osam.NONE, osam.NONE, 13, 13, osam.NONE, 13, osam.NONE, 13}

	inp := osam.CreateInputGraph(printGr, us, vs, ws)
	og := inp.CreateOSAMGraph()

	for i := 0; i < len(og.Vtcs); i++ {
		vtx := og.Vtcs[i]
		fmt.Printf("Vtx @ index %v: %v @ address %v \n", i, og.FakeRAM[vtx], vtx)
		for vtx != osam.NONE {
			vtx = og.FakeRAM[vtx].UP
			fmt.Printf("%v @ address %v \n", og.FakeRAM[vtx], vtx)
		}
	}

}

// ----------------------------------------------

func main() {
	testGraph()

	// testBSPBaseCase()
	// testBSP()
	// testSP()
	// testBasicSP()

	// OLD: OSAM-level program
	// a1 := os.Alloc()
	// fmt.Printf("[main] Read: %v \n", osam.GetData(os.Read(a1)))
	// a2 := os.Alloc()
	// os.Write(a2, "data_a2")
	// d2 := os.Read(a2)
	// fmt.Printf("[main] Read: %v \n", osam.GetData(d2))
}
