package mf

import (
	"log"
	"os"
)

func Create_file(file_name string) {
	_, err := os.Create(file_name)
	if err != nil {
		log.Println("Error: can't create file with name " + file_name)
		os.Exit(1)
	}
}

func Read_file(file_name string) string {
	content, err := os.ReadFile(file_name)
	if err != nil {
		log.Println("Error: couldn't read file " + file_name)
		os.Exit(1)
	}
	return string(content)
}

func Write_file(file_name string, new_text string) {
	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error: couldn't open file " + file_name)
		os.Exit(1)
	}
	defer file.Close()

	_, err = file.WriteString(new_text + "\n")
	if err != nil {
		log.Println("Error: couldn't write file " + file_name)
		os.Exit(1)
	}
}

