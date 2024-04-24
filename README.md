# Gestion de Plannings

This is a Go project for managing room reservations. It uses a MySQL database to store room and reservation data.

## Project Structure

The project is divided into two main packages:

1. `main`: This package contains the main application logic and user interface.
2. `fonction`: This package contains functions for interacting with the database.

## Database Setup

The project uses a MySQL database with the following tables:

1. `salles`: This table stores information about the rooms.
2. `reservations`: This table stores information about the reservations.

You can set up the database by running the following SQL commands:

```sql
CREATE DATABASE IF NOT EXISTS projetGolang;
USE projetGolang;

CREATE TABLE IF NOT EXISTS salles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    capacity INT NOT NULL,
    dispo INT NOT NULL
);

CREATE TABLE IF NOT EXISTS reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name TEXT NOT NULL,
    date TIMESTAMP NOT NULL,
    cours TEXT NOT NULL,
    description TEXT NOT NULL,
    salle_id INT NOT NULL,
    FOREIGN KEY (salle_id) REFERENCES salles(id)
);

INSERT INTO salles (name, capacity, dispo) VALUES 
('Salle 1', 10, 1),
('Salle 2', 25, 1),
('Salle 3', 10, 1),
('Salle 4', 15, 1),
('Salle 5', 18, 1),
('Salle 6', 50, 1),
('Salle 7', 11, 1),
('Salle 8', 20, 1),
('Salle 9', 10, 1),
('Salle 10', 25, 1);
