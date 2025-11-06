package utils

import "go.mongodb.org/mongo-driver/bson"

func BsonOr(b bson.D, conditions ...bson.D) bson.D {
	return append(b, bson.E{Key: "$or", Value: conditions})
}
