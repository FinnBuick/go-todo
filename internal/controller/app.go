package controller

import (
	"fmt"
	"log"

	"go-todo/internal/models"
)

type AppController struct {
	store Store
	ui    UI
}

type Store interface {
	GetTasks() ([]models.Task, error)
	AddTask(description string) (int64, error)
	ToggleTaskStatus(id int) error
	DeleteTask(id int) error
	Close()
}

type UI interface {
	Run() error
	Stop()
	RefreshList(tasks []models.Task)
	GetInputText() string
	ClearInput()
	FocusList()
	FocusInput()
	GetSelectedTaskID() (int, bool)
	GetItemCount() int
	ShowError(message string)
	ShowConfirmation(message string, onConfirm func())
}

func NewAppController(store Store) *AppController {
	return &AppController{
		store: store,
	}
}

func (c *AppController) SetUI(ui UI) {
	c.ui = ui
}

func (c *AppController) Start() error {
	if c.ui == nil {
		return fmt.Errorf("UI not initialised for controller")
	}
	log.Println("Loading and displaying tasks...")
	c.loadAndDisplayTasks()
	log.Println("Starting UI...")
	return c.ui.Run()
}

func (c *AppController) loadAndDisplayTasks() error {
	log.Println("Getting tasks from store...")
	tasks, err := c.store.GetTasks()
	if err != nil {
		log.Printf("Error loading tasks: %v", err)
		c.ui.ShowError(fmt.Sprintf("Failed to load tasks: %v", err))
		return err
	}
	log.Printf("Retrieved %d tasks, refreshing UI list...", len(tasks))
	c.ui.RefreshList(tasks)
	return nil
}

func (c *AppController) HandleAddTask() {
	description := c.ui.GetInputText()
	if description == "" {
		c.ui.ShowError(fmt.Sprintf("Task description cannot be empty"))
		return
	}

	_, err := c.store.AddTask(description)
	if err != nil {
		log.Printf("Error adding tasks: %v", err)
		c.ui.ShowError(fmt.Sprintf("Failed to add task: %v", err))
		return
	}

	c.ui.ClearInput()
	c.loadAndDisplayTasks()
	c.ui.FocusList()
}

func (c *AppController) HandleToggleTask() {
	taskID, selected := c.ui.GetSelectedTaskID()
	if !selected {
		log.Println("Toggle attempted on invalid or no selection.")
		return
	}

	err := c.store.ToggleTaskStatus(taskID)
	if err != nil {
		log.Printf("Error toggling task %d: %v", taskID, err)
		c.ui.ShowError(fmt.Sprintf("Failed to toggle task ID %d: %v", taskID, err))
		return
	}
	c.loadAndDisplayTasks()
}

func (c *AppController) HandleDeleteTask() {
	taskID, selected := c.ui.GetSelectedTaskID()
	if !selected {
		log.Println("Delete attempted on invalid or no selection.")
		return
	}

	confirmMsg := fmt.Sprintf("Are you sure you want to delete task ID %d?", taskID)

	c.ui.ShowConfirmation(confirmMsg, func() {
		err := c.store.DeleteTask(taskID)
		if err != nil {
			log.Printf("Error deleting task %d: %v", taskID, err)
			c.ui.ShowError(fmt.Sprintf("Failed to delete task ID %d: %v", taskID, err))
			return
		}
		c.loadAndDisplayTasks()
		if c.ui.GetItemCount() == 0 {
			c.ui.FocusInput()
		}
	})
}

func (c *AppController) HandleQuit() {
	c.ui.Stop()
}
