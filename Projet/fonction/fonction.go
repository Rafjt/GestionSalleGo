package fonction

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Variables globales pour la connexion à la base de données
const (
	dbUser     = "root"
	dbPassword = "1234.Azerty"
	dbName     = "projetgolang"
)

// Structure pour les salles
type Room struct {
	ID       int
	Name     string
	Capacity int
	Dispo    int
}

// Structure pour les réservations
type Reservation struct {
	ID          int
	Title       string
	Date        string
	Location    string
	Category    string
	Description string
}

// Fonction de gestion des erreurs
func errors(strer string, err error) {
	if err != nil {
		fmt.Println(strer, err)
	}
}

func RechercheRoom(db *sql.DB, date time.Time) (string, error) {
	// On ajoute 1h30 à la date pour obtenir la date de fin du cours
	endDate := date.Add(time.Hour*1 + time.Minute*30)

	startDateStr := date.Format("2006-01-02 15:04")
	endDateStr := endDate.Format("2006-01-02 15:04")

	// LA query .
	rows, err := db.Query("SELECT id FROM salles WHERE id NOT IN (SELECT salle_id FROM reservations WHERE date >= ? AND date <= ?)", startDateStr, endDateStr)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var availableRooms []string
	for rows.Next() {
		var roomID int
		if err := rows.Scan(&roomID); err != nil {
			return "", err
		}
		availableRooms = append(availableRooms, strconv.Itoa(roomID))
	}

	return strings.Join(availableRooms, ", "), nil
}

// Fonction de connexion à la base de données
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

	fmt.Fprintf(os.Stdout, "\033[1;32m%s\033[0m\n", "Connecté à la base de données MySQL!") // ligne de connexion bdd en vert

	return db, nil
}

// Fonction pour changer la disponibilité de la salle
func ChangeEtatSalle(db *sql.DB, salleID int) error {
	_, err := db.Exec("UPDATE salles SET dispo = 0 WHERE id = ?", salleID)
	errors("Erreur lors de la mise à jour de la salle: %v", err)

	return nil
}

// Fonction pour ajouter un événement
func AjoutEvenement(db *sql.DB, Reservation Reservation) error {
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

// Fonction pour mettre à jour un événement
func MajEvenement(db *sql.DB, Reservation Reservation, locationID int) error {

	query := "UPDATE reservations SET name=?, date=?, salle_id=?, cours=?, description=? WHERE id =" + strconv.Itoa(locationID) + ""
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(Reservation.Title, Reservation.Date, Reservation.Location, Reservation.Category, Reservation.Description)
	if err != nil {
		return err
	}

	ChangeEtatSalle(db, locationID)
	return nil
}

// Fonction pour créer une réservation
func CreationReservations(db *sql.DB) error {
	fmt.Println("----------------------------------------Créer son évènement---------------------------------------")

	var title, date, timeInput, location, category, description string

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Entrez le titre de la réservation: ")
	title, _ = reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Entrez la date (YYYY-MM-DD): ")
	date, _ = reader.ReadString('\n')
	date = strings.TrimSpace(date)

	fmt.Print("Entrez l'heure (HH:mm): ")
	timeInput, _ = reader.ReadString('\n')
	timeInput = strings.TrimSpace(timeInput)

	// pour éviter l'erreur de temps/date
	dateTimeInput := fmt.Sprintf("%s %s", date, timeInput)

	layout := "2006-01-02 15:04"
	thenafter, err := time.Parse(layout, dateTimeInput)
	errors("Erreur lors de l'analyse du temps (Veuillez entrer une date dans le format 'YYYY-MM-DD' ): %v", err)

	if thenafter.Before(time.Now()) {
		return fmt.Errorf("La date doit être dans le futur.")
	}

	// pour éviter les réservations après 20h45 et avant 8h
	hour, _, _ := thenafter.Clock()
	if hour < 8 || hour > 20 || (hour == 20 && thenafter.Minute() > 45) {
		return fmt.Errorf("L'heure doit être entre 8AM et 8:45PM.")
	}

	// Query the database to find available rooms at the given date and time
	availableRooms, err := RechercheRoom(db, thenafter)
	errors("Erreur lors de la récupération des salles disponibles: %v", err)

	fmt.Println("Salles disponibles: ", availableRooms)

	fmt.Print("Entrez la salle que vous souhaiter réserver: ")
	fmt.Scanln(&location)
	if !strings.Contains(availableRooms, location) {
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

	err = AjoutEvenement(db, newReservation)
	errors("Erreur lors de l'ajout de l'événement: %v", err)

	locationID, err := strconv.Atoi(location)
	errors("Erreur lors de la conversion de la salle en nombre: %v", err)

	err = ChangeEtatSalle(db, locationID)
	errors("Erreur lors de la mise à jour de la salle: %v", err)

	fmt.Println("Événement ajouté avec succès !")
	return nil
}

// Fonction pour visualiser les salles
func VisualiserSalles(db *sql.DB) error {
	fmt.Println("--------------------------------------Visualiser les Salles--------------------------------------")

	rows, err := db.Query("SELECT * FROM salles")
	errors("Erreur lors de la récupération des événements: %v", err)
	defer rows.Close()

	fmt.Printf("|%-5s | %-20s |\n", "numéro", "capacité")
	fmt.Println("-----------------------------------------------------")

	for rows.Next() {
		var Salle Room
		if err := rows.Scan(&Salle.ID, &Salle.Name, &Salle.Capacity, &Salle.Dispo); err != nil {
			return fmt.Errorf("Erreur lors de la lecture de la ligne d'événement: %v", err)
		}
		fmt.Printf("|%-5d | %-13d places |\n", Salle.ID, Salle.Capacity)
	}

	fmt.Println("------------------------------------------------------------------------------------------------")

	return nil
}

// Fonction pour supprimer une réservation
func SupprimerReservation(db *sql.DB, locationID int) error {
	result, err := db.Exec("DELETE FROM reservations WHERE id = ?", locationID)
	errors("Erreur lors de la suppression de l'événement: %v", err)

	rowsAffected, err := result.RowsAffected()
	errors("Erreur lors de la récupération du nombre de lignes affectées: %v", err)

	if rowsAffected == 0 {
		return fmt.Errorf("Aucun événement trouvé avec l'ID %d", locationID)
	}

	fmt.Printf("Événement avec l'ID %d supprimé avec succès.\n", locationID)
	ChangeEtatSalle(db, locationID)

	return nil
}

// Fonction pour visualiser les réservations
func VisualiserReservations(db *sql.DB) error {
	fmt.Println("--------------------------------------Visualiser les réservations-----------------------------------------------------------------------------------")

	rows, err := db.Query("SELECT id, name, date, cours, description, salle_id FROM reservations")
	errors("Erreur lors de la récupération des réservations: %v", err)
	defer rows.Close()

	fmt.Printf("|%-20s |%-20s | %-19s | %-20s | %-19s | %-19s |\n", "id", "nom", "date", "cours", "n de salle", "Description")
	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")

	for rows.Next() {
		var reservation Reservation
		if err := rows.Scan(&reservation.ID, &reservation.Title, &reservation.Date, &reservation.Category, &reservation.Description, &reservation.Location); err != nil {
			return fmt.Errorf("Erreur lors de la lecture de la ligne de réservation: %v", err)
		}
		fmt.Printf("|%-20d |%-20s | %-19s | %-20s | %-19s | %-19s |\n", reservation.ID, reservation.Title, reservation.Date, reservation.Category, reservation.Location, reservation.Description)
	}

	fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")

	return nil
}

// Fonction pour modifier une réservation
func ModifierReservation(db *sql.DB, locationID int) error {
	fmt.Println("--------------------------------------Modifier une réservation--------------------------------------")

	var title, date, timeInput, location, category, description string

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Entrez le titre de la réservation: ")
	title, _ = reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Entrez la date (YYYY-MM-DD): ")
	date, _ = reader.ReadString('\n')
	date = strings.TrimSpace(date)

	fmt.Print("Entrez l'heure (HH:mm): ")
	timeInput, _ = reader.ReadString('\n')
	timeInput = strings.TrimSpace(timeInput)

	// pour éviter l'erreur de temps/date
	dateTimeInput := fmt.Sprintf("%s %s", date, timeInput)

	layout := "2006-01-02 15:04"
	thenafter, err := time.Parse(layout, dateTimeInput)
	errors("Erreur lors de l'analyse du temps (Veuillez entrer une date dans le format 'YYYY-MM-DD' ): %v", err)

	if thenafter.Before(time.Now()) {
		return fmt.Errorf("La date doit être dans le futur.")
	}

	// pour éviter les réservations après 20h45 et avant 8h
	hour, _, _ := thenafter.Clock()
	if hour < 8 || hour > 20 || (hour == 20 && thenafter.Minute() > 45) {
		return fmt.Errorf("L'heure doit être entre 8AM et 8:45PM.")
	}

	// Query the database to find available rooms at the given date and time
	availableRooms, err := RechercheRoom(db, thenafter)
	errors("Erreur lors de la récupération des salles disponibles: %v", err)

	fmt.Println("Salles disponibles: ", availableRooms)

	fmt.Print("Entrez la salle que vous souhaiter réserver: ")
	fmt.Scanln(&location)
	if !strings.Contains(availableRooms, location) {
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

	err = MajEvenement(db, newReservation, locationID)
	errors("Erreur lors de la mise à jour de l'événement: %v", err)

	return nil
}
