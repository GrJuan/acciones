package main

import (
    "log"
    "net/http"
    "backend/api"
    "backend/db"
)

func main() {
	// Cadena de conexión a CockroachDB
	connectionString := "postgresql://root@localhost:26257/truora?sslmode=disable"

	// Iniciar la BD y crear tablas si no existen
	if err := db.InitDB(connectionString); err != nil {
		log.Fatalf("❌ Error al conectar a la base de datos: %v", err)
	}
	defer db.DB.Close() // Cierra la conexión global cuando el programa finalice

	// Iniciar el scheduler de actualización diaria.
	api.StartDailyUpdateScheduler()

	// Configurar endpoints
	http.HandleFunc("/api/stocks", api.GetStocksHandler)
	http.HandleFunc("/api/daily-recommendation", api.RecommendationHandler)
	http.HandleFunc("/ws", api.WebSocketHandler)

	log.Println("🚀 Servidor iniciado en http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("❌ Error al iniciar el servidor: %v", err)
	}
}

