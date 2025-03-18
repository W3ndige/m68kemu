package memory

import (
	"testing"
)

func TestMemory_ReadWriteByteAt(t *testing.T) {
	memory := NewMemory(0xFF)
	err := memory.WriteByteAt(0x16, 0xAA)
	if err != nil {
		t.Fatal(err)
	}

	value, err := memory.ReadByteAt(0x16)
	if err != nil {
		t.Fatal(err)
	}

	if value != 0xAA {
		t.Fatalf("expected 0xAA but got 0x%x", value)
	}
}

func TestMemory_OutOfRangeByteWrite(t *testing.T) {
	memory := NewMemory(0xF)
	err := memory.WriteByteAt(0x16, 0xAA)
	if err == nil {
		t.Fatalf("expected out of range error but got nil")
	}
}

func TestMemory_OutOfRangeByteWrite_2(t *testing.T) {
	memory := NewMemory(0xF)
	err := memory.WriteByteAt(0xF, 0xAA)
	if err == nil {
		t.Fatalf("expected out of range error but got nil")
	}
}

func TestMemory_OutOfRangeByteRead(t *testing.T) {
	memory := NewMemory(0xF)
	_, err := memory.ReadByteAt(0x16)
	if err == nil {
		t.Fatalf("expected out of range error but got nil")
	}
}

func TestMemory_OutOfRangeWordWrite(t *testing.T) {
	memory := NewMemory(0xF)
	err := memory.WriteWordAt(0xe, 0xAABB)
	if err == nil {
		t.Fatalf("expected out of range error but got nil")
	}
}
