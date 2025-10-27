[Memory Map](https://gbdev.io/pandocs/Memory_Map.html)

```
Game Boy Memory Map (16-bit address space = 64KB)

0xFFFF ┌─────────────────────┐
       │  Interrupt Enable   │  (1 byte)
0xFF80 ├─────────────────────┤
       │   High RAM (HRAM)   │  127 bytes - Super fast!
       │                     │  Used for critical code
0xFF00 ├─────────────────────┤
       │   I/O Registers     │  Joypad, Serial, Timer,
       │  (Memory-Mapped)    │  Audio, Video controls
0xFE00 ├─────────────────────┤
       │   OAM (Sprite)      │  Sprite attribute table
       │      Memory         │  (160 bytes)
0xE000 ├─────────────────────┤
       │  Echo RAM (Mirror)  │  Mirror of 0xC000-0xDDFF
0xC000 ├─────────────────────┤
       │   Work RAM (WRAM)   │  8KB - General purpose
       │                     │  Your program's variables
0xA000 ├─────────────────────┤
       │  External RAM       │  Cartridge RAM (if present)
       │   (Cart RAM)        │  Save games stored here!
0x8000 ├─────────────────────┤
       │   Video RAM (VRAM)  │  Tile data, background maps
       │                     │  What gets drawn on screen
0x4000 ├─────────────────────┤
       │  ROM Bank 01-NN     │  Switchable ROM banks
       │   (Switchable)      │  (for large games)
0x0000 └─────────────────────┘
       │  ROM Bank 00        │  Always visible
       │  (Fixed/Boot)       │  Game starts here!
       └─────────────────────┘
```