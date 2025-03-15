package models

// Stock representa una acción en la base de datos
type Stock struct {
    Ticker  string `json:"ticker"`  // Símbolo de la acción (ej: AAPL)
    Company string `json:"company"` // Nombre de la empresa
    Rating  string `json:"rating"`  // Calificación de la acción (ej: Buy, Sell, Hold)
}