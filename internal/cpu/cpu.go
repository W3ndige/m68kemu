package cpu

import (
	"errors"
	"m68kemu/internal/memory"
)

type CPU struct {
	dataRegisters    [8]uint32
	addressRegisters [8]uint32
	programCounter   uint32
	statusRegister   uint16

	memory *memory.Memory
}

func NewCPU(memory *memory.Memory) *CPU {
	return &CPU{memory: memory}
}

func (cpu *CPU) Init() {
	panic("not implemented")
}

func (cpu *CPU) GetDataRegister(index uint32) (uint32, error) {
	if index >= 8 {
		return 0, errors.New("data register index out of range")
	}
	return cpu.dataRegisters[index], nil
}

func (cpu *CPU) SetDataRegister(index uint32, value uint32) error {
	if index >= 8 {
		return errors.New("data register index out of range")
	}

	cpu.dataRegisters[index] = value
	return nil
}

func (cpu *CPU) GetAddressRegister(index uint32) (uint32, error) {
	if index >= 8 {
		return 0, errors.New("data register index out of range")
	}
	return cpu.addressRegisters[index], nil
}

func (cpu *CPU) SetAddressRegister(index uint32, value uint32) error {
	if index >= 8 {
		return errors.New("data register index out of range")
	}

	cpu.addressRegisters[index] = value
	return nil
}

func (cpu *CPU) GetStackPointer() uint32 {
	return cpu.addressRegisters[7]
}

func (cpu *CPU) SetStackPointer(value uint32) {
	cpu.addressRegisters[7] = value
}

func (cpu *CPU) GetStatusRegister() uint16 {
	return cpu.statusRegister
}

func (cpu *CPU) SetStatusRegister(value uint16) {
	cpu.statusRegister = value
}

func (cpu *CPU) GetConditionCodeRegister() uint8 {
	return uint8(cpu.statusRegister & 0xF)
}

func (cpu *CPU) SetConditionCodeRegister(value uint8) {
	panic("not implemented")
}

func (cpu *CPU) fetchAndExecute() (Instruction, error) {
	var instruction Instruction
	if uint(cpu.programCounter) >= cpu.memory.Size {
		return Instruction{name: "EOF"}, nil
	}

	opcode, err := cpu.memory.ReadWordAt(uint(cpu.programCounter))
	if err != nil {
		instruction.name = "EOF"
		return instruction, err
	}

	instruction.opcode = opcode
	instruction.address = uint(cpu.programCounter)

	switch opcode {
	case RESET:
		instruction.name = "RESET"
		cpu.programCounter += 2
	case NOP:
		instruction.name = "NOP"
		cpu.programCounter += 2
	case STOP:
		instruction.name = "STOP"
		instruction.operandSize = OP_SIZE_WORD
		_, err := instruction.ReadOperand(cpu.memory)
		if err != nil {
			return instruction, err
		}

		cpu.programCounter += 4
	case RTE:
		instruction.name = "RTE"
		cpu.programCounter += 2
	case RTS:
		instruction.name = "RTS"
		cpu.programCounter += 2
	case TRAPV:
		instruction.name = "TRAPV"
		cpu.programCounter += 2
	case RTR:
		instruction.name = "RTR"
		cpu.programCounter += 2
	default:
		higherOpcode := opcode >> 8
		switch higherOpcode {
		case 0x00:
			instruction.name = "ORI"
			err := instruction.parseORI(cpu.memory)
			if err != nil {
				return instruction, err
			}

			cpu.programCounter += 2 + uint32(instruction.operandSize+1)

		default:
			instruction.name = "INVALID"
			return instruction, errors.New("unknown opcode")
		}
	}

	return instruction, nil
}
