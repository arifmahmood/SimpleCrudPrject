package mongorepo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"simple-crud-project/model"
	"time"
)

type bsonUrl struct {
	model.Url `bson:"inline"`
	ID         primitive.ObjectID `bson:"_id"`
}

func (u *bsonUrl) toModel() *model.Url {
	u.Url.ID = u.ID.Hex()
	return &u.Url
}


type Url struct {
	db   *mongo.Database
	collectionName string
}


func (u *Url) collection() *mongo.Collection {
	return u.db.Collection(u.collectionName)
}


func NewUrl(db *mongo.Database, table string) *Url {
	return &Url{db, table}
}


func (u *Url) EnsureIndices(*model.Url) error {
	log.Println("Starting EnsureIndices")
	_, err := u.collection().Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "url", Value: 1}},
				Options: options.Index().SetUnique(true)},
		})
	log.Println("Completed EnsureIndices", err)
	return err
}

func (u *Url) Create(user *model.Url) error {
	log.Println("Starting Create", user)
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := u.collection().InsertOne(context.Background(), user)
	if err != nil {
		log.Println("Completed Create", err)
		if err, ok := err.(mongo.WriteException); ok {
			for _, err := range err.WriteErrors {
				if err.Code == 11000 {
					return err
				}
			}
		}
		return err
	}
	user.ID = result.InsertedID.(primitive.ObjectID).Hex()
	log.Println("Completed Create")
	return nil
}

func (u *Url) Fetch(urlName string) (*model.Url, error) {
	log.Println("Starting Fetch", urlName)

	result := u.collection().FindOne(context.Background(), bson.M{"url": urlName})
	if err := result.Err(); err != nil {
		log.Println("Completed Fetch", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var url bsonUrl
	if err := result.Decode(&url); err != nil {
		log.Println("Completed Fetch", err)
		return nil, err
	}
	log.Println("Completed Fetch")
	return url.toModel(), nil
}
func (u *Url) Delete(urlName string) (int64, error) {
	log.Println("Starting Deleting", urlName)

	result, err := u.collection().DeleteOne(context.Background(), bson.M{"url": urlName})
	if err != nil {
		log.Println("Deleting error", err)
		if err == mongo.ErrNoDocuments {
			return 0, nil
		}
		return 0, err
	}

	count := result.DeletedCount

	log.Println("Completed Fetch")
	return count, nil
}





