package main

import (
	"flag"
	"fmt"
	"os"

	"mygoprograms.com/todo"
)

const todoFileName = ".todo.json"

func main() {
	// parsing command line flags
	task := flag.String("task", "", "Task to be included in the ToDo List")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()
	l := &todo.List{}

	//use get command to read to do items from list
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	switch {
	case *list:
		// List current todo items
		for _, item := range *l {
			if !item.Done { // this makes sure we display only tasks not completed
				fmt.Println(item.Task)
			}
		}
	case *complete > 0:
		// complete to given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// save the new list
		if err := l.Save(todoFileName); err != nil {

			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *task != "":
		// Add the task
		l.Add(*task)
		// Save the list
		if err := l.Save(todoFileName); err != nil {

			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}

}
