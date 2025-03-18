package cpu

import (
	"fmt"
)

type OpcodeType uint16

const (
	ORI  OpcodeType = 0
	ANDI OpcodeType = iota

	RESET OpcodeType = 0x4e70
	NOP   OpcodeType = 0x4e71
	STOP  OpcodeType = 0x4e72
	RTE   OpcodeType = 0x4e73
	RTS   OpcodeType = 0x4e75
	TRAPV OpcodeType = 0x4e76
	RTR   OpcodeType = 0x4e77
)

type operandSize uint8

const (
	OpSizeByte operandSize = 0x00
	OpSizeWord operandSize = 0x01
	OpSizeLong operandSize = 0x02
)

func (opSize operandSize) isValid() bool {
	if opSize != OpSizeByte && opSize != OpSizeWord && opSize != OpSizeLong {
		return false
	}

	return true
}

func (opSize operandSize) getSizeInBytes() (uint8, error) {
	switch opSize {
	case OpSizeByte:
		return 1, nil
	case OpSizeWord:
		return 2, nil
	case OpSizeLong:
		return 4, nil
	default:
		return 0, fmt.Errorf("invalid operand size 0x%x", opSize)
	}
}

type RegisterType uint8

const (
	RegisterNone RegisterType = iota
	RegisterPC
	RegisterSP
	RegisterSR
	RegisterCCR
)
