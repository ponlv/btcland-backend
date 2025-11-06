package usercol

import (
	"api/internal/mongodb"
	bsonutil "api/internal/mongodb/utils"
	"api/internal/timer"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// create
func Create(ctx context.Context, data *User) (interface{}, error) {
	// get collection
	coll := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), data)

	// set createAt and updateAt
	data.CreatedAt = timer.Now()
	data.UpdatedAt = timer.Now()

	// create
	id, err := coll.CreateWithCtx(ctx, data)
	if err != nil {
		return nil, err
	}

	// end
	return id, nil
}

// Update update
func Update(ctx context.Context, data *User) (bool, error) {

	objID, err := primitive.ObjectIDFromHex(data.GetIDString())
	if err != nil {
		return false, err
	}

	// filter by project id
	filter := bsonutil.BsonAdd(nil, "_id", objID)

	data.UpdatedAt = timer.Now()

	// update
	update := bsonutil.BsonSetMap(nil,
		bsonutil.ConvertStructToBSONMap(
			data,
			&bsonutil.MappingOpts{RemoveID: true},
		),
	)

	// update
	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &User{})
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func FindWithFilter(ctx context.Context, filter primitive.D, ops *options.FindOptions) ([]*User, int64, error) {
	// get collection
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &User{})

	// get all
	var results []*User
	cursor, err := coll.Find(ctx, filter, ops)
	if err != nil {
		return nil, 0, err
	} else {
		if err = cursor.All(ctx, &results); err != nil {
			return nil, 0, err
		}
	}

	// get count
	count, err := coll.Count(filter)
	if err != nil {
		return nil, 0, err
	}

	// end
	return results, count, nil
}

// FindWithUserID find with userid
func FindWithUserID(ctx context.Context, userId string) (*User, error) {

	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	// filter by project id
	filter := bsonutil.BsonAdd(nil, "_id", objID)

	// end
	return FindWithCondition(ctx, filter)
}

func FindWithWallet(ctx context.Context, walletAddress string) (*User, error) {
	// generate filter
	filter := bsonutil.BsonAdd(nil, "wallet_address", walletAddress)

	// end
	return FindWithCondition(ctx, filter)
}

// FindByWalletAddress is an alias for FindWithWallet for better naming
func FindByWalletAddress(ctx context.Context, walletAddress string) (*User, error) {
	return FindWithWallet(ctx, walletAddress)
}

func FindWithCitizenId(ctx context.Context, citizenId string) (*User, error) {
	// generate filter
	filter := bsonutil.BsonAdd(nil, "citizen_id", citizenId)

	// end
	return FindWithCondition(ctx, filter)
}

func FindWithEmail(ctx context.Context, email string) (*User, error) {
	// generate filter
	filter := bsonutil.BsonAdd(nil, "email", email)

	// end
	return FindWithCondition(ctx, filter)
}

func FindUserByPhone(ctx context.Context, phone string, isVerify bool) (*User, error) {
	// generate filter
	filter := bsonutil.BsonAdd(nil, "phone_number", phone)
	filter = bsonutil.BsonAdd(filter, "is_verify_phone", isVerify)

	// end
	return FindWithCondition(ctx, filter)
}

// FindWithCondition find common
func FindWithCondition(ctx context.Context, filter interface{}, findOptions ...*options.FindOneOptions) (*User, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &User{})

	result := &User{}
	if err := coll.FirstWithCtx(ctx, filter, result, findOptions...); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateByID updates a user by ID with the provided update data
func UpdateByID(ctx context.Context, userID primitive.ObjectID, updateData bson.M) (*User, error) {
	// get collection
	coll := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &User{})

	// filter by user ID
	filter := bsonutil.BsonAdd(nil, "_id", userID)

	// prepare update with $set operator
	update := bsonutil.BsonSetMap(nil, updateData)

	// update the document
	_, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	// find and return the updated user
	return FindWithUserID(ctx, userID.Hex())
}

// FindByVeriffDocumentNumber finds a user by Veriff document number
func FindByVeriffDocumentNumber(ctx context.Context, documentNumber string) (*User, error) {
	// get collection
	coll := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &User{})

	// filter by veriff_document_number
	filter := bsonutil.BsonAdd(nil, "veriff_document_number", documentNumber)

	// find the user
	var user User
	err := coll.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func Collection() *mongodb.Collection {
	return mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &User{})
}
