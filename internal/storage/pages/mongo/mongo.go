package mongo

import (
	"context"
	"errors"
	pmongo "github.com/mnogokotin/golang-packages/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"log"
	"time"
	"work-routine-bot/internal/domain"
	"work-routine-bot/internal/storage/pages"
	"work-routine-bot/pkg/e"
)

type Storage struct {
	pages Pages
}

type Pages struct {
	*mongo.Collection
}

type Page struct {
	Url      string `bson:"url"`
	UserName string `bson:"username"`
}

func New(connectionString string, connectTimeout time.Duration, dbName string) Storage {
	mongo_, err := pmongo.New(connectionString, connectTimeout)
	if err != nil {
		log.Fatal(err)
	}

	pages_ := Pages{
		Collection: mongo_.Client.Database(dbName).Collection("pages"),
	}

	return Storage{
		pages: pages_,
	}
}

func (s Storage) Save(ctx context.Context, page *domain.Page) error {
	_, err := s.pages.InsertOne(ctx, Page{
		Url:      page.Url,
		UserName: page.UserName,
	})
	if err != nil {
		return e.Wrap("can't save page", err)
	}

	return nil
}

func (s Storage) PickRandom(ctx context.Context, userName string) (page *domain.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	pipe := bson.A{
		bson.M{"$sample": bson.M{"size": 1}},
	}

	cursor, err := s.pages.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}

	var p Page

	cursor.Next(ctx)

	err = cursor.Decode(&p)
	switch {
	case errors.Is(err, io.EOF):
		return nil, pages.ErrNoSavedPages
	case err != nil:
		return nil, err
	}

	return &domain.Page{
		Url:      p.Url,
		UserName: p.UserName,
	}, nil
}

func (s Storage) Remove(ctx context.Context, domainPage *domain.Page) error {
	_, err := s.pages.DeleteOne(ctx, toPage(domainPage).Filter())
	if err != nil {
		return e.Wrap("can't remove page", err)
	}

	return nil
}

func (s Storage) IsExists(ctx context.Context, domainPage *domain.Page) (bool, error) {
	count, err := s.pages.CountDocuments(ctx, toPage(domainPage).Filter())
	if err != nil {
		return false, e.Wrap("can't check if page exists", err)
	}

	return count > 0, nil
}

func toPage(p *domain.Page) Page {
	return Page{
		Url:      p.Url,
		UserName: p.UserName,
	}
}

func (p Page) Filter() bson.M {
	return bson.M{
		"url":      p.Url,
		"username": p.UserName,
	}
}
