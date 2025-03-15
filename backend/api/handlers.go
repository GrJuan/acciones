package api


import (
	"backend/db"
    "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
    "strconv"
	"github.com/gorilla/websocket"
)


// ----------------------------------------------------------------
// Lógica de WebSocket (SIN CAMBIOS)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // En producción, restringe esto
	},
}

var (
	connectedUsers = make(map[string]*websocket.Conn) // Mapa de conexiones
	userAvatars    = make(map[string]string)            // Mapa de avatares asignados a cada usuario
	mu             sync.Mutex                           // Mutex para sincronizar acceso
)

var avatars = []string{
	"https://media.istockphoto.com/id/1090878494/es/foto/retrato-de-joven-sonriente-a-hombre-guapo-en-camiseta-polo-azul-aislado-sobre-fondo-gris-de.jpg?s=612x612&w=0&k=20&c=dHFsDEJSZ1kuSO4wTDAEaGOJEF-HuToZ6Gt-E2odc6U=",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQF5TcVFjPc_Z0ZdLUAA2Df6uTrJL1C5Al4-w&s",
	"https://pymstatic.com/5844/conversions/personas-emocionales-wide_webp.webp",
	"https://cdn.prod.website-files.com/6452937893cd845b6181c39e/65cd10b25a74e4e74e77bc2c_Joosten-David-240117-8464-fave-final-BG13-linkedin.jpg",
	"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcT2mxuoA14QzeYIyKPQI6EJ-eSU4Yr74HtUPw&s",
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Intentamos obtener el user_id de la query
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		userID = generateUserID()
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error al actualizar a WebSocket:", err)
		return
	}
	defer conn.Close()

	// Seleccionar avatar (según tu lógica)
	avatar := avatars[rand.Intn(len(avatars))]
	mu.Lock()
	userAvatars[userID] = avatar
	connectedUsers[userID] = conn

	// Construir la lista actual de usuarios conectados
	currentUsers := []map[string]string{}
	for id := range connectedUsers {
		currentUsers = append(currentUsers, map[string]string{
			"id":     id,
			"avatar": userAvatars[id],
		})
	}
	mu.Unlock()

	welcomeMessage := map[string]interface{}{
		"event":        "user-connected",
		"id":           userID,
		"avatar":       avatar,
		"currentUsers": currentUsers,
	}

	if err := conn.WriteJSON(welcomeMessage); err != nil {
		log.Println("Error al enviar mensaje WebSocket:", err)
		return
	}

	broadcastUserConnected(userID, avatar)

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				broadcastUserDisconnected(userID)
				mu.Lock()
				delete(connectedUsers, userID)
				delete(userAvatars, userID)
				mu.Unlock()
				log.Printf("Usuario desconectado: %s", userID)
				return
			}

			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Println("Error al parsear el mensaje:", err)
				continue
			}

			if eventType, ok := msg["event"].(string); ok && eventType == "private-message" {
				recipient, ok := msg["recipient"].(string)
				if !ok {
					log.Println("Mensaje privado sin destinatario")
					continue
				}
				mu.Lock()
				if recConn, exists := connectedUsers[recipient]; exists {
					recConn.WriteJSON(msg)
				} else {
					log.Printf("Destinatario %s no conectado", recipient)
				}
				mu.Unlock()
			} else {
				log.Printf("Mensaje recibido de %s: %s", userID, message)
			}
		}
	}()

	select {}
}

func broadcastUserDisconnected(userID string) {
	mu.Lock()
	defer mu.Unlock()
	for id, conn := range connectedUsers {
		if id != userID {
			conn.WriteJSON(map[string]string{
				"event": "user-disconnected",
				"id":    userID,
			})
		}
	}
}

func generateUserID() string {
	return "user-" + time.Now().Format("20060102150405")
}

func broadcastUserConnected(userID, avatar string) {
	mu.Lock()
	defer mu.Unlock()
	for id, conn := range connectedUsers {
		if id != userID {
			conn.WriteJSON(map[string]string{
				"event":  "user-connected",
				"id":     userID,
				"avatar": avatar,
			})
		}
	}
}

// ----------------------------------------------------------------
// Integración con CockroachDB, token y caché para recomendación del día

// Stock representa la estructura recibida de la API de stocks.
// Stock representa la estructura recibida de la API de stocks.
type Stock struct {
	Ticker     string `json:"ticker"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Company    string `json:"company"`
	Action     string `json:"action"`
	Brokerage  string `json:"brokerage"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Time       string `json:"time"`
}

// Recommendation define la estructura de la recomendación del día.
type Recommendation struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	Reason     string `json:"reason"`
	Date       string `json:"date"`
}

// Variables para cachear la recomendación diaria.
var (
	cachedRecommendation *Recommendation
	cacheMutex           sync.RWMutex
)

func cleanPrice(price string) (float64, error) {
	price = strings.Replace(price, "$", "", -1)  // Elimina el símbolo $
	price = strings.Replace(price, ",", "", -1)  // Elimina comas
	return strconv.ParseFloat(price, 64)         // Convierte a float64
}


func GetStocksHandler(w http.ResponseWriter, r *http.Request) {
	// Permitir CORS (opcional, para desarrollo)
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	// Leer parámetros de consulta: page y limit
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	page := 1
	limit := 10
	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
	}
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}
	}
	offset := (page - 1) * limit

	// Consultar los stocks ordenados por fecha (descendente)
	rows, err := db.DB.Query(`
        SELECT ticker, company, target_from, target_to, rating_from, rating_to, time
        FROM stocks
        ORDER BY time DESC
        LIMIT $1 OFFSET $2
    `, limit, offset)
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var stocks []Stock
	for rows.Next() {
		var s Stock
		var t time.Time
		if err := rows.Scan(&s.Ticker, &s.Company, &s.TargetFrom, &s.TargetTo, &s.RatingFrom, &s.RatingTo, &t); err != nil {
			http.Error(w, "Error al leer resultados", http.StatusInternalServerError)
			return
		}
		s.Time = t.Format(time.RFC3339)
		stocks = append(stocks, s)
	}

	// Si se obtuvieron 'limit' registros, asumimos que hay más páginas.
	var nextPage string
	if len(stocks) == limit {
		nextPage = fmt.Sprintf("%d", page+1)
	} else {
		nextPage = ""
	}

	result := map[string]interface{}{
		"items":     stocks,
		"next_page": nextPage,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}


func calculateMLRecommendations(excludeTicker string) ([]Recommendation, error) {
	var query string
	var rows *sql.Rows
	var err error

	if excludeTicker != "" {
		query = `
			SELECT ticker, company, target_from, target_to, rating_from, rating_to, time
			FROM stocks
			WHERE ticker != $1
			ORDER BY (
				COALESCE(target_to, 0) - COALESCE(target_from, 0)
			) DESC
			LIMIT 5;
		`
		rows, err = db.DB.Query(query, excludeTicker)
	} else {
		query = `
			SELECT ticker, company, target_from, target_to, rating_from, rating_to, time
			FROM stocks
			WHERE ticker != $1
			ORDER BY (
				COALESCE(target_to, 0) - COALESCE(target_from, 0)
			) DESC
			LIMIT 5;
		`	
		rows, err = db.DB.Query(query)
	}
	if err != nil {
		log.Println("Error executing ML query: ", err)
		return nil, err
	}
	defer rows.Close()

	var mlRecs []Recommendation
	for rows.Next() {
		var rec Recommendation
		var t time.Time
		err := rows.Scan(&rec.Ticker, &rec.Company, &rec.TargetFrom, &rec.TargetTo, &rec.RatingFrom, &rec.RatingTo, &t)
		if err != nil {
			log.Println("Error scanning ML row: ", err)
			return nil, err
		}
		rec.Date = t.Format("2006-01-02")
		rec.Reason = "ML Recommendation"
		mlRecs = append(mlRecs, rec)
	}
	return mlRecs, nil
}

// fetchStocksFromAPI realiza una solicitud GET a la API de stocks usando el token y un parámetro de paginación.
func fetchStocksFromAPI(token string, nextPage string) ([]Stock, string, error) {
	apiURL := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
	if nextPage != "" {
		apiURL = fmt.Sprintf("%s?next_page=%s", apiURL, nextPage)
	}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, "", err
	}
	// Asegurarse de que el token incluya el prefijo "Bearer "
	if !strings.HasPrefix(token, "Bearer ") {
		token = "Bearer " + token
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("error en la respuesta: %d", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	var result struct {
		Items    []Stock `json:"items"`
		NextPage string  `json:"next_page"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, "", err
	}
	return result.Items, result.NextPage, nil
}

// updateStockData recorre todas las páginas de la API y actualiza/inserta los stocks en CockroachDB.
func updateStockData(token string) {
	var nextPage string = ""
	for {
		stocks, np, err := fetchStocksFromAPI(token, nextPage)
		if err != nil {
			log.Println("Error al obtener stocks de la API:", err)
			return
		}
		for _, stock := range stocks {
			// Limpiar los valores de target_from y target_to
			targetFrom, err := cleanPrice(stock.TargetFrom)
			if err != nil {
				log.Println("Error al convertir target_from:", err)
				targetFrom = 0
			}
		
			targetTo, err := cleanPrice(stock.TargetTo)
			if err != nil {
				log.Println("Error al convertir target_to:", err)
				targetTo = 0
			}
		
			_, err = db.DB.Exec(`
				INSERT INTO stocks (ticker, company, target_from, target_to, rating_from, rating_to, time)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				ON CONFLICT (ticker) DO UPDATE SET
					company = EXCLUDED.company,
					target_from = EXCLUDED.target_from,
					target_to = EXCLUDED.target_to,
					rating_from = EXCLUDED.rating_from,
					rating_to = EXCLUDED.rating_to,
					time = EXCLUDED.time
			`, stock.Ticker, stock.Company, targetFrom, targetTo, stock.RatingFrom, stock.RatingTo, stock.Time)
		
			if err != nil {
				log.Println("Error al insertar/actualizar stock:", err)
			}
		}
		
		if np == "" {
			break
		}
		nextPage = np
	}
}

// calculateDailyRecommendation consulta CockroachDB para obtener la recomendación del día.
// En este ejemplo, se selecciona un stock con rating_to "outperform" (en minúsculas).
func calculateDailyRecommendation() (*Recommendation, error) {
	var rec Recommendation
	var t time.Time

	// Intentamos obtener un stock con rating "outperform".
	row := db.DB.QueryRow(`
        SELECT ticker, company, target_from, target_to, rating_from, rating_to, time
        FROM stocks
        WHERE LOWER(rating_to) = 'outperform'
        ORDER BY time DESC
        LIMIT 1
    `)
	err := row.Scan(&rec.Ticker, &rec.Company, &rec.TargetFrom, &rec.TargetTo, &rec.RatingFrom, &rec.RatingTo, &t)
	if err == sql.ErrNoRows {
		// Fallback: seleccionar el stock con mayor diferencia entre targets.
		row = db.DB.QueryRow(`
          SELECT ticker, company, target_from, target_to, rating_from, rating_to, time
          FROM stocks
          ORDER BY (CAST(REPLACE(target_to, '$','') AS FLOAT) - CAST(REPLACE(target_from, '$','') AS FLOAT)) DESC
          LIMIT 1
      `)
		err = row.Scan(&rec.Ticker, &rec.Company, &rec.TargetFrom, &rec.TargetTo, &rec.RatingFrom, &rec.RatingTo, &t)
		if err != nil {
			return nil, err
		}
		rec.Reason = "Fallback: Mejor diferencia entre targets"
	} else if err != nil {
		return nil, err
	} else {
		rec.Reason = "Basado en rendimiento 'Outperform'"
	}
	rec.Date = t.Format("2006-01-02")
	return &rec, nil
}
// RecommendationHandler expone el endpoint para obtener la recomendación del día.
// Si ya existe una recomendación cacheada para el día actual, la retorna.
func RecommendationHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodOptions {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
        w.WriteHeader(http.StatusOK)
        return
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Token no provisto", http.StatusUnauthorized)
		return
	}

	cacheMutex.RLock()
	if cachedRecommendation != nil && cachedRecommendation.Date == time.Now().Format("2006-01-02") {
		dailyRec := cachedRecommendation
		cacheMutex.RUnlock()
		mlRecs, err := calculateMLRecommendations(dailyRec.Ticker)
		if err != nil {
			http.Error(w, "No se pudieron calcular las recomendaciones ML", http.StatusInternalServerError)
			return
		}
		result := map[string]interface{}{
			"daily_recommendation": dailyRec,
			"ml_recommendations":   mlRecs,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
		return
	}
	cacheMutex.RUnlock()

	updateStockData(token)
	dailyRec, err := calculateDailyRecommendation()
	if err != nil {
		log.Println("Error al calcular recomendación:", err)
		http.Error(w, "No se pudo calcular la recomendación", http.StatusInternalServerError)
		return
	}

	mlRecs, err := calculateMLRecommendations(dailyRec.Ticker)
	if err != nil {
		http.Error(w, "No se pudieron calcular las recomendaciones ML", http.StatusInternalServerError)
		return
	}

	cacheMutex.Lock()
	cachedRecommendation = dailyRec
	cacheMutex.Unlock()

	result := map[string]interface{}{
		"daily_recommendation": dailyRec,
		"ml_recommendations":   mlRecs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}



// StartDailyUpdateScheduler programa una actualización diaria de los stocks y la recomendación.
// Se utiliza un token de servicio (por ejemplo, definido en una variable de entorno) para la actualización.
func StartDailyUpdateScheduler() {
	token := os.Getenv("API_TOKEN")
	if token == "" {
		log.Println("No se ha definido API_TOKEN para actualización diaria")
		return
	}
	go func() {
		for {
			now := time.Now()
			// Calcula el siguiente instante a la medianoche.
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			sleepDuration := nextMidnight.Sub(now)
			log.Printf("Scheduler dormirá por %s hasta la siguiente actualización", sleepDuration)
			time.Sleep(sleepDuration)

			log.Println("Ejecutando actualización diaria...")
			updateStockData(token)
			rec, err := calculateDailyRecommendation()
			if err != nil {
				log.Println("Error al calcular la recomendación diaria:", err)
			} else {
				cacheMutex.Lock()
				cachedRecommendation = rec
				cacheMutex.Unlock()
				log.Println("Recomendación diaria actualizada:", rec)
			}
		}
	}()
}
