package formatter

import (
	"fmt"

	"github.com/dcilke/hz/pkg/g"
)

const (
	KeyMessage = "message"
	KeyMsg     = "msg"
)

var _ Formatter = (*Message)(nil)

type Message struct {
	noColor   bool
	formatKey Stringer
	keys      []string
}

func NewMessage(noColor bool, formatKey Stringer) Formatter {
	return &Message{
		noColor:   noColor,
		formatKey: formatKey,
		keys:      []string{KeyMessage, KeyMsg},
	}
}

func (f *Message) Format(m map[string]any) string {
	var message string
	var msg string
	if i, ok := m[KeyMessage]; ok {
		message = fmt.Sprintf("%s", i)
	}
	if i, ok := m[KeyMsg]; ok {
		msg = fmt.Sprintf("%s", i)
	}

	if ok, value := g.SameOrEmpty(message, msg); ok {
		if value == "" {
			return ""
		}
		return value
	}
	return kvJoin(
		f.formatKey(KeyMessage), message,
		f.formatKey(KeyMsg), msg,
	)
}

func (f *Message) ExcludeKeys() []string {
	return f.keys
}
