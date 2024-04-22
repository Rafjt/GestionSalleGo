CREATE DATABASE IF NOT EXISTS projetGolang;
use projetGolang;

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