package duke

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf8"
)

type Frame struct {
	IsFragment bool
	Opcode     byte
	Reserved   byte
	IsMasked   bool
	Length     uint64
	Payload    []byte
}

func NewCloseFrame(status uint16) *Frame {
	v := Frame{
		Opcode: 8,
		Length: 2,
		Payload: make([]byte, 2),
	}

	binary.BigEndian.PutUint16(v.Payload, status)
	return &v
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

func (self *Frame) Reason() []byte {
	return self.Payload[2:]
}

func (self *Frame) IsReasonValid() bool {
	return utf8.Valid(self.Reason())
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

func (self *Frame) Validate(status uint16) (uint16, error) {
	if !self.IsMasked {
		return ProtocolError.UInt16(), errors.New("protocol error: unmasked client frame")
	}

	if self.IsControl() && (self.Length > 125 || self.IsFragment) {
		return ProtocolError.UInt16(), errors.New("protocol error: all control frames MUST have a payload length of 125 bytes or less and MUST NOT be fragmented")
	}

	if self.IsReserved() {
		return ProtocolError.UInt16(), errors.New("protocol error: opcode " + fmt.Sprintf("%x", self.Opcode) + " is reserved")
	}

	if self.Reserved > 0 {
		return ProtocolError.UInt16(), errors.New("protocol error: RSV " + fmt.Sprintf("%x", self.Reserved) + " is reserved")
	}

	if self.Opcode == 1 && !self.IsFragment && !utf8.Valid(self.Payload) {
		return InvalidFramePayloadData.UInt16(), errors.New("wrong code: invalid UTF-8 text message")
	}

	if self.IsClose() {
		if self.Length >= 2 {
			code := self.Code()

			if code >= 5000 || IsCloseCodeUnassigned(code) {
				return ProtocolError.UInt16(), errors.New(ProtocolError.String() + " Wrong Code")
			}

			if self.Length > 2 && !self.IsReasonValid() {
				return InvalidFramePayloadData.UInt16(), errors.New(InvalidFramePayloadData.String() + " invalid UTF-8 reason message")
			}
		} else if self.Length != 0 {
			return ProtocolError.UInt16(), errors.New(ProtocolError.String() + " Wrong Code")
		}
	}

	return status, nil
}
