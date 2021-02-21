package utils

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

func CursorToList(context context.Context, cursor *mongo.Cursor, arr interface{}) error {
	valuePtr := reflect.ValueOf(arr)
	for cursor.Next(context) {
		elem := valuePtr.Elem()
		err := cursor.Decode(elem)
		if err != nil {
			return err
		}
		elem.Set(reflect.Append(elem, reflect.ValueOf(elem)))
	}
	return nil
}
