package processor

import "testing"

func TestNewRegisters(t *testing.T) {
	regs := NewRegisters()

	// Check initial values
	if regs.PC != 0x0000 {
		t.Errorf("PC: expected 0x0000, got 0x%04X", regs.PC)
	}
	if regs.SP != 0xFFFE {
		t.Errorf("SP: expected 0xFFFE, got 0x%04X", regs.SP)
	}
}

func TestRegisterPairs(t *testing.T) {
	regs := NewRegisters()

	// Test BC pair
	regs.B = 0x12
	regs.C = 0x34
	bc := regs.BC()
	if bc != 0x1234 {
		t.Errorf("BC: expected 0x1234, got 0x%04X", bc)
	}

	// Test setting BC
	regs.SetBC(0xABCD)
	if regs.B != 0xAB || regs.C != 0xCD {
		t.Errorf("SetBC: expected B=0xAB C=0xCD, got B=0x%02X C=0x%02X", regs.B, regs.C)
	}

	// Test HL pair
	regs.H = 0xFF
	regs.L = 0x80
	hl := regs.HL()
	if hl != 0xFF80 {
		t.Errorf("HL: expected 0xFF80, got 0x%04X", hl)
	}

	// Test setting HL
	regs.SetHL(0xC000)
	if regs.H != 0xC0 || regs.L != 0x00 {
		t.Errorf("SetHL: expected H=0xC0 L=0x00, got H=0x%02X L=0x%02X", regs.H, regs.L)
	}
}

func TestFlags(t *testing.T) {
	regs := NewRegisters()

	// Initially all flags should be clear
	if regs.GetFlagZ() || regs.GetFlagN() || regs.GetFlagH() || regs.GetFlagC() {
		t.Error("Flags should be initially clear")
	}

	// Set Zero flag
	regs.SetFlagZ(true)
	if !regs.GetFlagZ() {
		t.Error("Zero flag should be set")
	}
	if regs.F != FlagZ {
		t.Errorf("F register: expected 0x%02X, got 0x%02X", FlagZ, regs.F)
	}

	// Set Carry flag (Zero should still be set)
	regs.SetFlagC(true)
	if !regs.GetFlagZ() || !regs.GetFlagC() {
		t.Error("Zero and Carry flags should both be set")
	}

	// Clear Zero flag (Carry should still be set)
	regs.SetFlagZ(false)
	if regs.GetFlagZ() || !regs.GetFlagC() {
		t.Error("Only Carry flag should be set")
	}

	// Test SetFlags (set all at once)
	regs.SetFlags(true, false, true, false) // Z and H set
	if !regs.GetFlagZ() || regs.GetFlagN() || !regs.GetFlagH() || regs.GetFlagC() {
		t.Error("Z and H flags should be set, N and C clear")
	}
}

func TestFRegisterMasking(t *testing.T) {
	regs := NewRegisters()

	// Set AF with lower 4 bits set (should be masked)
	regs.SetAF(0x12FF) // Trying to set F to 0xFF

	if regs.A != 0x12 {
		t.Errorf("A: expected 0x12, got 0x%02X", regs.A)
	}

	// Lower 4 bits of F should be masked to 0
	if regs.F != 0xF0 {
		t.Errorf("F: expected 0xF0 (masked), got 0x%02X", regs.F)
	}
}

func TestAFPair(t *testing.T) {
	regs := NewRegisters()

	regs.A = 0x01
	regs.SetFlags(true, false, true, true) // Z=1, H=1, C=1 = 0xB0

	af := regs.AF()
	expected := uint16(0x01B0)
	if af != expected {
		t.Errorf("AF: expected 0x%04X, got 0x%04X", expected, af)
	}
}
