package gostashlg

import "time"

type Timestamp string

func NewTimestamp() Timestamp {
	str := time.Now().Format("2006-01-02 15:04:05")
	return Timestamp(str)
}

type Fields struct {
	Timestamp Timestamp `json:"@timestamp"`
	Level     Level     `json:"level"`
	Event     string    `json:"event"`
	Message   string    `json:"log_message"`
	Data      any       `json:"data,omitempty"`
}

func NewFields() *Fields {
	return &Fields{
		Timestamp: NewTimestamp(),
	}
}

func (f *Fields) SetLevel(level Level) *Fields {
	f.Level = level
	return f
}

func (f *Fields) SetEvent(e string) *Fields {
	f.Event = e
	return f
}

func (f *Fields) SetMessage(msg string) *Fields {
	f.Message = msg
	return f
}

func (f *Fields) SetData(data any) *Fields {
	f.Data = data
	return f
}

func (f *Fields) Get() Fields {
	return *f
}
