package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Загружаем конфигурацию
	cfg := LoadConfig()

	// Подключаемся к PostgreSQL
	db := InitDB(cfg)
	defer db.Close()

	// Создаём тестовый заказ
	testOrder := Order{
		OrderUID:        "test124",
		TrackNumber:     "TRACK124",
		Entry:           "entry124",
		Locale:          "en",
		CustomerID:      "cust124",
		DeliveryService: "DHL",
		Shardkey:        "124",
		SMID:            1,
		DateCreated:     time.Now(),
		OOFShard:        "shard124",
	}

	// Вставляем тестовый заказ в БД
	_, err := db.Exec(`INSERT INTO orders (order_uid, track_number, entry, locale, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		testOrder.OrderUID, testOrder.TrackNumber, testOrder.Entry, testOrder.Locale, testOrder.CustomerID,
		testOrder.DeliveryService, testOrder.Shardkey, testOrder.SMID, testOrder.DateCreated, testOrder.OOFShard,
	)
	if err != nil {
		log.Fatal("Ошибка вставки в Postgres:", err)
	}
	fmt.Println("Тестовый заказ добавлен в Postgres!")

	// Настраиваем Kafka продюсер
	kafkaURL := fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaURL},
		Topic:    "orders-topic",
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	// Сериализуем заказ в JSON
	orderJSON, err := json.Marshal(testOrder)
	if err != nil {
		log.Fatal("Ошибка сериализации заказа:", err)
	}

	// Отправляем сообщение в Kafka
	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(testOrder.OrderUID),
			Value: orderJSON,
		},
	)
	if err != nil {
		log.Fatal("Ошибка отправки в Kafka:", err)
	}
	fmt.Println("Тестовое сообщение отправлено в Kafka!")
}
