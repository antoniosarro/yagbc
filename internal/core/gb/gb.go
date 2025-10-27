// Package gb implements the Game Boy system emulation.
// It coordinates the CPU, memory, and other components to emulate
// the original Game Boy (DMG) hardware.
package gb

import (
	"github.com/antoniosarro/yagbc/internal/core/gb/memory"
	"github.com/antoniosarro/yagbc/internal/core/gb/processor"
)

// GameBoy represents the entire Game Boy system.
// It ties together all hardware components (CPU, memory, PPU, etc.)
type GameBoy struct {
	CPU    *processor.CPU
	Memory memory.Memory

	// TODO: Add more components (PPU, APU, Timers, etc.)
}

// NewGameBoy creates and initializes a new Game Boy system.
func NewGameBoy() *GameBoy {
	gb := &GameBoy{}

	// TODO: Initialize memory
	// TODO: Initialize CPU with memory reference

	return gb
}

// Step executes one machine cycle of the Game Boy.
// Returns the number of cycles that elapsed.
func (gb *GameBoy) Step() int {
	// TODO: Step the CPU
	// TODO: Step other components (PPU, timers, etc.)
	return 4 // Placeholder
}
