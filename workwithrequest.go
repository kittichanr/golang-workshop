package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Course struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Instructor string `json:"instructor"`
}

var CourseList []Course

func init() {
	CourseJSON := `[
		{
		  "id": 1,
		  "name": "Mathematics",
		  "price": 50,
		  "instructor": "John Doe"
		},
		{
		  "id": 2,
		  "name": "English Literature",
		  "price": 40,
		  "instructor": "Jane Smith"
		},
		{
		  "id": 3,
		  "name": "Computer Science",
		  "price": 60,
		  "instructor": "Mike Johnson"
		}
	  ]`
	err := json.Unmarshal([]byte(CourseJSON), &CourseList)

	if err != nil {
		log.Fatal(err)
	}
}

func courseHandler(w http.ResponseWriter, r *http.Request) {
	courseJSON, err := json.Marshal(CourseList)

	switch r.Method {
	case http.MethodGet:
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Header().Set("Content-type", "application/json")
		w.Write(courseJSON)
	case http.MethodPost:
		var newCourse Course

		bodyByte, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(bodyByte, &newCourse)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if newCourse.ID != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newCourse.ID = getNextID()
		CourseList = append(CourseList, newCourse)
		w.WriteHeader(http.StatusCreated)
		return
	}
}

func getNextID() int {
	highestID := -1
	for _, course := range CourseList {
		if highestID < course.ID {
			highestID = course.ID
		}
	}
	return highestID + 1
}

func courseServe() {
	http.HandleFunc("/course", courseHandler)
	http.ListenAndServe(":8080", nil)
}
