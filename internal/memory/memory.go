package memory

import (
	"errors"
	"fmt"
)

type Memory struct {
	Size    uint
	content []uint8
}

func NewMemory(size uint) *Memory {
	return &Memory{size, make([]uint8, size)}
}

func NewFromFile(filename string) error {
	panic("not implemented")
}

func NewMemoryFromArray(data []uint8) *Memory {
	return &Memory{
		Size:    uint(len(data)),
		content: data,
	}
}

func (memory *Memory) ReadByteAt(address uint) (uint8, error) {
	if address < 0 || address > memory.Size-1 {
		return 0, errors.New(fmt.Sprintf("address {:%x} out of range", address))
	}

	return memory.content[address], nil
}

func (memory *Memory) WriteByteAt(address uint, value uint8) error {
	if address < 0 || address > memory.Size-1 {
		return errors.New(fmt.Sprintf("address 0x%x out of range [0...0x%x]", address, memory.Size))
	}

	memory.content[address] = value
	return nil
}

func (memory *Memory) ReadWordAt(address uint) (uint16, error) {
	if address < 0 || address > memory.Size-2 {
		return 0, errors.New(fmt.Sprintf("address 0x%x out of range [0...0x%x]", address, memory.Size))
	}

	lowBytes := memory.content[address]
	highBytes := memory.content[address+1]

	return uint16(lowBytes)<<8 | uint16(highBytes), nil
}

func (memory *Memory) WriteWordAt(address uint, value uint16) error {
	if address < 0 || address > memory.Size-2 {
		return errors.New(fmt.Sprintf("address 0x%x out of range [0...0x%x]", address, memory.Size))
	}

	lowBytes := value & 0xFF
	highBytes := value >> 8

	memory.content[address] = uint8(lowBytes)
	memory.content[address+1] = uint8(highBytes)

	return nil
}
