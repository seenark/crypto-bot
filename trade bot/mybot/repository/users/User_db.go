package users

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepositoryDb struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserRepository(cl *mongo.Collection, ctx context.Context) UserRepository {
	return userRepositoryDb{
		collection: cl,
		ctx:        ctx,
	}
}

func (u userRepositoryDb) Create(user User) (string, error) {
	inserted, err := u.collection.InsertOne(u.ctx, user)
	if err != nil {
		return "", err
	}
	objId, ok := inserted.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("cannot translate inserted_id to obj_id")
	}
	return objId.Hex(), nil
}
func (u userRepositoryDb) GetAll() ([]User, error) {
	all := []User{}
	cur, err := u.collection.Find(u.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(u.ctx) {
		user := User{}
		err = cur.Decode(&user)
		if err != nil {
			continue
		}
		all = append(all, user)
	}
	return all, nil
}
func (u userRepositoryDb) GetUserById(id string) (*User, error) {
	user := User{}
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = u.collection.FindOne(u.ctx, bson.D{primitive.E{Key: "_id", Value: _id}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (u userRepositoryDb) UpdateUserById(id string, user User) (*User, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	opts := options.Update().SetUpsert(true)
	updateResult, err := u.collection.UpdateOne(u.ctx, bson.D{primitive.E{Key: "_id", Value: _id}}, user, opts)
	if err != nil {
		return nil, err
	}
	// update
	if updateResult.ModifiedCount > 0 {
		return &user, err
	}
	if updateResult.UpsertedCount > 0 {
		return &user, err
	}
	return nil, fmt.Errorf("no id matched")
}
func (u userRepositoryDb) DeleteUserById(id string) (int, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	deleteResult, err := u.collection.DeleteOne(u.ctx, bson.D{primitive.E{Key: "_id", Value: _id}})
	if err != nil {
		return 0, err
	}
	if deleteResult.DeletedCount < 0 {
		return 0, fmt.Errorf("cannot delete due to no id matched")
	}
	return int(deleteResult.DeletedCount), nil
}
