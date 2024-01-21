package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Aventurier représente la structure d'un aventurier
type Aventurier struct {
	ID          int    `json:"id"`
	Nom         string `json:"nom"`
	Classe      string `json:"classe"`
	Niveau      int    `json:"niveau"`
	PointVie    int    `json:"pointVie"`
	PointVieMax int    `json:"pointVieMax"`
	Defense     int    `json:"defense"`
	Attaque     int    `json:"attaque"`
	Vitesse     int    `json:"vitesse"`
	Avatar      string `json:"avatar"`
}

var idCounter int = 0
var aventuriers []Aventurier

func main() {
	loadAventuriersFromJSON()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/profil", profilHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/modify", modifyHandler)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/home.html", nil)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idParam := r.URL.Query().Get("id")
		idToCreate := parseInt(idParam)
		nom := r.FormValue("nom")
		classe := r.FormValue("classe")
		niveau := r.FormValue("niveau")

		niveauInt := parseInt(niveau)

		if nom != "" && classe != "" && niveauInt != -1 {
			for isIDUsed(idToCreate) {
				idToCreate++
			}
			var pointsVieMax int
			var défense int
			var attaque int
			var vitesse int
			var avatar string
			switch classe {
			case "guerrier":
				pointsVieMax = 100
				défense = 45
				attaque = 60
				vitesse = 30
				avatar = "guerrier.jpg"
			case "mage":
				pointsVieMax = 80
				défense = 40
				attaque = 65
				vitesse = 40
				avatar = "mage.jpg"
			case "archer":
				pointsVieMax = 90
				défense = 35
				attaque = 55
				vitesse = 50
				avatar = "archer.jpg"
			default:
				pointsVieMax = 0
				défense = 0
				attaque = 0
				vitesse = 0
				avatar = "avatar3.jpg"
			}
			aventurier := Aventurier{
				ID:          idToCreate,
				Nom:         nom,
				Classe:      classe,
				Niveau:      niveauInt,
				PointVie:    pointsVieMax,
				PointVieMax: pointsVieMax,
				Defense:     défense,
				Attaque:     attaque,
				Vitesse:     vitesse,
				Avatar:      avatar,
			}
			aventuriers = append(aventuriers, aventurier)
			saveAventuriersToJSON()
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else if r.Method == http.MethodGet {
		renderTemplate(w, "templates/create.html", nil)
	}
}

func profilHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "templates/profil.html", aventuriers)
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idParam := r.URL.Query().Get("id")
		idToDelete := parseInt(idParam)

		if idToDelete != -1 {
			index := findAventurierIndexByID(idToDelete)
			if index != -1 {
				aventuriers = append(aventuriers[:index], aventuriers[index+1:]...)
				saveAventuriersToJSON()
				idCounter -= 1

			}
		}
	}

	// Redirige vers la page de profil après la suppression
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}

func findAventurierIndexByID(id int) int {
	for i, aventurier := range aventuriers {
		if aventurier.ID == id {
			return i
		}
	}
	return -1
}
func modifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idParam := r.URL.Query().Get("id")
		idToModify := parseInt(idParam)

		if idToModify != -1 {
			index := findAventurierIndexByID(idToModify)
			if index != -1 && len(aventuriers) > 0 && index < len(aventuriers) {
				// Supprime l'aventurier existant
				aventuriers = append(aventuriers[:index], aventuriers[index+1:]...)
				saveAventuriersToJSON()
				idCounter -= 1

				// Redirige vers la page de création avec les informations pré-remplies
				http.Redirect(w, r, "/create", http.StatusSeeOther)
				return

			}
		}
	}

	// Redirige vers la page de profil si l'ID n'est pas valide
	http.Redirect(w, r, "/profil", http.StatusSeeOther)
}

func renderTemplate(w http.ResponseWriter, templateFile string, data interface{}) {
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loadAventuriersFromJSON() {
	file, err := ioutil.ReadFile("aventuriers.json")
	if err == nil {
		json.Unmarshal(file, &aventuriers)
	}
}

func saveAventuriersToJSON() {
	data, err := json.MarshalIndent(aventuriers, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = ioutil.WriteFile("aventuriers.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

func isIDUsed(id int) bool {
	return findAventurierIndexByID(id) != -1
}

func generateUniqueID() int {
	idCounter++
	return idCounter
}
