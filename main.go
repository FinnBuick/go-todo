package main

import (
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("todo_app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Application starting...")

	// 1. Init Database Store
	store, err := NewStore()
	if err != nil {
		log.Fatalf("Failed to initialise data store: %v", err)
	}
	defer store.Close()
	log.Println("Database store initialised.")

	// 2. Initialise Controller
	controller := NewAppController(store)

	// 3. Initialise UI
	ui := NewUI(controller)
	log.Println("UI initialised.")

	// 4. Set UI for the controller
	controller.SetUI(ui)
	log.Println("UI set for controller.")

	// 5. Start the application via the controller
	log.Println("Starting application controller...")
	if err := controller.Start(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}

	log.Println("Application stopped.")
}
