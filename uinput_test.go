package uinput

import (
	"os"
	"strings"
	"testing"
	"unsafe"
)

// use go test -v -run TestFFEffectMemoryLayout to see what went wrong
// and yes this was a problem so making a test for it
func TestFFEffectMemoryLayout(t *testing.T) {
  var i FFEffect
  t.Log("ff_effect\t size\t offset") 
  t.Logf("Type\t\t %v\t %v", unsafe.Sizeof(i.Type), unsafe.Offsetof(i.Type))
  t.Logf("ID\t\t %v\t %v", unsafe.Sizeof(i.ID), unsafe.Offsetof(i.ID))
  t.Logf("Direction\t %v\t %v", unsafe.Sizeof(i.Direction), unsafe.Offsetof(i.Direction))
  t.Logf("Trigger\t\t %v\t %v", unsafe.Sizeof(i.Trigger), unsafe.Offsetof(i.Trigger))
  t.Logf("Replay\t\t %v\t %v", unsafe.Sizeof(i.Replay), unsafe.Offsetof(i.Replay))
  t.Logf("u\t\t %v\t %v", unsafe.Sizeof(i.u), unsafe.Offsetof(i.u))
  t.Logf("Whole Size\t %v", unsafe.Sizeof(i))
  if unsafe.Offsetof(i.u) != 16 {
    t.Fatalf("Expected Union Offset in FFEffect to be 16\nActual: %v", unsafe.Offsetof(i.u))
  }
  if unsafe.Sizeof(i) != 48 {
    t.Fatalf("Expected Size of FFEffect to be 48\nActual: %v", unsafe.Sizeof(i))
  }
}

func TestValidateDevicePathEmptyPathPanics(t *testing.T) {
	expected := "device path must not be empty"
	err := validateDevicePath("")
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestValidateDevicePathInvalidPathPanics(t *testing.T) {
	path := "/some/bogus/path"
	err := validateDevicePath(path)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
	}
}

func TestValidateUinputNameEmptyNamePanics(t *testing.T) {
	expected := "device name may not be empty"
	err := validateUinputName(nil)
	if err.Error() != expected {
		t.Fatalf("Expected: %s\nActual: %s", expected, err)
	}
}

func TestFailedDeviceFileCreationGeneratesError(t *testing.T) {
	expected := "could not open device file"
	_, err := createDeviceFile("/root/testfile")
	if err == nil || err.Error() != expected {
		t.Fatalf("expected error, but got none")
	}
}

func TestNonExistentDeviceFileCausesError(t *testing.T) {
	expected := "failed to write uidev struct to device file:"
	_, err := createUsbDevice(nil, uinputUserDev{})
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("got '%v', but expected '%v'", err.Error(), expected)
	}
}
