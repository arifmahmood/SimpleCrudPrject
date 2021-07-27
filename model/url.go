package model

import "time"

type Url struct {
	ID        	string    	`json:"id" bson:"-"`
	Url      	string    	`json:"url" bson:"url"`
	Description string    	`json:"description" bson:"description"`
	Status   	string		`json:"status" bson:"status"`
	CreatedAt 	time.Time 	`json:"createdAt" bson:"createdAt"`
	UpdatedAt 	time.Time 	`json:"updatedAt" bson:"updatedAt"`
}
