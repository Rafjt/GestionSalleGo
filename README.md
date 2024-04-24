# Gestion des Réservations

Il s'agit d'un projet Go pour la gestion des réservations de salles. Il utilise une base de données MySQL pour stocker les données des salles et des réservations.

## Structure du Projet

Le projet est divisé en deux packages principaux :

1. `main` : Ce package contient la logique principale de l'application et l'interface utilisateur.
2. `fonction` : Ce package contient les fonctions pour interagir avec la base de données.

## Configuration de la Base de Données

Le projet utilise une base de données MySQL avec les tables suivantes :

1. `salles` : Cette table stocke les informations sur les salles.
2. `reservations` : Cette table stocke les informations sur les réservations.

Vous pouvez configurer la base de données en exécutant les commandes SQL suivantes :

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
```

Cela crée la base de données et les tables nécessaires, puis insère des données de salle fictives pour démarrer.

## Fonctions Principales

Les principales fonctions de l'application sont les suivantes :

1. `affichagePrincipal()`: Cette fonction affiche le menu principal à l'utilisateur et retourne le choix de l'utilisateur sous forme d'entier.
2. `choixUtilisateur(choix int)`: Cette fonction prend le choix de l'utilisateur comme argument et effectue l'action correspondante. Les actions comprennent la création d'une réservation, la visualisation des salles disponibles, la suppression d'une réservation, la visualisation des réservations et la sortie de l'application.
3. `main()`: Il s'agit de la fonction principale de l'application. Elle établit une connexion à la base de données et affiche continuellement le menu principal à l'utilisateur jusqu'à ce que l'utilisateur choisisse de quitter.

## Exécution du Projet

Pour exécuter le projet, accédez au répertoire du projet et exécutez la commande `go run main.go`.
