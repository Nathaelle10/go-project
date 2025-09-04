package main

import "fmt"

func main() {
	fmt.Println("Start :) ")

	go StartAPIServer()
	StartWebServer()
}
