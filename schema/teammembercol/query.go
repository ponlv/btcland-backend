package teammembercol

import (
	"api/internal/mongodb"
	bsonutil "api/internal/mongodb/utils"
	"api/internal/timer"
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create tạo mới team member relationship
func Create(ctx context.Context, data *TeamMember) (interface{}, error) {
	coll := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), data)

	data.CreatedAt = timer.Now()
	data.UpdatedAt = timer.Now()
	data.JoinedAt = timer.Now()
	data.IsDelete = false

	id, err := coll.CreateWithCtx(ctx, data)
	if err != nil {
		return nil, err
	}

	return id, nil
}

// FindByManagerID tìm tất cả nhân viên của một manager
func FindByManagerID(ctx context.Context, managerID string, ops *options.FindOptions) ([]*TeamMember, int64, error) {
	filter := bsonutil.BsonAdd(nil, "manager_id", managerID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	return FindWithFilter(ctx, filter, ops)
}

// FindByEmployeeID tìm manager của một employee
func FindByEmployeeID(ctx context.Context, employeeID string) (*TeamMember, error) {
	filter := bsonutil.BsonAdd(nil, "employee_id", employeeID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	return FindWithCondition(ctx, filter)
}

// FindByManagerAndEmployee tìm relationship giữa manager và employee
func FindByManagerAndEmployee(ctx context.Context, managerID, employeeID string) (*TeamMember, error) {
	filter := bsonutil.BsonAdd(nil, "manager_id", managerID)
	filter = bsonutil.BsonAdd(filter, "employee_id", employeeID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	return FindWithCondition(ctx, filter)
}

// FindWithCondition tìm với điều kiện
func FindWithCondition(ctx context.Context, filter interface{}, findOptions ...*options.FindOneOptions) (*TeamMember, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &TeamMember{})

	result := &TeamMember{}
	if err := coll.FirstWithCtx(ctx, filter, result, findOptions...); err != nil {
		return nil, err
	}

	return result, nil
}

// FindWithFilter tìm danh sách với filter
func FindWithFilter(ctx context.Context, filter primitive.D, ops *options.FindOptions) ([]*TeamMember, int64, error) {
	coll := mongodb.CollRead(os.Getenv("MONGODB_DATABASE"), &TeamMember{})

	var results []*TeamMember
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

// SoftDelete xóa mềm team member relationship
func SoftDelete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bsonutil.BsonAdd(nil, "_id", objID)
	update := bsonutil.BsonSet(nil, "is_delete", true)
	update = bsonutil.BsonSet(update, "deleted_at", timer.Now())
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &TeamMember{})
	_, err = collection.UpdateOne(ctx, filter, update)
	return err
}

// SoftDeleteByManagerAndEmployee xóa mềm relationship giữa manager và employee
func SoftDeleteByManagerAndEmployee(ctx context.Context, managerID, employeeID string) error {
	filter := bsonutil.BsonAdd(nil, "manager_id", managerID)
	filter = bsonutil.BsonAdd(filter, "employee_id", employeeID)
	filter = bsonutil.BsonAdd(filter, "is_delete", false)

	update := bsonutil.BsonSet(nil, "is_delete", true)
	update = bsonutil.BsonSet(update, "deleted_at", timer.Now())
	update = bsonutil.BsonSet(update, "updated_at", timer.Now())

	collection := mongodb.Coll(os.Getenv("MONGODB_DATABASE"), &TeamMember{})
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// CheckEmployeeBelongsToManager kiểm tra employee có thuộc team của manager không
func CheckEmployeeBelongsToManager(ctx context.Context, managerID, employeeID string) (bool, error) {
	_, err := FindByManagerAndEmployee(ctx, managerID, employeeID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

