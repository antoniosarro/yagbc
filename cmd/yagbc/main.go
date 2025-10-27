package main

import (
	"fmt"

	"github.com/antoniosarro/yagbc/internal/core/gb/memory"
	"github.com/antoniosarro/yagbc/internal/core/gb/processor"
)

func main() {
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║   Game Boy CPU - Instruction Showcase     ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	// Run multiple test programs
	testBasicArithmetic()
	testRegisterCopying()
	testJumpInstruction()
	testFlagBehavior()
}

// Test 1: Basic arithmetic
func testBasicArithmetic() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Test 1: Basic Arithmetic")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	mem := memory.NewBasicMemory()
	cpu := processor.NewCPU(mem)

	// Program: 10 + 20 = 30
	program := []byte{
		0x3E, 0x0A, // LD A, 10
		0x06, 0x14, // LD B, 20
		0x80, // ADD A, B
	}
	mem.LoadROM(program)

	fmt.Println("Program: LD A, 10 → LD B, 20 → ADD A, B")
	fmt.Println("Expected result: A = 30 (0x1E)")
	fmt.Println()

	cpu.Step()
	cpu.Step()
	cpu.Step()

	fmt.Printf("Result: A = %d (0x%02X)\n", cpu.Registers.A, cpu.Registers.A)

	if cpu.Registers.A == 30 {
		fmt.Println("✅ PASS")
	} else {
		fmt.Println("❌ FAIL")
	}
	fmt.Println()
}

// Test 2: Register copying
func testRegisterCopying() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Test 2: Register Copying")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	mem := memory.NewBasicMemory()
	cpu := processor.NewCPU(mem)

	// Program: Load values into B and C, then copy them to A
	program := []byte{
		0x06, 0x42, // LD B, 0x42
		0x78,       // LD A, B     (A should now be 0x42)
		0x0E, 0x99, // LD C, 0x99
		0x79, // LD A, C     (A should now be 0x99)
	}
	mem.LoadROM(program)

	fmt.Println("Program:")
	fmt.Println("  LD B, 0x42")
	fmt.Println("  LD A, B     (copy B to A)")
	fmt.Println("  LD C, 0x99")
	fmt.Println("  LD A, C     (copy C to A)")
	fmt.Println()

	cpu.Step() // LD B, 0x42
	fmt.Printf("After LD B, 0x42:  B=0x%02X\n", cpu.Registers.B)

	cpu.Step() // LD A, B
	fmt.Printf("After LD A, B:     A=0x%02X (copied from B)\n", cpu.Registers.A)

	cpu.Step() // LD C, 0x99
	fmt.Printf("After LD C, 0x99:  C=0x%02X\n", cpu.Registers.C)

	cpu.Step() // LD A, C
	fmt.Printf("After LD A, C:     A=0x%02X (copied from C)\n", cpu.Registers.A)

	if cpu.Registers.A == 0x99 && cpu.Registers.B == 0x42 && cpu.Registers.C == 0x99 {
		fmt.Println("✅ PASS - Register operations working correctly")
	} else {
		fmt.Println("❌ FAIL")
	}
	fmt.Println()
}

// Test 3: Jump instruction
func testJumpInstruction() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Test 3: Jump Instruction (JP)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	mem := memory.NewBasicMemory()
	cpu := processor.NewCPU(mem)

	// Program with a jump
	program := []byte{
		0x3E, 0x11, // 0x0000: LD A, 0x11
		0xC3, 0x08, 0x00, // 0x0002: JP 0x0008 (skip next instruction)
		0x3E, 0xFF, // 0x0005: LD A, 0xFF (this should be skipped!)
		0x3E, 0x22, // 0x0008: LD A, 0x22 (jump lands here)
	}
	mem.LoadROM(program)

	fmt.Println("Program:")
	fmt.Println("  0x0000: LD A, 0x11")
	fmt.Println("  0x0002: JP 0x0008      (jump forward)")
	fmt.Println("  0x0005: LD A, 0xFF     (skipped!)")
	fmt.Println("  0x0008: LD A, 0x22     (execution continues here)")
	fmt.Println()

	cpu.Step() // LD A, 0x11
	fmt.Printf("After LD A, 0x11:  A=0x%02X, PC=0x%04X\n", cpu.Registers.A, cpu.Registers.PC)

	cpu.Step() // JP 0x0008
	fmt.Printf("After JP 0x0008:   A=0x%02X, PC=0x%04X (jumped!)\n", cpu.Registers.A, cpu.Registers.PC)

	cpu.Step() // LD A, 0x22
	fmt.Printf("After LD A, 0x22:  A=0x%02X, PC=0x%04X\n", cpu.Registers.A, cpu.Registers.PC)

	// A should be 0x22, not 0xFF (which was skipped)
	if cpu.Registers.A == 0x22 && cpu.Registers.PC == 0x000A {
		fmt.Println("✅ PASS - Jump instruction working correctly")
	} else {
		fmt.Println("❌ FAIL")
	}
	fmt.Println()
}

// Test 4: Flag behavior with different additions
func testFlagBehavior() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("Test 4: CPU Flags (Z, N, H, C)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Test 4a: Zero flag
	fmt.Println("Test 4a: Zero Flag")
	fmt.Println("  0 + 0 should set Z flag")
	mem := memory.NewBasicMemory()
	cpu := processor.NewCPU(mem)
	program := []byte{
		0x3E, 0x00, // LD A, 0
		0x06, 0x00, // LD B, 0
		0x80, // ADD A, B
	}
	mem.LoadROM(program)
	cpu.Step()
	cpu.Step()
	cpu.Step()

	fmt.Printf("  Result: A=0x%02X, Z=%v, N=%v, H=%v, C=%v\n",
		cpu.Registers.A,
		cpu.Registers.GetFlagZ(),
		cpu.Registers.GetFlagN(),
		cpu.Registers.GetFlagH(),
		cpu.Registers.GetFlagC())

	if cpu.Registers.GetFlagZ() {
		fmt.Println("  ✅ Z flag correctly set")
	} else {
		fmt.Println("  ❌ Z flag should be set")
	}
	fmt.Println()

	// Test 4b: Half-carry flag
	fmt.Println("Test 4b: Half-Carry Flag")
	fmt.Println("  0x0F + 0x01 = 0x10 (carry from bit 3 to 4)")
	mem = memory.NewBasicMemory()
	cpu = processor.NewCPU(mem)
	program = []byte{
		0x3E, 0x0F, // LD A, 0x0F
		0x06, 0x01, // LD B, 0x01
		0x80, // ADD A, B
	}
	mem.LoadROM(program)
	cpu.Step()
	cpu.Step()
	cpu.Step()

	fmt.Printf("  Result: A=0x%02X, Z=%v, N=%v, H=%v, C=%v\n",
		cpu.Registers.A,
		cpu.Registers.GetFlagZ(),
		cpu.Registers.GetFlagN(),
		cpu.Registers.GetFlagH(),
		cpu.Registers.GetFlagC())

	if cpu.Registers.GetFlagH() {
		fmt.Println("  ✅ H flag correctly set")
	} else {
		fmt.Println("  ❌ H flag should be set")
	}
	fmt.Println()

	// Test 4c: Carry flag
	fmt.Println("Test 4c: Carry Flag")
	fmt.Println("  0xFF + 0x01 = 0x00 (carry from bit 7)")
	mem = memory.NewBasicMemory()
	cpu = processor.NewCPU(mem)
	program = []byte{
		0x3E, 0xFF, // LD A, 0xFF
		0x06, 0x01, // LD B, 0x01
		0x80, // ADD A, B
	}
	mem.LoadROM(program)
	cpu.Step()
	cpu.Step()
	cpu.Step()

	fmt.Printf("  Result: A=0x%02X, Z=%v, N=%v, H=%v, C=%v\n",
		cpu.Registers.A,
		cpu.Registers.GetFlagZ(),
		cpu.Registers.GetFlagN(),
		cpu.Registers.GetFlagH(),
		cpu.Registers.GetFlagC())

	if cpu.Registers.GetFlagC() && cpu.Registers.GetFlagZ() {
		fmt.Println("  ✅ C and Z flags correctly set (overflow to zero)")
	} else {
		fmt.Println("  ❌ C and Z flags should be set")
	}
	fmt.Println()

	// Test 4d: N flag (always 0 for ADD)
	fmt.Println("Test 4d: Subtract Flag (N)")
	fmt.Println("  N flag should always be 0 for ADD instructions")
	if !cpu.Registers.GetFlagN() {
		fmt.Println("  ✅ N flag correctly cleared (ADD is not subtraction)")
	} else {
		fmt.Println("  ❌ N flag should be 0 for ADD")
	}
	fmt.Println()

	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("All flag tests complete!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
