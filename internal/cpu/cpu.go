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
	return uint8(cpu.statusRegister & 0xFF)
}

func (cpu *CPU) SetConditionCodeRegister(value uint8) {
	cpu.statusRegister = (cpu.statusRegister & 0xFF00) | uint16(value)
}

func (cpu *CPU) GetRegisterValue(registerType RegisterType) uint32 {
	var registerValue uint32
	if registerType == RegisterSR {
		registerValue = uint32(cpu.GetStatusRegister())
	} else if registerType == RegisterCCR {
		registerValue = uint32(cpu.GetConditionCodeRegister())
	} else {
		panic("unimplmented")
	}

	return registerValue
}

func (cpu *CPU) fetchAndExecute() error {
	opcode, err := cpu.memory.ReadWordAt(uint(cpu.programCounter))
	if err != nil {
		return err
	}

	switch OpcodeType(opcode) {
	case RESET:
		cpu.programCounter += 2
	case NOP:
		cpu.programCounter += 2
	case STOP:
		_, err = cpu.memory.ReadWordAt(uint(cpu.programCounter + 2))
		if err != nil {
			return err
		}

		cpu.programCounter += 4
	case RTE:
		cpu.programCounter += 2
	case RTS:
		cpu.programCounter += 2
	case TRAPV:
		cpu.programCounter += 2
	case RTR:
		cpu.programCounter += 2
	default:
		higherOpcode := opcode >> 8
		switch OpcodeType(higherOpcode) {
		case ORI:
			err := cpu.ORI(opcode)
			if err != nil {
				return err
			}
		case ANDI:
			err := cpu.ANDI(opcode)
			if err != nil {
				return err
			}

		default:
			return errors.New("unknown opcode")
		}
	}

	return nil
}

func (cpu *CPU) ORI(opcode uint16) error {
	lowerOpcode := opcode & 0xFF
	switch lowerOpcode {
	case 0x3c:
		operand, err := cpu.memory.ReadByteAt(uint(cpu.programCounter + 2))
		if err != nil {
			return err
		}

		cpu.SetConditionCodeRegister(cpu.GetConditionCodeRegister() | operand)

		cpu.programCounter += 4
	case 0x7c:
		operand, err := cpu.memory.ReadWordAt(uint(cpu.programCounter + 2))
		if err != nil {
			return err
		}

		cpu.SetStatusRegister(cpu.GetStatusRegister() | operand)
		cpu.programCounter += 4
	default:
		panic("not implemented")
	}

	return nil
}

func (cpu *CPU) ANDI(opcode uint16) error {
	lowerOpcode := opcode & 0xFF
	switch lowerOpcode {
	case 0x3c:
		operand, err := cpu.memory.ReadByteAt(uint(cpu.programCounter + 2))
		if err != nil {
			return err
		}

		cpu.SetConditionCodeRegister(cpu.GetConditionCodeRegister() & operand)

		cpu.programCounter += 4
	case 0x7c:
		operand, err := cpu.memory.ReadWordAt(uint(cpu.programCounter + 2))
		if err != nil {
			return err
		}

		cpu.SetStatusRegister(cpu.GetStatusRegister() & operand)
		cpu.programCounter += 4
	default:
		panic("not implemented")
	}

	return nil
}
