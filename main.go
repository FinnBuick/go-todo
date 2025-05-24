package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const file = "todos.json"

type TodoList struct {
	tasks          []Task
	app            *tview.Application
	mainFlex       *tview.Flex
	list           *tview.List
	helpText       *tview.TextView
	inputField     *tview.InputField
	isShowingInput bool
}

func newTodoList() *TodoList {
	return &TodoList{
		app:      tview.NewApplication(),
		mainFlex: tview.NewFlex().SetDirection(tview.FlexRow),
		list:     tview.NewList().ShowSecondaryText(false),
		helpText: tview.NewTextView().
			SetText("Space: Toggle completed | Delete: Remove task | a: Add new | q: Quit").
			SetTextAlign(tview.AlignCenter),
		inputField:     tview.NewInputField().SetLabel("New task: "),
		isShowingInput: false,
	}
}

func (t *TodoList) loadTasks() error {
	tasks, err := LoadTasks(file)
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
		if task.Done {
			status = "[x] "
		}
		id := task.ID
		escapedStatus := tview.Escape(status)
		t.list.AddItem(fmt.Sprintf("%s%d: %s", escapedStatus, id, task.Description), "", rune(0), nil)
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

func (t *TodoList) addTask(text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}

	nextID := 1
	for _, task := range t.tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}

	newTask := NewTask(text, nextID)

	t.tasks = append(t.tasks, newTask)
	SaveTasks(t.tasks, file)
	t.refreshTodoList()
}

func (t *TodoList) deleteCurrentTask() {
	index := t.list.GetCurrentItem()
	if index >= len(t.tasks) {
		return
	}
	t.tasks = slices.Delete(t.tasks, index, index+1)
	SaveTasks(t.tasks, file)
	t.refreshTodoList()
}

func (t *TodoList) showInput() {
	if t.isShowingInput {
		return
	}
	t.isShowingInput = true
	t.mainFlex.AddItem(t.inputField, 1, 0, true)
	t.app.SetFocus(t.inputField)
}

func (t *TodoList) hideInput() {
	if !t.isShowingInput {
		return
	}
	t.isShowingInput = false
	t.mainFlex.RemoveItem(t.inputField)
	t.inputField.SetText("")
	t.app.SetFocus(t.list)
}

func (t *TodoList) setupKeyBindings() {
	t.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyBackspace2:
			t.deleteCurrentTask()
			return nil
		}

		switch event.Rune() {
		case ' ':
			index := t.list.GetCurrentItem()
			if index >= len(t.tasks) {
				return nil
			}
			t.tasks[index].Toggle()
			SaveTasks(t.tasks, file)
			t.refreshTodoList()
			return nil
		case 'j':
			index := t.list.GetCurrentItem()
			if index < t.list.GetItemCount()-1 {
				t.list.SetCurrentItem(index + 1)
			}
			return nil
		case 'k':
			index := t.list.GetCurrentItem()
			if index > 0 {
				t.list.SetCurrentItem(index - 1)
			}
			return nil
		case 'a':
			t.showInput()
			return nil
		}
		return event
	})

	// Set up input field to handle ESC to cancel and Enter to add
	t.inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			t.hideInput()
			return nil
		case tcell.KeyEnter:
			text := t.inputField.GetText()
			t.addTask(text)
			t.hideInput()
			return nil
		}
		return event
	})
}

func (t *TodoList) setupLayout() {
	t.mainFlex.AddItem(t.list, 0, 1, true).
		AddItem(t.helpText, 1, 0, false)

	t.inputField.
		SetFieldWidth(40).
		SetFieldBackgroundColor(tcell.ColorDefault).
		SetLabel("New task: ")
}

func (t *TodoList) run() error {
	t.setupLayout()
	t.setupKeyBindings()
	t.refreshTodoList()
	return t.app.SetRoot(t.mainFlex, true).EnableMouse(true).Run()
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
