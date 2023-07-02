package main

import "fmt"

func add(value1, value2 float64) {
	result := value1 + value2
	fmt.Println("result :", result)
}

func loop() {
	for i := 0; i < 10; i++ {
		fmt.Println("i =", i)
	}
}

func deferloop() {
	for j := 0; j < 10; j++ {
		defer fmt.Println("j =", j)
	}
}

func main() {
	// fmt.Println("Welcome to Calculator")

	// defer fmt.Println("End")
	// defer add(20, 10)
	// defer add(15, 15)
	// defer add(12, 12)

	loop()
	deferloop()
}
