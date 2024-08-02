package uinput

import (
	"syscall"
	"unsafe"
)

// types needed from uinput.h
const (
	uinputMaxNameSize = 80
	uiDevCreate       = 0x5501
	uiDevDestroy      = 0x5502
	uiDevSetup        = 0x405c5503
	// this is for 64 length buffer to store name
	// for another length generate using : (len << 16) | 0x8000552C
	uiGetSysname  = 0x8041552c

	uiSetEvBit    = 0x40045564
	uiSetKeyBit   = 0x40045565
	uiSetRelBit   = 0x40045566
	uiSetAbsBit   = 0x40045567
	uiSetFFBit    = 0x4004556b 

  uiBeginFFUpload = 0xc06855c8
  uiEndFFUpload   = 0x406855c9
  uiBeginFFErase  = 0xc00c55ca
  uiEndFFErase    = 0x400c55cb
  
	busUsb      = 0x03
)

// input event codes as specified in input-event-codes.h
const (
	evSyn     = 0x00
	evKey     = 0x01
	evRel     = 0x02
	evAbs     = 0x03
  evFF      = 0x15
	relX      = 0x0
	relY      = 0x1
	relHWheel = 0x6
	relWheel  = 0x8
	relDial   = 0x7

	absX     = 0x00
	absY     = 0x01
	absZ     = 0x02
	absRX    = 0x03
	absRY    = 0x04
	absRZ    = 0x05
	absHat0X = 0x10
	absHat0Y = 0x11

	absMtSlot       = 0x2f
	absMtTouchMajor = 0x30
	absMtPositionX  = 0x35
	absMtPositionY  = 0x36
	absMtTrackingId = 0x39

	synReport        = 0
	evMouseBtnLeft   = 0x110
	evMouseBtnRight  = 0x111
	evMouseBtnMiddle = 0x112
	evBtnTouch       = 0x14a
)

const (
	btnStateReleased = 0
	btnStatePressed  = 1
	absSize          = 64
)

// ff uinput consts
const (
  evUinput    = 0x0101
  uiFFUpload  = 1
  uiFFErase   = 2
)

type inputID struct {
	Bustype uint16
	Vendor  uint16
	Product uint16
	Version uint16
}

// translated to go from uinput.h
type uinputUserDev struct {
	Name       [uinputMaxNameSize]byte
	ID         inputID
	EffectsMax uint32
	Absmax     [absSize]int32
	Absmin     [absSize]int32
	Absfuzz    [absSize]int32
	Absflat    [absSize]int32
}

// translated to go from input.h
type inputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// ff-effect structs from input.h

type FFReplay struct {
  Length uint16
  Delay uint16
}

type FFTrigger struct {
  Button uint16
  Interval uint16
}

type FFEnvelope struct {
  AttackLength uint16
  AttackLevel uint16
  FadeLength uint16
  FadeLevel uint16
}

type FFConstantEffect struct {
  Level int16
  Envelope FFEnvelope
}

type FFRampEffect struct {
  StartLevel int16
  EndLevel int16
  Envelope FFEnvelope
}

type FFConditionEffect struct {
  RightSaturation uint16
  LeftSaturation uint16

  RightCoeff int16
  LeftCoeff int16

  Deadband uint16
  Center int16
}

type FFPeriodicEffect struct {
  Waveform uint16
  Period uint16
  Magnitude int16
  Offset int16
  Phase uint16

  Envelope FFEnvelope

  CustomLen uint32
  CustomData *int16 //not suse about this one original:	__s16 __user *custom_data;
}

type FFRumbleEffect struct {
  StrongMagnitude uint16
  WeakMagnitude uint16
}

// to access the values in U use the function Rumble() Periodic() etc
// check FFEffect.Type to call the correct function
// calling a function that doesn't match the type is undefined
type FFEffect struct {
  Type      uint16
  ID        int16
  Direction uint16
  Trigger   FFTrigger
  Replay    FFReplay
  // padding 
  _ uint16 
  /* 
  what U is supposed to mean
	  union {
	  struct ff_constant_effect constant;
	  struct ff_ramp_effect ramp;
	  struct ff_periodic_effect periodic; // this one is* the bigger with 192 bits
	  struct ff_condition_effect condition[2]; \/* One for each axis *\/ 
	  struct ff_rumble_effect rumble;
	 } u;
  */
  u         [32]byte 
}

func (ff *FFEffect) Rumble() FFRumbleEffect {
  return *(*FFRumbleEffect)(unsafe.Pointer(&ff.u[0]))
}

func (ff *FFEffect) Periodic() FFPeriodicEffect {
  return *(*FFPeriodicEffect)(unsafe.Pointer(&ff.u[0]))
}

func (ff *FFEffect) Ramp() FFRampEffect {
  return *(*FFRampEffect)(unsafe.Pointer(&ff.u[0]))
}

func (ff *FFEffect) Constant() FFConstantEffect {
  return *(*FFConstantEffect)(unsafe.Pointer(&ff.u[0]))
}

func (ff *FFEffect) Condition() [2]FFConditionEffect {
  var r [2]FFConditionEffect 
  r[0] = *(*FFConditionEffect)(unsafe.Pointer(&ff.u[0]))
  r[1] = *(*FFConditionEffect)(unsafe.Pointer(&ff.u[16]))
  return r
}

// uinput force-feedback structs from uinput.h

type UInputFFUpload struct {
  RequestID   uint32
  ReturnValue int32
  Effect      FFEffect
  OldEffect   FFEffect
}

type UInputFFErase struct {
  RequestID   uint32
  ReturnValue int32
  EffectID    uint32
}

