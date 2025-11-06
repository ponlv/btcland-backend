package userdevicecol

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
func Create(ctx context.Context, data *UserDevice) (interface{}, error) {
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

func DisableAllDevice(ctx context.Context, userId string) error {
	coll := mongodb.CollRead(mongodb.GetDatabaseName(), &UserDevice{})

	filter := bsonutil.BsonAdd(nil, "user_id", userId)
	update := bsonutil.BsonSet(nil, "is_current", false)

	_, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func FindWithDeviceId(ctx context.Context, userId, deviceId string) (*UserDevice, error) {

	// filter by project id
	filter := bsonutil.BsonAdd(nil, "user_id", userId)
	filter = bsonutil.BsonAdd(filter, "device_id", deviceId)

	// end
	return FindWithCondition(ctx, filter)
}

// FindWithCondition find common
func FindWithCondition(ctx context.Context, filter interface{}, findOptions ...*options.FindOneOptions) (*UserDevice, error) {
	coll := mongodb.CollRead(mongodb.GetDatabaseName(), &UserDevice{})

	result := &UserDevice{}
	if err := coll.FirstWithCtx(ctx, filter, result, findOptions...); err != nil {
		return nil, err
	}

	return result, nil
}

func FindWithFilter(ctx context.Context, filter primitive.D, ops *options.FindOptions) ([]*UserDevice, error) {
	// get collection
	coll := mongodb.CollRead(mongodb.GetDatabaseName(), &UserDevice{})

	// get all
	var results []*UserDevice
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
func Update(ctx context.Context, data *UserDevice) (bool, error) {

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
	collection := mongodb.Coll(mongodb.GetDatabaseName(), &UserDevice{})
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}
