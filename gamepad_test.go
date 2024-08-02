package uinput

import (
	"fmt"
	"io/ioutil"
	"os"
"testing"
)

// This test inputs the konami code
func TestInfiniteKonami(t *testing.T) {
  test := func(vg Gamepad) { 
    var err error
    for i := 0; i < 10; i++ {
      for j := 0; j < 3; j++ {
        err = vg.ButtonPress(ButtonDpadUp)
        if err != nil {
          t.Fatalf("Failed to send button press. Last error was: %s\n", err)
        }

        err = vg.ButtonPress(ButtonDpadDown)
        if err != nil {
          t.Fatalf("Failed to send button press. Last error was: %s\n", err)
        }

      }

      for j := 0; j < 3; j++ {
        err = vg.ButtonPress(ButtonDpadLeft)
        if err != nil {
          t.Fatalf("Failed to send button press. Last error was: %s\n", err)
        }

        err = vg.ButtonPress(ButtonDpadRight)
        if err != nil {
          t.Fatalf("Failed to send button press. Last error was: %s\n", err)
        }

      }

      err = vg.ButtonPress(ButtonSouth)
      if err != nil {
        t.Fatalf("Failed to send button press. Last error was: %s\n", err)
      }

      err = vg.ButtonPress(ButtonEast)
      if err != nil {
        t.Fatalf("Failed to send button press. Last error was: %s\n", err)
      }

      err = vg.ButtonPress(ButtonStart)
      if err != nil {
        t.Fatalf("Failed to send button press. Last error was: %s\n", err)
      }
    }
  }

  vg, err := CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)

	vg, err = CreateGamepadWithRumble("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF, 1)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)

}

// This test moves the axes around a bit
func TestAxisMovement(t *testing.T) {
  test := func(vg Gamepad) {
    var err error
    err = vg.LeftStickMove(0.2, 1.0)
    if err != nil {
      t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
    }

    err = vg.LeftStickMoveX(0.2)
    if err != nil {
      t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
    }

    err = vg.LeftStickMoveY(1)
    if err != nil {
      t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
    }

    err = vg.RightStickMove(0.2, 1.0)
    if err != nil {
      t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
    }

    err = vg.RightStickMoveX(0.2)
    if err != nil {
      t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
    }

    err = vg.RightStickMoveY(1)
    if err != nil {
      t.Fatalf("Failed to send axis event. Last error was: %s\n", err)
    }
  }

	vg, err := CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)
	vg, err = CreateGamepadWithRumble("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF, 1)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)
}

//This test press the triggers with some amount of force
func TestTriggerAxis(t *testing.T) {
  test := func(vg Gamepad) {
    var err error
    err = vg.LeftTriggerForce(1.0)
    if err != nil {
      t.Fatalf("Failed to send trigger axis event. Last error was: %s\n", err)
    }

    err = vg.RightTriggerForce(-1.0)
    if err != nil {
      t.Fatalf("Failed to send trigger axis event. Last error was: %s\n", err)
    }

    err = vg.LeftTriggerForce(-0.2)
    if err != nil {
      t.Fatalf("Failed to send trigger axis event. Last error was: %s\n", err)
    }

    err = vg.RightTriggerForce(0.2)
    if err != nil {
      t.Fatalf("Failed to send trigger axis event. Last error was: %s\n", err)
    }
  }

	vg, err := CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)
	vg, err = CreateGamepadWithRumble("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF, 1)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)
}

func TestHatMovement(t *testing.T) {
  test := func(vg Gamepad) {
    var err error
    err = vg.HatPress(HatUp)
    if err != nil {
      t.Fatalf("Failed to move hat up")
    }
    err = vg.HatRelease(HatUp)
    if err != nil {
      t.Fatalf("Failed to release hat")
    }
    err = vg.HatPress(HatRight)
    if err != nil {
      t.Fatalf("Failed to move hat right")
    }
    err = vg.HatRelease(HatRight)
    if err != nil {
      t.Fatalf("Failed to release hat")
    }
    err = vg.HatPress(HatDown)
    if err != nil {
      t.Fatalf("Failed to move hat down")
    }
    err = vg.HatRelease(HatDown)
    if err != nil {
      t.Fatalf("Failed to release hat")
    }
    err = vg.HatPress(HatLeft)
    if err != nil {
      t.Fatalf("Failed to move hat left")
    }
    err = vg.HatRelease(HatLeft)
    if err != nil {
      t.Fatalf("Failed to release hat")
    }
  }

	vg, err := CreateGamepad("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)
	vg, err = CreateGamepadWithRumble("/dev/uinput", []byte("Hot gophers in your area"), 0xDEAD, 0xBEEF, 1)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(vg)
}

func TestGamepadCreationFailsOnEmptyPath(t *testing.T) {
	expected := "device path must not be empty"

  test := func(err error) {
    if err.Error() != expected {
      t.Fatalf("Expected: %s\nActual: %s", expected, err)
    }
  }

	_, err := CreateGamepad("", []byte("Gamepad"), 0xDEAD, 0xBEEF)
  test(err)
	_, err = CreateGamepadWithRumble("", []byte("Gamepad"), 0xDEAD, 0xBEEF, 1)
  test(err)
}

func TestGamepadCreationFailsOnNonExistentPathName(t *testing.T) {
	path := "/some/bogus/path"

  test := func(err error) {
    if !os.IsNotExist(err) {
      t.Fatalf("Expected: os.IsNotExist error\nActual: %s", err)
    }
  }

	_, err := CreateGamepad(path, []byte("Gamepad"), 0xDEAD, 0xBEEF)
  test(err)
	_, err = CreateGamepadWithRumble(path, []byte("Gamepad"), 0xDEAD, 0xBEEF, 1)
  test(err)
}

func TestGamepadCreationFailsOnWrongPathName(t *testing.T) {
	expected := "failed to register virtual gamepad device: failed to close device: inappropriate ioctl for device"

  test := func(err error) {
    if err == nil || !(expected == err.Error()) {
      t.Fatalf("Expected: %s\nActual: %s", expected, err)
    }
  }

	file, err := ioutil.TempFile(os.TempDir(), "uinput-gamepad-test-")
	if err != nil {
		t.Fatalf("Failed to setup test. Unable to create tempfile: %v", err)
	}
	defer file.Close()

	_, err = CreateGamepad(file.Name(), []byte("GamepadDevice"), 0xDEAD, 0xBEEF)
  test(err)
	_, err = CreateGamepadWithRumble(file.Name(), []byte("GamepadDevice"), 0xDEAD, 0xBEEF, 1)
  test(err)
}

func TestGamepadCreationFailsIfNameIsTooLong(t *testing.T) {
	name := "adsfdsferqewoirueworiuejdsfjdfa;ljoewrjeworiewuoruew;rj;kdlfjoeai;jfewoaifjef;das"
	expected := fmt.Sprintf("device name %s is too long (maximum of %d characters allowed)", name, uinputMaxNameSize)

  test := func(err error) {
    if err.Error() != expected {
      t.Fatalf("Expected: %s\nActual: %s", expected, err)
    }
  }

	_, err := CreateGamepad("/dev/uinput", []byte(name), 0xDEAD, 0xBEEF)
  test(err)
	_, err = CreateGamepadWithRumble("/dev/uinput", []byte(name), 0xDEAD, 0xBEEF, 1)
  test(err)
}

func TestGamepadButtonEventsFailOnClosedDevice(t *testing.T) {
  test := func(gamepad Gamepad) {
    var err error
    _ = gamepad.Close()

    err = gamepad.ButtonUp(1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
    err = gamepad.ButtonDown(1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
    err = gamepad.ButtonPress(1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
  }

	gamepad, err := CreateGamepad("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(gamepad)
	gamepad, err = CreateGamepadWithRumble("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF, 1)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(gamepad)
}

func TestGamepadHatEventsFailOnClosedDevice(t *testing.T) {
  test := func(gamepad Gamepad) {
    var err error 
    _ = gamepad.Close()

    err = gamepad.HatPress(HatUp)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
    err = gamepad.HatRelease(HatUp)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
  }
  
	gamepad, err := CreateGamepad("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(gamepad)
	gamepad, err = CreateGamepadWithRumble("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF, 1)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(gamepad)
}

func TestGamepadMoveToFailsOnClosedDevice(t *testing.T) {
  test := func(gamepad Gamepad) {
    var err error
    _ = gamepad.Close()

    err = gamepad.LeftStickMoveX(1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
    err = gamepad.LeftStickMoveY(1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
    err = gamepad.LeftStickMove(1, 1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }

    err = gamepad.RightStickMoveX(1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
    err = gamepad.RightStickMoveY(1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
    err = gamepad.RightStickMove(1, 1)
    if err == nil {
      t.Fatalf("Expected error due to closed device, but no error was returned.")
    }
  }

	gamepad, err := CreateGamepad("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(gamepad)
	gamepad, err = CreateGamepadWithRumble("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF, 1)
	if err != nil {
		t.Fatalf("Failed to create the virtual gamepad. Last error was: %s\n", err)
	}
  test(gamepad)
}

func TestGamepadWith0EffectsMax(t *testing.T) {
  expected := "effectsMax is below the minimum value of 1, use CreateGamepad if you don't want rumble support"

  test := func(err error) {
    if err.Error() != expected {
      t.Fatalf("Expected: %s\nActual: %s", expected, err)
    }
  }

  _, err := CreateGamepadWithRumble("/dev/uinput", []byte("Test Gamepad"), 0xDEAD, 0xBEEF, 0)
  test(err)
}

//TODO to test if rumble is working we need to send a rumble event to the gamepad
