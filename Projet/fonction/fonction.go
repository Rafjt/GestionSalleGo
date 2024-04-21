package fonction

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUser     = "root"
	dbPassword = "1234.Azerty"
	dbName     = "projetGolang"
)

type Event struct {
	ID          int
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

func AddEvent(db *sql.DB, event Event) error {
	query := "INSERT INTO events (title, date, location, category, description) VALUES (?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(event.Title, event.Date, event.Location, event.Category, event.Description)
	if err != nil {
		return err
	}

	return nil
}

func CreationEvenement(db *sql.DB) error {
	fmt.Println("----------------------------------------Créer son évènement---------------------------------------")

	var title, date, timeInput, location, category, description string

	fmt.Print("Entrez le titre de l'événement: ")
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
		return fmt.Errorf("Erreur lors de l'analyse du temps: %v", err)
	}

	fmt.Print("Entrez le lieu: ")
	fmt.Scanln(&location)

	for {
		fmt.Print("Choisissez une catégorie (Professionnel, Personnel, Loisir): ")
		fmt.Scanln(&category)

		if category == "Professionnel" || category == "Personnel" || category == "Loisir" {
			break
		}

		fmt.Println("Choix invalide.")
	}

	fmt.Print("Entrez une brève description: ")
	fmt.Scanln(&description)

	fmt.Println("--------------------------------------------------------------------------------------------------")

	newEvent := Event{
		Title:       title,
		Date:        thenafter.Format("2006-01-02 15:04"),
		Location:    location,
		Category:    category,
		Description: description,
	}

	err = AddEvent(db, newEvent)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'ajout de l'événement: %v", err)
	}

	fmt.Println("Événement ajouté avec succès !")
	return nil
}

func VisualiserEvenements(db *sql.DB) error {
	fmt.Println("--------------------------------------Visualiser les événements--------------------------------------")

	rows, err := db.Query("SELECT * FROM events")
	if err != nil {
		return fmt.Errorf("Erreur lors de la récupération des événements: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-5s | %-20s | %-19s | %-30s | %-15s | %-50s\n", "ID", "Titre", "Date", "Lieu", "Catégorie", "Description")
	fmt.Println("--------------------------------------------------------------------------------------------------")

	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.Title, &event.Date, &event.Location, &event.Category, &event.Description); err != nil {
			return fmt.Errorf("Erreur lors de la lecture de la ligne d'événement: %v", err)
		}

		fmt.Printf("%-5d | %-20s | %-19s | %-30s | %-15s | %-50s\n", event.ID, event.Title, event.Date, event.Location, event.Category, event.Description)
	}

	fmt.Println("--------------------------------------------------------------------------------------------------")

	return nil
}

func SupprimerEvenement(db *sql.DB, eventID int) error {
	result, err := db.Exec("DELETE FROM events WHERE id = ?", eventID)
	if err != nil {
		return fmt.Errorf("Erreur lors de la suppression de l'événement: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Erreur lors de la récupération du nombre de lignes affectées: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Aucun événement trouvé avec l'ID %d", eventID)
	}

	fmt.Printf("Événement avec l'ID %d supprimé avec succès.\n", eventID)

	return nil
}

func AfficherEvenement(db *sql.DB, eventID int) error {
	fmt.Println("--------------------------------------------------------------------------------------------------")
	row := db.QueryRow("SELECT * FROM events WHERE id = ?", eventID)

	var event Event
	err := row.Scan(&event.ID, &event.Title, &event.Date, &event.Location, &event.Category, &event.Description)
	if err != nil {
		return fmt.Errorf("Erreur lors de la récupération de l'événement: %v", err)
	}

	fmt.Printf("ID: %d\nTitre: %s\nDate: %s\nLieu: %s\nCatégorie: %s\nDescription: %s\n",
		event.ID, event.Title, event.Date, event.Location, event.Category, event.Description)

	fmt.Println("--------------------------------------------------------------------------------------------------")

	return nil
}
