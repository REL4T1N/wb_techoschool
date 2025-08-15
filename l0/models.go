package main

import "time"

type Order struct {
	OrderUID        string    `json:"order_uid"`
	TrackNumber     string    `json:"track_number"`
	Entry           string    `json:"entry"`
	Locale          string    `json:"locale"`
	CustomerID      string    `json:"customer_id"`
	DeliveryService string    `json:"delivery_service"`
	Shardkey        string    `json:"shardkey"`
	SMID            int       `json:"sm_id"`
	DateCreated     time.Time `json:"date_created"`
	OOFShard        string    `json:"oof_shard"`
}

type Delivery struct {
	OrderUID string `json:"order_uid"`
	// добавить остальные поля
}

type Payment struct {
	OrderUID string `json:"order_uid"`
	// добавить остальные поля
}

type Item struct {
	OrderUID string `json:"order_uid"`
	ChrtID   int    `json:"chrt_id"`
	// добавить остальные поля
}
