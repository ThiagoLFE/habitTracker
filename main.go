package main

import (
	log "log/slog"
	"net/http"
)

type Habit struct {
	ID     int
	Title  string
	Streak int
}

type HabitList struct {
	Items  []Habit
	NextID int
}

func createHabitList() HabitList {
	return HabitList{Items: []Habit{}, NextID: 1}
}

var myHabits HabitList = createHabitList()

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", loadPage)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Info("Server open on the http://localhost:8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Error("server failed", "error", err)
	}
}

func loadPage(w http.ResponseWriter, r *http.Request) {

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
