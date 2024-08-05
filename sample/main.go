package main

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/kpango/glg"
	"github.com/randyardiansyah25/gostashlg"
)

func main() {
	useDefine()
	var input string
	// glg.Error("cobain deh")
	// glg.CustomLog("INFO", "hello--world")
	// glg.Error(TestError())

	err := ErrLevel3()
	//el := handleErrorLog(err)
	glg.Error(err)
	fmt.Scanln(&input)
}

func TestError() error {
	file, line := GetCallerInfo()
	s := fmt.Sprintf("Called from file: %s, line: %d\n", file, line)
	return errors.New(s)
}

func ErrLevel3() error {
	err := ErrLevel2()
	return AddError(err, "error level3")
	//return errors.Join(err, errors.New("err level3"))
	//return errors.Join(err,errors.New("this third error"))
}

func ErrLevel2() error {
	err := ErrLevel1()
	return AddError(err, "error level2")
	// return errors.Join(err, errors.New("err level2"))

	// return errors.Join(err, errors.New("this second error"))
}

func ErrLevel1() error {
	return errors.New("This first error")
}

func handleErrorLog(err error) string {
	// if (errors.Is(err, &errors.Error{})) {
	// 	return err.(*errors.Error).ErrorStack()
	// } else {
	// 	return errors.New(err).ErrorStack()
	// }
	return ""
}

func AddError(er error, message string) error {
	file, line := GetCallerInfo()
	s := fmt.Sprintf("%s\n\t%s:%d", message, file, line)
	return errors.Join(er, errors.New(s))
}

func GetCallerInfo() (string, int) {
	// `2` adalah level call stack (1 adalah fungsi GetCallerInfo, 2 adalah fungsi pemanggil)
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "Unknown", 0
	}

	//file = filepath.Base(file)
	return file, line
}

func useDefine() {

	body := struct {
		Type       string
		RemoteAddr string
		Body       string
	}{
		Type:       "RECV",
		RemoteAddr: "172.0.0.1:23994",
		Body:       "Ini contoh body, kalo http bisa aja concate dari header dan body",
	}
	log, _ := gostashlg.UseDefine(gostashlg.Define{
		Template: gostashlg.NewTemplate().
			Add(gostashlg.LOG, "{{.Timestamp}} [{{.Level}}] {{.Data.Type}}, FROM: {{.Data.RemoteAddr}}, {{.Event}}, {{.Message}}, Data:\n{{.Data.Body}}"),
	})
	field := gostashlg.NewFields().
		SetLevel(gostashlg.LOG).
		SetEvent("Test").
		SetMessage("Nyoba log").
		SetData(body).
		Get()

	log.Put(field,
		false)
}

func useDefault() {
	log, _ := gostashlg.UseDefault()
	field := gostashlg.NewFields().
		SetLevel(gostashlg.INFO).
		SetEvent("Test").
		SetMessage("Nyoba log").
		SetData("Detail:\nInformasi Detail").
		Get()

	log.Put(field,
		false)

	var input string
	fmt.Scanln(&input)
}
