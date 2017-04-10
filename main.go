package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	//"image"
	"log"
	"os"
	"vaporspec/vm"
)

func main() {
	// Set up flag handling.
	fileName := flag.String("f", "", "Input file")
	scale := flag.Int("s", 1.0, "Display scale")
	romName := flag.String("r", "", "ROM file")

	// Read in all flag values
	flag.Parse()

	// Exit with error if file doesn't open
	file := readVaporFile(*fileName, false, true)

	var rom []uint16
	if len(*romName) > 0 {
		// This is not a required field
		rom = readVaporFile(*romName, true, true)
	}

	// Create the display
	//width, height := 256*(*scale), 192*(*scale)
	//m := image.NewRGBA(image.Rect(0, 0, width, height))
	machine := vm.CreateVM(file, rom)

	fmt.Printf("File: %v, ROM: %v, Scale: %v\n", len(file), len(rom), scale)

	machine.Run()

	return
}

func readVaporFile(fileName string, isRom, printStatus bool) []uint16 {
	rom, err := os.Open(fileName)
	defer rom.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Read the first 2 bytes for the instruction count
	b := make([]byte, 2)
	n, err := rom.Read(b)

	if n != 2 {
		fmt.Println("Fuckin' wut.")
		os.Exit(1)
	}

	size := binary.LittleEndian.Uint16(b)

	// Read the declared number of instructions
	b = make([]byte, size*2)
	numRead, err := rom.Read(b)

	if numRead != int(size*2) {
		fmt.Printf("%s:\tRead %v, expected %v\n", fileName, numRead, size*2)
		os.Exit(1)
	}

	// Convert the instructions into a uint16 slice
	u := make([]uint16, size)

	for i := range u {
		u[i] = uint16(binary.BigEndian.Uint16(b[i*2 : (i+1)*2]))
	}

	if printStatus {
		fmt.Printf("Read %v bytes from %s.\n", n, fileName)
	}

	return u
}
