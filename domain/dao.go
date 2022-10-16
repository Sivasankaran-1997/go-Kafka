package domain

import (
	"context"
	"kafka/utils"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (user *User) Create() (*mongo.InsertOneResult, *utils.Resterr) {
	usersC := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	emailCount, _ := usersC.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if emailCount > 0 {
		return nil, utils.BadRequest("Email Already Register")
	}

	result, err := usersC.InsertOne(ctx, user)

	if err != nil {
		restErr := utils.InternalErr("can't insert user to the database.")
		return nil, restErr
	}

	return result, nil
}

func (user *User) FindUser() *utils.Resterr {

	userC := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	// filter := bson.M{
	// 	"$and": []bson.M{ // you can try this in []interface
	// 		bson.M{"email": user.Email},
	// 		bson.M{"password": user.Password}},
	// }

	filter := bson.M{"email": user.Email}
	err := userC.FindOne(ctx, filter).Decode(&user)

	defer cancel()

	if err != nil {
		return utils.NotFound("Email is Not Found")
	}

	return nil

}

func (user *User) Delete() (*mongo.DeleteResult, *utils.Resterr) {

	userC := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	filter := bson.M{"email": user.Email}
	result, err := userC.DeleteOne(ctx, filter)
	defer cancel()
	if result.DeletedCount == 0 {
		return nil, utils.BadRequest("No Record Found")
	}

	if err != nil {
		return nil, utils.NotFound("Email is Not Found")
	}

	return result, nil

}

func (user *User) Update() (*mongo.UpdateResult, *utils.Resterr) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)

	userC := db.Collection("users")

	filter := bson.M{"email": user.Email}

	updateValue := bson.M{"$set": bson.M{"name": user.Name}}

	opts := options.Update().SetUpsert(true)

	result, err := userC.UpdateOne(ctx, filter, updateValue, opts)

	defer cancel()

	if result.ModifiedCount == 0 {
		return nil, utils.BadRequest("not modified")
	}

	if err != nil {
		return nil, utils.InternalErr("Data not Updated")
	}

	return result, nil

}
