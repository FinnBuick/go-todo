package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UI struct {
	app     *tview.Application
	list    *tview.List
	input   *tview.InputField
	details *tview.TextView
	pages   *tview.Pages
	flex    *tview.Flex

	controller *AppController
}

const helpText = `[yellow]Controls:
[green]Tab:[white] Cycle Focus | [green]Enter (in list):[white] Toggle Done | [green]d (in list):[white] Delete
[green]Enter (in input):[white] Add Task | [green]Esc (in input):[white] Focus List | [green]q:[white] Quit`

func NewUI(controller *AppController) *UI {
	ui := &UI{
		app:        tview.NewApplication(),
		controller: controller,
	}

	ui.list = tview.NewList().ShowSecondaryText(false)
	ui.list.SetBorder(true).SetTitle("To-Do List")
	ui.list.SetSelectedFocusOnly(true)

	ui.input = tview.NewInputField().SetLabel("New Task: ").SetFieldWidth(0)
	ui.input.SetBorder(true)

	ui.details = tview.NewTextView().SetDynamicColors(true).SetScrollable(true).SetChangedFunc(func() {
		ui.app.Draw()
	})
	ui.details.SetBorder(true).SetTitle("Help / Info")
	fmt.Fprint(ui.details, helpText)

	leftPanel := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(ui.list, 0, 1, true).AddItem(ui.input, 3, 0, false)

	ui.flex = tview.NewFlex().AddItem(leftPanel, 0, 2, true).AddItem(ui.details, 0, 1, false)

	ui.pages = tview.NewPages().AddPage("main", ui.flex, true, true)

	ui.setupKeybindings()

	return ui
}

func (ui *UI) Run() error {
	return ui.app.SetRoot(ui.pages, true).EnableMouse(true).Run()
}

func (ui *UI) Stop() {
	ui.app.Stop()
}

func (ui *UI) RefreshList(tasks []Task) {
	currentSelection := ui.list.GetCurrentItem()
	ui.list.Clear()

	if len(tasks) == 0 {
		ui.list.AddItem(
			"No tasks yet!", "Press Tab then Enter in input field to add one.", 0, nil)
		ui.app.Draw()
		return
	}

	for _, task := range tasks {
		prefix := " [ ] "
		if task.Done {
			prefix = "[lime][✔][white] "
		}
		mainText := fmt.Sprintf("%s%s", prefix, task.Description)
		ui.list.AddItem(mainText, strconv.Itoa(task.ID), 0, func() {
			ui.controller.HandleToggleTask()
		})
	}

	if currentSelection >= 0 && currentSelection < ui.list.GetCurrentItem() {
		ui.list.SetCurrentItem(currentSelection)
	} else if ui.list.GetCurrentItem() > 0 {
		ui.list.SetCurrentItem(0)
	}
	ui.app.Draw()
}

func (ui *UI) GetSelectedTaskID() (int, bool) {
	if ui.list.GetItemCount() == 0 {
		return 0, false
	}
	index := ui.list.GetCurrentItem()
	if index < 0 {
		return 0, false
	}
	_, idStr := ui.list.GetItemText(index)
	taskID, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, false
	}
	return taskID, true
}

func (ui *UI) GetInputText() string {
	return strings.TrimSpace(ui.input.GetText())
}

func (ui *UI) ClearInput() {
	ui.input.SetText("")
}

func (ui *UI) FocusList() {
	ui.app.SetFocus(ui.list)
}

func (ui *UI) FocusInput() {
	ui.app.SetFocus(ui.input)
}

func (ui *UI) ShowConfirmation(message string, onConfirm func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Confirm", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.pages.RemovePage("confirmModal")
			if buttonLabel == "Confirm" {
				onConfirm()
			}
			ui.app.SetFocus(ui.list)
		})
	ui.pages.AddPage("confirmModal", modal, true, true)
}

func (ui *UI) ShowError(message string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("[red]Error:\n%s", message)).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.pages.RemovePage("errorModal")
			// Decide where to focus after error.
			if ui.app.GetFocus() == ui.input {
				ui.FocusInput()
			} else {
				ui.FocusList()
			}
		})
	ui.pages.AddPage("errorModal", modal, false, true)
}

func (ui *UI) setupKeybindings() {
	// List-specific keybindings
	ui.list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter: // Already handled by list.SetSelectedFunc
			ui.controller.HandleToggleTask()
			return nil
		case tcell.KeyRune:
			switch event.Rune() {
			case 'd':
				ui.controller.HandleDeleteTask()
				return nil
			}
		}
		return event
	})

	// Input field keybindings
	ui.input.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			ui.controller.HandleAddTask()
		case tcell.KeyEscape:
			ui.FocusList()
		}
	})

	// Global application keybindings
	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Check if a modal is active. If so, let the modal handle input.
		if name, _ := ui.pages.GetFrontPage(); name != "main" {
			return event // Pass event to modal
		}

		switch event.Key() {
		case tcell.KeyTab:
			if ui.input.HasFocus() {
				ui.FocusList()
			} else {
				ui.FocusInput()
			}
			return nil
		case tcell.KeyBacktab: // Shift+Tab
			if ui.input.HasFocus() {
				ui.FocusList()
			} else {
				ui.FocusInput()
			}
			return nil
		case tcell.KeyRune:
			if event.Rune() == 'q' {
				ui.controller.HandleQuit()
				return nil
			}
		}
		return event
	})
}
