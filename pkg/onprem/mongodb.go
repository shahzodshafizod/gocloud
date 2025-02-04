package onprem

/*
Notes:
Tables in MongoDB are called Collections.
Rows are Documents
Columns are Fields
*/

import (
	"context"
	"os"

	"github.com/pkg/errors"
	"github.com/shahzodshafizod/gocloud/pkg"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type nosql struct {
	// client    *mongo.Client
	database *mongo.Database
}

func NewNoSQL() (pkg.NoSQL, error) {
	opts := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	if username != "" && password != "" {
		opts = opts.SetAuth(options.Credential{Username: username, Password: password})
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, errors.Wrap(err, "mongo.Connect")
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, errors.Wrap(err, "cl.Ping")
	}

	var database = client.Database(os.Getenv("MONGODB_DATABASE"))

	return &nosql{database: database}, nil
}

func (n *nosql) Insert(ctx context.Context, collection string, item pkg.Map) (string, error) {
	result, err := n.database.
		Collection(collection).
		InsertOne(ctx, item)
	if err != nil {
		return "", n.parseError(err, "n.database.InsertOne")
	}
	return n.parseID(result.InsertedID), nil
}

func (n *nosql) GetItem(ctx context.Context, collection string, filter pkg.Map) (pkg.Map, error) {
	n.setID(filter)
	var document pkg.Map
	var err = n.database.
		Collection(collection).
		FindOne(ctx, bson.M(filter)).
		Decode(&document)
	if err != nil {
		return nil, n.parseError(err, "n.database.Decode")
	}
	document["id"] = n.parseID(document["_id"])
	return document, err
}

func (n *nosql) GetItems(ctx context.Context, collection string, filter pkg.Map) ([]pkg.Map, error) {
	n.setID(filter)
	cursor, err := n.database.
		Collection(collection).
		Find(ctx, bson.M(filter) /*options.Find().SetSort(bson.D{{Key: "regdate", Value: -1}})*/)
	if err != nil {
		return nil, n.parseError(err, "n.collection.Find")
	}
	var documents []pkg.Map
	err = cursor.All(ctx, &documents)
	if err != nil {
		return nil, errors.Wrap(err, "cursor.All")
	}
	for idx := range documents {
		documents[idx]["id"] = n.parseID(documents[idx]["_id"])
	}
	return documents, nil
}

func (n *nosql) Update(ctx context.Context, collection string, filter pkg.Map, update pkg.Map) (pkg.Map, error) {
	n.setID(filter)
	var client = n.database.Collection(collection)
	result, err := client.
		UpdateOne(ctx, bson.M(filter), bson.M{"$set": bson.M(update)})
	if err != nil {
		return nil, n.parseError(err, "n.collection.UpdateOne")
	}
	if result.ModifiedCount == 0 {
		return nil, pkg.ErrNoRowsAffected
	}

	var document pkg.Map
	err = client.FindOne(ctx, bson.M(filter)).Decode(&document)
	if err != nil {
		return nil, n.parseError(err, "n.collection.FindOne")
	}
	document["id"] = n.parseID(document["_id"])
	return document, nil
}

func (n *nosql) setID(filter pkg.Map) {
	if _, ok := filter["id"].(string); ok {
		filter["_id"], _ = bson.ObjectIDFromHex(filter["id"].(string))
		delete(filter, "id")
	}
}

func (n *nosql) parseError(err error, msg string) error {
	if err == mongo.ErrNoDocuments {
		return pkg.ErrNoRows
	}
	if mongo.IsDuplicateKeyError(err) {
		return pkg.ErrDuplicate
	}
	return errors.Wrap(err, msg)
}

func (n *nosql) parseID(id any) string {
	if id == nil {
		return ""
	}
	if objectID, ok := id.(bson.ObjectID); ok {
		return objectID.Hex()
	}
	return ""
}
