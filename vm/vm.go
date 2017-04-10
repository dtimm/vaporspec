package vm

import (
	"fmt"
	"os"
)

// VaporVM An instance of a VaporSpec VM
type VaporVM struct {
	Code   []uint16
	Rom    []uint16
	memory [256][256]uint16
	regs   [16]uint16 // Hard coded register count
	pc     uint16     // program counter
}

type Instruction struct {
	opcode uint16
	arg0   uint16
	arg1   uint16
	arg2   uint16
}

const JumpSegmentSize = 256

const (
	EXT      = iota
	EXT_HALT = iota
	EXT_CPY  = iota
	EXT_NOT  = iota
	EXT_LSL  = iota
	EXT_LSR  = iota
	EXT_JMP  = iota
	EXT_NOP  = iota
	ADD      = iota
	SUB      = iota
	ADDC     = iota
	SUBC     = iota
	CMP      = iota
	JLT      = iota
	JGT      = iota
	JEQ      = iota
	LDR      = iota
	STR      = iota
	LRC      = iota
	AND      = iota
	OR       = iota
	XOR      = iota
)

// CreateVM creates an initialized vaporspec VM.
func CreateVM(code, rom []uint16) VaporVM {
	var v VaporVM
	v.Code = code
	v.Rom = rom
	return v
}

// Run begins execution of code loaded in the VM
func (v *VaporVM) Run() {
	for i := range v.Code {
		i := decode(v.Code[i])

		fmt.Printf("Execute: %4X\n", i.opcode)
	}
}

func decode(u uint16) Instruction {
	var inst Instruction
	clean := uint16(0x000F)
	inst.opcode = u >> 12 & clean
	inst.arg0 = u >> 8 & clean
	inst.arg1 = u >> 4 & clean
	inst.arg2 = u & clean
	return inst
}

// Executes an instruction
func (v *VaporVM) exec(instr *Instruction) {
	switch instr.opcode {
	case EXT:

		switch instr.arg0 {
		case EXT_HALT:
			fmt.Printf("Exiting at halt instruction\n")
			os.Exit(0)
		case EXT_CPY:
			v.regs[instr.arg1] = v.regs[instr.arg2]
		case EXT_NOT:
			v.regs[instr.arg1] = ^(v.regs[instr.arg2])
		case EXT_LSL:
			v.regs[instr.arg0] = v.regs[instr.arg0] << v.regs[instr.arg1]
		case EXT_LSR:
			v.regs[instr.arg0] = v.regs[instr.arg0] >> v.regs[instr.arg1]
		case EXT_JMP:
			v.pc = (v.regs[instr.arg1] * JumpSegmentSize) + v.regs[instr.arg2] - 1
		case EXT_NOP:
			// No operation
		}

	case ADD:
		v.regs[instr.arg0] = v.regs[instr.arg1] + v.regs[instr.arg2]
	case SUB:
		v.regs[instr.arg0] = v.regs[instr.arg1] - v.regs[instr.arg2]
	case ADDC:
		v.regs[instr.arg0] += ((instr.arg1 << 4) & 0x00F0) + instr.arg2
	case SUBC:
		v.regs[instr.arg0] -= ((instr.arg1 << 4) & 0x00F0) + instr.arg2
	case CMP:
		if v.regs[instr.arg1] < v.regs[instr.arg2] {
			v.regs[instr.arg0] = 0
		} else if v.regs[instr.arg1] > v.regs[instr.arg2] {
			v.regs[instr.arg0] = 2
		} else {
			v.regs[instr.arg0] = 1
		}
	case JLT:
		if v.regs[instr.arg0] == 0 {
			v.pc = (v.regs[instr.arg1] * JumpSegmentSize) + v.regs[instr.arg2] - 1
		}
	case JGT:
		if v.regs[instr.arg0] == 2 {
			v.pc = (v.regs[instr.arg1] * JumpSegmentSize) + v.regs[instr.arg2] - 1
		}
	case JEQ:
		if v.regs[instr.arg0] == 1 {
			v.pc = (v.regs[instr.arg1] * JumpSegmentSize) + v.regs[instr.arg2] - 1
		}
	case LDR:
		v.regs[instr.arg0] = v.memory[v.regs[instr.arg1]][v.regs[instr.arg2]]
	case STR:
		if v.regs[instr.arg1] < 128 { // Segment is not part of ROM
			v.memory[v.regs[instr.arg1]][v.regs[instr.arg2]] = v.regs[instr.arg0]
		} else {
			fmt.Printf("Attempted illegal write to ROM\n")
			os.Exit(1)
		}
	case LRC:
		v.regs[instr.arg0] = ((instr.arg1 << 4) & 0x00F0) + instr.arg2
	case AND:
		v.regs[instr.arg0] = v.regs[instr.arg1] & v.regs[instr.arg2]
	case OR:
		v.regs[instr.arg0] = v.regs[instr.arg1] | v.regs[instr.arg2]
	case XOR:
		v.regs[instr.arg0] = v.regs[instr.arg1] ^ v.regs[instr.arg2]
	}
}
