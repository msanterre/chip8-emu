package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var (
	Opcodes    [123]uint16
	RegistersX [8]uint16
	PC         uint16
	RegisterI  uint16
	Memory     [3584]byte
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: chip8 <path_to_rom>")
		return
	}

	fmt.Printf("\n\nStart Execution:\n\n")

	LoadOpcodes()
	Run()
}

func Run() {
	PC = 0

	for {
		// RunOpcode[PC]
	}
}

func LoadOpcodes() {
	data, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		fmt.Println("Could not load ROM.")
	}

	fmt.Printf("ROM is %d bytes\n", len(data))

	for i := 0; i < len(data); i += 2 {
		if i%16 == 0 {
			fmt.Println()
		}
		opcode := (uint16(data[i]) << 8) | uint16(data[i+1])
		Opcodes[i/2] = opcode

		fmt.Printf("%x ", opcode)
	}
}

func RunOpcode(opcode uint16) {
	switch opcode & 0xf000 {
	case 0x1000:
		switch opcode {
		case 0x00e0:
			DisplayClear()
		case 0x0ee0:
			SubReturn()
		default:
			CallRcaProgram(opcode)
		}
	case 0x2000:
		CallSubroutine(opcode & 0x0fff)
	case 0x3000:
		SkipIfEqual(opcode&0x0f00, opcode&0x00ff)
	case 0x4000:
		SkipIfNotEqual(opcode&0x0f00, opcode&0x00ff)
	case 0x5000:
		SkipIfVXEqualV(opcode&0x0f00, opcode&0x00f0)
	case 0x6000:
		SetVX(opcode&0x0f00, opcode&0x00ff)
	case 0x7000:
		AddVX(opcode&0x0f00, opcode&0x00ff)
	case 0x8000:
		registerX := opcode & 0x0f00
		registerY := opcode & 0x00f0

		switch opcode & 0x0001 {
		case 0x0001:
			SetVXToY(registerX, registerY)
		case 0x0002:
			SetVXToXorY(registerX, registerY)
		case 0x0003:
			SetVXToXandY(registerX, registerY)
		case 0x0004:
			AddVYToVX(registerX, registerY)
		case 0x0005:
			SubstractVYFromVX(registerX, registerY)
		case 0x0006:
			ShiftRightVX(registerX)
		case 0x0007:
			SetVXToVYMinusVX(registerX, registerY)
		case 0x000e:
			ShiftLeftVX(registerX)
		}
	case 0x9000:
		SkipIfVXNotEqualVY(opcode&0x0f00, opcode&0x00f0)
	case 0xa000:
		SetIToAddr(opcode & 0x0fff)
	case 0xb000:
		JumpToAddrPlusV0(opcode & 0x0fff)
	case 0xc000:
		SetVXRandomAndVal(opcode & 0x0fff)
	case 0xd000:
		Draw(opcode&0x0f00, opcode&0x00f0, opcode&0x000f)
	case 0xe000:
		registerX := opcode & 0x0f00
		switch opcode & 0x00ff {
		case 0x0007:
			SkipIfVXPressed(registerX)
		case 0x000a:
			SkipIfVXUnpressed(registerX)
		}
	case 0xf000:
		registerX := opcode & 0x0f00
		switch opcode & 0x00ff {
		case 0x0015:
			SetDelayTimerToVX(registerX)
		case 0x0018:
			SetSoundTimerToVX(registerX)
		case 0x001E:
			AddVXTOI(registerX)
		case 0x0029:
			SetIToSpriteAddrInVX(registerX)
		case 0x0033:
			SetBCD(registerX)
		case 0x0055:
			RegDump(registerX)
		case 0x0065:
			RegLoad(registerX)
		}
	}
}

//////////////////
// Operations ///
////////////////

// Calls RCA 1802 program at address NNN. Not necessary for most ROMs.

func CallRcaProgram(addr uint16) { // 0NNN

}

// Clears the screen.

func DisplayClear() { // 00E0

}

// Returns from a subroutine.

func SubReturn() { // 00EE

}

// Jumps to address NNN.

func Goto(addr uint16) { // 1NNN

}

// Calls subroutine at NNN.

func CallSubroutine(addr uint16) { // 2NNN

}

// Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)

func SkipIfEqual(registerX, addr uint16) { //3XNN

}

// Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)

func SkipIfNotEqual(registerX, addr uint16) { // 4XNN

}

// Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)

func SkipIfVXEqualV(registerX, registerY uint16) { // 5XNN

}

// Sets VX to NN.

func SetVX(registerX, value uint16) { // 6XNN

}

// Adds NN to VX.

func AddVX(registerX, value uint16) { // 7XNN

}

// Sets VX to the value of VY.

func SetVXToY(registerX, registerY uint16) { // 8XY0

}

// Sets VX to VX or VY. (Bitwise OR operation)

func SetVXToXorY(registerX, registerY uint16) { // 8XY1

}

// Sets VX to VX and VY. (Bitwise AND operation)

func SetVXToXandY(registerX, registerY uint16) { // 8XY2

}

// Sets VX to VX xor VY.

func SetVXToXxorY(registerX, registerY uint16) { // 8XY3

}

// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.

func AddVYToVX(registerX, registerY uint16) { // 8XY4

}

// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.

func SubstractVYFromVX(registerX, registerY uint16) { //8XY5

}

// Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.[2]

func ShiftRightVX(registerX uint16) { // 8XY6

}

// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.

func SetVXToVYMinusVX(registerX, registerY uint16) { // 8XY7

}

// Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift.[2]

func ShiftLeftVX(registerX uint16) { // 8XYE

}

// Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)

func SkipIfVXNotEqualVY(registerX, registerY uint16) { // 9XY0

}

// Sets I to the address NNN.

func SetIToAddr(addr uint16) { // ANNN

}

// Jumps to the address NNN plus V0.
func JumpToAddrPlusV0(addr uint16) { // BNNN

}

// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and
func SetVXRandomAndVal(value uint16) { // CXNN

}

// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
// Each row of 8 pixels is read as bit-coded starting from memory location I; I value doesn’t change after the execution of this instruction.
// As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn’t happen
func Draw(registerX, registerY, val uint16) { // DXYN

}

// Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
func SkipIfVXPressed(registerX uint16) { // EX9N

}

// Skips the next instruction if the key stored in VX isn't pressed. (Usually the next instruction is a jump to skip a code block)
func SkipIfVXUnpressed(registerX uint16) { // EXA1

}

// Sets VX to the value of the delay timer.
func SetVXToDelayTimer(registerX uint16) { // FX07

}

// A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
func WaitAndSetKeypressToVX(registerX uint16) { // FX0A

}

// Sets the delay timer to VX.
func SetDelayTimerToVX(registerX uint16) { // FX15

}

// Sets the sound timer to VX.
func SetSoundTimerToVX(registerX uint16) { // FX18

}

// Adds VX to I.
func AddVXTOI(registerX uint16) { // FX1E

}

// Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
func SetIToSpriteAddrInVX(registerX uint16) { // FX29

}

// Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2.
// (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.)
func SetBCD(registerX uint16) { // FX33

}

// Stores V0 to VX (including VX) in memory starting at address I.
func RegDump(registerX uint16) { // FX55

}

// Fills V0 to VX (including VX) with values from memory starting at address I.
func RegLoad(registerX uint16) { // FX65

}
