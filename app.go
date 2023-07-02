package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

var Cache = make(map[string]Order)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost dbname=testdb user=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("cant connect to db")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("cant connect to db")
		return nil, err
	}

	return db, nil
}

func cacheInit(db *sql.DB) {
	rows, _ := db.Query("select * from WB_Order")
	defer rows.Close()

	for rows.Next() {
		var curr_order Order

		var items_raw_data interface{}

		rows.Scan(
			&curr_order.OrderUID,
			&curr_order.TrackNumber,
			&curr_order.Entry,
			&curr_order.Delivery.Name,
			&curr_order.Delivery.Phone,
			&curr_order.Delivery.Zip,
			&curr_order.Delivery.City,
			&curr_order.Delivery.Address,
			&curr_order.Delivery.Region,
			&curr_order.Delivery.Email,
			&curr_order.Payment.Transaction,
			&curr_order.Payment.RequestID,
			&curr_order.Payment.Currency,
			&curr_order.Payment.Provider,
			&curr_order.Payment.Amount,
			&curr_order.Payment.PaymentDt,
			&curr_order.Payment.Bank,
			&curr_order.Payment.DeliveryCost,
			&curr_order.Payment.GoodsTotal,
			&curr_order.Payment.CustomFee,
			&items_raw_data,
			&curr_order.Locale,
			&curr_order.InternalSignature,
			&curr_order.CustomerID,
			&curr_order.DeliveryService,
			&curr_order.Shardkey,
			&curr_order.SmID,
			&curr_order.DateCreated,
			&curr_order.OofShard,
		)

		json.Unmarshal(items_raw_data.([]byte), &curr_order.Items)

		Cache[curr_order.OrderUID] = curr_order
	}

}

func JetStreamSet(nc *nats.Conn, db *sql.DB) {
	js, _ := nc.JetStream()

	js.AddStream(&nats.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})

	js.Subscribe("ORDERS.*", func(m *nats.Msg) {
		// здесь нужно будет обрабатывать получение данных
		var curr_order Order

		err := json.Unmarshal(m.Data, &curr_order) // данный код отсекает все невалидные данные.
		if err != nil {
			log.Println("invalid data ingored")
			return
		} else {
			log.Println("data added successfully")
		}
		_, ok := Cache[curr_order.OrderUID]
		if ok {
			return
		}
		Cache[curr_order.OrderUID] = curr_order

		stmt, err := db.Prepare("INSERT INTO WB_Order VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)")
		if err != nil {
			log.Fatal(err)
		}

		items_raw, _ := json.Marshal(curr_order.Items)

		_, err = stmt.Exec(
			curr_order.OrderUID,
			curr_order.TrackNumber,
			curr_order.Entry,
			curr_order.Delivery.Name,
			curr_order.Delivery.Phone,
			curr_order.Delivery.Zip,
			curr_order.Delivery.City,
			curr_order.Delivery.Address,
			curr_order.Delivery.Region,
			curr_order.Delivery.Email,
			curr_order.Payment.Transaction,
			curr_order.Payment.RequestID,
			curr_order.Payment.Currency,
			curr_order.Payment.Provider,
			curr_order.Payment.Amount,
			curr_order.Payment.PaymentDt,
			curr_order.Payment.Bank,
			curr_order.Payment.DeliveryCost,
			curr_order.Payment.GoodsTotal,
			curr_order.Payment.CustomFee,
			items_raw,
			curr_order.Locale,
			curr_order.InternalSignature,
			curr_order.CustomerID,
			curr_order.DeliveryService,
			curr_order.Shardkey,
			curr_order.SmID,
			curr_order.DateCreated,
			curr_order.OofShard,
		)
		if err != nil {
			log.Fatal(err)
		}
	}, nats.Durable("uniqueApp"))
}

func HTTPServerInit() http.Server {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<h1>Введите данные UID <h1> <form action="/get"><label for="fname">UID:</label><br><input type="text" id="uid" name="uid"><br><br><input type="submit" value="Получить"></form> `)

	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("uid")

		val, ok := Cache[id] // нужно это вернуть в виде json
		if ok {
			raw_data, _ := json.MarshalIndent(val, "", "\t")
			w.Write(raw_data)
		} else {
			fmt.Fprintf(w, "<h1> No such UID <h1>")
		}
	})

	httpServer := http.Server{
		Addr: ":8080",
	}

	return httpServer
}

func main() {
	// создание + инициализация объекта БД
	db, err := initDB()
	if err != nil {
		return
	}

	// восстановление кэша из БД
	cacheInit(db)

	// подписка на канал JetStream
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Close()
	JetStreamSet(nc, db)

	// создание http server
	httpServer := HTTPServerInit()

	//graceful shutdown
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-idleConnectionsClosed
	log.Printf("\nBye bye")
}
