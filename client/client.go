package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Valor string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	file, err := os.OpenFile("cotacao.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	cotacao, err := getCotacao(ctx)
	if err != nil {
		log.Printf("Error getting cotacao: %v", err)
		return
	}

	txt := fmt.Sprintf("Dólar: %v\n", cotacao.Valor)
	_, err = file.Write([]byte(txt))
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return
	}

	log.Printf("Cotacao received: %v", cotacao.Valor)

}

func getCotacao(ctx context.Context) (*Cotacao, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Client request timeout exceeded (300ms)")
		}
		log.Printf("Error creating request: %v", err)
		return nil, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Client request timeout exceeded (300ms)")
		}
		log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Client request timeout exceeded (300ms)")
		}
		log.Println("Error reading response body:", err)
		return nil, err
	}
	log.Printf("Return Body: %v", string(body))

	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Client request timeout exceeded (300ms)")
		}
		log.Println("Error unmarshaling JSON:", err)
		return nil, err
	}

	return &cotacao, nil
}
