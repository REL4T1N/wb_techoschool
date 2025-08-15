package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func main() {
	// --- Подключение к Postgres ---
	psqlInfo := "host=localhost port=5433 user=order_user password=order_pass dbname=orders_db sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Ошибка подключения к Postgres:", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Postgres недоступен:", err)
	}
	fmt.Println("Postgres подключён!")

	// --- Вставка тестового заказа ---
	_, err = db.Exec(`INSERT INTO orders(order_uid, track_number, entry, locale, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
                      VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
                      ON CONFLICT DO NOTHING`,
		"test123", "TRACK123", "WBIL", "en", "cust_test", "meest", "9", 99, time.Now(), "1")
	if err != nil {
		log.Fatal("Ошибка вставки заказа:", err)
	}
	fmt.Println("Тестовый заказ добавлен в Postgres!")

	// --- Подключение к Kafka с ожиданием готовности ---
	for {
		conn, err := kafka.Dial("tcp", "localhost:9092")
		if err == nil {
			conn.Close()
			break
		}
		fmt.Println("Ждем Kafka...")
		time.Sleep(2 * time.Second)
	}
	fmt.Println("Kafka готова!")

	// --- Подключение к Kafka ---
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "orders-topic",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	// --- Отправка тестового сообщения ---
	msg := kafka.Message{
		Key:   []byte("test123"),
		Value: []byte(`{"order_uid":"test123","status":"created"}`),
	}

	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Fatal("Ошибка отправки в Kafka:", err)
	}

	fmt.Println("Тестовое сообщение отправлено в Kafka!")
}
