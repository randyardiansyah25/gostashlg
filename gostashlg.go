package gostashlg

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"os"
)

type LogHandlerFunc func()

type Define struct {
	Template *Template
	Fun      LogHandlerFunc
}

type logItem struct {
	Field     Fields
	PushStash bool
}

func UseDefault() (l LoggerEngine, e error) {
	return use(NewTemplate(), nil)
}

func UseDefine(d Define) (l LoggerEngine, e error) {
	return use(d.Template, d.Fun)
}

func use(tmpl *Template, fun LogHandlerFunc) (l LoggerEngine, er error) {

	defaultTemplate, er := template.New("logTemplate").Parse(LogTemplate)
	if er != nil {
		return l, fmt.Errorf("error parsing template: %v", er)
	}

	defineTemplate := make(map[Level]*template.Template, 0)

	for k, v := range tmpl.pattern {
		templateItem, er := template.New(string(k)).Parse(v)
		if er != nil {
			return l, fmt.Errorf("error parsing template for level - %s: %v", k, er)
		}
		defineTemplate[k] = templateItem
	}

	l = LoggerEngine{
		defaultTemplate: defaultTemplate,
		definedTemplate: defineTemplate,
		fun:             fun,
		itemChan:        make(chan logItem),
		isUseStash:      false,
	}

	stashUrl := os.Getenv("logstash.host")
	parsedURL, err := url.Parse(stashUrl)

	if err == nil || parsedURL.Scheme != "" || parsedURL.Host != "" {
		l.isUseStash = true
	}

	go l.observe()
	return
}

type LoggerEngine struct {
	definedTemplate map[Level]*template.Template
	defaultTemplate *template.Template
	fun             LogHandlerFunc
	itemChan        chan logItem
	isUseStash      bool
}

func (l *LoggerEngine) Put(f Fields, putToStash bool) {
	l.itemChan <- logItem{
		Field:     f,
		PushStash: putToStash,
	}
}

func (l *LoggerEngine) observe() {
	for i := range l.itemChan {
		go l.exec(i)
	}
}

func (l *LoggerEngine) exec(item logItem) {
	var tmpl *template.Template

	tmpl, ok := l.definedTemplate[item.Field.Level]
	if !ok {
		tmpl = l.defaultTemplate
	}

	var out bytes.Buffer
	_ = tmpl.Execute(&out, item.Field)

	logStr := out.String()

	fmt.Println(logStr)

	if l.fun != nil {
		l.fun()
	}
}
