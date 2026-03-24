package main

import "database/sql"

type Habit struct {
	ID        int64
	Title     string
	Done      bool
	CreatedAt string
}

func createHabit(db *sql.DB, title string) error {
	_, err := db.Exec(
		`INSERT INTO habits (title, done) VALUES (?, ?)`,
		title,
		false,
	)
	return err
}

func listHabits(db *sql.DB) ([]Habit, error) {
	rows, err := db.Query(`
		SELECT id, title, done, created_at
		FROM habits
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []Habit

	for rows.Next() {
		var h Habit
		if err := rows.Scan(&h.ID, &h.Title, &h.Done, &h.CreatedAt); err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return habits, nil
}
