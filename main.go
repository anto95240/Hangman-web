package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"piscine/piscine"
	"text/template"
)

type Test struct {
	Att  int
	Word string
	Jose string
	Rep  []string
	Win  []piscine.Score
}

const port = ":5520"

var attempt int
var UdScore []rune
var pick string
var boolean = true
var rep []string
var Name string
var level string
var winners []piscine.Score

type Score struct {
	Name   string
	Points int
}

// Redirection vers la page cible lorsque jeu terminé.
func Redirect(w http.ResponseWriter, r *http.Request) {

	if boolean {
		boolean = false
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		new := Test{Att: attempt, Word: string(UdScore), Jose: piscine.Check(attempt), Rep: rep}
		tmpl := template.Must(template.ParseFiles("./pageshtml/game.html"))
		tmpl.Execute(w, new)
	}
}

func main() {
	http.HandleFunc("/", Redirect)

	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pageshtml/index.html")
	})

	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) {
		boolean = true

		if level == "EASY" {
			winners = append(winners, piscine.Score{Name: Name, Points: attempt})
		} else if level == "NORMAL" {
			winners = append(winners, piscine.Score{Name: Name, Points: attempt * 2})
		} else {
			winners = append(winners, piscine.Score{Name: Name, Points: attempt * 3})
		}
		winners = piscine.ScoreJoueur(winners)
		new := Test{Win: winners}
		tmpl := template.Must(template.ParseFiles("./pageshtml/win.html"))
		tmpl.Execute(w, new)
	})

	//Redict for loose
	http.HandleFunc("/loose", func(w http.ResponseWriter, r *http.Request) {
		boolean = true
		new := Test{Word: pick}
		tmpl := template.Must(template.ParseFiles("./pageshtml/loose.html"))
		tmpl.Execute(w, new)
	})

	//Redict for game.html
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		letter := piscine.ToUpper(r.Form.Get("field2"))
		deja := true
		for i := range rep {
			if rep[i] == letter {
				deja = false
			}
		}

		udd := attempt
		UdScore, attempt = piscine.Compare(UdScore, attempt, pick, letter)
		if deja && udd != attempt {
			rep = append(rep, letter)
		}

		if attempt <= 0 && r.Method == "POST" {
			boolean = true
			http.Redirect(w, r, "/loose", http.StatusSeeOther)
		}
		if string(UdScore) == pick && r.Method == "POST" {
			boolean = true
			http.Redirect(w, r, "/win", http.StatusSeeOther)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)

	})

	http.HandleFunc("/hangman", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./pageshtml/hangman.html")
	})

	http.HandleFunc("/choix", func(w http.ResponseWriter, r *http.Request) {
		UdScore = []rune{}
		if len(os.Args) == 1 {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			rep = []string{}
			piscine.Repetition = []string{}
			level = r.Form.Get("w")
			pick = piscine.ToUpper(piscine.Random(level))
			Name = r.Form.Get("nom_utilisateur")
			attempt = 10

			for range pick {
				UdScore = append(UdScore, '_')
			}

			for v := 0; v < len(pick)/2-1; v++ {
				random := rand.Intn(len(pick))
				if UdScore[random] == '_' {
					UdScore[random] = rune(pick[random])
				} else {
					v--
				}
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./pageshtml/game.html")
	// })

	// Serve asset files (CSS and JS)
	http.Handle("/asset/", http.StripPrefix("/asset/", http.FileServer(http.Dir("asset"))))

	// Start server
	fmt.Printf("\nServer started at http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
