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

	helpText := tview.NewTextView().
		SetText("Space: Toggle completed | Delete: Remove task | a: Add new | q: Quit").
		SetTextAlign(tview.AlignCenter)

	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(todoList, 0, 1, true).
		AddItem(helpText, 1, 0, false)

	pages.AddPage("list", mainLayout, true, true)

	refreshTodoList(todoList, tasks, app, pages)

	todoList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ' ':
			index := todoList.GetCurrentItem()
			if index >= len(tasks) {
				return nil
			}
			tasks[index].Toggle()
			todo.SaveTasks(tasks, file)
			refreshTodoList(todoList, tasks, app, pages)
			return nil

		case 'j':
			current := todoList.GetCurrentItem()
			if current < todoList.GetItemCount()-2 {
				todoList.SetCurrentItem(current + 1)
			}
			return nil

		case 'k':
			current := todoList.GetCurrentItem()
			if current > 0 {
				todoList.SetCurrentItem(current - 1)
			}
			return nil

		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
