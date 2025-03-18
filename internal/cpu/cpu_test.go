package cpu

import (
	"m68kemu/internal/memory"
	"testing"
)

func TestORItoCCRInstructions(t *testing.T) {
	// Third opcode is an ANDI instruction, used only to clear the SR register
	byteArray := [...]uint8{
		0x00, 0x3c, 0x12, 0x00,
		0x00, 0x3c, 0x44, 0x00,
		0x01, 0x7c, 0x00, 0x00,
		0x00, 0x7c, 0x12, 0x34,
	}

	tests := []struct {
		registerType RegisterType
		value        uint32
	}{
		{RegisterCCR, 0x12},
		{RegisterCCR, 0x12 | 0x44},
		{RegisterSR, 0x00},
		{RegisterSR, 0x1234},
	}

	memoryBlock := memory.NewMemoryFromArray(byteArray[:])

	cpu := NewCPU(memoryBlock)

	for i, test := range tests {
		err := cpu.fetchAndExecute()
		if err != nil {
			t.Fatalf("%d TestORItoCCRInstructions: %v", i, err)
		}

		if cpu.GetRegisterValue(test.registerType) != test.value {
			t.Fatalf(
				"%d TestORItoCCRInstructions: expected 0x%x, got 0x%x", i, test.value, cpu.GetRegisterValue(test.registerType))
		}
	}

}
