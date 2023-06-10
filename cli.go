package main

import (
	"flag"
	"log"
)

type jobFunc func()

var (
	fPtr   jobFunc
	exists bool
)

func main() {
	job := flag.String("job", "all", "Job to perform. REQUIRED")
	flag.Parse()

	tasks := map[string]func(){
		"all":    allAssignments,
		"grades": allGrades,
	}
	if fPtr, exists = tasks[*job]; !exists {
		flag.PrintDefaults()
		log.Fatalf("invalid job: %v", *job)
	}
	fPtr()

}
