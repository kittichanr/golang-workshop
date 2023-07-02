package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func coursesHandler(w http.ResponseWriter, r *http.Request) {
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

func courseHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, "course/")
	ID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusFound)
		return
	}

	course, listItemIndex := findID(ID)
	if course == nil {
		http.Error(w, fmt.Sprintf("no course with id %d", ID), http.StatusFound)
	}
	switch r.Method {
	case http.MethodGet:
		courseJSON, err := json.Marshal(course)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(courseJSON)
	case http.MethodPut:
		var updateCourse Course
		byteBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err = json.Unmarshal(byteBody, &updateCourse)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if updateCourse.ID != ID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		course = &updateCourse
		CourseList[listItemIndex] = *course
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func findID(ID int) (*Course, int) {
	for i, course := range CourseList {
		if course.ID == ID {
			return &course, i
		}
	}
	return nil, 0
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

func main() {
	http.HandleFunc("/course/", courseHandler)
	http.HandleFunc("/course", coursesHandler)
	http.ListenAndServe(":8080", nil)
}
