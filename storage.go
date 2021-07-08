package main

import (
	"context"
	"database/sql"
	"io"
	"os"

	"cloud.google.com/go/storage"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

var bucketName = "storagedb"
var objectPath = "storagedb.sqlite3"

func objectHandle(ctx context.Context) (*storage.ObjectHandle, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	obj := client.Bucket(bucketName).Object(objectPath)
	return obj, nil
}

func uploadFile(ctx context.Context) error {
	obj, err := objectHandle(ctx)
	if err != nil {
		return err
	}

	w := obj.NewWriter(ctx)

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func downloadFile(ctx context.Context) error {
	obj, err := objectHandle(ctx)
	if err != nil {
		return err
	}

	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	f, err := os.Create(objectPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	return nil
}

func initDb() (*gorp.DbMap, error) {
	db, err := sql.Open("sqlite3", objectPath)
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	return dbmap, err
}
