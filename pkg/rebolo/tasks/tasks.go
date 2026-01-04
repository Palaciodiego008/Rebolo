package tasks

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Task represents a runnable task
type Task struct {
	Name        string
	Description string
	Handler     func(args []string) error
}

var (
	tasks   = make(map[string]*Task)
	tasksMu sync.RWMutex
	app     interface{} // Reference to Application for tasks that need it
)

// Register registers a new task
func Register(name, description string, handler func(args []string) error) {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	if _, exists := tasks[name]; exists {
		panic(fmt.Sprintf("task %s already registered", name))
	}

	tasks[name] = &Task{
		Name:        name,
		Description: description,
		Handler:     handler,
	}
}

// List returns all registered tasks sorted by name
func List() []*Task {
	tasksMu.RLock()
	defer tasksMu.RUnlock()

	result := make([]*Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, task)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// Get returns a task by name
func Get(name string) (*Task, error) {
	tasksMu.RLock()
	defer tasksMu.RUnlock()

	task, exists := tasks[name]
	if !exists {
		return nil, fmt.Errorf("task %s not found", name)
	}

	return task, nil
}

// Run executes a task by name with the given arguments
func Run(name string, args []string) error {
	task, err := Get(name)
	if err != nil {
		return err
	}

	return task.Handler(args)
}

// SetApp sets the application reference for tasks that need it
func SetApp(a interface{}) {
	app = a
}

// GetApp returns the application reference
func GetApp() interface{} {
	return app
}

// PrintList prints all available tasks
func PrintList() {
	tasks := List()

	if len(tasks) == 0 {
		fmt.Println("No tasks available")
		return
	}

	fmt.Println("Available tasks:")
	fmt.Println()

	maxNameLen := 0
	for _, task := range tasks {
		if len(task.Name) > maxNameLen {
			maxNameLen = len(task.Name)
		}
	}

	for _, task := range tasks {
		padding := strings.Repeat(" ", maxNameLen-len(task.Name))
		desc := task.Description
		if desc == "" {
			desc = "No description"
		}
		fmt.Printf("  %s%s  %s\n", task.Name, padding, desc)
	}
}

// RunFromArgs runs a task from command line arguments
func RunFromArgs(args []string) error {
	if len(args) == 0 {
		PrintList()
		return nil
	}

	taskName := args[0]
	taskArgs := args[1:]

	return Run(taskName, taskArgs)
}

// DefaultTasks registers default tasks
func DefaultTasks() {
	Register("secret", "Generate a cryptographically secure secret key", func(args []string) error {
		// Generate 64 random bytes
		b := make([]byte, 64)
		_, err := rand.Read(b)
		if err != nil {
			return err
		}

		// Print as base64
		fmt.Println(base64.URLEncoding.EncodeToString(b))
		return nil
	})
}
