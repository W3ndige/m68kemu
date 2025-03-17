package cpu

import (
	"m68kemu/internal/memory"
	"testing"
)

func TestStaticInstructions(t *testing.T) {
	byteArray := [...]uint8{
		0x4e, 0x70, 0x4e, 0x71,
		0x4e, 0x72, 0x13, 0x37,
		0x4e, 0x73, 0x4e, 0x75,
		0x4e, 0x76, 0x4e, 0x77,
	}

	instructions := [...]Instruction{
		{name: "RESET"},
		{name: "NOP"},
		{name: "STOP"},
		{name: "RTE"},
		{name: "RTS"},
		{name: "TRAPV"},
		{name: "RTR"},
	}

	memoryBlock := memory.NewMemoryFromArray(byteArray[:])

	cpu := NewCPU(memoryBlock)

	instruction, err := cpu.fetchAndExecute()
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	for {
		if i >= len(instructions) {
			break
		}

		if instruction.name != instructions[i].name {
			t.Fatalf("expected instruction name: %s, got %s", instructions[i].name, instruction.name)
		}

		instruction, err = cpu.fetchAndExecute()
		if err != nil {
			t.Fatal(err)
		}

		i++
	}

	if i != len(instructions) {
		t.Fatalf("unexpected number of instructions: %d", i)
	}

}

func TestORIInstructions(t *testing.T) {
	byteArray := [...]uint8{
		0x00, 0x3c, 0x12,
		0x00, 0x7c, 0x13, 0x37,
		0x00, 0x01, 0xff,
	}

	instructions := [...]Instruction{
		{name: "ORI", operandSize: OP_SIZE_BYTE, operand: 0x12, addressingMode: ADDR_MODE_IMMEDIATE},
		{name: "ORI", operandSize: OP_SIZE_WORD, operand: 0x1337, addressingMode: ADDR_MODE_IMMEDIATE},
		{name: "ORI", operandSize: OP_SIZE_BYTE, operand: 0xff, addressingMode: ADDR_MODE_DATA_REGISTER, registerField: 0x1},
	}

	memoryBlock := memory.NewMemoryFromArray(byteArray[:])

	cpu := NewCPU(memoryBlock)

	instruction, err := cpu.fetchAndExecute()
	if err != nil {
		t.Fatal(err)
	}

	i := 0
	for {
		if i >= len(instructions) {
			break
		}

		if instruction.name != instructions[i].name {
			t.Fatalf("expected instruction name: %s, got %s", instructions[i].name, instruction.name)
		}

		if instruction.operand != instructions[i].operand {
			t.Fatalf("expected instruction operand: 0x%x, got 0x%x", instructions[i].operand, instruction.operand)
		}

		if instruction.addressingMode != instructions[i].addressingMode {
			t.Fatalf("expected addressingMode 0x%x, got 0x%x", instructions[i].addressingMode, instructions[i].addressingMode)
		}

		instruction, err = cpu.fetchAndExecute()
		if err != nil {
			t.Fatal(err)
		}

		i++
	}

	if i != len(instructions) {
		t.Fatalf("unexpected number of instructions: %d", i)
	}

}
