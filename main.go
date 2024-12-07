package main

import (
	"fmt"
	"html/template"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

type PageData struct {
	Title    string
	Numbers  []Number
	Complete bool
}

type Number struct {
	Value    int
	Letter   string
	DrawTime string
}

var data PageData
var prevNums []int

func main() {
	// Initialize data
	data = PageData{
		Title:   "Bingo",
		Numbers: make([]Number, 0, 75),
	}
	prevNums = make([]int, 0, 75)

	// Set up routes
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/draw", handleDraw)
	http.HandleFunc("/reset", handleReset)

	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "bingo", data)
}

func handleDraw(w http.ResponseWriter, r *http.Request) {
	if len(prevNums) >= 75 {
		data.Complete = true
	} else {
		newNum := getNextNumber()
		data.Numbers = append(data.Numbers, Number{
			Value:    newNum,
			Letter:   getBingoLetter(newNum),
			DrawTime: time.Now().Format("15:04:05"),
		})
	}

	// Use temporary redirect instead of see other
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func handleReset(w http.ResponseWriter, r *http.Request) {
	data.Numbers = make([]Number, 0, 75)
	data.Complete = false
	prevNums = make([]int, 0, 75)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getNextNumber() int {
	var n int
	for {
		n = rand.IntN(75) + 1
		if !contains(prevNums, n) {
			break
		}
	}
	prevNums = append(prevNums, n)
	return n
}

func contains(arr []int, n int) bool {
	for _, v := range arr {
		if v == n {
			return true
		}
	}
	return false
}

func getBingoLetter(n int) string {
	switch {
	case n <= 15:
		return "B"
	case n <= 30:
		return "I"
	case n <= 45:
		return "N"
	case n <= 60:
		return "G"
	default:
		return "O"
	}
}
