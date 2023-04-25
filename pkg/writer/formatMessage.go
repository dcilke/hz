package writer

import "fmt"

const (
	KeyMessage = "message"
	KeyMsg     = "msg"
)

type messageFormatter struct {
	noColor   bool
	formatKey Stringer
	keys      []string
}

func newMessageFormatter(noColor bool, formatKey Stringer) Formatter {
	return &messageFormatter{
		noColor:   noColor,
		formatKey: formatKey,
		keys:      []string{KeyMessage, KeyMsg},
	}
}

func (f *messageFormatter) Format(m map[string]any, _ string) string {
	var message string
	var msg string
	if i, ok := m[KeyMessage]; ok {
		message = fmt.Sprintf("%s", i)
	}
	if i, ok := m[KeyMsg]; ok {
		msg = fmt.Sprintf("%s", i)
	}

	if ok, value := sameOrEmpty(message, msg); ok {
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

func (f *messageFormatter) ExcludeKeys() []string {
	return f.keys
}
