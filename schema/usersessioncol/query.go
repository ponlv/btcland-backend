package usersessioncol

import (
	"api/internal/mongodb"
	bsonutil "api/internal/mongodb/utils"
	"api/internal/timer"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

// create
func Create(ctx context.Context, data *UserSession) (interface{}, error) {
	// get collection
	coll := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), data)

	// set createAt and updateAt
	data.CreatedAt = timer.Now()

	// create
	id, err := coll.CreateWithCtx(ctx, data)
	if err != nil {
		return nil, err
	}

	// end
	return id, nil
}

func FindWithId(ctx context.Context, id string) (*UserSession, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// filter by project id
	filter := bsonutil.BsonAdd(nil, "_id", objID)

	// end
	return FindWithCondition(ctx, filter)
}

func FindWithDeviceId(ctx context.Context, userId string) (*UserSession, error) {

	// filter by project id
	filter := bsonutil.BsonAdd(nil, "user_id", userId)
	filter = bsonutil.BsonNotEqual(filter, "device_token", "")

	// end
	return FindWithCondition(ctx, filter)
}

func RemoveAllSession(ctx context.Context, userId string) error {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &UserSession{})

	filter := bsonutil.BsonAdd(nil, "user_id", userId)
	update := bsonutil.BsonSet(nil, "is_delete", true)

	_, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// FindWithCondition find common
func FindWithCondition(ctx context.Context, filter interface{}, findOptions ...*options.FindOneOptions) (*UserSession, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &UserSession{})

	result := &UserSession{}
	if err := coll.FirstWithCtx(ctx, filter, result, findOptions...); err != nil {
		return nil, err
	}

	return result, nil
}

func FindWithFilter(ctx context.Context, filter primitive.D, ops *options.FindOptions) ([]*UserSession, error) {
	// get collection
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &UserSession{})

	// get all
	var results []*UserSession
	cursor, err := coll.Find(ctx, filter, ops)
	if err != nil {
		return nil, err
	} else {
		if err = cursor.All(ctx, &results); err != nil {
			return nil, err
		}
	}

	// end
	return results, nil
}

// Update update
func Update(ctx context.Context, data *UserSession) (bool, error) {

	objID, err := primitive.ObjectIDFromHex(data.GetIDString())
	if err != nil {
		return false, err
	}

	// filter by project id
	filter := bsonutil.BsonAdd(nil, "_id", objID)

	// update
	update := bsonutil.BsonSetMap(nil,
		bsonutil.ConvertStructToBSONMap(
			data,
			&bsonutil.MappingOpts{RemoveID: true},
		),
	)

	// update
	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &UserSession{})
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}
