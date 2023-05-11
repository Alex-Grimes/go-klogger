package keylogger

import (
	"testing"
)

func TestFileDescriptor(t *testing.T) {
	k := &Keylogger{}

	err := k.Close()
	if err != nil {
		t.Error("Closing empty file descriptor should not yield an error", err)
		return
	}
}

func TestBufferParser(t *testing.T) {
	k := &Keylogger{}

	input, err := k.eventFromBuffer([]byte{138, 180, 84, 92, 0, 0, 0, 0, 62, 75, 0, 0, 0, 0, 0, 0, 4, 0, 4, 0, 30, 0, 0, 0})
	if err != nil {
		t.Error("Parsing buffer should not yield an error", err)
		return
	}

	if input == nil {
		t.Error("Parsing buffer should not yield a nil input")
		return
	}

	if input.KeyString() != "3" {
		t.Errorf("wrong input key. got %v, expected %v", input.KeyString(), "3")
		return
	}

	if input.Type != EvMsc {
		t.Errorf("wrong input type. expected %v", input.Type)
		return
	}
}
