package workconfirmationcol

import (
	"api/internal/mongodb"
	bsonutil "api/internal/mongodb/utils"
	"api/internal/timer"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create tạo mới đơn xác nhận công tác
func Create(ctx context.Context, data *WorkConfirmation) (interface{}, error) {
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

// Update cập nhật đơn xác nhận công tác
func Update(ctx context.Context, data *WorkConfirmation) (bool, error) {
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

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindByID tìm đơn theo ID
func FindByID(ctx context.Context, id string) (*WorkConfirmation, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	return FindWithCondition(ctx, filter)
}

// FindWithCondition tìm đơn với điều kiện
func FindWithCondition(ctx context.Context, filter interface{}, findOptions ...*options.FindOneOptions) (*WorkConfirmation, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})

	result := &WorkConfirmation{}
	if err := coll.FirstWithCtx(ctx, filter, result, findOptions...); err != nil {
		return nil, err
	}

	return result, nil
}

// FindWithFilter tìm danh sách đơn với filter
func FindWithFilter(ctx context.Context, filter primitive.D, ops *options.FindOptions) ([]*WorkConfirmation, int64, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})

	// Thêm điều kiện không bị xóa
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	var results []*WorkConfirmation
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

// FindByCreatedBy tìm đơn theo người tạo
func FindByCreatedBy(ctx context.Context, userID string, ops *options.FindOptions) ([]*WorkConfirmation, int64, error) {
	filter := bsonutil.BsonAdd(nil, "created_by", userID)
	return FindWithFilter(ctx, filter, ops)
}

// FindByStatus tìm đơn theo trạng thái
func FindByStatus(ctx context.Context, status WorkConfirmationStatus, ops *options.FindOptions) ([]*WorkConfirmation, int64, error) {
	filter := bsonutil.BsonAdd(nil, "status", status)
	return FindWithFilter(ctx, filter, ops)
}

// UpdateStatus cập nhật trạng thái đơn
func UpdateStatus(ctx context.Context, id string, status WorkConfirmationStatus) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	update := bsonutil.BsonSet(nil, "status", status)
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// ApproveByManager xác nhận bởi quản lý
func ApproveByManager(ctx context.Context, id string, managerID string, comment string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	approval := ApprovalInfo{
		ApprovedBy: managerID,
		ApprovedAt: timer.Now(),
		Comment:    comment,
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	update := bsonutil.BsonSet(nil, "manager_approval", approval)
	update = bsonutil.BsonSet(update, "status", StatusPendingLeader)
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// ApproveByLeader xác nhận bởi lãnh đạo
func ApproveByLeader(ctx context.Context, id string, leaderID string, comment string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	approval := ApprovalInfo{
		ApprovedBy: leaderID,
		ApprovedAt: timer.Now(),
		Comment:    comment,
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	update := bsonutil.BsonSet(nil, "leader_approval", approval)
	update = bsonutil.BsonSet(update, "status", StatusApproved)
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// Reject từ chối đơn
func Reject(ctx context.Context, id string, rejectedBy string, reason string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	rejection := RejectionInfo{
		RejectedBy: rejectedBy,
		RejectedAt: timer.Now(),
		Reason:     reason,
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	update := bsonutil.BsonSet(nil, "rejection", rejection)
	update = bsonutil.BsonSet(update, "status", StatusRejected)
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDelete xóa mềm đơn
func SoftDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	update := bsonutil.BsonSet(nil, "is_delete", true)
	update = bsonutil.BsonSet(update, "deleted_at", timer.Now())
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// FindByIDs tìm nhiều đơn theo danh sách ID
func FindByIDs(ctx context.Context, ids []string) ([]*WorkConfirmation, error) {
	var objectIDs []primitive.ObjectID
	for _, id := range ids {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objID)
	}

	if len(objectIDs) == 0 {
		return []*WorkConfirmation{}, nil
	}

	filter := primitive.D{
		{Key: "_id", Value: primitive.D{{Key: "$in", Value: objectIDs}}},
		{Key: "is_delete", Value: false},
	}

	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
	var results []*WorkConfirmation
	cursor, err := coll.Find(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// Collection trả về collection
func Collection() *mongodb.Collection {
	return mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &WorkConfirmation{})
}

