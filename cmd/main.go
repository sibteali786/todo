package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"mygoprograms.com/todo"
)

var todoFileName = ".todo.json"

// check if user defined env variable exists
func main() {
	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for Fun\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2024\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Define command-line flags
	add := flag.Bool("add", false, "Add task to the ToDo List")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	delete := flag.Int("delete", 0, "Item to be deleted")
	verbose := flag.Bool("verbose", false, "Enable verbose output with date/time")

	// Parse the flags
	flag.Parse()

	// Create a new to-do list
	l := &todo.List{}

	// Retrieve the to-do items
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	switch {
	case *list:
		fmt.Print(l.String(*verbose))
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// Add a new task
		tasks, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, task := range tasks {
			l.Add(task)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *delete > 0:
		if err := l.Delete(*delete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid option, show usage message
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}

}
func getTask(r io.Reader, args ...string) ([]string, error) {
	if len(args) > 0 {
		// Join arguments and split by '\n' to handle multiline input
		input := strings.Join(args, " ")
		lines := strings.Split(input, "\n")

		var tasks []string
		for _, line := range lines {
			task := strings.TrimSpace(line)
			if len(task) > 0 {
				tasks = append(tasks, task)
			}
		}

		if len(tasks) == 0 {
			return nil, fmt.Errorf("No valid tasks provided")
		}

		return tasks, nil
	}

	s := bufio.NewScanner(r)
	var tasks []string
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) > 0 {
			tasks = append(tasks, line)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, fmt.Errorf("No tasks entered")
	}
	return tasks, nil
}
