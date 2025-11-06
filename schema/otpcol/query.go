package otpcol

import (
	"api/internal/mongodb"
	bsonutil "api/internal/mongodb/utils"
	"api/internal/timer"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// create
func Create(ctx context.Context, data *OTP) (interface{}, error) {
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
func Update(ctx context.Context, data *OTP) (bool, error) {

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
	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &OTP{})
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

func FindUserByPhone(ctx context.Context, phone string) (*OTP, error) {
	// generate filter
	filter := bsonutil.BsonAdd(nil, "phone", phone)
	filter = bsonutil.BsonGreaterThan(filter, "expire_at", timer.Now())

	// Sort by created_at in descending order to get the latest OTP
	findOptions := options.FindOne().SetSort(bsonutil.BsonAdd(nil, "created_at", -1))

	// end
	return FindWithCondition(ctx, filter, findOptions)
}

func FindUserWithUserID(ctx context.Context, userID string, otpType string) (*OTP, error) {
	// generate filter
	filter := bsonutil.BsonAdd(nil, "user_id", userID)
	filter = bsonutil.BsonGreaterThan(filter, "expire_at", timer.Now())
	filter = bsonutil.BsonAdd(filter, "type", otpType)

	// Sort by created_at in descending order to get the latest OTP
	findOptions := options.FindOne().SetSort(bsonutil.BsonAdd(nil, "created_at", -1))

	// end
	return FindWithCondition(ctx, filter, findOptions)
}

func FindOTP(ctx context.Context, id string) (*OTP, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// filter by project id
	filter := bsonutil.BsonAdd(nil, "_id", objID)

	// end
	return FindWithCondition(ctx, filter)
}

// FindWithCondition find common
func FindWithCondition(ctx context.Context, filter interface{}, findOptions ...*options.FindOneOptions) (*OTP, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &OTP{})

	result := &OTP{}
	if err := coll.FirstWithCtx(ctx, filter, result, findOptions...); err != nil {
		return nil, err
	}

	return result, nil
}

func Collection() *mongodb.Collection {
	return mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &OTP{})
}
