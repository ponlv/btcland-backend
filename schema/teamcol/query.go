package teamcol

import (
	"api/internal/mongodb"
	bsonutil "api/internal/mongodb/utils"
	"api/internal/timer"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create tạo mới team
func Create(ctx context.Context, data *Team) (interface{}, error) {
	coll := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), data)

	data.CreatedAt = timer.Now()
	data.UpdatedAt = timer.Now()
	data.IsDelete = false

	id, err := coll.CreateWithCtx(ctx, data)
	if err != nil {
		return nil, err
	}

	return id, nil
}

// Update cập nhật team
func Update(ctx context.Context, data *Team) (bool, error) {
	objID, err := primitive.ObjectIDFromHex(data.GetIDString())
	if err != nil {
		return false, err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	data.UpdatedAt = timer.Now()

	update := bsonutil.BsonSetMap(nil,
		bsonutil.ConvertStructToBSONMap(
			data,
			&bsonutil.MappingOpts{RemoveID: true},
		),
	)

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &Team{})
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindByID tìm team theo ID
func FindByID(ctx context.Context, id string) (*Team, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	return FindWithCondition(ctx, filter)
}

// FindWithCondition tìm team với điều kiện
func FindWithCondition(ctx context.Context, filter interface{}, findOptions ...*options.FindOneOptions) (*Team, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &Team{})

	result := &Team{}
	if err := coll.FirstWithCtx(ctx, filter, result, findOptions...); err != nil {
		return nil, err
	}

	return result, nil
}

// FindWithFilter tìm danh sách team với filter
func FindWithFilter(ctx context.Context, filter primitive.D, ops *options.FindOptions) ([]*Team, int64, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &Team{})

	// Thêm điều kiện không bị xóa
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	var results []*Team
	cursor, err := coll.Find(ctx, filter, ops)
	if err != nil {
		return nil, 0, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	count, err := coll.Count(filter)
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// FindByManagerID tìm team theo manager ID
func FindByManagerID(ctx context.Context, managerID string, ops *options.FindOptions) ([]*Team, int64, error) {
	filter := bsonutil.BsonAdd(nil, "manager_id", managerID)
	return FindWithFilter(ctx, filter, ops)
}

// SoftDelete xóa mềm team
func SoftDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	update := bsonutil.BsonSet(nil, "is_delete", true)
	update = bsonutil.BsonSet(update, "deleted_at", timer.Now())
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &Team{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateManagerID cập nhật manager ID cho team
func UpdateManagerID(ctx context.Context, teamID string, managerID string) error {
	objID, err := primitive.ObjectIDFromHex(teamID)
	if err != nil {
		return err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)
	update := bsonutil.BsonSet(nil, "manager_id", managerID)
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &Team{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// Collection trả về collection
func Collection() *mongodb.Collection {
	return mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &Team{})
}

