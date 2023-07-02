package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

const coursePath = "courses"
const basePath = "/api"

type Course struct {
	CourseID   int    `json:"id"`
	CourseName string `json:"name"`
	Price      int    `json:"price"`
	ImageURL   string `json:"imageUrl"`
}

func setUpDB() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/coursedb")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func getCourseList() ([]Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := db.QueryContext(ctx, `SELECT
	courseid, coursename, price, image_url FROM courseonline
	`)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer results.Close()
	courses := make([]Course, 0)
	for results.Next() {
		var course Course
		results.Scan(&course.CourseID, &course.CourseName, &course.Price, &course.ImageURL)

		courses = append(courses, course)
	}
	return courses, nil
}

func insertProduct(course Course) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	results, err := db.ExecContext(ctx, `INSERT INTO courseonline (
		courseid, coursename, price, image_url
	) VALUES (?, ?, ?, ?);`, course.CourseID, course.CourseName, course.Price, course.ImageURL)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	insertID, err := results.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return int(insertID), nil
}

func getCourse(courseid int) (*Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, `SELECT
	courseid, coursename, price, image_url FROM courseonline
	WHERE courseid = ?;`, courseid)

	course := &Course{}
	err := row.Scan(
		&course.CourseID,
		&course.CourseName,
		&course.Price,
		&course.ImageURL,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return course, nil
}

func removeCourse(courseID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, `DELETE FROM courseonline WHERE courseid = ?;`, courseID)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func handlerCourses(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		courseList, err := getCourseList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		j, err := json.Marshal(courseList)
		if err != nil {
			log.Fatal(err)
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		var course Course
		err := json.NewDecoder(r.Body).Decode(&course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		CourseID, err := insertProduct(course)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf(`{"courseid":%d}`, CourseID)))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handlerCourse(w http.ResponseWriter, r *http.Request) {
	urlPathSegment := strings.Split(r.URL.Path, fmt.Sprintf("%s/", coursePath))
	if len(urlPathSegment[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	courseID, err := strconv.Atoi(urlPathSegment[len(urlPathSegment)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		course, err := getCourse(courseID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if course == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		j, err := json.Marshal(course)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = w.Write(j)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodDelete:
		err := removeCourse(courseID)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Content-Type, Content-Length,Authorization,X-Cross-origin")

		handler.ServeHTTP(w, r)
	})
}

func setupRoutes(apiBasePath string) {
	coursesHandler := http.HandlerFunc(handlerCourses)
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, coursePath), corsMiddleware(coursesHandler))

	courseHandler := http.HandlerFunc(handlerCourse)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, coursePath), corsMiddleware(courseHandler))
}

func main() {
	setUpDB()
	setupRoutes(basePath)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
