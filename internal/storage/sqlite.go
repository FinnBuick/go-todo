package storage

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"go-todo/internal/models"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

const dbFileName = "tasks.db"

type Store struct {
	db *sql.DB
}

func NewStore() (*Store, error) {
	dbPath := filepath.Join(".", dbFileName)

	d, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := d.Ping(); err != nil {
		d.Close()
		return nil, fmt.Errorf("failed to connnect to database: %w", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		done INTEGER DEFAULT 0 CHECK(done in (0,1))
	);`
	if _, err = d.Exec(createTableSQL); err != nil {
		d.Close()
		return nil, fmt.Errorf("failed to create tasks table: %w", err)
	}

	return &Store{db: d}, nil
}

func (s *Store) Close() {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			log.Printf("Error closing databse: %v\n", err)
		}
	}
}

func (s *Store) GetTasks() ([]models.Task, error) {
	rows, err := s.db.Query("SELECT id, description, done FROM tasks ORDER BY id ASC")
	if err != nil {
		return nil, fmt.Errorf("querying tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		var doneInt int
		if err := rows.Scan(&t.ID, &t.Description, &doneInt); err != nil {
			return nil, fmt.Errorf("scanning task row: %w", err)
		}
		t.Done = (doneInt == 1)
		tasks = append(tasks, t)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("after iteration task rows: %w", err)
	}
	return tasks, nil
}

func (s *Store) AddTask(description string) (int64, error) {
	res, err := s.db.Exec("INSERT INTO tasks (description) VALUES (?)", description)
	if err != nil {
		return 0, fmt.Errorf("inserting task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert id: %w", err)
	}
	return id, nil
}

func (s *Store) ToggleTaskStatus(id int) error {
	var currentStatus bool
	err := s.db.QueryRow("SELECT done FROM tasks WHERE id = ?", id).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task with ID %d not found", id)
		}
		return fmt.Errorf("querying task status for toggle: %w", err)
	}

	_, err = s.db.Exec("UPDATE tasks SET done = ? WHERE id = ?", !currentStatus, id)
	if err != nil {
		return fmt.Errorf("updating task status: %w", err)
	}
	return nil
}

func (s *Store) DeleteTask(id int) error {
	res, err := s.db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected by delete: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found for deletion", id)
	}
	return nil
}
