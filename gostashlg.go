package gostashlg

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/imroc/req/v3"
	"github.com/parnurzeal/gorequest"
	"github.com/randyardiansyah25/glg"
	"golang.org/x/sync/singleflight"
)

const (
	FORMAT_YMD = "20060102"
)

var (
	sfGroup singleflight.Group
	lSync   sync.Mutex
)

type logItem struct {
	Field       Fields
	PushToStash bool
}

func UseDefault() (l LoggerEngine, e error) {
	return use(NewTemplate())
}

func UseDefine(template *Template) (l LoggerEngine, e error) {
	return use(template)
}

func use(tmpl *Template) (l LoggerEngine, er error) {

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
		itemChan:        make(chan logItem),
		isUseStash:      false,
	}
	gorequest.New()
	stashUrl := os.Getenv("logstash.host")
	parsedURL, err := url.Parse(stashUrl)

	if err == nil && parsedURL.Scheme != "" && parsedURL.Host != "" {
		l.isUseStash = true
		l.stashUrl = stashUrl
	}

	if l.isUseStash {
		l.client = req.C() //.DevMode()
	}

	l.prepareLogFile()
	go l.observe()
	return
}

type LoggerEngine struct {
	definedTemplate map[Level]*template.Template
	defaultTemplate *template.Template
	itemChan        chan logItem
	isUseStash      bool
	stashUrl        string
	LastSuffix      string
	client          *req.Client
}

func (l *LoggerEngine) Write(f Fields) {
	l.doWrite(f, true)
}

func (l *LoggerEngine) WriteOnly(f Fields) {
	l.doWrite(f, false)
}

func (l *LoggerEngine) doWrite(f Fields, putToStash bool) {
	l.itemChan <- logItem{
		Field:       f,
		PushToStash: putToStash,
	}
}

func (l *LoggerEngine) observe() {
	for i := range l.itemChan {
		go l.exec(i)
	}
}

func (l *LoggerEngine) exec(item logItem) {
	message := l.composeOutput(item.Field)
	l.printLog(item.Field.Timestamp, item.Field.Level, message)
	if l.isUseStash && item.PushToStash {
		l.put(item)
	}
}

func (l *LoggerEngine) prepareLogFile() {
	sfGroup.Do("prepare_log_file", func() (interface{}, error) {
		lSync.Lock()
		l.LastSuffix = time.Now().Format(FORMAT_YMD)
		logFl := glg.FileWriter(fmt.Sprintf("log/app_%s.log", l.LastSuffix), 0660)

		glg.Get().
			SetMode(glg.BOTH).
			AddLevelWriter(glg.DEBG, logFl).
			AddLevelWriter(glg.INFO, logFl).
			AddLevelWriter(glg.LOG, logFl).
			AddLevelWriter(glg.PRINT, logFl).
			AddLevelWriter(glg.TRACE, logFl).
			AddLevelWriter(glg.OK, logFl).
			AddLevelWriter(glg.WARN, logFl).
			AddLevelWriter(glg.ERR, logFl).
			AddLevelWriter(glg.FAIL, logFl).
			AddLevelWriter(glg.FATAL, logFl)

		lSync.Unlock()
		return nil, nil
	})
}

func (l *LoggerEngine) printLog(timestamp string, level Level, msg string) {
	_ = glg.CustomTimestampLog([]byte(timestamp), string(level), msg)
}

func (l *LoggerEngine) put(item logItem) {
	c := l.client.Clone()
	resp, er := c.R().SetBody(item.Field).Post(l.stashUrl)
	if er != nil {
		_ = glg.Warn("Unable to put log to stash,", er)
	} else {
		if !resp.IsSuccessState() {
			_ = glg.Warn("Problem when put log to stash,", er)
		}
	}
}

func (l *LoggerEngine) composeOutput(f Fields) string {
	var tmpl *template.Template

	tmpl, ok := l.definedTemplate[f.Level]
	if !ok {
		tmpl = l.defaultTemplate
	}

	currentSuffix := time.Now().Format(FORMAT_YMD)
	if l.LastSuffix != currentSuffix {
		l.prepareLogFile()
	}

	var out bytes.Buffer
	_ = tmpl.Execute(&out, f)

	return out.String()
}
