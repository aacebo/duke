package duke

type CloseCode uint16
const (
	NormalClosure			CloseCode = 1000
	GoingAway				CloseCode = 1001
	ProtocolError			CloseCode = 1002
	UnsupportedData			CloseCode = 1003
	Reserved				CloseCode = 1004
	NoStatusRcvd			CloseCode = 1005
	AbnormalClosure			CloseCode = 1006
	InvalidFramePayloadData	CloseCode = 1007
	PolicyViolation			CloseCode = 1008
	MessageTooBig			CloseCode = 1009
	MandatoryExt			CloseCode = 1010
	InternalError			CloseCode = 1011
	ServiceRestart			CloseCode = 1012
	TryAgainLater			CloseCode = 1013
	BadGateway				CloseCode = 1014
	TLSHandshake			CloseCode = 1015
)

func (self CloseCode) UInt16() uint16 {
	return uint16(self)
}

func (self CloseCode) String() string {
	switch self {
	case NormalClosure:
		return "Normal Closure"
	case GoingAway:
		return "Going Away"
	case ProtocolError:
		return "Protocol Error"
	case UnsupportedData:
		return "Unsupported Data"
	case Reserved:
		return "Reserved"
	case NoStatusRcvd:
		return "No Status Received"
	case AbnormalClosure:
		return "Abnormal Closure"
	case InvalidFramePayloadData:
		return "Invalid Frame Payload Data"
	case PolicyViolation:
		return "Policy Violation"
	case MessageTooBig:
		return "Message Too Big"
	case MandatoryExt:
		return "Mandatory Ext"
	case InternalError:
		return "Internal Error"
	case ServiceRestart:
		return "Service Restart"
	case TryAgainLater:
		return "Try Again Later"
	case BadGateway:
		return "Bad Gateway"
	case TLSHandshake:
		return "TLS Handshake"
	}

	return "Unexpected Error"
}

func IsCloseCodeUnassigned(code uint16) bool {
	return ((code >= 1016 && code <= 2999) ||
		(code >= 3001 && code <= 3999))
}
