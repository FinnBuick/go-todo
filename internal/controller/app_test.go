package controller

import (
	"go-todo/internal/models"
	"testing"
)

type MockStore struct {
	AddTaskCalls          int64
	ToggleTaskStatusCalls int
}

func (ms *MockStore) GetTasks() ([]models.Task, error) {
	return []models.Task{}, nil
}
func (ms *MockStore) AddTask(description string) (int64, error) {
	ms.AddTaskCalls++
	return ms.AddTaskCalls, nil
}
func (ms *MockStore) ToggleTaskStatus(id int) error {
	ms.ToggleTaskStatusCalls++
	return nil
}
func (ms *MockStore) DeleteTask(id int) error {
	return nil
}
func (ms *MockStore) Close() {}

type MockUI struct {
	GetInputTextCalls      int
	ClearInputCalls        int
	FocusListCalls         int
	GetSelectedTaskIDCalls int
	SelectedTaskID         int
	ShowErrorMsg           string
	inputText              string
}

func (mu *MockUI) Run() error {
	return nil
}
func (mu *MockUI) Stop() {
}
func (mu *MockUI) RefreshList(tasks []models.Task) {
}
func (mu *MockUI) GetInputText() string {
	mu.GetInputTextCalls++
	return mu.inputText
}
func (mu *MockUI) ClearInput() {
	mu.ClearInputCalls++
}
func (mu *MockUI) FocusList() {
	mu.FocusListCalls++
}
func (mu *MockUI) FocusInput() {
}
func (mu *MockUI) GetSelectedTaskID() (int, bool) {
	mu.GetSelectedTaskIDCalls++
	return mu.SelectedTaskID, true
}
func (mu *MockUI) GetItemCount() int {
	return 1
}
func (mu *MockUI) ShowError(message string) {
	mu.ShowErrorMsg = message
}
func (mu *MockUI) ShowConfirmation(message string, onConfirm func()) {
}

func setupTest(inputText string, selectedTaskID int) (*MockStore, *MockUI, *AppController) {
	mockStore := &MockStore{}
	mockUI := &MockUI{}
	mockUI.inputText = inputText
	mockUI.SelectedTaskID = selectedTaskID
	controller := NewAppController(mockStore)
	controller.SetUI(mockUI)
	return mockStore, mockUI, controller
}

func TestHandleAddTask_EmptyDescription(t *testing.T) {
	// Arrange
	mockStore, mockUI, controller := setupTest("", 1)

	// Act
	controller.HandleAddTask()

	// Assert
	if mockUI.GetInputTextCalls != 1 {
		t.Errorf("GetInputText called incorrect number of times, got=%d, expected=%d", mockUI.GetInputTextCalls, 1)
	}

	if mockUI.ShowErrorMsg != "Task description cannot be empty" {
		t.Errorf("ui error message incorrect, got=%s", mockUI.ShowErrorMsg)
	}

	if mockStore.AddTaskCalls != 0 {
		t.Errorf("AddTask should not be called, got=%d", mockStore.AddTaskCalls)
	}
}

func TestHandleAddTask_ValidDescription(t *testing.T) {
	// Arrange
	mockStore, mockUI, controller := setupTest("Valid description", 1)

	// Act
	controller.HandleAddTask()

	// Assert
	if mockUI.GetInputTextCalls != 1 {
		t.Errorf("GetInputText called incorrect number of times, got=%d, expected=%d", mockUI.GetInputTextCalls, 1)
	}

	if mockUI.ShowErrorMsg != "" {
		t.Errorf("ui error message incorrect, got=%s, expected=%s", mockUI.ShowErrorMsg, "")
	}

	if mockStore.AddTaskCalls != 1 {
		t.Errorf("AddTask should be called once, got=%d", mockStore.AddTaskCalls)
	}

	if mockUI.ClearInputCalls != 1 {
		t.Errorf("ClearInput should be called once, got=%d", mockUI.ClearInputCalls)
	}

	if mockUI.FocusListCalls != 1 {
		t.Errorf("FocusList should be called once, got=%d", mockUI.FocusListCalls)
	}
}

func TestHandleToggleTask_NoExistingTask(t *testing.T) {
	// Arrange
	mockStore, mockUI, controller := setupTest("", 0)

	// Act
	controller.HandleToggleTask()

	// Assert
	if mockUI.GetSelectedTaskIDCalls != 1 {
		t.Errorf("GetSelectedTaskID should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}

	if mockUI.SelectedTaskID != 0 {
		t.Errorf("SelectedTaskID incorrect, got=%d", mockUI.SelectedTaskID)
	}

	if mockStore.ToggleTaskStatusCalls != 1 {
		t.Errorf("ToggleTaskStatus should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}
}

func TestHandleToggleTask_ExistingTask(t *testing.T) {
	// Arrange
	mockStore, mockUI, controller := setupTest("Description", 1)

	// Act
	controller.HandleToggleTask()

	// Assert
	if mockUI.GetSelectedTaskIDCalls != 1 {
		t.Errorf("GetSelectedTaskID should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}

	if mockStore.ToggleTaskStatusCalls != 1 {
		t.Errorf("ToggleTaskStatus should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}
}
