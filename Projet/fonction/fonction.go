package fonction

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser     = "root"
	dbPassword = "1234.Azerty"
	dbName     = "projetGolang"
)

type Room struct {
	ID       int
	Name     string
	Capacity int
	Dispo    int
}

type Reservation struct {
	Title       string
	Date        string
	Location    string
	Category    string
	Description string
}

func ConnexionBdd() (*sql.DB, error) {
	dbURI := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", dbUser, dbPassword, dbName)

	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connecté à la base de données MySQL!")

	return db, nil
}

func MakeSalleIndispobible(db *sql.DB, salleID int) error {
	_, err := db.Exec("UPDATE salles SET dispo = 0 WHERE id = ?", salleID)
	if err != nil {
		return fmt.Errorf("Erreur lors de la mise à jour de la salle: %v", err)
	}

	return nil
}

func AddEvent(db *sql.DB, Reservation Reservation) error {
	query := "INSERT INTO reservations (name, date, salle_id, cours, description) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(Reservation.Title, Reservation.Date, Reservation.Location, Reservation.Category, Reservation.Description)
	if err != nil {
		return err
	}

	return nil
}

func CreationReservations(db *sql.DB) error {
	fmt.Println("----------------------------------------Créer son évènement---------------------------------------")

	var title, date, timeInput, location, category, description string

	fmt.Print("Entrez le titre de la réservation: ")
	fmt.Scanln(&title)

	fmt.Print("Entrez la date (YYYY-MM-DD): ")
	fmt.Scanln(&date)

	fmt.Print("Entrez l'heure (HH:mm): ")
	fmt.Scanln(&timeInput)

	// pour éviter l'erreur de temps/date
	dateTimeInput := fmt.Sprintf("%s %s", date, timeInput)

	layout := "2006-01-02 15:04"
	thenafter, err := time.Parse(layout, dateTimeInput)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'analyse du temps (Veuillez entrer une date dans le format 'YYYY-MM-DD' ): %v", err)

	}

	fmt.Print("Entrez la salle que vous souhaiter réserver (1-10): ")
	fmt.Scanln(&location)
	if location != "1" && location != "2" && location != "3" && location != "4" && location != "5" && location != "6" && location != "7" && location != "8" && location != "9" && location != "10" {
		return fmt.Errorf("Salle invalide.")
	}

	for {
		fmt.Print("Choisissez la matière enseignée (LangageC, Go, Python): ")
		fmt.Scanln(&category)

		if category == "LangageC" || category == "Go" || category == "Python" {
			break
		}

		fmt.Println("Choix invalide.")
	}

	fmt.Print("Entrez une brève description: ")
	fmt.Scanln(&description)

	fmt.Println("--------------------------------------------------------------------------------------------------")

	newReservation := Reservation{
		Title:       title,
		Date:        thenafter.Format("2006-01-02 15:04"),
		Location:    location,
		Category:    category,
		Description: description,
	}

	err = AddEvent(db, newReservation)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'ajout de l'événement: %v", err)
	}

	// query := "SELECT id FROM reservations WHERE name = ? AND date = ? AND salle_id = ? AND cours = ? AND description = ?"
	// row := db.QueryRow(query, title, thenafter.Format("2006-01-02 15:04"), location, category, description)

	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	return fmt.Errorf("Erreur lors de la récupération de l'ID de la réservation: %v", err)
	// }

	locationID, err := strconv.Atoi(location)
	if err != nil {
		return fmt.Errorf("Erreur lors de la conversion de la salle en nombre: %v", err)
	}

	// Pass the integer ID to the function
	err = MakeSalleIndispobible(db, locationID)
	if err != nil {
		return fmt.Errorf("Erreur lors de la mise à jour de la salle: %v", err)
	}

	fmt.Println("Événement ajouté avec succès !")
	return nil
}

func VisualiserSalles(db *sql.DB) error {
	fmt.Println("--------------------------------------Visualiser les Salles--------------------------------------")

	rows, err := db.Query("SELECT * FROM salles WHERE dispo = 1")
	if err != nil {
		return fmt.Errorf("Erreur lors de la récupération des événements: %v", err)
	}
	defer rows.Close()

	fmt.Printf("|%-5s | %-20s | %-19s |\n", "numéro", "capacité", "disponibilité")
	fmt.Println("-----------------------------------------------------")

	for rows.Next() {
		var Salle Room
		if err := rows.Scan(&Salle.ID, &Salle.Name, &Salle.Capacity, &Salle.Dispo); err != nil {
			return fmt.Errorf("Erreur lors de la lecture de la ligne d'événement: %v", err)
		}
		fmt.Printf("|%-5d | %-20d | %-19d |\n", Salle.ID, Salle.Capacity, Salle.Dispo)
	}

	fmt.Println("-----------------------------------------------------")

	return nil
}

func SupprimerReservation(db *sql.DB, locationID int) error {
	result, err := db.Exec("DELETE FROM reservations WHERE id = ?", locationID)
	if err != nil {
		return fmt.Errorf("Erreur lors de la suppression de l'événement: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Erreur lors de la récupération du nombre de lignes affectées: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Aucun événement trouvé avec l'ID %d", locationID)
	}

	fmt.Printf("Événement avec l'ID %d supprimé avec succès.\n", locationID)

	return nil
}

func VisualiserReservations(db *sql.DB) error {
	fmt.Println("--------------------------------------Visualiser les réservations--------------------------------------")

	rows, err := db.Query("SELECT name, date, cours, description, salle_id FROM reservations")
	if err != nil {
		return fmt.Errorf("Erreur lors de la récupération des réservations: %v", err)
	}
	defer rows.Close()

	fmt.Printf("|%-20s | %-19s | %-20s | %-19s | %-19s |\n", "name", "date", "cours", "n de salle", "Description")
	fmt.Println("-------------------------------------------------------------------------------------------------------------------")

	for rows.Next() {
		var reservation Reservation
		if err := rows.Scan(&reservation.Title, &reservation.Date, &reservation.Category, &reservation.Description, &reservation.Location); err != nil {
			return fmt.Errorf("Erreur lors de la lecture de la ligne de réservation: %v", err)
		}
		fmt.Printf("|%-20s | %-19s | %-20s | %-19s | %-19s |\n", reservation.Title, reservation.Date, reservation.Category, reservation.Location, reservation.Description)
	}

	fmt.Println("-------------------------------------------------------------------------------------------------------------------")

	return nil
}

// func ModifierReservation(db *sql.DB, locationID int) error {
// 	fmt.Println("--------------------------------------Modifier une réservation--------------------------------------")

// 	var title, date, timeInput, location, category, description string

// 	fmt.Print("Entrez le titre de la réservation: ")
// 	fmt.Scanln(&title)

// 	fmt.Print("Entrez la date (YYYY-MM-DD): ")
// 	fmt.Scanln(&date)

// 	fmt.Print("Entrez l'heure (HH:mm): ")
// 	fmt.Scanln(&timeInput)

// 	// pour éviter l'erreur de temps/date
// 	dateTimeInput := fmt.Sprintf("%s %s", date, timeInput)

// 	layout := "2006-01-02 15:04"
// 	thenafter, err := time.Parse(layout, dateTimeInput)
// 	if err != nil {
// 		return fmt.Errorf("Erreur lors de l'analyse du temps (Veuillez entrer une date dans le format 'YYYY-MM-DD' ): %v", err)

// 	}

// 	fmt.Print("Entrez la salle que vous souhaiter réserver (1-10): ")
// 	fmt.Scanln(&location)
// 	if location != "1" && location != "2" && location != "3" && location != "4" && location != "5" && location != "6" && location != "7" && location != "8" && location != "9" && location != "10" {
// 		return fmt.Errorf("Salle invalide.")
// 	}

// 	for {
// 		fmt.Print("Choisissez la matière enseignée (LangageC, Go, Python): ")
// 		fmt.Scanln(&category)

// 		if category == "LangageC" || category == "Go" || category == "Python" {
// 			break
// 		}

// 		fmt.Println("Choix invalide.")
// 	}

// 	fmt.Print("Entrez une brève description: ")
// 	fmt.Scanln(&description)

// 	fmt.Println("--------------------------------------------------------------------------------------------------")

// 	newReservation := Reservation{
// 		Title:       title,
// 		Date:        thenafter.Format("2006-01-02 15:04"),
// 		Location:    location,
// 		Category:    category,
// 		Description: description,
// 	}

// 	err = UpdateEvent(db, newReservation, locationID)
// 	if err != nil {
// 		return fmt.Errorf("Erreur lors de la mise à jour de l'événement: %v", err)
// 	}

// 	locationID, err := strconv.Atoi(location)
// 	if err != nil {
// 		return fmt.Errorf("Erreur lors de la conversion de la salle en nombre: %v", err)
// 	}
