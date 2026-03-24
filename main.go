package main

import (
	"errors"
	log "log/slog"
	"net/http"
	"strings"
)

type Habit struct {
	ID    int
	Title string
}

type HabitList struct {
	Items  []Habit
	NextID int
}

var (
	errTitleRequired     = errors.New("Title required")
	errLengthOverflow120 = errors.New("Text is too large")
)

func formattedError(err error) string {
	switch err {
	case errTitleRequired:
		return "Title is required"
	case errLengthOverflow120:
		return "Max of 120 characters"
	default:
		return "Unexpected Error"
	}
}

func createHabitList() HabitList {
	return HabitList{Items: []Habit{}, NextID: 1}
}

func (hl *HabitList) add(title string) error {
	title = strings.TrimSpace(title)

	if len(title) == 0 {
		return errTitleRequired
	}

	if len(title) > 120 {
		return errLengthOverflow120
	}

	hl.Items = append(hl.Items, Habit{
		ID:    hl.NextID,
		Title: title,
	})

	hl.NextID++
	return nil
}

var myHabits HabitList = createHabitList()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", loadHomePage)
	mux.HandleFunc("/habits", loadCreateHabit)
	mux.HandleFunc("/habits/create", createHabit)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Info("Server open on the http://localhost:8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Error("server failed", "error", err)
	}
}

func createHabit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Form not found", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")

	if err := myHabits.add(title); err != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if renderErr := BaseLayout("habit", FormHabit(true, formattedError(err), title)).Render(r.Context(), w); renderErr != nil {
			http.Error(w, "Error to render page", http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func loadHomePage(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", `text/html; charset=utf-8`)

	if err := BaseLayout("Habit tracker", Home(myHabits.Items)).Render(r.Context(), w); err != nil {
		http.Error(w, "failed to render page", http.StatusInternalServerError)
		return
	}

}

func loadCreateHabit(w http.ResponseWriter, r *http.Request) {
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
