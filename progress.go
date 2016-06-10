package main

// Status is a structure for communicating the status of a project.
import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

// ProjectState is the state of the project
type ProjectState int

// These are the states that a project can exist in.
const (
	StateQueued     ProjectState = iota // Project is queued.
	StateProcessing                     // Fetching project.
	StateSuccess                        // Project has successfully fetched.
	StateFailed                         // Project has failed to fetch.
)

// String representation of a ProjectState.
func (ps ProjectState) String() string {
	var state string

	switch ps {
	case StateQueued:
		state = "queued"
	case StateProcessing:
		state = "processing"
	case StateSuccess:
		state = "success"
	case StateFailed:
		state = "failed"
	}

	return state
}

// Status communicates the state of a single propject.
type Status struct {
	Project *Project
	State   ProjectState
}

// SortStatus represents the state of a project.
type SortStatus []*Status

func (c SortStatus) Len() int      { return len(c) }
func (c SortStatus) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c SortStatus) Less(i, j int) bool {
	// Is c[i] less than c[j]?
	if (c[i].Project.Sticky && c[j].Project.Sticky) || (!c[i].Project.Sticky && !c[j].Project.Sticky) {
		// Either both c[i] and c[j] are sticky, or neither are. In either case
		// they can be compared normally.
		return c[i].Project.Name < c[j].Project.Name
	}

	return !c[j].Project.Sticky
}

// Progress tracks the progress of the entire source code assembly.
type Progress struct {
	C               chan string
	projectStatuses SortStatus

	// Needed to ensure exclusive access to the projectStates map with
	// concurrent goroutines.
	mutex *sync.Mutex
}

// NewProgress creates a new progress.
func NewProgress(c chan string) *Progress {
	p := new(Progress)
	p.C = c
	p.mutex = &sync.Mutex{}

	return p
}

// AddStatus adds a Status to Progress. Updates the Progress on the state of an
// individual project.
func (p *Progress) AddStatus(status Status) {
	p.mutex.Lock()

	var found bool
	for _, projectStatus := range p.projectStatuses {
		if projectStatus.Project.Name == status.Project.Name {
			// We found an existing project status to update.
			projectStatus.State = status.State
			found = true
		}
	}

	if !found {
		// We didn't find an existing status for the project.
		p.projectStatuses = append(p.projectStatuses, &status)
	}

	p.mutex.Unlock()
}

// Add creates a Status from the string and state and adds it to the Progress.
func (p *Progress) Add(project *Project, state ProjectState) *Progress {
	status := Status{Project: project, State: state}
	p.AddStatus(status)
	return p
}

// QueueProjects queues a slice of projects.
func (p *Progress) QueueProjects(projects []*Project) {
	for _, project := range projects {
		status := Status{project, StateQueued}
		p.AddStatus(status)
	}
}

// A helper function to get the length of the longest project name
func (p *Progress) longestProjectNameLength() int {
	length := 0
	for _, status := range p.projectStatuses {
		l := utf8.RuneCountInString(status.Project.Name)
		if l > length {
			length = l
		}
	}

	return length
}

// Report the status of projects
func (p *Progress) Report() {
	// Find the longest project name
	length := p.longestProjectNameLength()

	// Calculate the row format
	rowElements := []string{"| %-", strconv.Itoa(length), "s | %-9s | %10s |\n"}
	rowFormat := strings.Join(rowElements, "")

	// Seperator format
	sepElements := []string{"| %-", strconv.Itoa(length), "s | %-9s   %10s |\n"}
	sepFormat := strings.Join(sepElements, "")
	sepString, capString := "", ""
	for i := 0; i < length; i++ {
		sepString += "-"
		capString += "_"
	}
	cap := fmt.Sprintf("__%s___________________________\n", capString)

	report := fmt.Sprintf(rowFormat, "Projects", "Blocking", "Status")
	report += fmt.Sprintf(sepFormat, sepString, "---------", "----------")
	p.mutex.Lock()
	sort.Sort(p.projectStatuses)
	for _, status := range p.projectStatuses {
		blocking := "No"

		switch {
		case status.Project.Blocking && status.State == StateProcessing:
			blocking = "BLOCKING"
		case status.Project.Blocking && status.State != StateProcessing:
			blocking = "unblocked"
		}

		report += fmt.Sprintf(rowFormat, status.Project.Name, blocking, status.State)
	}
	p.mutex.Unlock()

	report = fmt.Sprintf("%s%s%s", cap, report, cap)
	p.C <- report
}
