package messages

type MessageType string

type CommandType string

const (
	Request MessageType = "request"
	Reply   MessageType = "reply"
	Cancel  MessageType = "cancel"
	Flood   MessageType = "flood"
	Echo    MessageType = "echo"
	Short   MessageType = "short"

	Ask      CommandType = "ask"
	Detect   CommandType = "detect"
	Liberate CommandType = "liberate"
	Unknown  CommandType = "unknown"
)

func (m MessageType) String() string {
	return string(m)
}
