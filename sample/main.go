package main

import (
	"fmt"

	"github.com/randyardiansyah25/gostashlg"
)

func main() {
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
	log, _ := gostashlg.UseDefine(gostashlg.Define{
		Template: gostashlg.NewTemplate().
			Add(gostashlg.LOG, "{{.Data.Type}}, FROM: {{.Data.RemoteAddr}}, {{.Event}}, {{.Message}}, Data:\n{{.Data.Body}}"),
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
