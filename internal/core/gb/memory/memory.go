// Package memory implements the Game Boy memory system.
package memory

import "fmt"

// Game Boy Memory Map (16-bit address space = 64KB)

// Memory is the interface that all memory implementations must satisfy.
type Memory interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

// BasicMemory is a simple implementation of the Game Boy memory system.
// This is a simplified version for learning - it only includes:
//   - ROM area (0x0000-0x7FFF): 32KB
//   - WRAM (0xC000-0xDFFF): 8KB
//   - HRAM (0xFF80-0xFFFE): 127 bytes
//
// Other regions will return 0xFF (common behavior for unmapped memory).
type BasicMemory struct {
	// ROM - Read Only Memory (game code)
	// In a real Game Boy, this comes from the cartridge
	rom [0x8000]uint8 // 32KB: 0x0000-0x7FFF

	// WRAM - Work RAM (general purpose RAM)
	wram [0x2000]uint8 // 8KB: 0xC000-0xDFFF

	// HRAM - High RAM (fast RAM on CPU die)
	hram [0x7F]uint8 // 127 bytes: 0xFF80-0xFFFE

	// TODO Phase 2: Add VRAM, OAM, I/O registers, etc.
}

// NewBasicMemory creates a new BasicMemory instance.
// All memory is initialized to 0x00.
func NewBasicMemory() *BasicMemory {
	return &BasicMemory{}
	// Arrays are zero-initialized in Go, so all bytes start at 0x00
}

func (m *BasicMemory) Read(addr uint16) uint8 {
	switch {
	// ROM Area: 0x0000 - 0x7FFF (32KB)
	case addr <= 0x7FFF:
		return m.rom[addr]

	// WRAM: 0xC000 - 0xDFFF (8KB)
	case addr >= 0xC000 && addr <= 0xDFFF:
		// Subtract base address to get array index
		return m.wram[addr-0xC000]

	// Echo RAM: 0xE000 - 0xFDFF (mirror of 0xC000-0xDDFF)
	// The Game Boy mirrors WRAM here (except last 512 bytes)
	case addr >= 0xE000 && addr <= 0xFDFF:
		// Mirror of WRAM: redirect the read
		return m.wram[addr-0xE000]

	// HRAM: 0xFF80 - 0xFFFE (127 bytes)
	case addr >= 0xFF80 && addr <= 0xFFFE:
		return m.hram[addr-0xFF80]

	// Unmapped regions return 0xFF
	// This is typical behavior when reading from empty space
	default:
		return 0xFF
	}
}

// Write stores a byte at the given 16-bit address.
// This implements the Memory interface.
func (m *BasicMemory) Write(addr uint16, val uint8) {
	switch {
	// ROM Area: 0x0000 - 0x7FFF
	// ROM is READ-ONLY, but we allow writes for loading programs
	// In a real Game Boy, writes here control memory banking
	case addr <= 0x7FFF:
		m.rom[addr] = val

	// WRAM: 0xC000 - 0xDFFF (8KB)
	case addr >= 0xC000 && addr <= 0xDFFF:
		m.wram[addr-0xC000] = val

	// Echo RAM: 0xE000 - 0xFDFF (writes go to WRAM)
	case addr >= 0xE000 && addr <= 0xFDFF:
		m.wram[addr-0xE000] = val

	// HRAM: 0xFF80 - 0xFFFE (127 bytes)
	case addr >= 0xFF80 && addr <= 0xFFFE:
		m.hram[addr-0xFF80] = val

	// Writes to unmapped regions are ignored
	// (In a real emulator, we might log these for debugging)
	default:
		// Silently ignore for now
	}
}

// LoadROM loads a byte slice into ROM starting at address 0x0000.
// This is a helper function for testing - in Phase 2, we'll load
// actual ROM files from cartridges.
func (m *BasicMemory) LoadROM(data []byte) error {
	if len(data) > len(m.rom) {
		return fmt.Errorf("ROM data too large: %d bytes (max %d)", len(data), len(m.rom))
	}

	copy(m.rom[:], data)
	return nil
}

// DirectWrite writes directly to ROM without bounds checking.
// ONLY use this for setting up test programs!
// In a real Game Boy, ROM comes from the cartridge and can't be written.
func (m *BasicMemory) DirectWrite(addr uint16, val uint8) {
	if addr <= 0x7FFF {
		m.rom[addr] = val
	}
}
