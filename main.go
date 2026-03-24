package main

import (
	"database/sql"
	"errors"
	log "log/slog"
	"net/http"
)

type HabitList struct {
	Items  []Habit
	NextID int
}

type App struct {
	db *sql.DB
}

var (
	errTitleRequired     = errors.New("Title required")
	errLengthOverflow120 = errors.New("Text is too large")
)

func messageForError(err error) string {
	switch err {
	case errTitleRequired:
		return "Title is required"
	case errLengthOverflow120:
		return "Max of 120 characters"
	default:
		return "Unexpected Error"
	}
}

func main() {
	db, err := openDB()
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer db.Close()

	if err := initSchema(db); err != nil {
		log.Error(err.Error())
		return
	}

	app := &App{db: db}

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.loadHomePage)
	mux.HandleFunc("/habits", app.loadCreateHabit)
	mux.HandleFunc("/habits/create", app.createHabitHandle)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Info("Server open on the http://localhost:8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Error("server failed", "error", err)
	}
}

func (a *App) createHabitHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form not found", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	if err := createHabit(a.db, title); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if renderErr := BaseLayout("habit", FormHabit(true, messageForError(err), title)).Render(r.Context(), w); renderErr != nil {
			http.Error(w, "Error to render page", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) loadHomePage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	habits, err := listHabits(a.db)
	if err != nil {
		http.Error(w, "failed to load habits", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", `text/html; charset=utf-8`)

	if err := BaseLayout("Habit tracker", Home(habits)).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}

}

func (a *App) loadCreateHabit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", `text/html; charset=utf-8`)

	if err := BaseLayout("habit", FormHabit(true, "", "")).Render(r.Context(), w); err != nil {
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
		return
	}

}
