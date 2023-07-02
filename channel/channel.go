package main

import (
	"fmt"
	"time"
)

func process1(c chan string, data string) {
	c <- data
}

func main() {
	ch := make(chan string)
	go process1(ch, "Data1")

	ch2 := make(chan string, 1)
	process1(ch2, "Data2")

	fmt.Println(<-ch)
	fmt.Println(<-ch2)
	time.Sleep(5 * time.Second)
}
