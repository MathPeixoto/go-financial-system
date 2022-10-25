package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Transaction struct {
	SourceId      string  `json:"sourceId"`
	Value         float32 `json:"value"`
	Date          string  `json:"date"`
	DestinationId string  `json:"destinationId"`
}

type TransactionDb struct {
	ID            bson.ObjectId `bson:"_id" json:"id"`
	SourceId      string        `json:"sourceId"`
	Value         float32       `json:"value"`
	Date          string        `json:"date"`
	DestinationId string        `json:"destinationId"`
}
