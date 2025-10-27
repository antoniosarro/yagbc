package processor

// Registers represents the Sharp SM83 CPU register set.
//
// The Game Boy CPU has:
//   - 8 eight-bit registers: A, F, B, C, D, E, H, L
//   - 2 sixteen-bit registers: SP, PC
//
// The 8-bit registers can be accessed as pairs:
//   - AF (Accumulator + Flags)
//   - BC (general purpose)
//   - DE (general purpose)
//   - HL (often used for memory addressing)
type Registers struct {
	// 8-bit registers
	A uint8 // Accumulator - primary register for arithmetic
	F uint8 // Flags - holds CPU state (Z, N, H, C flags)
	B uint8 // General purpose
	C uint8 // General purpose
	D uint8 // General purpose
	E uint8 // General purpose
	H uint8 // High byte of HL pair (often used for addresses)
	L uint8 // Low byte of HL pair

	// 16-bit registers
	SP uint16 // Stack Pointer - points to top of stack
	PC uint16 // Program Counter - points to next instruction
}

// Flag bit positions in the F register
// The F register layout: [Z N H C 0 0 0 0]
// Only the top 4 bits are used; lower 4 bits are always 0
const (
	FlagZ uint8 = 0b10000000 // Zero flag (bit 7)
	FlagN uint8 = 0b01000000 // Subtract flag (bit 6)
	FlagH uint8 = 0b00100000 // Half-carry flag (bit 5)
	FlagC uint8 = 0b00010000 // Carry flag (bit 4)
)

// NewRegisters creates a Registers instance with Game Boy power-on values.
//
// Initial register state after boot ROM (which we're skipping):
//
//	AF = 0x01B0  (A=0x01, flags set to 0xB0)
//	BC = 0x0013
//	DE = 0x00D8
//	HL = 0x014D
//	SP = 0xFFFE  (top of memory)
//	PC = 0x0100  (where game code starts)
//
// For now, we'll use simpler defaults and set PC=0x0000.
func NewRegisters() *Registers {
	return &Registers{
		A:  0x00,
		F:  0x00,
		B:  0x00,
		C:  0x00,
		D:  0x00,
		E:  0x00,
		H:  0x00,
		L:  0x00,
		SP: 0xFFFE, // Stack grows down from top of memory
		PC: 0x0000, // Start at beginning of ROM
	}
}

// ========================================
// Register Pair Getters (16-bit values)
// ========================================

// AF returns the AF register pair (A in high byte, F in low byte).
// Result: 0xAF = (A << 8) | F
func (r *Registers) AF() uint16 {
	return uint16(r.A)<<8 | uint16(r.F)
}

// BC returns the BC register pair (B in high byte, C in low byte).
func (r *Registers) BC() uint16 {
	return uint16(r.B)<<8 | uint16(r.C)
}

// DE returns the DE register pair (D in high byte, E in low byte).
func (r *Registers) DE() uint16 {
	return uint16(r.D)<<8 | uint16(r.E)
}

// HL returns the HL register pair (H in high byte, L in low byte).
// This is the most commonly used pair for memory addressing.
func (r *Registers) HL() uint16 {
	return uint16(r.H)<<8 | uint16(r.L)
}

// ========================================
// Register Pair Setters (16-bit values)
// ========================================

// SetAF sets the AF register pair from a 16-bit value.
// High byte goes to A, low byte goes to F.
// Note: Lower 4 bits of F are masked to 0 (hardware behavior).
func (r *Registers) SetAF(value uint16) {
	r.A = uint8(value >> 8)   // High byte
	r.F = uint8(value) & 0xF0 // Low byte, mask lower 4 bits
}

// SetBC sets the BC register pair from a 16-bit value.
func (r *Registers) SetBC(value uint16) {
	r.B = uint8(value >> 8) // High byte
	r.C = uint8(value)      // Low byte
}

// SetDE sets the DE register pair from a 16-bit value.
func (r *Registers) SetDE(value uint16) {
	r.D = uint8(value >> 8) // High byte
	r.E = uint8(value)      // Low byte
}

// SetHL sets the HL register pair from a 16-bit value.
func (r *Registers) SetHL(value uint16) {
	r.H = uint8(value >> 8) // High byte
	r.L = uint8(value)      // Low byte
}

// ========================================
// Flag Getters (individual flag checks)
// ========================================

// GetFlag returns true if the specified flag is set.
// Usage: r.GetFlag(FlagZ) to check if Zero flag is set.
func (r *Registers) GetFlag(flag uint8) bool {
	return (r.F & flag) != 0
}

// GetFlagZ returns true if the Zero flag is set.
// Zero flag: Set when the result of an operation is zero.
func (r *Registers) GetFlagZ() bool {
	return r.GetFlag(FlagZ)
}

// GetFlagN returns true if the Subtract flag is set.
// Subtract flag: Set if the last operation was a subtraction.
func (r *Registers) GetFlagN() bool {
	return r.GetFlag(FlagN)
}

// GetFlagH returns true if the Half-carry flag is set.
// Half-carry: Set when carry occurs from bit 3 to bit 4.
func (r *Registers) GetFlagH() bool {
	return r.GetFlag(FlagH)
}

// GetFlagC returns true if the Carry flag is set.
// Carry flag: Set when carry occurs from bit 7 (overflow).
func (r *Registers) GetFlagC() bool {
	return r.GetFlag(FlagC)
}

// ========================================
// Flag Setters (set/clear individual flags)
// ========================================

// SetFlag sets the specified flag to the given value.
// If value is true, the flag is set (1), otherwise cleared (0).
func (r *Registers) SetFlag(flag uint8, value bool) {
	if value {
		r.F |= flag // Set the bit (OR operation)
	} else {
		r.F &^= flag // Clear the bit (AND NOT operation)
	}
}

// SetFlagZ sets or clears the Zero flag.
func (r *Registers) SetFlagZ(value bool) {
	r.SetFlag(FlagZ, value)
}

// SetFlagN sets or clears the Subtract flag.
func (r *Registers) SetFlagN(value bool) {
	r.SetFlag(FlagN, value)
}

// SetFlagH sets or clears the Half-carry flag.
func (r *Registers) SetFlagH(value bool) {
	r.SetFlag(FlagH, value)
}

// SetFlagC sets or clears the Carry flag.
func (r *Registers) SetFlagC(value bool) {
	r.SetFlag(FlagC, value)
}

// SetFlags sets all four flags at once.
// This is more efficient than setting them individually.
func (r *Registers) SetFlags(z, n, h, c bool) {
	r.F = 0 // Clear all flags first

	if z {
		r.F |= FlagZ
	}
	if n {
		r.F |= FlagN
	}
	if h {
		r.F |= FlagH
	}
	if c {
		r.F |= FlagC
	}
}
