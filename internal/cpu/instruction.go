package cpu

import (
	"errors"
	"fmt"
	"m68kemu/internal/memory"
)

const (
	RESET = 0x4e70
	NOP   = 0x4e71
	STOP  = 0x4e72
	RTE   = 0x4e73
	RTS   = 0x4e75
	TRAPV = 0x4e76
	RTR   = 0x4e77

	OP_SIZE_BYTE = 0x00
	OP_SIZE_WORD = 0x01
	OP_SIZE_LONG = 0x02

	ADDR_MODE_DATA_REGISTER    = 0x0
	ADDR_MODE_ADDRESS_REGISTER = 0x1
	ADDR_MODE_IMMEDIATE        = 0x7
)

type Instruction struct {
	address        uint
	opcode         uint16
	name           string
	operand        uint32
	operandSize    uint16
	addressingMode uint8
	registerField  uint8
}

func (instr *Instruction) ReadOperand(memory *memory.Memory) (uint32, error) {
	if instr.operandSize != OP_SIZE_BYTE && instr.operandSize != OP_SIZE_WORD && instr.operandSize != OP_SIZE_LONG {
		return 0, errors.New(fmt.Sprintf("operand size 0x%x invalid", instr.operandSize))
	}

	if instr.operandSize == OP_SIZE_BYTE {
		operand, err := memory.ReadByteAt(instr.address + 2)
		if err != nil {
			return 0, err
		}
		instr.operand = uint32(operand)
	} else if instr.operandSize == OP_SIZE_WORD {
		operand, err := memory.ReadWordAt(instr.address + 2)
		if err != nil {
			return 0, err
		}

		instr.operand = uint32(operand)
	} else {
		return 0, errors.New("operand size is invalid")
	}

	return instr.operand, nil
}

func (instr *Instruction) parseORI(memory *memory.Memory) error {
	if instr.opcode == 0 {
		return errors.New("invalid opcode for ORI")
	}

	lowerOpcode := instr.opcode & 0xff
	switch lowerOpcode {
	case 0x3c:
		instr.operandSize = OP_SIZE_BYTE
		instr.addressingMode = ADDR_MODE_IMMEDIATE
		_, err := instr.ReadOperand(memory)
		if err != nil {
			return err
		}
	case 0x7c:
		instr.operandSize = OP_SIZE_WORD
		instr.addressingMode = ADDR_MODE_IMMEDIATE
		_, err := instr.ReadOperand(memory)
		if err != nil {
			return err
		}
	default:
		instr.operandSize = lowerOpcode >> 6
		_, err := instr.ReadOperand(memory)
		if err != nil {
			return err
		}

		instr.addressingMode = uint8((lowerOpcode >> 3) & 0b111)
		instr.registerField = uint8(lowerOpcode & 0b111)
	}

	return nil
}
