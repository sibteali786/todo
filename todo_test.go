package todo_test

import (
	"testing"

	"mygoprograms.com/todo"
)

func TestAdd(t *testing.T) {
	l := todo.List{}
	taskName := "New york"
	l.Add(taskName)
	if l[0].Task != taskName {
		t.Errorf("Expected %q, got %q instead.", taskName, l[0].Task)
	}
}
