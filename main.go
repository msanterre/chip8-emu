package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	PROGRAM_POS    uint16 = 0x200
	SCREEN_WIDTH   uint16 = 64
	SCREEN_HEIGHT  uint16 = 32
	DISPLAY_WIDTH  int32  = 640
	DISPLAY_HEIGHT int32  = 320
)

var (
	Surface    *sdl.Surface
	Window     *sdl.Window
	V          [16]uint16
	VI         uint16
	PC         uint16
	SP         byte
	Stack      [16]uint16
	Memory     [4096]byte
	GFX        [SCREEN_HEIGHT * SCREEN_WIDTH]byte
	DelayTimer byte
	SoundTimer byte
	GfxFlag    bool
	SoundFlag  bool
	Sprites    = [80]byte{
		0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2
		0xf0, 0x10, 0xf0, 0x10, 0xf0, // 3
		0x90, 0x90, 0xf0, 0x10, 0x10, // 4
		0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
		0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
		0xf0, 0x10, 0x20, 0x40, 0x40, // 7
		0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
		0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
		0xf0, 0x90, 0xf0, 0x90, 0x90, // A
		0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
		0xf0, 0x80, 0x80, 0x80, 0xf0, // C
		0xe0, 0x90, 0x90, 0x90, 0xe0, // D
		0xf0, 0x80, 0xf0, 0x80, 0xf0, // E
		0xf0, 0x80, 0xf0, 0x80, 0x80} // F
	KeyPositions = [16]uint8{
		sdl.K_0, sdl.K_1, sdl.K_2, sdl.K_3,
		sdl.K_4, sdl.K_5, sdl.K_6, sdl.K_7,
		sdl.K_8, sdl.K_9, sdl.K_a, sdl.K_b,
		sdl.K_c, sdl.K_d, sdl.K_e, sdl.K_f,
	}
)

func main() {
	fmt.Println("Chip-8 Emulator")

	if len(os.Args) != 2 {
		fmt.Println("Usage: chip8 <path_to_rom>")
		return
	}

	fmt.Printf("\n\nStart Execution:\n\n")

	LoadSprites()
	LoadROM()

	if err := InitializeDrawingSurface(); err != nil {
		fmt.Println("Could not initialize drawing surface.")
		os.Exit(-1)
	}
	defer Close()

	Run()
}

func InitializeDrawingSurface() error {
	var err error

	sdl.Init(sdl.INIT_EVERYTHING)

	Window, err = sdl.CreateWindow("Chip-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, DISPLAY_WIDTH, DISPLAY_HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}

	Surface, err = Window.GetSurface()
	if err != nil {
		return err
	}

	ClearScreen()
	Window.UpdateSurface()

	return nil
}

func Close() {
	Window.Destroy()
	sdl.Quit()
}

func Run() {
	PC = PROGRAM_POS

	for {
		sdl.PumpEvents()

		opcode := (uint16(Memory[PC]) << 8) | uint16(Memory[PC+1])
		RunOpcode(opcode)
		if GfxFlag {
			DrawScreen()
			Window.UpdateSurface()
			GfxFlag = false
		}
		UpdateTimers()
		sdl.Delay(2)
	}
}

func ClearScreen() {
	// Clear GFX values
	for i := range GFX {
		GFX[i] = 0
	}
	DrawScreen()
	PC += 2
}

func DrawScreen() {
	pixelWidth := int32(DISPLAY_WIDTH / int32(SCREEN_WIDTH))
	pixelHeight := int32(DISPLAY_HEIGHT / int32(SCREEN_HEIGHT))

	for y := uint16(0); y < SCREEN_HEIGHT; y++ {
		for x := uint16(0); x < SCREEN_WIDTH; x++ {
			position := y*SCREEN_WIDTH + x
			pixel := GFX[position]

			rect := sdl.Rect{int32(x) * pixelWidth, int32(y) * pixelHeight, pixelWidth, pixelHeight}

			if pixel == 1 {
				// fmt.Printf("x:%d\ny:%d\n\n", x, y)
				Surface.FillRect(&rect, 0xffffffff)
			} else {
				Surface.FillRect(&rect, 0)
			}
		}
	}
}

func LoadSprites() {
	copy(Memory[:], Sprites[:])
}

func UpdateTimers() {
	if DelayTimer > 0 {
		fmt.Println("Delay:", DelayTimer)
		DelayTimer--
	}

	if SoundTimer > 0 {
		fmt.Println("Sound:", SoundTimer)
		SoundTimer--
	}
}

func LoadROM() {
	data, err := os.ReadFile(os.Args[1])

	if err != nil {
		fmt.Println("Could not load ROM.")
		os.Exit(-1)
	}

	fmt.Printf("ROM is %d bytes\n", len(data))

	for i := uint16(0); i < uint16(len(data)); i += 2 {
		if i%16 == 0 {
			fmt.Println()
		}
		opcode := (uint16(data[i]) << 8) | uint16(data[i+1])

		Memory[i+PROGRAM_POS] = data[i]
		Memory[i+1+PROGRAM_POS] = data[i+1]

		fmt.Printf("%02x ", opcode)
	}
	fmt.Printf("\n\n")
}

func RunOpcode(opcode uint16) {
	fmt.Printf("Opcode: %02x\n\n", opcode)
	switch opcode & 0xf000 {
	case 0x0000:
		switch opcode {
		case 0x00e0:
			DisplayClear()
		case 0x00ee:
			SubReturn()
		default:
			CallRcaProgram(opcode)
		}
	case 0x1000:
		Jump(opcode & 0x0fff)
	case 0x2000:
		CallSubroutine(opcode & 0x0fff)
	case 0x3000:
		SkipIfEqual(opcode&0x0f00>>8, opcode&0x00ff)
	case 0x4000:
		SkipIfNotEqual(opcode&0x0f00>>8, opcode&0x00ff)
	case 0x5000:
		SkipIfVXEqualVY(opcode&0x0f00>>8, opcode&0x00f0>>4)
	case 0x6000:
		SetVX(opcode&0x0f00>>8, opcode&0x00ff)
	case 0x7000:
		AddVX(opcode&0x0f00>>8, opcode&0x00ff)
	case 0x8000:
		vx := opcode & 0x0f00 >> 8
		vy := opcode & 0x00f0 >> 4
		switch opcode & 0x000f {
		case 0x0000:
			SetVXToY(vx, vy)
		case 0x0001:
			SetVXToXorY(vx, vy)
		case 0x0002:
			SetVXToXandY(vx, vy)
		case 0x0003:
			SetVXToXxorY(vx, vy)
		case 0x0004:
			AddVYToVX(vx, vy)
		case 0x0005:
			SubstractVYFromVX(vx, vy)
		case 0x0006:
			ShiftRightVX(vx)
		case 0x0007:
			SetVXToVYMinusVX(vx, vy)
		case 0x000e:
			ShiftLeftVX(vx)
		}
	case 0x9000:
		SkipIfVXNotEqualVY(opcode&0x0f00>>8, opcode&0x00f0>>4)
	case 0xa000:
		SetIToAddr(opcode & 0x0fff)
	case 0xb000:
		JumpToAddrPlusV0(opcode & 0x0fff)
	case 0xc000:
		SetVXRandomAndVal(opcode&0x0f00>>8, opcode&0x00ff)
	case 0xd000:
		Draw(opcode&0x0f00>>8, opcode&0x00f0>>4, opcode&0x000f)
	case 0xe000:
		vx := opcode & 0x0f00 >> 8
		switch opcode & 0x00ff {
		case 0x009e:
			SkipIfVXPressed(vx)
		case 0x00a1:
			SkipIfVXUnpressed(vx)
		}
	case 0xf000:
		vx := opcode & 0x0f00 >> 8
		switch opcode & 0x00ff {
		case 0x0007:
			SetVXToDelayTimer(vx)
		case 0x000a:
			WaitAndSetKeypressToVX(vx)
		case 0x0015:
			SetDelayTimerToVX(vx)
		case 0x0018:
			SetSoundTimerToVX(vx)
		case 0x001E:
			AddVXToI(vx)
		case 0x0029:
			SetIToSpriteAddrInVX(vx)
		case 0x0033:
			SetBCD(vx)
		case 0x0055:
			RegDump(vx)
		case 0x0065:
			RegLoad(vx)
		}
	}
}

//////////////////
// Operations ///
////////////////

// Calls RCA 1802 program at address NNN. Not necessary for most ROMs.

func CallRcaProgram(addr uint16) { // 0NNN
	log.Fatal("Call RCA Program called")
}

// Clears the screen.
func DisplayClear() { // 00E0
	ClearScreen()
}

// Returns from a subroutine.
func SubReturn() { // 00EE
	PC = Stack[SP]
	SP -= 1
	PC += 2
}

// Jumps to address
func Jump(addr uint16) { // 1NNN
	PC = addr
}

// Calls subroutine at
func CallSubroutine(addr uint16) { // 2NNN
	SP += 1
	Stack[SP] = PC
	PC = addr
}

// Skips the next instruction if VX equals NN. (Usually the next instruction is a jump to skip a code block)
func SkipIfEqual(vx, value uint16) { //3XNN
	if V[vx] == value {
		PC += 4
	} else {
		PC += 2
	}
}

// Skips the next instruction if VX doesn't equal NN. (Usually the next instruction is a jump to skip a code block)
func SkipIfNotEqual(vx, value uint16) { // 4XNN
	if V[vx] != value {
		PC += 4
	} else {
		PC += 2
	}
}

// Skips the next instruction if VX equals VY. (Usually the next instruction is a jump to skip a code block)
func SkipIfVXEqualVY(vx, vy uint16) { // 5XNN
	if V[vx] == V[vy] {
		PC += 4
	} else {
		PC += 2
	}
}

// Sets VX to NN.
func SetVX(vx, value uint16) { // 6XNN
	V[vx] = value
	PC += 2
}

// Adds NN to VX.
func AddVX(vx, value uint16) { // 7XNN
	V[vx] += value
	PC += 2
}

// Sets VX to the value of VY.
func SetVXToY(vx, vy uint16) { // 8XY0
	V[vx] = V[vy]
	PC += 2
}

// Sets VX to VX or VY. (Bitwise OR operation)
func SetVXToXorY(vx, vy uint16) { // 8XY1
	V[vx] = V[vx] | V[vy]
	PC += 2
}

// Sets VX to VX and VY. (Bitwise AND operation)
func SetVXToXandY(vx, vy uint16) { // 8XY2
	V[vx] = V[vx] & V[vy]
	PC += 2
}

// Sets VX to VX xor VY.
func SetVXToXxorY(vx, vy uint16) { // 8XY3
	V[vx] = V[vx] ^ V[vy]
	PC += 2
}

// Adds VY to VX. VF is set to 1 when there's a carry, and to 0 when there isn't.
func AddVYToVX(vx, vy uint16) { // 8XY4
	add := V[vx] + V[vy]

	V[vx] = add & 0xff
	V[0xf] = (add >> 8) & 0xf

	PC += 2
}

// VY is subtracted from VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
func SubstractVYFromVX(vx, vy uint16) { // 8XY5
	var sub int32
	sub = int32(V[vx]) - int32(V[vy])

	if sub < 0 {
		sub = 0
		V[0xF] = 0
	} else {
		V[0xF] = 1
	}
	V[vx] = uint16(sub)

	PC += 2
}

// Shifts VX right by one. VF is set to the value of the least significant bit of VX before the shift.[2]
func ShiftRightVX(vx uint16) { // 8XY6
	if V[vx]&0x000f == 0x0001 {
		V[0xf] = 1
	} else {
		V[0xf] = 0
	}
	V[vx] /= 2
	PC += 2
}

// Sets VX to VY minus VX. VF is set to 0 when there's a borrow, and 1 when there isn't.
func SetVXToVYMinusVX(vx, vy uint16) { // 8XY7
	V[vx] = V[vy] - V[vx]
	PC += 2
}

// Shifts VX left by one. VF is set to the value of the most significant bit of VX before the shift.[2]
func ShiftLeftVX(vx uint16) { // 8XYE
	V[vx] = V[vx] << 1
	PC += 2
}

// Skips the next instruction if VX doesn't equal VY. (Usually the next instruction is a jump to skip a code block)
func SkipIfVXNotEqualVY(vx, vy uint16) { // 9XY0
	if V[vx] != V[vy] {
		PC += 4
	} else {
		PC += 2
	}
}

// Sets I to the address NNN.
func SetIToAddr(addr uint16) { // ANNN
	VI = addr
	PC += 2
}

// Jumps to the address NNN plus V0.
func JumpToAddrPlusV0(addr uint16) { // BNNN
	PC = addr + V[0]
}

// Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
func SetVXRandomAndVal(vx, value uint16) { // CXNN
	V[vx] = uint16(rand.Int()) & value & 0xff
	PC += 2
}

// Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
// Each row of 8 pixels is read as bit-coded starting from memory location I; I value doesn’t change after the execution of this instruction.
// As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn’t happen
func Draw(vx, vy, height uint16) { // DXYN
	var pixel byte

	x := V[vx]
	y := V[vy]

	V[0xf] = 0
	for yline := uint16(0); yline < height && yline+y < SCREEN_HEIGHT; yline++ {
		pixel = Memory[VI+yline]

		for xline := uint16(0); xline < 8; xline++ {
			if pixel&(0x80>>xline) != 0 {
				offset := (x + xline + ((y + yline) * SCREEN_WIDTH))

				if GFX[offset] == 1 {
					// VF is set to 1 if any screen pixels are flipped from
					// set to unset when the sprite is drawn, and to 0 if
					// that doesn't happen.
					V[0xF] = 1
				}
				GFX[offset] ^= 1
			}
		}
	}

	GfxFlag = true
	PC += 2
}

// Skips the next instruction if the key stored in VX is pressed. (Usually the next instruction is a jump to skip a code block)
func SkipIfVXPressed(vx uint16) { // EX9N
	i := V[vx]

	if sdl.GetKeyboardState()[KeyPositions[i]] == 1 {
		PC += 4
	} else {
		PC += 2
	}
}

// Skips the next instruction if the key stored in VX isn't pressed. (Usually the next instruction is a jump to skip a code block)
func SkipIfVXUnpressed(vx uint16) { // EXA1
	i := V[vx]

	if sdl.GetKeyboardState()[KeyPositions[i]] == 0 {
		PC += 4
	} else {
		PC += 2
	}
}

// Sets VX to the value of the delay timer.
func SetVXToDelayTimer(vx uint16) { // FX07
	V[vx] = uint16(DelayTimer)
	PC += 2
}

// A key press is awaited, and then stored in VX. (Blocking Operation. All instruction halted until next key event)
func WaitAndSetKeypressToVX(vx uint16) { // FX0A
	for {
		event := sdl.WaitEvent()
		switch t := event.(type) {
		case *sdl.KeyboardEvent:
			key := t.Keysym.Sym
			chipKey := GetChipKey(int(key))
			fmt.Println("Key down:", key, "Chippy: ", chipKey)

			if chipKey >= 0 {
				V[vx] = uint16(chipKey)
				PC += 2
				return
			}
		}
	}
}

// Sets the delay timer to VX.
func SetDelayTimerToVX(vx uint16) { // FX15
	DelayTimer = byte(V[vx])
	fmt.Println(DelayTimer)
	PC += 2
}

// Sets the sound timer to VX.
func SetSoundTimerToVX(vx uint16) { // FX18
	SoundTimer = byte(V[vx])
	PC += 2
}

// Adds VX to I.
func AddVXToI(vx uint16) { // FX1E
	VI += V[vx]
	PC += 2
}

// Sets I to the location of the sprite for the character in VX. Characters 0-F (in hexadecimal) are represented by a 4x5 font.
func SetIToSpriteAddrInVX(vx uint16) { // FX29
	VI = V[vx] * 5
	PC += 2
}

// Stores the binary-coded decimal representation of VX, with the most significant of three digits at the address in I, the middle digit at I plus 1, and the least significant digit at I plus 2.
func SetBCD(vx uint16) { // FX33
	val := V[vx]

	Memory[VI] = byte(val / 100 % 10)
	Memory[VI+1] = byte(val / 10 % 10)
	Memory[VI+2] = byte(val % 10)

	PC += 2
}

// Stores V0 to VX (including VX) in memory starting at address I.
func RegDump(vx uint16) { // FX55
	for i := uint16(0); i <= vx; i++ {
		Memory[VI+i] = byte(V[i])
	}
	PC += 2
}

// Fills V0 to VX (including VX) with values from memory starting at address I.
func RegLoad(vx uint16) { // FX65
	for i := uint16(0); i <= vx; i++ {
		V[i] = uint16(Memory[VI+i])
	}
	PC += 2
}

// /
// Utility
// /
func GetChipKey(sdlKey int) int {
	for i, key := range KeyPositions {
		if sdlKey == int(key) {
			return i
		}
	}
	return -1
}
