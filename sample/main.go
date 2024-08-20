package main

import (
	"fmt"
	"os"

	"github.com/randyardiansyah25/gostashlg"
)

func main() {
	os.Setenv("logstash.host", "http://localhost:5044")
	useDefine()
	var input string
	fmt.Scanln(&input)
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

	bErr := struct {
		Body string
	}{
		Body: "ini adalah error",
	}
	log, _ := gostashlg.UseDefine(gostashlg.NewTemplate().
		Add(gostashlg.LOG, "{{.Data.Type}}, FROM:4 {{.Data.RemoteAddr}}, {{.Event}}, {{.Message}}, Data:\n{{.Data.Body}}").
		Add(gostashlg.ERROR, "Data:\n{{.Data.Body}}"),
	)

	_ = gostashlg.NewFields().
		SetIdentifierName("myapp").
		SetLevel(gostashlg.LOG).
		SetEvent("Test").
		SetMessage("Nyoba log").
		SetData(body).
		Get()

	//  log.Write(field)

	fieldErr := gostashlg.NewFields().
		SetIdentifierName("myapp").
		SetLevel(gostashlg.ERROR).
		SetEvent("Test error").
		SetMessage("pesan error").
		SetData(bErr).
		Get()

	log.Write(fieldErr)
}

func useDefault() {
	log, _ := gostashlg.UseDefault()
	field := gostashlg.NewFields().
		SetLevel(gostashlg.INFO).
		SetEvent("Test").
		SetMessage("Nyoba log").
		SetData("Detail:\nInformasi Detail").
		Get()

	log.Write(field)

	var input string
	fmt.Scanln(&input)
}
