# Go Todo

A simple, terminal-based todo application built in Go with a clean TUI (Terminal User Interface) and SQLite persistence.

## Features

- **Terminal User Interface**: Clean, keyboard-driven interface using `tview`
- **Persistent Storage**: SQLite database for task persistence
- **Task Management**: Add, toggle completion, and delete tasks
- **Real-time Updates**: Immediate UI updates with database synchronization
- **Logging**: Comprehensive logging for debugging and monitoring

## Prerequisites

- Go 1.23.4 or later
- SQLite3 (included via go-sqlite3 driver)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/FinnBuick/go-todo
cd go-todo
```

2. Install dependencies:
```bash
go mod download
```

3. Build the application:
```bash
go build -o go-todo
```

## Usage

### Running the Application

```bash
./go-todo
```

### Controls

- **Tab**: Cycle focus between input field and task list
- **Enter** (in input field): Add new task
- **Enter** (in task list): Toggle task completion status
- **d** (in task list): Delete selected task
- **Esc** (in input field): Focus back to task list
- **q**: Quit application

### Interface Layout

The application features a three-panel layout:
- **Task List**: Displays all tasks with completion status
- **Input Field**: For adding new tasks
- **Help Panel**: Shows available controls and shortcuts

## Architecture

The application follows a clean architecture pattern with separation of concerns:

```
go-todo/
├── main.go              # Application entry point
├── internal/
│   ├── models/          # Data models
│   │   └── task.go      # Task struct and methods
│   ├── storage/         # Database layer
│   │   └── sqlite.go    # SQLite implementation
│   ├── controller/      # Business logic
│   │   ├── app.go       # Main controller
│   │   └── app_test.go  # Controller tests
│   └── ui/              # User interface
│       └── tui.go       # Terminal UI implementation
├── go.mod               # Go module definition
└── tasks.db             # SQLite database (created on first run)
```

### Components

- **Models**: Define the core data structures (`Task`)
- **Storage**: Handle database operations with SQLite
- **Controller**: Manage application logic and coordinate between UI and storage
- **UI**: Provide terminal-based user interface using `tview`

## Dependencies

- [`github.com/rivo/tview`](https://github.com/rivo/tview) - Terminal UI framework
- [`github.com/gdamore/tcell/v2`](https://github.com/gdamore/tcell) - Terminal handling
- [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) - SQLite driver

## Database Schema

The application uses a simple SQLite schema:

```sql
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	description TEXT NOT NULL,
	done INTEGER DEFAULT 0 CHECK(done in (0,1)),
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);
```

## Logging

Application logs are written to `todo_app.log` for debugging purposes. The log includes:
- Application startup/shutdown events
- Database operations
- UI interactions
- Error conditions

## Development

### Running Tests

```bash
go test ./...
```

### Building for Production

```bash
go build -ldflags "-s -w" -o go-todo
```
