package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// DB es la conexión global a la base de datos.
var DB *sql.DB

// InitDB conecta a la BD, la asigna a `DB` y crea las tablas si no existen.
func InitDB(connString string) error {
	var err error
	DB, err = sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	// Verificar conexión
	if err := DB.Ping(); err != nil {
		return err
	}

	// Crear tablas si no existen
	query := `
	CREATE TABLE IF NOT EXISTS stocks (
	    id SERIAL PRIMARY KEY,
	    ticker STRING NOT NULL UNIQUE,
	    company STRING NOT NULL,
	    target_from DECIMAL(10,2) NULL,
	    target_to DECIMAL(10,2) NULL,
	    rating_from STRING NULL,
	    rating_to STRING NULL,
	    time TIMESTAMPTZ NOT NULL DEFAULT now()
	);
	
	CREATE TABLE IF NOT EXISTS recommendations (
	    id SERIAL PRIMARY KEY,
	    ticker STRING NOT NULL,
	    company STRING NOT NULL,
	    target_from DECIMAL(10,2) NULL,
	    target_to DECIMAL(10,2) NULL,
	    rating_from STRING NULL,
	    rating_to STRING NULL,
	    reason STRING NOT NULL,
	    date TIMESTAMP NOT NULL DEFAULT now()
	);
	`

	_, err = DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("✅ Tablas verificadas/creadas correctamente")
	return nil
}
