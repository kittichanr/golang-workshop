package main

import "os"

func check(e error) {
	if e != nil {
		panic(e)
	}

}

func main() {
	data1 := []byte("hello\nworld")

	err := os.WriteFile("data.txt", data1, 0644)
	check(err)

	f, err := os.Create("employeeName")
	check(err)

	defer f.Close()

	data2 := []byte("Sira\nManee")
	os.WriteFile("employeeName.txt", data2, 0644)
}
