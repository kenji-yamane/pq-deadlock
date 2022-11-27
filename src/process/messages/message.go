package messages

import (
	"encoding/json"
	"fmt"
	"github.com/kenji-yamane/pq-deadlock/src/clock"
	"github.com/kenji-yamane/pq-deadlock/src/customerror"
)

type Message struct {
	SenderId    int     `json:"id"`
	Text        string  `json:"type"`
	ClockStr    string  `json:"clock_str"`
	Weight      float64 `json:"weight"`
	InitiatorId int     `json:"initiator_id"`
	InitiatedAt int     `json:"initiated_at"`
}

func BuildMessage(myId int, c clock.LogicalClock, messageType MessageType, weight float64, initiatorId int, initiatedAt int) string {
	m := Message{
		SenderId:    myId,
		Text:        messageType.String(),
		ClockStr:    c.GetClockStr(),
		Weight:      weight,
		InitiatorId: initiatorId,
		InitiatedAt: initiatedAt,
	}
	mStr, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		customerror.CheckError(fmt.Errorf("error serializing"))
	}
	return string(mStr)
}

func ParseMessage(msg string) (Message, error) {
	var msgSerializer Message
	err := json.Unmarshal([]byte(msg), &msgSerializer)
	return msgSerializer, err
}
