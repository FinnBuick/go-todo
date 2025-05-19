package main

import (
	"fmt"
	"go-todo/todo"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const file = "todos.json"

type TodoList struct {
	app      *tview.Application
	pages    *tview.Pages
	list     *tview.List
	tasks    []todo.Task
	helpText *tview.TextView
	form     *tview.Form
}

func newTodoList() *TodoList {
	return &TodoList{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),
		list:  tview.NewList().ShowSecondaryText(false),
		helpText: tview.NewTextView().
			SetText("Space: Toggle completed | Delete: Remove task | a: Add new | q: Quit").
			SetTextAlign(tview.AlignCenter),
		form: tview.NewForm(),
	}
}

func (t *TodoList) loadTasks() error {
	tasks, err := todo.LoadTasks(file)
	if err != nil {
		return err
	}
	t.tasks = tasks
	return nil
}

func (t *TodoList) refreshTodoList() {
	currentIndex := t.list.GetCurrentItem()

	t.list.Clear()

	for _, task := range t.tasks {
		status := "[ ] "
		if task.Completed {
			status = "[x] "
		}

		id := task.ID
		escapedStatus := tview.Escape(status)
		t.list.AddItem(fmt.Sprintf("%s%d: %s", escapedStatus, id, task.Text), "", rune(0), nil)
	}

	t.list.AddItem("Quit", "Press to exit", 'q', func() {
		t.app.Stop()
	})

	if currentIndex >= 0 && currentIndex < t.list.GetItemCount() {
		t.list.SetCurrentItem(currentIndex)
	} else if t.list.GetItemCount() > 0 {
		t.list.SetCurrentItem(0)
	}
}

func (t *TodoList) setupKeyBindings() {
	t.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case ' ':
			index := t.list.GetCurrentItem()
			if index >= len(t.tasks) {
				return nil
			}
			t.tasks[index].Toggle()
			todo.SaveTasks(t.tasks, file)
			t.refreshTodoList()
			return nil
		case 'j':
			current := t.list.GetCurrentItem()
			if current < t.list.GetItemCount()-2 {
				t.list.SetCurrentItem(current + 1)
			}
			return nil
		case 'k':
			current := t.list.GetCurrentItem()
			if current > 0 {
				t.list.SetCurrentItem(current - 1)
			}
			return nil

		}
		return event
	})
}

func (t *TodoList) setupLayout() {
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.list, 0, 1, true).
		AddItem(t.helpText, 1, 0, false)

	t.pages.AddPage("list", mainLayout, true, true)
}

func (t *TodoList) run() error {
	t.setupLayout()
	t.setupKeyBindings()
	t.refreshTodoList()
	return t.app.SetRoot(t.pages, true).EnableMouse(true).Run()
}

func main() {
	app := newTodoList()

	if err := app.loadTasks(); err != nil {
		fmt.Println("Error loading tasks:", err)
		return
	}

	if err := app.run(); err != nil {
		panic(err)
	}
}
