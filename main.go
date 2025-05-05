package main

import (
	"flag"
	"fmt"
	"go-todo/todo"
)

// Commands
// list - lists the existing todos
// add "Buy milk" - adds a todo
// done - marks a todos as completed

func main() {
	add := flag.String("add", "", "Task to add")
	list := flag.Bool("list", false, "List tasks")
	done := flag.Int("done", -1, "Mark task as done")
	flag.Parse()

	const file = "todos.json"
	tasks, err := todo.LoadTasks(file)
	if err != nil {
		fmt.Println("Error loading task:", err)
		return
	}

	switch {
	case *add != "":
		fmt.Println(tasks)
		t := todo.Task{ID: len(tasks) + 1, Text: *add}
		tasks = append(tasks, t)
		todo.SaveTasks(tasks, file)
		fmt.Println("Added task:", t.Text)
	case *list:
		for _, t := range tasks {
			status := ""
			if t.Completed {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s\n", status, t.ID, t.Text)
		}
	case *done > 0:
		for i, t := range tasks {
			if t.ID == *done {
				tasks[i].Completed = true
				todo.SaveTasks(tasks, file)
				fmt.Println("Marked as done:", t.Text)
				break
			}
		}
	default:
		fmt.Println("Usage: -add \"task\" | -list | -done N")
	}
}
