package main

import (
	"fmt"
	"os"
	"gestionProjetGolang/fonction"
	"database/sql"
)

var db *sql.DB

func affichagePrincipal() int {
	var inputUtilisateur int
	fmt.Fprintf(os.Stdout, "\033[0;31m%s\033[0m\n", "Bienvenue dans le Système de gestion de plannings") // ligne de bienvenue en rouge
	fmt.Println("-----------------------------------------------------")
	fmt.Println("1. Créer un nouvel événement")
	fmt.Println("2. Visualiser les événements")
	fmt.Println("3. Modifier un événement")
	fmt.Println("4. Supprimer un événement")
	fmt.Println("5. Rechercher un événement")
	fmt.Println("6. Quitter")
	fmt.Println("Choisissez une option :")
	fmt.Scanln(&inputUtilisateur)
	return inputUtilisateur
}

func choixUtilisateur(choix int) {
	switch choix {
	case 1:
		err := fonction.CreationEvenement(db)
		if err != nil {
			fmt.Println("Erreur lors de la création de l'événement:", err)
		}
	case 2:
		fonction.VisualiserEvenements(db)
	case 3:
		// modification à faire
	case 4:
		var eventID int
        fmt.Print("Entrez l'ID de l'événement à supprimer: ")
        fmt.Scanln(&eventID)
        err := fonction.SupprimerEvenement(db, eventID)
        if err != nil {
            fmt.Println(err)
        }
	case 5:
        var eventID int
        fmt.Print("Entrez le nom de l'événement à afficher: ")
        fmt.Scanln(&eventID)
        err := fonction.AfficherEvenement(db, eventID)
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
