package mf

import (
	"log"
	"os"
)

func create_file(file_name string) {
	_, err := os.Create(file_name)
	if err != nil {
		log.Println("Error: can't create file with name " + file_name)
		os.Exit(1)
	}
}

func read_file(file_name string) string {
	content, err := os.ReadFile(file_name)
	if err != nil {
		log.Println("Error: couldn't read file " + file_name)
		os.Exit(1)
	}
	return string(content)
}

func write_file(file_name string, new_text string) {
	file, err := os.OpenFile(file_name, os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Error: couldn't open file " + file_name)
		os.Exit(1)
	}
	_, err = file.WriteString(new_text)
	if err != nil {
		log.Println("Error: couldn't write file " + file_name)
		os.Exit(1)
	}
}
