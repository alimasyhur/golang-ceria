package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	f1, err := os.Open("./file0.go")
	if err != nil {
		fmt.Println("Unable to open file: ", err)
		os.Exit(1)
	}

	defer f1.Close()

	f2, err := os.Create("./file0.bkp")
	if err != nil {
		fmt.Println("Unable to create file: ", err)
		os.Exit(1)
	}

	defer f2.Close()

	f3, err := io.Copy(f2, f1)
	if err != nil {
		fmt.Println("Unable to Copy file: ", err)
		os.Exit(1)
	}

	fmt.Printf("Copied %d bytes from %s to %s\n",
		f3, f1.Name(), f2.Name())

	//Openning File
	open1, err := os.OpenFile("./file0.go", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("Unable to open file: ", err)
		os.Exit(1)
	}

	defer open1.Close()
}
