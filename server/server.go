package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Exchange struct {
	ID     int    `gorm:"primaryKey"`
	Usdbrl Usdbrl `json:"USDBRL" gorm:"embedded"`
}

type CurrentExchangeValue struct {
	Bid string `json:"bid"`
}

type Usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

var db *gorm.DB

func main() {
	var err error

	log.Printf("Starting Database connection")
	db, err = gorm.Open(sqlite.Open("exchange.db"), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database:", err)
		return
	}

	db.AutoMigrate(&Exchange{})
	log.Printf("Database connected successfully")

	log.Printf("Initializing HTTP server")
	http.HandleFunc("/cotacao", ExchangeHandler)
	log.Printf("Server Init -> Started at port :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Failed to start HTTP server:", err)
		return
	}
}

func ExchangeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Printf("Handling /cotacao request")
	exchange, err := getExchange(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch exchange rate", http.StatusInternalServerError)
		return
	}

	_, err = saveExchangeRate(ctx, exchange)
	if err != nil {
		http.Error(w, "Failed to save exchange rate", http.StatusInternalServerError)
		return
	}

	exchangeRate := CurrentExchangeValue{
		Bid: exchange.Usdbrl.Bid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("Current Exchange Value: %+v\n", exchangeRate)
	json.NewEncoder(w).Encode(exchangeRate)
}

func getExchange(ctx context.Context) (*Exchange, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	log.Printf("Fetching exchange rate from external API")
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("API request timeout exceeded (200ms)")
		}
		log.Println("Error creating request:", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("API request timeout exceeded (200ms)")
		}
		log.Println("Error making request:", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("API request timeout exceeded (200ms)")
		}
		log.Println("Error reading response body:", err)
		return nil, err
	}

	var exchange Exchange
	err = json.Unmarshal(body, &exchange)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("API request timeout exceeded (200ms)")
		}
		log.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	log.Printf("Exchange Rate: %+v\n", exchange.Usdbrl)
	return &exchange, nil

}

func saveExchangeRate(ctx context.Context, exchange *Exchange) (*Exchange, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	log.Printf("Persisting exchange rate: %+v\n", exchange)

	err := db.WithContext(ctx).Create(&exchange).Error
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Database operation timeout exceeded (10ms)")
		}
		log.Println("Failed to save exchange rate:", err)
		return nil, err
	}
	log.Printf("Exchange rate saved successfully: %+v\n", exchange)

	return exchange, nil
}
