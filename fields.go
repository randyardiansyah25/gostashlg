package gostashlg

import "time"

func NewTimestamp() string {
	str := time.Now().Format("2006-01-02 15:04:05")
	return str
}

type Fields struct {
	IdentifierName string      `json:"service"`
	Timestamp      string      `json:"service_timestamp"`
	Level          Level       `json:"level"`
	Event          string      `json:"event"`
	Message        string      `json:"log_message"`
	Data           interface{} `json:"data,omitempty"`
}

func NewFields() *Fields {
	return &Fields{
		Timestamp: NewTimestamp(),
	}
}

func (f *Fields) SetIdentifierName(id string) *Fields {
	f.IdentifierName = id
	return f
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

func (f *Fields) SetData(data interface{}) *Fields {
	f.Data = data
	return f
}

func (f *Fields) Get() Fields {
	return *f
}
