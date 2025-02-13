package repository

import (
	"SangXanh/pkg/common/query"
	"SangXanh/pkg/model"
	"context"
	"github.com/samber/do/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection[T model.IModel] interface {
	Create(ctx context.Context, data T) error
	Update(ctx context.Context, data T) error
	FindOne(ctx context.Context, filter bson.M) (T, error)
	FindMany(ctx context.Context, filter bson.M, p *query.Pagination) ([]T, error)
	DeleteOne(ctx context.Context, filter bson.M) error
	DeleteMany(ctx context.Context, filter bson.M) error
}

type collection[T model.IModel] struct {
	db *mongo.Collection
}

func newCollection(di do.Injector, name string) *mongo.Collection {
	db := do.MustInvoke[*mongo.Database](di)
	return db.Collection(name)
}

func (r *collection[T]) Create(ctx context.Context, data T) error {
	result, err := r.db.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	data.SetID(result.InsertedID.(primitive.ObjectID))
	return nil
}

func (r *collection[T]) Update(ctx context.Context, data T) error {
	_, err := r.db.UpdateOne(ctx,
		bson.M{
			"_id": data.GetID(),
		},
		bson.M{
			"$set": data,
		},
	)
	return err
}

func (r *collection[T]) FindOne(ctx context.Context, filter bson.M) (T, error) {
	var t T
	if err := r.db.FindOne(ctx, filter).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}

func (r *collection[T]) FindMany(ctx context.Context, filter bson.M, p *query.Pagination) ([]T, error) {
	p.Correct()
	ts := make([]T, 0)
	total, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}
	p.SetTotal(total)

	c, err := r.db.Find(ctx, filter, options.Find().SetLimit(p.Limit).SetSkip((p.Page-1)*p.Limit).SetSort(p.GetSort()))
	if err != nil {
		return nil, err
	}
	if err := c.All(ctx, &ts); err != nil {
		return nil, err
	}
	return ts, nil
}

func (r *collection[T]) DeleteOne(ctx context.Context, filter bson.M) error {
	_, err := r.db.DeleteOne(ctx, filter)
	return err
}

func (r *collection[T]) DeleteMany(ctx context.Context, filter bson.M) error {
	_, err := r.db.DeleteMany(ctx, filter)
	return err
}
