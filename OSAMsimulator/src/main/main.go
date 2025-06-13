package main

import (
	"fmt"
	"src/osam"
)

func main() {
	or := osam.CreateORAM(12, true)
	os := osam.CreateOSAM(or)
	a1 := os.Alloc()
	fmt.Printf("[main] Read: %v \n", osam.GetData(os.Read(a1)))
	a2 := os.Alloc()
	os.Write(a2, "data_a2")
	d2 := os.Read(a2)
	fmt.Printf("[main] Read: %v \n", osam.GetData(d2))
}
