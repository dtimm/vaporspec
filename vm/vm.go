package vm

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// VaporVM An instance of a VaporSpec VM
type VaporVM struct {
	Code   []uint16
	Rom    []uint16
	memory [256][256]uint16
	regs   [16]uint16 // Hard coded register count
	pc     uint16     // program counter
}

type instruction struct {
	opcode uint16
	arg0   uint16
	arg1   uint16
	arg2   uint16
}

const jumpSegmentSize = 256

const (
	ext     = 0x0
	extHalt = 0x0
	extCpy  = 0x1
	extNot  = 0x2
	extLsl  = 0x3
	extLsr  = 0x4
	extJmp  = 0x5
	extNop  = 0x6
	add     = 0x1
	sub     = 0x2
	addc    = 0x3
	subc    = 0x4
	cmp     = 0x5
	jlt     = 0x6
	jgt     = 0x7
	jeq     = 0x8
	ldr     = 0x9
	str     = 0xA
	lrc     = 0xB
	and     = 0xC
	or      = 0xD
	xor     = 0xE
)

// NewVaporVM creates an initialized vaporspec VM.
func NewVaporVM(code, rom []uint16) VaporVM {
	var v VaporVM
	v.Code = code
	v.Rom = rom
	return v
}

// Run begins execution of code loaded in the VM
func (v *VaporVM) Run() {
	key := make(chan int)

	go getKeys(key)
	tick := time.Tick(16 * time.Millisecond)

MainLoop:
	for {
		select {
		case pressed := <-key:
			// Do this shit
			if pressed == 113 {
				break MainLoop
			}
			fmt.Printf("Key: %d\n", pressed)
		case <-tick:
			// update display!
			v.updateDisplay()
		default:
			// Run execute steps.
			//fmt.Printf("v.pc=%X\nreg=%v", v.pc, v.regs)
			i := decode(v.Code[v.pc])
			//fmt.Printf("Execute: %X%X%X%X\n", i.opcode, i.arg0, i.arg1, i.arg2)
			v.exec(i)
			v.pc++

			// execute every 2us
			time.Sleep(2 * time.Microsecond)
		}
	}
}

func getKeys(ch chan int) {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadByte()
		value := int(input)
		if value == 113 {
			ch <- value
			break
		}
		ch <- value
	}
}

func decode(u uint16) instruction {
	var inst instruction
	clean := uint16(0x000F)
	inst.opcode = u >> 12 & clean
	inst.arg0 = u >> 8 & clean
	inst.arg1 = u >> 4 & clean
	inst.arg2 = u & clean
	return inst
}

func (v *VaporVM) updateDisplay() {

}

// Executes an instruction
func (v *VaporVM) exec(instr instruction) {
	switch instr.opcode {
	case ext:
		switch instr.arg0 {
		case extHalt:
			fmt.Printf("Exiting at halt instruction\n")
			os.Exit(0)
		case extCpy:
			v.regs[instr.arg1] = v.regs[instr.arg2]
		case extNot:
			v.regs[instr.arg1] = ^(v.regs[instr.arg2])
		case extLsl:
			v.regs[instr.arg0] = v.regs[instr.arg0] << v.regs[instr.arg1]
		case extLsr:
			v.regs[instr.arg0] = v.regs[instr.arg0] >> v.regs[instr.arg1]
		case extJmp:
			v.pc = (v.regs[instr.arg1] * jumpSegmentSize) + v.regs[instr.arg2] - 1
		case extNop:
			// No operation
		}
	case add:
		v.regs[instr.arg0] = v.regs[instr.arg1] + v.regs[instr.arg2]
	case sub:
		v.regs[instr.arg0] = v.regs[instr.arg1] - v.regs[instr.arg2]
	case addc:
		v.regs[instr.arg0] += ((instr.arg1 << 4) & 0x00F0) + instr.arg2
	case subc:
		v.regs[instr.arg0] -= ((instr.arg1 << 4) & 0x00F0) + instr.arg2
	case cmp:
		if v.regs[instr.arg1] < v.regs[instr.arg2] {
			v.regs[instr.arg0] = 0
		} else if v.regs[instr.arg1] > v.regs[instr.arg2] {
			v.regs[instr.arg0] = 2
		} else {
			v.regs[instr.arg0] = 1
		}
	case jlt:
		if v.regs[instr.arg0] == 0 {
			v.pc = (v.regs[instr.arg1] * jumpSegmentSize) + v.regs[instr.arg2] - 1
		}
	case jgt:
		if v.regs[instr.arg0] == 2 {
			v.pc = (v.regs[instr.arg1] * jumpSegmentSize) + v.regs[instr.arg2] - 1
		}
	case jeq:
		if v.regs[instr.arg0] == 1 {
			v.pc = (v.regs[instr.arg1] * jumpSegmentSize) + v.regs[instr.arg2] - 1
		}
	case ldr:
		v.regs[instr.arg0] = v.memory[v.regs[instr.arg1]][v.regs[instr.arg2]]
	case str:
		if v.regs[instr.arg1] < 128 { // Segment is not part of ROM
			v.memory[v.regs[instr.arg1]][v.regs[instr.arg2]] = v.regs[instr.arg0]
		} else {
			fmt.Printf("Attempted illegal write to ROM\n")
			os.Exit(1)
		}
	case lrc:
		v.regs[instr.arg0] = ((instr.arg1 << 4) & 0x00F0) + instr.arg2
	case and:
		v.regs[instr.arg0] = v.regs[instr.arg1] & v.regs[instr.arg2]
	case or:
		v.regs[instr.arg0] = v.regs[instr.arg1] | v.regs[instr.arg2]
	case xor:
		v.regs[instr.arg0] = v.regs[instr.arg1] ^ v.regs[instr.arg2]
	}
}
