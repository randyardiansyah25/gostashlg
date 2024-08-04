package main

import (
	"fmt"

	"github.com/randyardiansyah25/gostashlg"
)

func main() {
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
