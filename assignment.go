package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Assignment struct {
	Name            string
	Due_at          time.Time
	Lock_at         time.Time
	Points_Possible float32
}

var waitGroup sync.WaitGroup

func main() {
	courses := map[string]int{
		"Team Software Project":           48420,
		"Workplace Technology and Skills": 48425,
	}
	client := http.Client{Timeout: 5 * time.Second}
	waitGroup.Add(len(courses))
	for name, courseId := range courses {
		go printAssignment(name, courseId, client)
	}
	waitGroup.Wait()
}

func printAssignment(name string, courseId int, client http.Client) {
	if res, err := getAssignmentForCourse(courseId, client); err == nil {
		if assignmentF := formatAssignments(res); len(assignmentF) != 0 {
			fmt.Println(name)
			fmt.Println(assignmentF)
		}
	} else {
		fmt.Println(fmt.Errorf("failed to retrieve assignments for %s: %v", name, err))
	}
	waitGroup.Done()
}
func getAssignmentForCourse(course_id int, client http.Client) ([]Assignment, error) {
	url := fmt.Sprintf("https://ucc.instructure.com/api/v1/courses/%d/assignments", course_id)
	// replace this token with your own from canvas
	token := "not-a-real-token"
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
