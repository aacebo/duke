package duke

import (
	"bytes"
	"encoding/binary"
)

type Frame struct {
	IsFragment bool
	Opcode     byte
	Reserved   byte
	IsMasked   bool
	Length     uint64
	Payload    []byte
}

func (self *Frame) Ping() *Frame {
	self.Opcode = 9
	return self
}

func (self *Frame) Pong() *Frame {
	self.Opcode = 10
	return self
}

// Get Text Payload
func (self *Frame) String() string {
	return string(self.Payload)
}

func (self *Frame) Code() uint16 {
	return binary.BigEndian.Uint16(self.Payload[:2])
}

// IsControl checks if the frame is a control frame identified by opcodes where the most significant bit of the opcode is 1
func (self *Frame) IsControl() bool {
	return self.Opcode&0x08 == 0x08
}

func (self *Frame) IsContinuation() bool {
	return self.Opcode == 0
}

func (self *Frame) IsText() bool {
	return self.Opcode == 1
}

func (self *Frame) IsBinary() bool {
	return self.Opcode == 2
}

func (self *Frame) IsClose() bool {
	return self.Opcode == 8
}

func (self *Frame) IsPing() bool {
	return self.Opcode == 9
}

func (self *Frame) IsPong() bool {
	return self.Opcode == 10
}

func (self *Frame) IsReserved() bool {
	return self.Opcode > 10 || (self.Opcode >= 3 && self.Opcode <= 7)
}

func (self *Frame) CloseCode() uint16 {
	var code uint16
	binary.Read(bytes.NewReader(self.Payload), binary.BigEndian, &code)
	return code
}
