package main

import (
	"database/sql"
	"fmt"
	"gestionProjetGolang/fonction"
	"os"
)

var db *sql.DB

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

func choixUtilisateur(choix int) {
	switch choix {
	case 1:
		err := fonction.CreationReservations(db)
		if err != nil {
			fmt.Println("Erreur lors de la création de l'événement:", err)
		}
	case 2:
		fonction.VisualiserSalles(db)
	case 3:
		var LocationID int
		fmt.Print("Entrez l'ID de l'événement à supprimer: ")
		fmt.Scanln(&LocationID)
		err := fonction.SupprimerReservation(db, LocationID)
		if err != nil {
			fmt.Println(err)
		}
	case 4:
		err := fonction.VisualiserReservations(db)
		if err != nil {
			fmt.Println(err)
		}
	// case 5:
	// 	var eventID int
	// 	fmt.Print("Entrez le nom de l'événement à afficher: ")
	// 	fmt.Scanln(&eventID)
	// 	err := fonction.AfficherEvenement(db, eventID)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	case 6:
		fmt.Println("Merci de votre visite !")
		os.Exit(0)
	default:
		fmt.Println("Choix non valide. Veuillez choisir une option de 1 à 6.")
	}

}

func main() {
	var err error
	db, err = fonction.ConnexionBdd()
	if err != nil {
		fmt.Println("Erreur lors de la connexion à la base de données:", err)
		return
	}
	defer db.Close()
	for {
		choix := affichagePrincipal()
		choixUtilisateur(choix)
	}
}
