package main_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	fmt.Println("Building tool...")
	if os.Getenv("TODO_FILENAME") != "" {
		fileName = os.Getenv("TODO_FILENAME")
	}
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	build := exec.Command("go", "build", "-o", binName)
	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
	}

	fmt.Println("Running tests...")
	result := m.Run()
	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Remove(fileName)

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	task := "test task number 1"

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmdPath := filepath.Join(dir, binName)

	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task) // command path points to binary file like ./todo ( in mac )
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

	})
	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdSTDIN, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		io.WriteString(cmdSTDIN, task2)
		cmdSTDIN.Close()
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("AddMultipleTasksFromMultilineInput", func(t *testing.T) {
		multilineTasks := "task3\ntask4\ntask5"
		cmd := exec.Command(cmdPath, "-add", multilineTasks)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list") // empty arguments mean list all tasks
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf(" 1: %s\n 2: %s\n 3: task3\n 4: task4\n 5: task5\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}

		t.Run("ListVerboseTasks", func(t *testing.T) {
			cmd := exec.Command(cmdPath, "-list", "-verbose")
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatal(err)
			}

			if !strings.Contains(string(out), "Created at:") {
				t.Errorf("Expected verbose output with 'Created at:', got %q instead\n", string(out))
			}
		})
	})
	deleteTaskNum := "3"
	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", deleteTaskNum)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}

		// Now list tasks to verify deletion
		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}

		// Expected output should include all tasks except task3
		expected := fmt.Sprintf(" 1: %s\n 2: %s\n 3: task4\n 4: task5\n", task, task2)
		if expected != string(out) {
			t.Errorf("Expected %q, got %q instead\n", expected, string(out))
		}
	})
}
