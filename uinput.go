/*
Package uinput is a pure go package that provides access to the userland input device driver uinput on linux systems.
Virtual keyboard devices as well as virtual mouse input devices may be created using this package.
The keycodes and other event definitions, that are available and can be used to trigger input events,
are part of this package ("Key1" for number 1, for example).

In order to use the virtual keyboard, you will need to follow these three steps:

 1. Initialize the device
    Example: vk, err := CreateKeyboard("/dev/uinput", "Virtual Keyboard")

 2. Send Button events to the device
    Example (print a single D):
    err = vk.KeyPress(uinput.KeyD)

    Example (keep moving right by holding down right arrow key):
    err = vk.KeyDown(uinput.KeyRight)

    Example (stop moving right by releasing the right arrow key):
    err = vk.KeyUp(uinput.KeyRight)

 3. Close the device
    Example: err = vk.Close()

A virtual mouse input device is just as easy to create and use:

 1. Initialize the device:
    Example: vm, err := CreateMouse("/dev/uinput", "DangerMouse")

 2. Move the cursor around and issue click events
    Example (move mouse right):
    err = vm.MoveRight(42)

    Example (move mouse left):
    err = vm.MoveLeft(42)

    Example (move mouse up):
    err = vm.MoveUp(42)

    Example (move mouse down):
    err = vm.MoveDown(42)

    Example (trigger a left click):
    err = vm.LeftClick()

    Example (trigger a right click):
    err = vm.RightClick()

 3. Close the device
    Example: err = vm.Close()

If you'd like to use absolute input events (move the cursor to specific positions on screen), use the touch pad.
Note that you'll need to specify the size of the screen area you want to use when you initialize the
device. Here are a few examples of how to use the virtual touch pad:

 1. Initialize the device:
    Example: vt, err := CreateTouchPad("/dev/uinput", "DontTouchThis", 0, 1024, 0, 768)

 2. Move the cursor around and issue click events
    Example (move cursor to the top left corner of the screen):
    err = vt.MoveTo(0, 0)

    Example (move cursor to the position x: 100, y: 250):
    err = vt.MoveTo(100, 250)

    Example (trigger a left click):
    err = vt.LeftClick()

    Example (trigger a right click):
    err = vt.RightClick()

 3. Close the device
    Example: err = vt.Close()
*/
package uinput

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

func validateDevicePath(path string) error {
	if path == "" {
		return errors.New("device path must not be empty")
	}
	_, err := os.Stat(path)
	return err
}

func validateUinputName(name []byte) error {
	if name == nil || len(name) == 0 {
		return errors.New("device name may not be empty")
	}
	if len(name) > uinputMaxNameSize {
		return fmt.Errorf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)
	}
	return nil
}

func toUinputName(name []byte) (uinputName [uinputMaxNameSize]byte) {
	var fixedSizeName [uinputMaxNameSize]byte
	copy(fixedSizeName[:], name)
	return fixedSizeName
}

func createDeviceFile(path string) (fd *os.File, err error) {
  // Needs to be read and write for force-feedback support
	deviceFile, err := os.OpenFile(path, syscall.O_RDWR|syscall.O_NONBLOCK, 0660) 
	if err != nil {
		return nil, errors.New("could not open device file")
	}
	return deviceFile, err
}

func registerDevice(deviceFile *os.File, evType uintptr) error {
	err := ioctl(deviceFile, uiSetEvBit, evType)
	if err != nil {
		defer deviceFile.Close()
		err = releaseDevice(deviceFile)
		if err != nil {
			return fmt.Errorf("failed to close device: %v", err)
		}
		return fmt.Errorf("invalid file handle returned from ioctl: %v", err)
	}
	return nil
}

func createUsbDevice(deviceFile *os.File, dev uinputUserDev) (fd *os.File, err error) {
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, dev)
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to write user device buffer: %v", err)
	}
	_, err = deviceFile.Write(buf.Bytes())
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to write uidev struct to device file: %v", err)
	}

	err = ioctl(deviceFile, uiDevCreate, uintptr(0))
	if err != nil {
		_ = deviceFile.Close()
		return nil, fmt.Errorf("failed to create device: %v", err)
	}

	time.Sleep(time.Millisecond * 200)

	return deviceFile, err
}

func closeDevice(deviceFile *os.File) (err error) {
	err = releaseDevice(deviceFile)
	if err != nil {
		return fmt.Errorf("failed to close device: %v", err)
	}
	return deviceFile.Close()
}

func releaseDevice(deviceFile *os.File) (err error) {
	return ioctl(deviceFile, uiDevDestroy, uintptr(0))
}

func fetchSyspath(deviceFile *os.File) (string, error) {
	sysInputDir := "/sys/devices/virtual/input/"
	// 64 for name + 1 for null byte
	path := make([]byte, 65)
	err := ioctl(deviceFile, uiGetSysname, uintptr(unsafe.Pointer(&path[0])))

	sysInputDir = sysInputDir + string(path)
	return sysInputDir, err
}

// Note that mice and touch pads do have buttons as well. Therefore, this function is used
// by all currently available devices and resides in the main source file.
func sendBtnEvent(deviceFile *os.File, keys []int, btnState int) (err error) {
	for _, key := range keys {
		buf, err := inputEventToBuffer(inputEvent{
			Time:  syscall.Timeval{Sec: 0, Usec: 0},
			Type:  evKey,
			Code:  uint16(key),
			Value: int32(btnState)})
		if err != nil {
			return fmt.Errorf("key event could not be set: %v", err)
		}
		_, err = deviceFile.Write(buf)
		if err != nil {
			return fmt.Errorf("writing btnEvent structure to the device file failed: %v", err)
		}
	}
	return syncEvents(deviceFile)
}

// Currently only used for force-feedback support
// if the above is no longer true the code will need to change
// to allow for consuming events in multiple places
//
// if nothing was read and no errors occoured both returns will be nil
func readEvent(deviceFile *os.File) (*inputEvent, error) {
  var err error
  buf := make([]byte, 24)
  n, err := deviceFile.Read(buf)
  if err != nil {
    return nil, fmt.Errorf("reading input event from device file failed: %v", err)
  }
  if n == 0 {
    return nil, nil
  }
  iev, err := inputEventFromBuffer(buf) 

  if err != nil {
    return nil, fmt.Errorf("device file read failed on input event from buffer: %v", err)
  }
  return iev, nil
}

// Expose this function in your device for force-feedback support
// you only need this function if on device creation uinputUserDev.EffectMax is > 0 
// this function blocks* 
// the callback return will be placed into upload.ReturnValue
//
// Note on blocking: 
// if for some reason this function doesn't block it will return nil 
// and the callback will not be called but it is not a error state 
// on my machine i got mixed results while testing. this is likely caused by 
// this lib using a old version of go.
//
// IMPORTANT:
// on some old kernel versions telling you support force-feedback without 
// handling the callback it will hang 
// on newer versions the callback will timeout after 30 seconds without a hang
// so on device creation give the option to add force-feedback support
// 
// Read linux/uinput.h for how this callback works
func forceFeedbackCallback(deviceFile *os.File, callback func(upload *UInputFFUpload, erase *UInputFFErase) int32) error {
  var err error
  var ie *inputEvent = nil

  // readEvent should block but I got mixed results 
  // leaving like this because in the case it doesn't block
  // the user can deal with it better 
  // 
  // try to read an event
  ie, err = readEvent(deviceFile)
  if err != nil {
    return err
  }

  // return early in case of no events
  if ie == nil {
    return nil
  }

  // not on my watch
  if callback == nil {
    return fmt.Errorf("callback was nil")
  }

  // only handle the events we care about
  if ie.Type == evUinput {
    switch ie.Code {
    case uiFFUpload:
      var ffUp = UInputFFUpload{}
      ffUp.RequestID = uint32(ie.Value)
      err = ioctl(deviceFile, uiBeginFFUpload, uintptr(unsafe.Pointer(&ffUp)))
      if err != nil {
        return fmt.Errorf("begin ff upload ioctl failed: %v", err)
      }
      ffUp.ReturnValue = callback(&ffUp, nil)
      err = ioctl(deviceFile, uiEndFFUpload, uintptr(unsafe.Pointer(&ffUp)))
      if err != nil {
        return fmt.Errorf("end ff upload ioctl failed: %v", err)
      }
    case uiFFErase:
      var ffErs = UInputFFErase{}
      ffErs.RequestID = uint32(ie.Value)
      err = ioctl(deviceFile, uiBeginFFErase, uintptr(unsafe.Pointer(&ffErs)))
      if err != nil {
        return fmt.Errorf("begin ff erase ioctl failed: %v", err)
      }
      ffErs.ReturnValue = callback(nil, &ffErs)
      err = ioctl(deviceFile, uiEndFFErase, uintptr(unsafe.Pointer(&ffErs)))
      if err != nil {
        return fmt.Errorf("end ff erase ioctl failed: %v", err)
      }
    }
  }

  return nil
}

func syncEvents(deviceFile *os.File) (err error) {
	buf, err := inputEventToBuffer(inputEvent{
		Time:  syscall.Timeval{Sec: 0, Usec: 0},
		Type:  evSyn,
		Code:  uint16(synReport),
		Value: 0})
	if err != nil {
		return fmt.Errorf("writing sync event failed: %v", err)
	}
	_, err = deviceFile.Write(buf)
	return err
}

func inputEventToBuffer(iev inputEvent) (buffer []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 24))
	err = binary.Write(buf, binary.LittleEndian, iev)
	if err != nil {
		return nil, fmt.Errorf("failed to write input event to buffer: %v", err)
	}
	return buf.Bytes(), nil
}

// make sure the buffer capacity is 24
// and don't use buffer after this
func inputEventFromBuffer(buffer []byte) (_ *inputEvent, err error) {
	buf := bytes.NewBuffer(buffer)
  iev := inputEvent{}
	err = binary.Read(buf, binary.LittleEndian, &iev)
	if err != nil {
		return nil, fmt.Errorf("failed to read buffer to input event: %v", err)
	}
	return &iev, nil

}

// original function taken from: https://github.com/tianon/debian-golang-pty/blob/master/ioctl.go
func ioctl(deviceFile *os.File, cmd, ptr uintptr) error {
	_, _, errorCode := syscall.Syscall(syscall.SYS_IOCTL, deviceFile.Fd(), cmd, ptr)
	if errorCode != 0 {
		return errorCode
	}
	return nil
}
