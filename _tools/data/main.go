package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	bytes, err := os.ReadFile("saka.png")
	if err != nil {
		log.Fatal("fialed to open: ", err)
	}

	fmt.Println(bytes)
}
