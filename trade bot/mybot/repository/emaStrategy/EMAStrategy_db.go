package emaStrategy

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type emaStrategyRepositoryDB struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewEmaStraRepository(collection *mongo.Collection, ctx context.Context) EMAStrategyRepository {
	return &emaStrategyRepositoryDB{
		collection: collection,
		ctx:        ctx,
	}
}

func (e *emaStrategyRepositoryDB) Create(ema EMAStrategy) (*EMAStrategy, error) {
	inserted, err := e.collection.InsertOne(e.ctx, ema)
	if err != nil {
		return nil, err
	}
	objId, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}
	ema.Id = objId.Hex()
	return &ema, nil
}

func (e *emaStrategyRepositoryDB) GetAll() ([]EMAStrategy, error) {
	all := []EMAStrategy{}
	cur, err := e.collection.Find(e.ctx, bson.M{})
	if err != nil {
		return all, err
	}
	for cur.Next(e.ctx) {
		ema := EMAStrategy{}
		err = cur.Decode(&ema)
		if err != nil {
			continue
		}
		all = append(all, ema)
	}
	return all, nil
}

func (e *emaStrategyRepositoryDB) GetById(id string) (*EMAStrategy, error) {
	ema := EMAStrategy{}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &ema, err
	}
	err = e.collection.FindOne(e.ctx, bson.M{"_id": _id}).Decode(&ema)
	if err != nil {
		return &ema, err
	}
	return &ema, nil
}

func (e *emaStrategyRepositoryDB) UpdateById(id string, ema EMAStrategy) (*EMAStrategy, error) {
	newEma := EMAStrategy{}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &ema, err
	}
	err = e.collection.FindOne(e.ctx, bson.M{"_id": _id}).Decode(&newEma)
	if err != nil {
		return &ema, err
	}
	opts := options.Update().SetUpsert(true)
	filter := bson.D{primitive.E{Key: "_id", Value: _id}}
	// update := bson.D{{"$set", bson.D{{"email", "newemail@example.com"}}}}

	result, err := e.collection.UpdateOne(e.ctx, filter, newEma, opts)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
		return &ema, nil
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
		objId, ok := result.UpsertedID.(primitive.ObjectID)
		if !ok {
			return nil, err
		}
		ema.Id = objId.Hex()
		return &ema, nil
	}
	return &ema, nil
}
func (e *emaStrategyRepositoryDB) DeleteById(id string) (int, error) {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}

	deleteResult, err := e.collection.DeleteOne(e.ctx, bson.D{primitive.E{Key: "_id", Value: _id}})
	if err != nil {
		return 0, err
	}
	if deleteResult.DeletedCount == 0 {
		return 0, fmt.Errorf("no id matched, so we can't delete")
	}
	return int(deleteResult.DeletedCount), nil
}
