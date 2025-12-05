package main

import (
	"log"
	"os"

	"github.com/rafinhacuri/Deku/controllers"
	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

type goRes struct {
	Message string `json:"message"`
}

func main() {
	converter := typescriptify.New()

	converter.Add(goRes{})
	converter.Add(controllers.AuthResponse{})
	converter.Add(controllers.Salary{})
	converter.Add(controllers.SalaryResponse{})

	converter.BackupDir = ""
	converter.CreateInterface = true

	os.MkdirAll("../../shared/types/", 0755)

	err := converter.ConvertToFile("../../shared/types/goServer.d.ts")
	if err != nil {
		log.Print(err)
	}
}
