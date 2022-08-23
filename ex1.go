package main

import "os"

func main() {
	data, err := os.ReadFile("Training Program") // read test data
	if err != nil {
		panic(err)
	}
	defer data.Close()

}
