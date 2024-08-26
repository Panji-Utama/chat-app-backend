// models/message.go
package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Sender    string             `bson:"sender" json:"sender"`
    Recipient string             `bson:"recipient" json:"recipient"`
    Content   string             `bson:"content" json:"content"`
    Timestamp primitive.DateTime `bson:"timestamp" json:"timestamp"`
}
