# YAGBC - Yet Another GameBoy Color Emulator

A Game Boy / Game Boy Color emulator written in Go as a learning project. This emulator is being built from the ground up to understand the intricacies of emulation, computer architecture, and the Game Boy hardware.

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-in%20development-yellow.svg)]()

## 🎯 Project Goals

- **Learning First**: Every component is implemented with a focus on understanding the "why" behind each decision
- **Clean Code**: Well-documented, tested, and organized Go code
- **Accuracy**: Striving for cycle-accurate emulation of Game Boy hardware
- **Progressive Development**: Building in phases, from basic CPU to full system emulation

## 🚀 Current Status

### Phase 1: Foundation (Complete)

**What's Working:**
- ✅ Project structure with clean separation of concerns
- ✅ Memory system with proper address space mapping
- ✅ CPU register implementation (8-bit and 16-bit)
- ✅ Basic instruction execution (fetch-decode-execute cycle)
- ✅ 8 fundamental opcodes implemented
- ✅ CPU flag system (Z, N, H, C)
- ✅ Comprehensive test suite

**Implemented Opcodes:**
- `0x00` - NOP (No Operation)
- `0x3E` - LD A, n (Load immediate into A)
- `0x06` - LD B, n (Load immediate into B)
- `0x0E` - LD C, n (Load immediate into C)
- `0x78` - LD A, B (Copy B to A)
- `0x79` - LD A, C (Copy C to A)
- `0x80` - ADD A, B (Add B to A with flags)
- `0xC3` - JP nn (Unconditional jump)

### Key Components

**CPU (`internal/core/gb/processor/`)**
- Implements the Sharp SM83 processor (8080-like with Z80 influences)
- Handles instruction fetch, decode, and execution
- Manages registers, flags, and program counter

**Memory (`internal/core/gb/memory/`)**
- Implements Game Boy's 64KB address space
- Supports ROM, WRAM, HRAM regions
- Interface-based design for future expansion (MBCs, I/O registers)

**Game Boy System (`internal/core/gb/gb.go`)**
- Coordinates CPU, memory, and future components
- Will manage timing and synchronization

## 📖 Learning Resources

This project is built using these excellent resources:

- **[Pan Docs](https://gbdev.io/pandocs/)** - The most comprehensive Game Boy technical documentation
- **[Game Boy CPU Opcodes](https://www.pastraiser.com/cpu/gameboy/gameboy_opcodes.html)** - Complete opcode reference
- **[Game Boy Development Community](https://gbdev.io/)** - Tools, docs, and community

## 🧪 Example Usage

Here's a simple example of running a program on the emulator:
```go
package main

import (
    "github.com/antoniosarro/yagbc/internal/core/gb/processor"
    "github.com/antoniosarro/yagbc/internal/core/gb/memory"
)

func main() {
    // Create memory and CPU
    mem := memory.NewBasicMemory()
    cpu := processor.NewCPU(mem)

    // Load a simple program: 5 + 3 = 8
    program := []byte{
        0x3E, 0x05, // LD A, 5
        0x06, 0x03, // LD B, 3
        0x80,       // ADD A, B
    }
    mem.LoadROM(program)

    // Execute the program
    cpu.Step() // LD A, 5
    cpu.Step() // LD B, 3
    cpu.Step() // ADD A, B

    // Result: A = 8
    println("Result:", cpu.Registers.A) // Output: 8
}
```

## 🤝 Contributing

This is primarily a learning project, but contributions, suggestions, and feedback are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 Technical Details

### CPU Architecture

The Game Boy uses a **Sharp SM83** CPU, which is:
- Similar to Intel 8080 and Zilog Z80
- 8-bit data bus, 16-bit address bus
- Clock speed: 4.194304 MHz
- 8 general-purpose 8-bit registers (can be paired into 16-bit)
- 4 flag bits: Zero (Z), Subtract (N), Half-Carry (H), Carry (C)

### Memory Map
```
0xFFFF          Interrupt Enable Register
0xFF80 - 0xFFFE High RAM (HRAM) - 127 bytes
0xFF00 - 0xFF7F I/O Registers
0xFE00 - 0xFE9F Sprite Attribute Table (OAM)
0xE000 - 0xFDFF Echo RAM (mirror of 0xC000-0xDDFF)
0xC000 - 0xDFFF Work RAM (WRAM) - 8KB
0xA000 - 0xBFFF External RAM (Cartridge RAM)
0x8000 - 0x9FFF Video RAM (VRAM) - 8KB
0x4000 - 0x7FFF ROM Bank 01-NN (switchable)
0x0000 - 0x3FFF ROM Bank 00 (fixed)
```

### Instruction Cycle

Every instruction follows the fetch-decode-execute cycle:

1. **Fetch**: Read opcode byte at Program Counter (PC)
2. **Decode**: Look up instruction in opcode table
3. **Execute**: Perform the operation
4. **Update**: Increment PC, update flags, consume cycles

## 🐛 Known Issues

- Only basic opcodes implemented (8 of 256)
- No interrupt handling yet
- No graphics output
- No ROM loading from files
- Memory bank controllers not implemented

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

*Built with ❤️ and a lot of ☕ as a learning journey into emulation*