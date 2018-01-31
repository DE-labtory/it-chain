package main

import (
	"os"
	"log"
	"fmt"
	"time"
)

func AppendFile() {
	now := time.Now().Format("2006-01-02 15:04:05");
	file, err := os.OpenFile("version", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	len, err := file.WriteString("file io test " + now)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func main()  {
	AppendFile()
}