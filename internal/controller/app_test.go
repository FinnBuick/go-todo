package controller

import (
	"errors"
	"go-todo/internal/models"
	"io"
	"log"
	"os"
	"strings"
	"testing"
)

// TestMain disables logging for all tests
func TestMain(m *testing.M) {
	// Disable logging during tests
	log.SetOutput(io.Discard)

	// Run tests
	code := m.Run()

	// Exit with the same code as the tests
	os.Exit(code)
}

type MockStore struct {
	AddTaskCalls          int64
	ToggleTaskStatusCalls int
	DeleteTaskCalls       int
	GetTasksCalls         int
	CloseCalls            int

	// Control behavior
	GetTasksError   error
	AddTaskError    error
	ToggleTaskError error
	DeleteTaskError error
	TasksToReturn   []models.Task
}

func (ms *MockStore) GetTasks() ([]models.Task, error) {
	ms.GetTasksCalls++
	if ms.GetTasksError != nil {
		return nil, ms.GetTasksError
	}
	return ms.TasksToReturn, nil
}

func (ms *MockStore) AddTask(description string) (int64, error) {
	ms.AddTaskCalls++
	if ms.AddTaskError != nil {
		return 0, ms.AddTaskError
	}
	return ms.AddTaskCalls, nil
}

func (ms *MockStore) ToggleTaskStatus(id int) error {
	ms.ToggleTaskStatusCalls++
	return ms.ToggleTaskError
}

func (ms *MockStore) DeleteTask(id int) error {
	ms.DeleteTaskCalls++
	return ms.DeleteTaskError
}

func (ms *MockStore) Close() {
	ms.CloseCalls++
}

type MockUI struct {
	GetInputTextCalls        int
	ClearInputCalls          int
	FocusListCalls           int
	FocusInputCalls          int
	GetSelectedTaskIDCalls   int
	GetSelectedTaskTextCalls int
	GetItemCountCalls        int
	ShowErrorCalls           int
	ShowConfirmationCalls    int
	RefreshListCalls         int
	StopCalls                int
	RunCalls                 int

	// Control behavior
	SelectedTaskID       int
	SelectedTaskText     string
	TaskSelected         bool
	ItemCount            int
	ShowErrorMsg         string
	ShowConfirmationMsg  string
	inputText            string
	RunError             error
	ConfirmationCallback func()
	TasksReceived        []models.Task
}

func (mu *MockUI) Run() error {
	mu.RunCalls++
	return mu.RunError
}

func (mu *MockUI) Stop() {
	mu.StopCalls++
}

func (mu *MockUI) RefreshList(tasks []models.Task) {
	mu.RefreshListCalls++
	mu.TasksReceived = tasks
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
	mu.FocusInputCalls++
}

func (mu *MockUI) GetSelectedTaskID() (int, bool) {
	mu.GetSelectedTaskIDCalls++
	return mu.SelectedTaskID, mu.TaskSelected
}

func (mu *MockUI) GetSelectedTaskText() (string, bool) {
	mu.GetSelectedTaskTextCalls++
	return mu.SelectedTaskText, mu.TaskSelected
}

func (mu *MockUI) GetItemCount() int {
	mu.GetItemCountCalls++
	return mu.ItemCount
}

func (mu *MockUI) ShowError(message string) {
	mu.ShowErrorCalls++
	mu.ShowErrorMsg = message
}

func (mu *MockUI) ShowConfirmation(message string, onConfirm func()) {
	mu.ShowConfirmationCalls++
	mu.ShowConfirmationMsg = message
	mu.ConfirmationCallback = onConfirm
}

func setupTest(inputText string, selectedTaskID int, taskSelected bool) (*MockStore, *MockUI, *AppController) {
	mockStore := &MockStore{}
	mockUI := &MockUI{
		inputText:      inputText,
		SelectedTaskID: selectedTaskID,
		TaskSelected:   taskSelected,
		ItemCount:      1,
	}
	controller := NewAppController(mockStore)
	controller.SetUI(mockUI)
	return mockStore, mockUI, controller
}

// Test Controller Creation and Initialization
func TestNewAppController(t *testing.T) {
	mockStore := &MockStore{}
	controller := NewAppController(mockStore)

	if controller.store != mockStore {
		t.Error("Controller store not set correctly")
	}

	if controller.ui != nil {
		t.Error("Controller UI should be nil initially")
	}
}

func TestSetUI(t *testing.T) {
	mockStore := &MockStore{}
	mockUI := &MockUI{}
	controller := NewAppController(mockStore)

	controller.SetUI(mockUI)

	if controller.ui != mockUI {
		t.Error("Controller UI not set correctly")
	}
}

// Test Start Method
func TestStart_NoUI(t *testing.T) {
	mockStore := &MockStore{}
	controller := NewAppController(mockStore)

	err := controller.Start()

	if err == nil {
		t.Error("Expected error when UI not set")
	}

	expectedMsg := "UI not initialised for controller"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestStart_Success(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 0, false)
	mockStore.TasksToReturn = []models.Task{
		{ID: 1, Description: "Test task", Done: false},
	}

	err := controller.Start()

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if mockStore.GetTasksCalls != 1 {
		t.Errorf("GetTasks should be called once, got %d", mockStore.GetTasksCalls)
	}

	if mockUI.RefreshListCalls != 1 {
		t.Errorf("RefreshList should be called once, got %d", mockUI.RefreshListCalls)
	}

	if mockUI.RunCalls != 1 {
		t.Errorf("Run should be called once, got %d", mockUI.RunCalls)
	}

	if len(mockUI.TasksReceived) != 1 {
		t.Errorf("Expected 1 task to be passed to UI, got %d", len(mockUI.TasksReceived))
	}
}

func TestStart_GetTasksError(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 0, false)
	mockStore.GetTasksError = errors.New("database error")

	err := controller.Start()

	if err != nil {
		t.Errorf("Start should not return error even if GetTasks fails: %v", err)
	}

	if mockUI.ShowErrorCalls != 1 {
		t.Errorf("ShowError should be called once, got %d", mockUI.ShowErrorCalls)
	}

	if !strings.Contains(mockUI.ShowErrorMsg, "Failed to load tasks") {
		t.Errorf("Expected error message about loading tasks, got '%s'", mockUI.ShowErrorMsg)
	}
}

func TestStart_UIRunError(t *testing.T) {
	_, mockUI, controller := setupTest("", 0, false)
	mockUI.RunError = errors.New("UI error")

	err := controller.Start()

	if err == nil {
		t.Error("Expected error from UI.Run to be propagated")
	}

	if err.Error() != "UI error" {
		t.Errorf("Expected 'UI error', got '%s'", err.Error())
	}
}

// Test HandleAddTask
func TestHandleAddTask_EmptyDescription(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 1, true)

	controller.HandleAddTask()

	if mockUI.GetInputTextCalls != 1 {
		t.Errorf("GetInputText called incorrect number of times, got=%d, expected=1", mockUI.GetInputTextCalls)
	}

	if mockUI.ShowErrorMsg != "Task description cannot be empty" {
		t.Errorf("Expected empty description error, got='%s'", mockUI.ShowErrorMsg)
	}

	if mockStore.AddTaskCalls != 0 {
		t.Errorf("AddTask should not be called, got=%d", mockStore.AddTaskCalls)
	}

	if mockUI.ClearInputCalls != 0 {
		t.Errorf("ClearInput should not be called on error, got=%d", mockUI.ClearInputCalls)
	}
}

func TestHandleAddTask_ValidDescription(t *testing.T) {
	mockStore, mockUI, controller := setupTest("Valid description", 1, true)
	mockStore.TasksToReturn = []models.Task{
		{ID: 1, Description: "Valid description", Done: false},
	}

	controller.HandleAddTask()

	if mockUI.GetInputTextCalls != 1 {
		t.Errorf("GetInputText called incorrect number of times, got=%d, expected=1", mockUI.GetInputTextCalls)
	}

	if mockUI.ShowErrorMsg != "" {
		t.Errorf("Expected no error message, got='%s'", mockUI.ShowErrorMsg)
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

	// Verify tasks are reloaded
	if mockStore.GetTasksCalls != 1 {
		t.Errorf("GetTasks should be called to reload tasks, got=%d", mockStore.GetTasksCalls)
	}

	if mockUI.RefreshListCalls != 1 {
		t.Errorf("RefreshList should be called to update UI, got=%d", mockUI.RefreshListCalls)
	}
}

func TestHandleAddTask_StoreError(t *testing.T) {
	mockStore, mockUI, controller := setupTest("Valid description", 1, true)
	mockStore.AddTaskError = errors.New("store error")

	controller.HandleAddTask()

	if mockStore.AddTaskCalls != 1 {
		t.Errorf("AddTask should be called once, got=%d", mockStore.AddTaskCalls)
	}

	if mockUI.ShowErrorCalls != 1 {
		t.Errorf("ShowError should be called once, got=%d", mockUI.ShowErrorCalls)
	}

	if !strings.Contains(mockUI.ShowErrorMsg, "Failed to add task") {
		t.Errorf("Expected add task error message, got='%s'", mockUI.ShowErrorMsg)
	}

	if mockUI.ClearInputCalls != 0 {
		t.Errorf("ClearInput should not be called on error, got=%d", mockUI.ClearInputCalls)
	}
}

func TestHandleAddTask_WhitespaceOnlyDescription(t *testing.T) {
	mockStore, _, controller := setupTest("   ", 1, true)

	controller.HandleAddTask()

	// Note: Current implementation doesn't trim whitespace
	// This test documents current behavior - you might want to add trimming
	if mockStore.AddTaskCalls != 1 {
		t.Errorf("AddTask called with whitespace-only description, got=%d", mockStore.AddTaskCalls)
	}
}

// Test HandleToggleTask
func TestHandleToggleTask_NoSelection(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 0, false)

	controller.HandleToggleTask()

	if mockUI.GetSelectedTaskIDCalls != 1 {
		t.Errorf("GetSelectedTaskID should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}

	if mockStore.ToggleTaskStatusCalls != 0 {
		t.Errorf("ToggleTaskStatus should not be called when no selection, got=%d", mockStore.ToggleTaskStatusCalls)
	}
}

func TestHandleToggleTask_ValidSelection(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 1, true)
	mockStore.TasksToReturn = []models.Task{
		{ID: 1, Description: "Test task", Done: true},
	}

	controller.HandleToggleTask()

	if mockUI.GetSelectedTaskIDCalls != 1 {
		t.Errorf("GetSelectedTaskID should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}

	if mockStore.ToggleTaskStatusCalls != 1 {
		t.Errorf("ToggleTaskStatus should be called once, got=%d", mockStore.ToggleTaskStatusCalls)
	}

	// Verify tasks are reloaded after toggle
	if mockStore.GetTasksCalls != 1 {
		t.Errorf("GetTasks should be called to reload tasks, got=%d", mockStore.GetTasksCalls)
	}

	if mockUI.RefreshListCalls != 1 {
		t.Errorf("RefreshList should be called to update UI, got=%d", mockUI.RefreshListCalls)
	}
}

func TestHandleToggleTask_StoreError(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 1, true)
	mockStore.ToggleTaskError = errors.New("toggle error")

	controller.HandleToggleTask()

	if mockStore.ToggleTaskStatusCalls != 1 {
		t.Errorf("ToggleTaskStatus should be called once, got=%d", mockStore.ToggleTaskStatusCalls)
	}

	if mockUI.ShowErrorCalls != 1 {
		t.Errorf("ShowError should be called once, got=%d", mockUI.ShowErrorCalls)
	}

	if !strings.Contains(mockUI.ShowErrorMsg, "Failed to toggle task ID 1") {
		t.Errorf("Expected toggle error message, got='%s'", mockUI.ShowErrorMsg)
	}

	// Should not reload tasks on error
	if mockStore.GetTasksCalls != 0 {
		t.Errorf("GetTasks should not be called on toggle error, got=%d", mockStore.GetTasksCalls)
	}
}

// Test HandleDeleteTask
func TestHandleDeleteTask_NoSelection(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 0, false)

	controller.HandleDeleteTask()

	if mockUI.GetSelectedTaskIDCalls != 1 {
		t.Errorf("GetSelectedTaskID should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}

	if mockUI.ShowConfirmationCalls != 0 {
		t.Errorf("ShowConfirmation should not be called when no selection, got=%d", mockUI.ShowConfirmationCalls)
	}

	if mockStore.DeleteTaskCalls != 0 {
		t.Errorf("DeleteTask should not be called when no selection, got=%d", mockStore.DeleteTaskCalls)
	}
}

func TestHandleDeleteTask_ValidSelection_Confirmed(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 5, true)
	mockStore.TasksToReturn = []models.Task{}
	mockUI.ItemCount = 0

	controller.HandleDeleteTask()

	if mockUI.GetSelectedTaskIDCalls != 1 {
		t.Errorf("GetSelectedTaskID should be called once, got=%d", mockUI.GetSelectedTaskIDCalls)
	}

	if mockUI.ShowConfirmationCalls != 1 {
		t.Errorf("ShowConfirmation should be called once, got=%d", mockUI.ShowConfirmationCalls)
	}

	expectedMsg := "Are you sure you want to delete task ID 5?"
	if mockUI.ShowConfirmationMsg != expectedMsg {
		t.Errorf("Expected confirmation message '%s', got '%s'", expectedMsg, mockUI.ShowConfirmationMsg)
	}

	// Simulate user confirming deletion
	if mockUI.ConfirmationCallback != nil {
		mockUI.ConfirmationCallback()
	}

	if mockStore.DeleteTaskCalls != 1 {
		t.Errorf("DeleteTask should be called once after confirmation, got=%d", mockStore.DeleteTaskCalls)
	}

	// Verify tasks are reloaded
	if mockStore.GetTasksCalls != 1 {
		t.Errorf("GetTasks should be called to reload tasks, got=%d", mockStore.GetTasksCalls)
	}

	// Should focus input when no items left
	if mockUI.FocusInputCalls != 1 {
		t.Errorf("FocusInput should be called when no items left, got=%d", mockUI.FocusInputCalls)
	}
}

func TestHandleDeleteTask_ValidSelection_WithItemsRemaining(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 1, true)
	mockStore.TasksToReturn = []models.Task{
		{ID: 2, Description: "Remaining task", Done: false},
	}
	mockUI.ItemCount = 1

	controller.HandleDeleteTask()

	// Simulate user confirming deletion
	if mockUI.ConfirmationCallback != nil {
		mockUI.ConfirmationCallback()
	}

	if mockStore.DeleteTaskCalls != 1 {
		t.Errorf("DeleteTask should be called once after confirmation, got=%d", mockStore.DeleteTaskCalls)
	}

	// Should not focus input when items remain
	if mockUI.FocusInputCalls != 0 {
		t.Errorf("FocusInput should not be called when items remain, got=%d", mockUI.FocusInputCalls)
	}
}

func TestHandleDeleteTask_StoreError(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 1, true)
	mockStore.DeleteTaskError = errors.New("delete error")

	controller.HandleDeleteTask()

	// Simulate user confirming deletion
	if mockUI.ConfirmationCallback != nil {
		mockUI.ConfirmationCallback()
	}

	if mockStore.DeleteTaskCalls != 1 {
		t.Errorf("DeleteTask should be called once, got=%d", mockStore.DeleteTaskCalls)
	}

	if mockUI.ShowErrorCalls != 1 {
		t.Errorf("ShowError should be called once, got=%d", mockUI.ShowErrorCalls)
	}

	if !strings.Contains(mockUI.ShowErrorMsg, "Failed to delete task ID 1") {
		t.Errorf("Expected delete error message, got='%s'", mockUI.ShowErrorMsg)
	}

	// Should not reload tasks on error
	if mockStore.GetTasksCalls != 0 {
		t.Errorf("GetTasks should not be called on delete error, got=%d", mockStore.GetTasksCalls)
	}
}

// Test HandleQuit
func TestHandleQuit(t *testing.T) {
	_, mockUI, controller := setupTest("", 0, false)

	controller.HandleQuit()

	if mockUI.StopCalls != 1 {
		t.Errorf("Stop should be called once, got=%d", mockUI.StopCalls)
	}
}

// Test loadAndDisplayTasks (private method, tested through public methods)
func TestLoadAndDisplayTasks_Success(t *testing.T) {
	mockStore, mockUI, controller := setupTest("", 0, false)
	expectedTasks := []models.Task{
		{ID: 1, Description: "Task 1", Done: false},
		{ID: 2, Description: "Task 2", Done: true},
	}
	mockStore.TasksToReturn = expectedTasks

	// Test through Start method
	controller.Start()

	if mockStore.GetTasksCalls != 1 {
		t.Errorf("GetTasks should be called once, got=%d", mockStore.GetTasksCalls)
	}

	if mockUI.RefreshListCalls != 1 {
		t.Errorf("RefreshList should be called once, got=%d", mockUI.RefreshListCalls)
	}

	if len(mockUI.TasksReceived) != 2 {
		t.Errorf("Expected 2 tasks to be passed to UI, got %d", len(mockUI.TasksReceived))
	}

	// Verify the actual tasks passed
	for i, task := range expectedTasks {
		if i < len(mockUI.TasksReceived) {
			received := mockUI.TasksReceived[i]
			if received.ID != task.ID || received.Description != task.Description || received.Done != task.Done {
				t.Errorf("Task %d mismatch: expected %+v, got %+v", i, task, received)
			}
		}
	}
}

// Integration-style tests
func TestAddTaskWorkflow_CompleteFlow(t *testing.T) {
	mockStore, mockUI, controller := setupTest("Buy groceries", 0, false)

	// First, start the controller to load initial tasks
	mockStore.TasksToReturn = []models.Task{}
	controller.Start()

	// Reset call counts after start
	mockStore.GetTasksCalls = 0
	mockUI.RefreshListCalls = 0

	// Now add a task
	mockStore.TasksToReturn = []models.Task{
		{ID: 1, Description: "Buy groceries", Done: false},
	}

	controller.HandleAddTask()

	// Verify the complete workflow
	if mockStore.AddTaskCalls != 1 {
		t.Errorf("Expected AddTask to be called once, got %d", mockStore.AddTaskCalls)
	}

	if mockStore.GetTasksCalls != 1 {
		t.Errorf("Expected GetTasks to be called once to reload, got %d", mockStore.GetTasksCalls)
	}

	if mockUI.ClearInputCalls != 1 {
		t.Errorf("Expected ClearInput to be called once, got %d", mockUI.ClearInputCalls)
	}

	if mockUI.RefreshListCalls != 1 {
		t.Errorf("Expected RefreshList to be called once, got %d", mockUI.RefreshListCalls)
	}

	if mockUI.FocusListCalls != 1 {
		t.Errorf("Expected FocusList to be called once, got %d", mockUI.FocusListCalls)
	}

	if len(mockUI.TasksReceived) != 1 {
		t.Errorf("Expected 1 task in UI, got %d", len(mockUI.TasksReceived))
	}
}
