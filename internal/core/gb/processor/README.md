```
┌─────────────────────────────────────────────────────┐
│           Game Boy CPU Registers (SM83)             │
├─────────────────────────────────────────────────────┤
│                                                     │
│  8-BIT REGISTERS (can be used individually):        │
│                                                     │
│     A  ┌──────┐  Accumulator - main work register   │
│        │      │  Most arithmetic happens here       │
│        └──────┘                                     │
│                                                     │
│     F  ┌──────┐  Flags - stores CPU state           │
│        │Z N H C│  Z=Zero, N=Subtract, H=Half-carry  │
│        └──────┘  C=Carry                            │
│                                                     │
│     B  ┌──────┐  General purpose                    │
│        │      │                                     │
│        └──────┘                                     │
│                                                     │
│     C  ┌──────┐  General purpose                    │
│        │      │                                     │
│        └──────┘                                     │
│                                                     │
│     D  ┌──────┐  General purpose                    │
│        │      │                                     │
│        └──────┘                                     │
│                                                     │
│     E  ┌──────┐  General purpose                    │
│        │      │                                     │
│        └──────┘                                     │
│                                                     │
│     H  ┌──────┐  High byte of HL (often used        │
│        │      │  for memory addressing)             │
│        └──────┘                                     │
│                                                     │
│     L  ┌──────┐  Low byte of HL                     │
│        │      │                                     │
│        └──────┘                                     │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  REGISTER PAIRS (two 8-bit regs = one 16-bit):      │
│                                                     │
│    AF  ┌──────┬──────┐  Accumulator + Flags         │
│        │  A   │  F   │                              │
│        └──────┴──────┘                              │
│                                                     │
│    BC  ┌──────┬──────┐  General purpose pair        │
│        │  B   │  C   │                              │
│        └──────┴──────┘                              │
│                                                     │
│    DE  ┌──────┬──────┐  General purpose pair        │
│        │  D   │  E   │                              │
│        └──────┴──────┘                              │
│                                                     │
│    HL  ┌──────┬──────┐  Memory pointer pair         │
│        │  H   │  L   │  (High/Low for addresses)    │
│        └──────┴──────┘                              │
│                                                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  16-BIT REGISTERS (separate, not pairs):            │
│                                                     │
│    SP  ┌────────────────┐  Stack Pointer            │
│        │                │  Points to top of stack   │
│        └────────────────┘  (grows downward!)        │
│                                                     │
│    PC  ┌────────────────┐  Program Counter          │
│        │                │  Points to next opcode    │
│        └────────────────┘  (where we are in code)   │
│                                                     │
└─────────────────────────────────────────────────────┘
```

## 🔑 Key Concepts

### 1. Register Pairs

You can access 8-bit registers **individually** OR **as 16-bit pairs**:
```
Example: BC register pair

As individual 8-bit registers:
B = 0x12
C = 0x34

As a 16-bit pair:
BC = 0x1234

┌─────────┬─────────┐
│ B (0x12)│ C (0x34)│  <- Two separate 8-bit values
└─────────┴─────────┘
     │         │
     └────┬────┘
          │
    BC = 0x1234          <- One 16-bit value
```

**Why is this useful?**
- 8-bit: `LD A, B` (copy B to A)
- 16-bit: `LD BC, 0x1234` (load 16-bit value into BC pair)
- 16-bit: `LD A, (BC)` (read from memory address pointed to by BC)

### 2. The Flag Register (F)

The F register is **special** - it holds 4 status flags that represent the CPU's current state:
```
F Register (8 bits, but only top 4 are used)
┌───┬───┬───┬───┬───┬───┬───┬───┐
│ Z │ N │ H │ C │ 0 │ 0 │ 0 │ 0 │
└───┴───┴───┴───┴───┴───┴───┴───┘
  7   6   5   4   3   2   1   0    <- Bit positions

Lower 4 bits are always 0!
```

**What each flag means:**

- **Z (Zero)** - Bit 7: Set to 1 if the last operation resulted in zero
  - Example: `SUB A, A` (subtract A from itself) → Result = 0 → Z flag = 1
  
- **N (Subtract)** - Bit 6: Set to 1 if the last operation was a subtraction
  - Used by DAA (Decimal Adjust) instruction
  
- **H (Half-Carry)** - Bit 5: Set to 1 if carry occurred from bit 3 to bit 4
  - Example: `0x0F + 0x01 = 0x10` → Lower nibble overflowed → H = 1
  - Also used by DAA for BCD arithmetic
  
- **C (Carry)** - Bit 4: Set to 1 if carry occurred from bit 7 (overflow)
  - Example: `0xFF + 0x01 = 0x00` (with carry) → C = 1
  - Also used for rotates and shifts

**Why flags matter:**
They enable **conditional instructions**:
- `JR Z, label` - Jump if Zero flag is set
- `CALL C, addr` - Call function if Carry flag is set

### 3. Stack Pointer (SP)

Points to the **top of the stack** in memory. The stack:
- Grows **downward** (toward lower addresses)
- Used for function calls, storing return addresses
- Typically starts at 0xFFFE (top of memory)
```
Stack grows DOWN:

0xFFFE  ┌─────┐  <- SP initially points here
        │     │
0xFFFC  ├─────┤  <- After PUSH, SP = 0xFFFC
        │data │
0xFFFA  ├─────┤  <- After another PUSH, SP = 0xFFFA
        │more │
        ├─────┤
        │     │
```

### 4. Fetch-Decode-Execute Cycle

```
    ┌─────────────────────────────────────┐
    │                                     │
    │    THE CPU'S INFINITE LOOP          │
    │                                     │
    └─────────────────────────────────────┘
              │
              ▼
       ┌─────────────┐
       │   FETCH     │  1. Read the opcode byte at PC
       │             │  2. Increment PC
       └──────┬──────┘
              │
              ▼
       ┌─────────────┐
       │   DECODE    │  3. Look up what this opcode means
       │             │  4. Determine how many bytes it needs
       └──────┬──────┘
              │
              ▼
       ┌─────────────┐
       │   EXECUTE   │  5. Perform the operation
       │             │  6. Update registers/memory/flags
       └──────┬──────┘
              │
              ▼
       ┌─────────────┐
       │   REPEAT    │  Loop forever!
       └──────┬──────┘
              │
              └─────────────┐
                            │
              ┌─────────────┘
              │
              ▼
       (back to FETCH)
```

### Real Example

Let's trace through a simple program:
```
Memory at 0x0000:  [0x3E] [0x42] [0x06] [0x03] [0x80]
                     ↑      ↑      ↑      ↑      ↑
                     │      │      │      │      │
Opcodes:          LD A,n    42   LD B,n   03   ADD A,B
```

**Step-by-step execution:**
```
PC = 0x0000
┌──────────┐
│  FETCH   │  Read byte at PC (0x0000) → opcode = 0x3E
│          │  PC becomes 0x0001
└──────────┘

┌──────────┐
│  DECODE  │  0x3E = "LD A, n" (load immediate into A)
│          │  This is a 2-byte instruction (opcode + data)
└──────────┘

┌──────────┐
│ EXECUTE  │  Read next byte at PC (0x0001) → 0x42
│          │  Set A = 0x42
│          │  PC becomes 0x0002
└──────────┘

PC = 0x0002
┌──────────┐
│  FETCH   │  Read byte at PC (0x0002) → opcode = 0x06
└──────────┘

┌──────────┐
│  DECODE  │  0x06 = "LD B, n" (load immediate into B)
└──────────┘

┌──────────┐
│ EXECUTE  │  Read next byte at PC (0x0003) → 0x03
│          │  Set B = 0x03
│          │  PC becomes 0x0004
└──────────┘

PC = 0x0004
┌──────────┐
│  FETCH   │  Read byte at PC (0x0004) → opcode = 0x80
└──────────┘

┌──────────┐
│  DECODE  │  0x80 = "ADD A, B" (add B to A)
└──────────┘

┌──────────┐
│ EXECUTE  │  A = A + B = 0x42 + 0x03 = 0x45
│          │  Update flags (Z=0, N=0, H=0, C=0)
│          │  PC becomes 0x0005
└──────────┘

Final result: A = 0x45, B = 0x03
```

### What is an Opcode?

An **opcode** (operation code) is a single byte that tells the CPU what to do:
```
0x3E = "Load the next byte into register A"
0x80 = "Add register B to register A"
0x00 = "Do nothing (NOP)"
```

The Game Boy has **256 possible opcodes** (0x00 to 0xFF), plus 256 more "CB-prefixed" opcodes (we'll add those later).

### Immediate Values

Some instructions need **extra data** that comes in the bytes following the opcode:
```
Memory:     [0x3E] [0x42]
             ↑      ↑
             │      └─ This is the "immediate value" (n)
             └──────── This is the opcode

Instruction: LD A, n  (Load immediate value into A)
Bytes: 2
Result: A = 0x42
```

### CPU Cycles

Every instruction takes **time** to execute, measured in **CPU cycles** (or "T-states"):
```
Instruction    | Bytes | Cycles | What it does
---------------|-------|--------|---------------------------
NOP            |   1   |   4    | Do nothing
LD A, n        |   2   |   8    | Load immediate into A
ADD A, B       |   1   |   4    | Add B to A
JP nn          |   3   |  16    | Jump to 16-bit address