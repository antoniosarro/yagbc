package memory

import "testing"

func TestBasicMemoryROM(t *testing.T) {
	mem := NewBasicMemory()

	// Write to ROM area
	mem.Write(0x0100, 0x42)

	// Read it back
	val := mem.Read(0x0100)
	if val != 0x42 {
		t.Errorf("Expected 0x42, got 0x%02X", val)
	}
}

func TestBasicMemoryWRAM(t *testing.T) {
	mem := NewBasicMemory()

	// Write to WRAM
	mem.Write(0xC000, 0xAB)
	mem.Write(0xDFFF, 0xCD) // Last byte of WRAM

	// Read back
	if mem.Read(0xC000) != 0xAB {
		t.Errorf("WRAM start: expected 0xAB, got 0x%02X", mem.Read(0xC000))
	}
	if mem.Read(0xDFFF) != 0xCD {
		t.Errorf("WRAM end: expected 0xCD, got 0x%02X", mem.Read(0xDFFF))
	}
}

func TestBasicMemoryHRAM(t *testing.T) {
	mem := NewBasicMemory()

	// Write to HRAM
	mem.Write(0xFF80, 0x11)
	mem.Write(0xFFFE, 0x22) // Last byte of HRAM

	if mem.Read(0xFF80) != 0x11 {
		t.Errorf("HRAM start: expected 0x11, got 0x%02X", mem.Read(0xFF80))
	}
	if mem.Read(0xFFFE) != 0x22 {
		t.Errorf("HRAM end: expected 0x22, got 0x%02X", mem.Read(0xFFFE))
	}
}

func TestEchoRAM(t *testing.T) {
	mem := NewBasicMemory()

	// Write to WRAM
	mem.Write(0xC100, 0x55)

	// Read from Echo RAM (should mirror WRAM)
	echo := mem.Read(0xE100)
	if echo != 0x55 {
		t.Errorf("Echo RAM: expected 0x55, got 0x%02X", echo)
	}

	// Write to Echo RAM
	mem.Write(0xE200, 0x66)

	// Should affect WRAM
	wram := mem.Read(0xC200)
	if wram != 0x66 {
		t.Errorf("Echo write to WRAM: expected 0x66, got 0x%02X", wram)
	}
}

func TestUnmappedRegions(t *testing.T) {
	mem := NewBasicMemory()

	// Reading unmapped regions should return 0xFF
	unmapped := []uint16{0x8000, 0x9FFF, 0xA000, 0xFE00, 0xFF00}

	for _, addr := range unmapped {
		val := mem.Read(addr)
		if val != 0xFF {
			t.Errorf("Unmapped 0x%04X: expected 0xFF, got 0x%02X", addr, val)
		}
	}
}

func TestLoadROM(t *testing.T) {
	mem := NewBasicMemory()

	// Create a small test program
	program := []byte{0x00, 0x3E, 0x42, 0xC3} // Some opcodes

	err := mem.LoadROM(program)
	if err != nil {
		t.Fatalf("LoadROM failed: %v", err)
	}

	// Verify the data was loaded
	for i, expected := range program {
		actual := mem.Read(uint16(i))
		if actual != expected {
			t.Errorf("Address 0x%04X: expected 0x%02X, got 0x%02X", i, expected, actual)
		}
	}
}
