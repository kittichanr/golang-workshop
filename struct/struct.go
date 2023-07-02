package main

import "fmt"

type employee struct {
	employeeID   string
	employeeName string
	phone        string
}

func main() {
	employeeList := []employee{}

	employee1 := employee{
		employeeID:   "101",
		employeeName: "John",
		phone:        "010000000",
	}

	employee2 := employee{
		employeeID:   "102",
		employeeName: "Alice",
		phone:        "020000000",
	}
	employee3 := employee{
		employeeID:   "103",
		employeeName: "Bob",
		phone:        "030000000",
	}
	employeeList = append(employeeList, employee1, employee2, employee3)

	fmt.Println("employee = ", employeeList)
}
