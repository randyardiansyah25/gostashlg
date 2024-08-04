package gostashlg

const (
	LogTemplate = "{{.Timestamp}} [{{.Level}}] {{.Event}}, {{.Message}}, {{.Data}}"
)

type Template struct {
	pattern map[Level]string
}

func NewTemplate() (t *Template) {
	return &Template{}
}

func (t *Template) Add(level Level, template string) *Template {
	t.pattern[level] = template
	return t
}
