package main

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/matiasbellaud/hangman"
)

// initialisation de la structure des données
type HangmanStruct struct {
	Letter         string
	Word           string
	DisplayHangman string
	Essai          string
	TotalEssai     int
	ToFind         string
	ListLetter     []string
}

func main() {
	Play()

	// chemin faire les assets
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))

	//ouverture du port 4000
	log.Print("Listening on :4000...")
	err := http.ListenAndServe(":4000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Play() {

	//initialisation des variables
	difficulty := "2"
	language := "fr"
	var essai int
	var totalEssai int
	var listLetter []string
	toFind := hangman.RandomWord(difficulty, language)
	word := hangman.RevealRandomLetter(toFind)

	//fonction de la page home
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//initialisation des variables
		tmpl := template.Must(template.ParseFiles("template/home.html"))
		essai = 10
		totalEssai = 0
		listLetter = listLetter[:0]
		data := HangmanStruct{}

		// formulaire de la difficulté pour la validité des inputs
		if r.FormValue("chooseDifficulty") == "1" || r.FormValue("chooseDifficulty") == "2" || r.FormValue("chooseDifficulty") == "3" {
			difficulty = r.FormValue("chooseDifficulty")
		}
		// formulaire de la langue pour la validité des inputs
		if r.FormValue("chooseLanguage") == "fr" || r.FormValue("chooseLanguage") == "eng" {
			language = r.FormValue("chooseLanguage")
		}
		errHome := tmpl.Execute(w, data)
		if errHome != nil {
			log.Fatal(errHome)
		}
	})
	// fonction de la page de win
	http.HandleFunc("/win", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("template/win.html"))
		data := HangmanStruct{
			ToFind: toFind,
		}
		errWin := tmpl.Execute(w, data)
		if errWin != nil {
			log.Fatal(errWin)
		}
	})
	// fonction de la page de loose
	http.HandleFunc("/loose", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("template/loose.html"))
		data := HangmanStruct{
			ToFind: toFind,
		}
		errLoose := tmpl.Execute(w, data)
		if errLoose != nil {
			log.Fatal(errLoose)
		}
	})
	// fonction de la page hangman
	http.HandleFunc("/hangman", func(w http.ResponseWriter, r *http.Request) {
		//initialisation des variables
		if totalEssai == 0 {
			toFind = hangman.RandomWord(difficulty, language)
			word = hangman.RevealRandomLetter(toFind)
		}
		tmpl := template.Must(template.ParseFiles("template/hangman.html"))
		input := r.FormValue("chooseLetter")
		newWord, _, _ := hangman.PrintWord(toFind, word, listLetter, input)
		var essai_str string
		letterUsed := false
		validLetter, inputLetter := hangman.VerifInput(input, word)

		// pour empecher le compteur d'aller en négatif si on revient en arriere
		if essai <= 0 {
			essai = 0
		}

		// vérification d'input faux
		if !validLetter {
			if totalEssai != 0 {
				input = "lettre non valide"
			}
		} else if newWord == word {
			for _, letterInList := range listLetter {
				if input == letterInList {
					letterUsed = true
				}
			}
			// si lettre déjà utilisée
			if letterUsed {
				input = "lettre déjà utilisée"
			} else {
				// si lettre n'est pas dans le mot
				input = ""
				essai = essai - 1
				listLetter = append(listLetter, inputLetter)
			}
		} else {
			input = ""
			listLetter = append(listLetter, inputLetter)
		}

		word = newWord
		essai_str = strconv.Itoa(essai)

		// pour afficher les données
		data := HangmanStruct{
			Letter:     input,
			Word:       word,
			Essai:      essai_str,
			ToFind:     toFind,
			ListLetter: listLetter,
		}

		essai_str = strconv.Itoa(10 - essai)                                 // Modifier essai pour aller dans le sens des images
		data.DisplayHangman = "/assets/hangman/hangman" + essai_str + ".jpg" // displayHangman est égale aux images
		totalEssai++
		errHangman := tmpl.Execute(w, data)
		if errHangman != nil {
			log.Fatal(errHangman)
		}
	})
}
