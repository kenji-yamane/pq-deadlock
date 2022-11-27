package puppeteer

import (
	"github.com/kenji-yamane/pq-deadlock/src"
	"github.com/kenji-yamane/pq-deadlock/src/process/messages"
	"time"
)

type Command struct {
	Id    int
	Line  string
	Await time.Duration
}

func (c *Command) ToMessage() messages.Message {
	return messages.Message{
		SenderId: src.PuppeteerId,
		Text:     c.Line,
	}
}
