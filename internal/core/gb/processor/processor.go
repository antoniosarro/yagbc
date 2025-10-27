package processor

import "github.com/antoniosarro/yagbc/internal/core/gb/memory"

// CPU represents the Sharp SM83 processor used in the Game Boy.
type CPU struct {
	Registers *Registers    // CPU registers (A, B, C, D, E, F, H, L, SP, PC)
	Memory    memory.Memory // Memory interface for reading/writing
	Halted    bool          // Is the CPU halted? (from HALT instruction)

	// Debug/stats
	TotalCycles uint64 // Total cycles executed (for debugging)
}

// NewCPU creates a new CPU instance connected to the given memory.
func NewCPU(mem memory.Memory) *CPU {
	return &CPU{
		Registers:   NewRegisters(),
		Memory:      mem,
		Halted:      false,
		TotalCycles: 0,
	}
}

// Step executes one CPU instruction (fetch-decode-execute cycle).
// Returns the number of cycles the instruction took.
func (cpu *CPU) Step() int {
	// If halted, do nothing (but still consume cycles)
	if cpu.Halted {
		return 4 // NOP-equivalent
	}

	// FETCH: Read the opcode at PC
	opcode := cpu.fetchByte()

	// DECODE & EXECUTE: Look up and execute the instruction
	instruction := opcodeTable[opcode]

	// Execute the instruction
	instruction.Execute(cpu)

	// Track total cycles (for debugging/stats)
	cpu.TotalCycles += uint64(instruction.Cycles)

	return instruction.Cycles
}

// fetchByte reads the byte at PC and increments PC.
// This is used to read the opcode and any immediate operands.
func (cpu *CPU) fetchByte() uint8 {
	value := cpu.Memory.Read(cpu.Registers.PC)
	cpu.Registers.PC++
	return value
}

// fetchWord reads a 16-bit value at PC (little-endian) and increments PC by 2.
// Little-endian means: low byte first, then high byte.
// Example: bytes [0x34, 0x12] = 0x1234
func (cpu *CPU) fetchWord() uint16 {
	low := cpu.fetchByte()  // Read low byte first
	high := cpu.fetchByte() // Read high byte second
	return uint16(high)<<8 | uint16(low)
}
