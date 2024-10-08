package gostashlg

const (
	LogTemplate = "{{.Event}}, {{.Message}}, {{.Data}}"
)

type Template struct {
	pattern map[Level]string
}

func NewTemplate() (t *Template) {
	return &Template{
		pattern: make(map[Level]string, 0),
	}
}

func (t *Template) Add(level Level, template string) *Template {
	t.pattern[level] = template
	return t
}
