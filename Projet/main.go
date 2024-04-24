package main

import (
	"database/sql"
	"fmt"
	"gestionProjetGolang/fonction"
	"os"
)

// func de gestion des erreurs
func errors(strer string, err error) {
	if err != nil {
		fmt.Println(strer, err)
	}
}

// variable globale pour la connexion à la base de données
var db *sql.DB

// fonction pour l'affichage principal
func affichagePrincipal() int {
	var inputUtilisateur int
	fmt.Fprintf(os.Stdout, "\033[0;31m%s\033[0m\n", "Bienvenue dans le Système de gestion de plannings") // ligne de bienvenue en rouge
	fmt.Println("-----------------------------------------------------")
	fmt.Println("1. Créer une réservations")
	fmt.Println("2. Visualiser les salle disponibles")
	fmt.Println("3. Supprimer une réservation")
	fmt.Println("4. Visualiser les réservations")
	fmt.Println("5. Modifier une réservation")
	fmt.Println("6. Quitter")
	fmt.Println("Choisissez une option :")
	fmt.Scanln(&inputUtilisateur)
	return inputUtilisateur
}

// fonction pour le choix de l'utilisateur
func choixUtilisateur(choix int) {
	switch choix {
	case 1:
		err := fonction.CreationReservations(db)
		errors("Erreur lors de la création de l'événement:", err)
	case 2:
		fonction.VisualiserSalles(db)
	case 3:
		var LocationID int
		fmt.Print("Entrez l'ID de l'événement à supprimer: ")
		fmt.Scanln(&LocationID)
		err := fonction.SupprimerReservation(db, LocationID)
		errors("", err)
	case 4:
		err := fonction.VisualiserReservations(db)
		errors("", err)
	case 5:
		var eventID int
		fmt.Print("Entrez l'id de l'événement à modifier: ")
		fmt.Scanln(&eventID)
		err := fonction.ModifierReservation(db, eventID)
		errors("", err)
		if err != nil {
			fmt.Println(err)
		}
	case 6:
		fmt.Println("Merci de votre visite !")
		os.Exit(0)
	default:
		fmt.Println("Choix non valide. Veuillez choisir une option de 1 à 6.")
	}

}

// fonction main
func main() {
	var err error
	db, err = fonction.ConnexionBdd()
	errors("Erreur lors de la connexion à la base de données:", err)
	defer db.Close()
	for {
		choix := affichagePrincipal()
		choixUtilisateur(choix)
	}
}
