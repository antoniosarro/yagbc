package processor

import (
	"testing"

	"github.com/antoniosarro/yagbc/internal/core/gb/memory"
)

// Helper function to create a CPU with a test program loaded
func setupCPU(program []byte) *CPU {
	mem := memory.NewBasicMemory()
	mem.LoadROM(program)
	return NewCPU(mem)
}

func TestOpNOP(t *testing.T) {
	// Program: NOP
	cpu := setupCPU([]byte{0x00})

	initialPC := cpu.Registers.PC
	cycles := cpu.Step()

	// NOP should take 4 cycles and advance PC by 1
	if cycles != 4 {
		t.Errorf("Expected 4 cycles, got %d", cycles)
	}
	if cpu.Registers.PC != initialPC+1 {
		t.Errorf("PC should advance by 1, got PC=%d", cpu.Registers.PC)
	}
}

func TestOpLD_A_n(t *testing.T) {
	// Program: LD A, 0x42
	cpu := setupCPU([]byte{0x3E, 0x42})

	cycles := cpu.Step()

	if cycles != 8 {
		t.Errorf("Expected 8 cycles, got %d", cycles)
	}
	if cpu.Registers.A != 0x42 {
		t.Errorf("Expected A=0x42, got A=0x%02X", cpu.Registers.A)
	}
	if cpu.Registers.PC != 2 {
		t.Errorf("Expected PC=2, got PC=%d", cpu.Registers.PC)
	}
}

func TestOpLD_B_n(t *testing.T) {
	// Program: LD B, 0x12
	cpu := setupCPU([]byte{0x06, 0x12})

	cpu.Step()

	if cpu.Registers.B != 0x12 {
		t.Errorf("Expected B=0x12, got B=0x%02X", cpu.Registers.B)
	}
}

func TestOpLD_C_n(t *testing.T) {
	// Program: LD C, 0x34
	cpu := setupCPU([]byte{0x0E, 0x34})

	cpu.Step()

	if cpu.Registers.C != 0x34 {
		t.Errorf("Expected C=0x34, got C=0x%02X", cpu.Registers.C)
	}
}

func TestOpLD_A_B(t *testing.T) {
	// Program: LD B, 0x99; LD A, B
	cpu := setupCPU([]byte{0x06, 0x99, 0x78})

	cpu.Step() // LD B, 0x99
	cpu.Step() // LD A, B

	if cpu.Registers.A != 0x99 {
		t.Errorf("Expected A=0x99, got A=0x%02X", cpu.Registers.A)
	}
	if cpu.Registers.B != 0x99 {
		t.Errorf("B should remain 0x99, got B=0x%02X", cpu.Registers.B)
	}
}

func TestOpLD_A_C(t *testing.T) {
	// Program: LD C, 0xAB; LD A, C
	cpu := setupCPU([]byte{0x0E, 0xAB, 0x79})

	cpu.Step() // LD C, 0xAB
	cpu.Step() // LD A, C

	if cpu.Registers.A != 0xAB {
		t.Errorf("Expected A=0xAB, got A=0x%02X", cpu.Registers.A)
	}
}

func TestOpADD_A_B(t *testing.T) {
	// Program: LD A, 0x05; LD B, 0x03; ADD A, B
	cpu := setupCPU([]byte{0x3E, 0x05, 0x06, 0x03, 0x80})

	cpu.Step() // LD A, 0x05
	cpu.Step() // LD B, 0x03
	cpu.Step() // ADD A, B

	if cpu.Registers.A != 0x08 {
		t.Errorf("Expected A=0x08, got A=0x%02X", cpu.Registers.A)
	}

	// No flags should be set for this simple add
	if cpu.Registers.GetFlagZ() || cpu.Registers.GetFlagH() || cpu.Registers.GetFlagC() {
		t.Error("No flags should be set for 5 + 3")
	}
	if cpu.Registers.GetFlagN() {
		t.Error("N flag should be 0 for addition")
	}
}

func TestOpADD_A_B_Zero(t *testing.T) {
	// Test: ADD resulting in zero sets Z flag
	// Program: LD A, 0x00; LD B, 0x00; ADD A, B
	cpu := setupCPU([]byte{0x3E, 0x00, 0x06, 0x00, 0x80})

	cpu.Step() // LD A, 0x00
	cpu.Step() // LD B, 0x00
	cpu.Step() // ADD A, B

	if cpu.Registers.A != 0x00 {
		t.Errorf("Expected A=0x00, got A=0x%02X", cpu.Registers.A)
	}
	if !cpu.Registers.GetFlagZ() {
		t.Error("Z flag should be set when result is zero")
	}
}

func TestOpADD_A_B_Carry(t *testing.T) {
	// Test: ADD with overflow sets C flag
	// Program: LD A, 0xFF; LD B, 0x01; ADD A, B
	cpu := setupCPU([]byte{0x3E, 0xFF, 0x06, 0x01, 0x80})

	cpu.Step() // LD A, 0xFF
	cpu.Step() // LD B, 0x01
	cpu.Step() // ADD A, B

	// 0xFF + 0x01 = 0x00 (with carry)
	if cpu.Registers.A != 0x00 {
		t.Errorf("Expected A=0x00, got A=0x%02X", cpu.Registers.A)
	}
	if !cpu.Registers.GetFlagZ() {
		t.Error("Z flag should be set (result is zero)")
	}
	if !cpu.Registers.GetFlagC() {
		t.Error("C flag should be set (overflow occurred)")
	}
	if !cpu.Registers.GetFlagH() {
		t.Error("H flag should be set (half-carry occurred)")
	}
}

func TestOpADD_A_B_HalfCarry(t *testing.T) {
	// Test: ADD with half-carry sets H flag
	// Program: LD A, 0x0F; LD B, 0x01; ADD A, B
	cpu := setupCPU([]byte{0x3E, 0x0F, 0x06, 0x01, 0x80})

	cpu.Step() // LD A, 0x0F
	cpu.Step() // LD B, 0x01
	cpu.Step() // ADD A, B

	// 0x0F + 0x01 = 0x10 (half-carry from bit 3 to 4)
	if cpu.Registers.A != 0x10 {
		t.Errorf("Expected A=0x10, got A=0x%02X", cpu.Registers.A)
	}
	if !cpu.Registers.GetFlagH() {
		t.Error("H flag should be set (half-carry occurred)")
	}
	if cpu.Registers.GetFlagC() {
		t.Error("C flag should NOT be set (no full overflow)")
	}
}

func TestOpJP_nn(t *testing.T) {
	// Program: JP 0x0150
	// Bytes: 0xC3, 0x50, 0x01 (little-endian!)
	cpu := setupCPU([]byte{0xC3, 0x50, 0x01})

	cycles := cpu.Step()

	if cycles != 16 {
		t.Errorf("Expected 16 cycles, got %d", cycles)
	}
	if cpu.Registers.PC != 0x0150 {
		t.Errorf("Expected PC=0x0150, got PC=0x%04X", cpu.Registers.PC)
	}
}
