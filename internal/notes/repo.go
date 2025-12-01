package notes

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("note not found")

type Repo struct { col *mongo.Collection }

func NewRepo(db *mongo.Database) (*Repo, error) {
	col := db.Collection("notes")

	// 1) Уникальный индекс по title
	_, err := col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil { return nil, err }

	// 2) Текстовый индекс (⭐): title + content
	_, err = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
	})
	if err != nil { return nil, err }

	// 3) TTL-индекс (⭐): автоудаление по expiresAt
	_, err = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "expiresAt", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil { return nil, err }

	return &Repo{col: col}, nil
}

func (r *Repo) Create(ctx context.Context, title, content string, expiresAt *time.Time) (Note, error) {
	now := time.Now()
	n := Note{Title: title, Content: content, CreatedAt: now, UpdatedAt: now, ExpiresAt: expiresAt}
	res, err := r.col.InsertOne(ctx, n)
	if err != nil { return Note{}, err }
	n.ID = res.InsertedID.(primitive.ObjectID)
	return n, nil
}

func (r *Repo) ByID(ctx context.Context, idHex string) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { return Note{}, ErrNotFound }
	var n Note
	if err := r.col.FindOne(ctx, bson.M{"_id": oid}).Decode(&n); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return Note{}, ErrNotFound }
		return Note{}, err
	}
	return n, nil
}

// List: текстовый поиск через $text. Если q пустой — вернём всё.
func (r *Repo) List(ctx context.Context, q string, limit, skip int64) ([]Note, error) {
	filter := bson.M{}
	if q != "" {
		filter = bson.M{"$text": bson.M{"$search": q}}
	}
	opts := options.Find().SetLimit(limit).SetSkip(skip).SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil { return nil, err }
	defer cur.Close(ctx)

	var out []Note
	for cur.Next(ctx) {
		var n Note
		if err := cur.Decode(&n); err != nil { return nil, err }
		out = append(out, n)
	}
	return out, cur.Err()
}

func (r *Repo) Update(ctx context.Context, idHex string, title, content *string, ttlMinutes *int) (Note, error) {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { return Note{}, ErrNotFound }

	set := bson.M{"updatedAt": time.Now()}
	if title != nil   { set["title"] = *title }
	if content != nil { set["content"] = *content }
	if ttlMinutes != nil {
		if *ttlMinutes <= 0 {
			set["expiresAt"] = nil
		} else {
			v := time.Now().Add(time.Duration(*ttlMinutes) * time.Minute)
			set["expiresAt"] = v
		}
	}

	after := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated Note
	if err := r.col.FindOneAndUpdate(ctx, bson.M{"_id": oid}, bson.M{"$set": set}, after).Decode(&updated); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return Note{}, ErrNotFound }
		return Note{}, err
	}
	return updated, nil
}

func (r *Repo) Delete(ctx context.Context, idHex string) error {
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil { return ErrNotFound }
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil { return err }
	if res.DeletedCount == 0 { return ErrNotFound }
	return nil
}
