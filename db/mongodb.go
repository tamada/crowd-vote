package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoStore struct {
	uri      string
	dbName   string
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoStore(uri, dbName string) *MongoStore {
	return &MongoStore{uri: uri, dbName: dbName}
}

func (s *MongoStore) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(s.uri))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	s.client = client
	s.database = client.Database(s.dbName)

	// Create indexes
	_, _ = s.database.Collection("locations").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "group", Value: 1}, {Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	_, _ = s.database.Collection("votes").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "group", Value: 1}, {Key: "location_id", Value: 1}, {Key: "timestamp", Value: 1}},
	})

	return nil
}

func (s *MongoStore) Close() error {
	if s.client != nil {
		return s.client.Disconnect(context.Background())
	}
	return nil
}

func (s *MongoStore) ListLocations(group string) ([]Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.database.Collection("locations").Find(ctx, bson.M{"group": group})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var locations []Location
	if err := cursor.All(ctx, &locations); err != nil {
		return nil, err
	}
	if locations == nil {
		return []Location{}, nil
	}
	return locations, nil
}

func (s *MongoStore) CreateLocation(l Location) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.database.Collection("locations").InsertOne(ctx, l)
	return err
}

func (s *MongoStore) GetLocation(group, id string) (Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var l Location
	err := s.database.Collection("locations").FindOne(ctx, bson.M{"group": group, "id": id}).Decode(&l)
	return l, err
}

func (s *MongoStore) CreateVote(v Vote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.database.Collection("votes").InsertOne(ctx, v)
	return err
}

func (s *MongoStore) CalculateCrowdRate(group, id string, since time.Time) (CrowdRate, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "group", Value: group},
		{Key: "location_id", Value: id},
		{Key: "timestamp", Value: bson.D{{Key: "$gt", Value: since}}},
	}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: nil},
		{Key: "avgLevel", Value: bson.D{{Key: "$avg", Value: "$level"}}},
		{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
	}}}

	cursor, err := s.database.Collection("votes").Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return CrowdRate{}, err
	}
	defer cursor.Close(ctx)

	rate := CrowdRate{
		Group:      group,
		LocationID: id,
		Since:      since,
	}

	if cursor.Next(ctx) {
		var result struct {
			AvgLevel float64 `bson:"avgLevel"`
			Count    int     `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return rate, err
		}
		rate.Votes = result.AvgLevel
		rate.Total = result.Count
	}

	return rate, nil
}

func (s *MongoStore) ListLocationIDs(group string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := s.database.Collection("locations").Find(ctx, bson.M{"group": group}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ids []string
	for cursor.Next(ctx) {
		var result struct {
			ID string `bson:"id"`
		}
		if err := cursor.Decode(&result); err != nil {
			continue
		}
		ids = append(ids, result.ID)
	}
	return ids, nil
}
