package main

import (
	"fmt"
	"go-todo/todo"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const file = "todos.json"

func main() {

	tasks, err := todo.LoadTasks(file)
	if err != nil {
		fmt.Println("Error loading task:", err)
		return
	}

	app := tview.NewApplication()

	pages := tview.NewPages()

	// Create list to display tasks
	todoList := tview.NewList().
		SetHighlightFullLine(true).
		ShowSecondaryText(false)

	refreshTodoList := func(list *tview.List, tasks []todo.Task, app *tview.Application, pages *tview.Pages) {
		list.Clear()

		for _, t := range tasks {
			status := "[ ] "
			if t.Completed {
				status = "[x] "
			}

			id := t.ID
			list.AddItem(fmt.Sprintf("%s%d: %s", status, id, t.Text), "", rune(0), nil)
		}

		list.AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})
	}

	// Create a help text view
	helpText := tview.NewTextView().
		SetText("Space: Toggle completed | Delete: Remove task | a: Add new | q: Quit").
		SetTextAlign(tview.AlignCenter)

	// Create the main layout
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(todoList, 0, 1, true).
		AddItem(helpText, 1, 0, false)

	pages.AddPage("list", mainLayout, true, true)

	// initialise the list
	refreshTodoList(todoList, tasks, app, pages)

	// Key handler
	todoList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ' ':
			index := todoList.GetCurrentItem()
			if index >= len(tasks) {
				return nil // "Quit" item
			}
			tasks[index].Toggle()
			todo.SaveTasks(tasks, file) // Save changes
			refreshTodoList(todoList, tasks, app, pages)
			return nil
		}
		return event
	})

	// Run the application
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
