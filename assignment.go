package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Assignment struct {
	Id              int
	Name            string
	Due_at          time.Time
	Lock_at         time.Time
	Points_Possible float32
	Points          float32
}
type Grade struct {
	Score float32
}

var (
	// modules you are enrolled in
	courses = map[string]int{
		"C-Programming for Microcontrollers": 48520,
		"Networks and data communications":   48480,
		"Theory of Computation":              48496,
		"Advanced Programming with Java":     48465,
		"Software Engineering":               48475,
		"Ethical Hacking and Web Security":   48505,
	}
	//change this to your university's canvas url
	universityUrl = "https://ucc.instructure.com"
	// canvas user id
	userId = 12345
	// canvas api token
	token     = "not-a-real-token"
	w         sync.WaitGroup
	waitGroup sync.WaitGroup

	// userId, courseId
	assignmentUrl = universityUrl + "/api/v1/users/%d/courses/%d/assignments"
	// courseId, assignmentId, userId
	gradesUrl = universityUrl + "/api/v1/courses/%d/assignments/%d/submissions/%d"
)

// prints all assignments for course
func allAssignments() {
	client := http.Client{Timeout: 5 * time.Second}
	waitGroup.Add(len(courses))
	for name, courseId := range courses {
		go printAssignment(name, courseId, client)
	}
	waitGroup.Wait()
}

func printAssignment(name string, courseId int, client http.Client) {
	defer waitGroup.Done()

	if res, err := getAssignmentForCourse(courseId, userId, client); err == nil {
		if assignmentF := formatAssignments(res); len(assignmentF) != 0 {
			fmt.Printf("%s\n\n%s", name, assignmentF)
		}
	} else {
		fmt.Println(fmt.Errorf("failed to retrieve assignments for %s: %v", name, err))
	}

}

// retrives all assignments for a given course
func getAssignmentForCourse(courseId, userId int, client http.Client) ([]Assignment, error) {
	url := fmt.Sprintf(assignmentUrl, userId, courseId)
	// replace this token with your own from canvas

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed  to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %v", err)
	}
	defer resp.Body.Close()
	z := []Assignment{}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to  read request: %v", err)
	}
	if err = json.Unmarshal(b, &z); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal response: %v", err)
	}
	return z, nil
}

// format assignments with time left until due, will not print assignments with due date in the past
func formatAssignments(resp []Assignment) string {
	if len(resp) == 0 {
		return ""
	}
	res := ""
	sort.Slice(resp[:], func(i, j int) bool {
		return resp[i].Due_at.String() > resp[j].Due_at.String()
	})
	//this is the time object unmarshal will produce if no due_at or lock_at field specified for the assignment
	empty, _ := time.Parse("2006-01-02T15:04:05Z", "0001-01-01T00:00:00Z")
	for _, a := range resp {
		if a.Due_at == empty {
			if a.Lock_at == empty {
				continue
			}
			a.Due_at = a.Lock_at
		}
		if a.Lock_at == empty {
			a.Lock_at = a.Due_at
		}
		diff := a.Due_at.Local().Sub(time.Now())
		if diff.Hours() < 0 {
			continue
		}
		res += fmt.Sprintf("%s\nPoints Possible: %.2f\nDue at: %s\nLock at: %s\nTime left: %s\n\n", a.Name, a.Points_Possible, a.Due_at.Local(), a.Lock_at.Local(), timeLeft(diff))
	}
	return res
}
func timeLeft(timeRemaining time.Duration) string {
	days := int(timeRemaining.Hours()) / 24
	hours := int(timeRemaining.Hours()) % 24
	minutes := int(timeRemaining.Minutes()) % 60
	return fmt.Sprintf("Days: %d, Hours: %d, Minutes: %d", days, hours, minutes)
}

// get the score for a given assignment
func getAssignmentGrade(courseId int, a *Assignment, client http.Client) {
	url := fmt.Sprintf(gradesUrl, courseId, a.Id, userId)
	defer w.Done()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	z := new(Grade)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(b, &z); err != nil {
		return
	}
	a.Points = z.Score
}

// gets grades for all assignments in course
func printGradesForCourse(courseId int, courseName string, client http.Client) {
	defer waitGroup.Done()
	assignments, err := getAssignmentForCourse(courseId, userId, client)
	if err != nil {
		log.Fatalf("failed to get assignment for course: %v", courseId)
	}
	w.Add(len(assignments))
	for i := range assignments {
		go getAssignmentGrade(courseId, &assignments[i], client)
	}
	w.Wait()
	fmt.Println(formatGrades(courseName, assignments))

}

// format assignments/grades for printing
func formatGrades(courseName string, assignments []Assignment) string {
	result := fmt.Sprintf("Course: %s\n", courseName)
	for _, a := range assignments {
		result = fmt.Sprintf("%vName: %v, Grade: %v/%v Percent: %v\n", result, a.Name, a.Points, a.Points_Possible, (a.Points/a.Points_Possible)*100)
	}
	return result
}

// retrieves all grades for all assignments
func allGrades() {
	client := http.Client{Timeout: 5 * time.Second}
	waitGroup.Add(len(courses))
	for name, courseId := range courses {
		go printGradesForCourse(courseId, name, client)
	}
	waitGroup.Wait()
}
