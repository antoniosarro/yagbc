package processor

import "fmt"

// Opcode represents a single CPU instruction.
type Opcode struct {
	Mnemonic string     // Human-readable name (e.g., "LD A, B")
	Bytes    int        // Number of bytes (including opcode)
	Cycles   int        // Number of CPU cycles
	Execute  func(*CPU) // Function that performs the operation
}

// opcodeTable maps each opcode byte (0x00-0xFF) to its implementation.
var opcodeTable [256]Opcode

// init initializes the opcode table.
// This runs automatically when the package is imported.
func init() {
	initOpcodes()
}

// initOpcodes registers all opcodes in the lookup table.
func initOpcodes() {
	// Initialize all opcodes as "UNKNOWN" first
	for i := range 256 {
		opcodeTable[i] = Opcode{
			Mnemonic: fmt.Sprintf("UNKNOWN_0x%02X", i),
			Bytes:    1,
			Cycles:   4,
			Execute:  opUnknown,
		}
	}

	// 0x00: NOP - No Operation
	opcodeTable[0x00] = Opcode{
		Mnemonic: "NOP",
		Bytes:    1,
		Cycles:   4,
		Execute:  opNOP,
	}

	// 0x3E: LD A, n - Load immediate 8-bit value into A
	opcodeTable[0x3E] = Opcode{
		Mnemonic: "LD A, n",
		Bytes:    2,
		Cycles:   8,
		Execute:  opLD_A_n,
	}

	// 0x06: LD B, n - Load immediate 8-bit value into B
	opcodeTable[0x06] = Opcode{
		Mnemonic: "LD B, n",
		Bytes:    2,
		Cycles:   8,
		Execute:  opLD_B_n,
	}

	// 0x0E: LD C, n - Load immediate 8-bit value into C
	opcodeTable[0x0E] = Opcode{
		Mnemonic: "LD C, n",
		Bytes:    2,
		Cycles:   8,
		Execute:  opLD_C_n,
	}

	// 0x78: LD A, B - Copy register B into A
	opcodeTable[0x78] = Opcode{
		Mnemonic: "LD A, B",
		Bytes:    1,
		Cycles:   4,
		Execute:  opLD_A_B,
	}

	// 0x79: LD A, C - Copy register C into A
	opcodeTable[0x79] = Opcode{
		Mnemonic: "LD A, C",
		Bytes:    1,
		Cycles:   4,
		Execute:  opLD_A_C,
	}

	// 0x80: ADD A, B - Add B to A
	opcodeTable[0x80] = Opcode{
		Mnemonic: "ADD A, B",
		Bytes:    1,
		Cycles:   4,
		Execute:  opADD_A_B,
	}

	// 0xC3: JP nn - Jump to 16-bit address
	opcodeTable[0xC3] = Opcode{
		Mnemonic: "JP nn",
		Bytes:    3,
		Cycles:   16,
		Execute:  opJP_nn,
	}
}

// ============================================================
// OPCODE IMPLEMENTATIONS
// ============================================================

// opUnknown is called for unimplemented opcodes.
// In a real emulator, this would be an error, but for learning
// we'll just do nothing (like NOP).
func opUnknown(cpu *CPU) {
	// For now, treat unknown opcodes as NOP
	// In Phase 2, we might want to panic or log this
}

// ============================================================
// 0x00: NOP - No Operation
// ============================================================
// Does nothing. Just wastes 4 cycles.
// Used for timing, padding, or as a placeholder.
//
// Flags: None affected
// Cycles: 4
// Bytes: 1
func opNOP(cpu *CPU) {
	// Literally do nothing!
	// The cycle count is handled by Step()
}

// ============================================================
// 0x3E: LD A, n - Load immediate into A
// ============================================================
// Loads the next byte (immediate value) into register A.
//
// Example:
//
//	Memory: [0x3E] [0x42]
//	Result: A = 0x42
//
// Flags: None affected
// Cycles: 8
// Bytes: 2
func opLD_A_n(cpu *CPU) {
	cpu.Registers.A = cpu.fetchByte()
}

// ============================================================
// 0x06: LD B, n - Load immediate into B
// ============================================================
// Loads the next byte (immediate value) into register B.
//
// Flags: None affected
// Cycles: 8
// Bytes: 2
func opLD_B_n(cpu *CPU) {
	cpu.Registers.B = cpu.fetchByte()
}

// ============================================================
// 0x0E: LD C, n - Load immediate into C
// ============================================================
// Loads the next byte (immediate value) into register C.
//
// Flags: None affected
// Cycles: 8
// Bytes: 2
func opLD_C_n(cpu *CPU) {
	cpu.Registers.C = cpu.fetchByte()
}

// ============================================================
// 0x78: LD A, B - Copy B to A
// ============================================================
// Copies the value in register B to register A.
// B remains unchanged.
//
// Flags: None affected
// Cycles: 4
// Bytes: 1
func opLD_A_B(cpu *CPU) {
	cpu.Registers.A = cpu.Registers.B
}

// ============================================================
// 0x79: LD A, C - Copy C to A
// ============================================================
// Copies the value in register C to register A.
// C remains unchanged.
//
// Flags: None affected
// Cycles: 4
// Bytes: 1
func opLD_A_C(cpu *CPU) {
	cpu.Registers.A = cpu.Registers.C
}

// ============================================================
// 0x80: ADD A, B - Add B to A
// ============================================================
// Adds register B to register A and stores result in A.
// B remains unchanged.
//
// Flags affected:
//
//	Z: Set if result is zero
//	N: Reset (0) - this is an addition
//	H: Set if carry from bit 3 to bit 4
//	C: Set if carry from bit 7 (overflow)
//
// Cycles: 4
// Bytes: 1
func opADD_A_B(cpu *CPU) {
	a := cpu.Registers.A
	b := cpu.Registers.B
	result := a + b

	// Calculate flags
	cpu.Registers.SetFlagZ(result == 0)                // Zero flag
	cpu.Registers.SetFlagN(false)                      // Addition, so N=0
	cpu.Registers.SetFlagH((a&0x0F)+(b&0x0F) > 0x0F)   // Half-carry
	cpu.Registers.SetFlagC(uint16(a)+uint16(b) > 0xFF) // Carry

	// Store result
	cpu.Registers.A = result
}

// ============================================================
// 0xC3: JP nn - Jump to 16-bit address
// ============================================================
// Unconditional jump. Sets PC to the 16-bit address specified
// in the next two bytes (little-endian).
//
// Example:
//
//	Memory: [0xC3] [0x50] [0x01]
//	Result: PC = 0x0150
//
// Flags: None affected
// Cycles: 16
// Bytes: 3
func opJP_nn(cpu *CPU) {
	addr := cpu.fetchWord() // Read 16-bit address (little-endian)
	cpu.Registers.PC = addr // Jump to that address
}
