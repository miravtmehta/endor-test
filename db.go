package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	ulid "github.com/oklog/ulid/v2"
)

type ObjectDB interface {
	// Store will store the object in the data store. The object will have
	// name and kind, and the Store method should create a unique ID.
	Store(ctx context.Context, object Object) error
	// GetObjectByID will retrieve the object with the provided ID.
	GetObjectByID(ctx context.Context, id string) (Object, error)
	// GetObjectByName will retrieve the object with the given name.
	GetObjectByName(ctx context.Context, name string) ([]Object, error)
	// ListObjects will return a list of all objects of the given kind.
	ListObjects(ctx context.Context, kind string) ([]Object, error)
	// DeleteObject will delete the object.
	DeleteObject(ctx context.Context, id string) error
}

type RedisDb struct {
	client *redis.Client
}

var _ ObjectDB = (*RedisDb)(nil)

func (r *RedisDb) Store(ctx context.Context, object Object) error {
	if object.GetName() == "" || object.GetKind() == "" {
		return fmt.Errorf("name / kind not found")
	}
	id := ulid.Make().String()
	object.SetID(id)

	key := object.GetID() + ":" + object.GetName() + ":" + object.GetKind()
	value, err := json.Marshal(object)
	if err != nil {
		return err
	}

	fmt.Printf("storing object with key %s: value %s\n", key, value)
	err = r.client.WithContext(ctx).Set(key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil

}

func (r *RedisDb) GetObjectByID(ctx context.Context, id string) (Object, error) {
	pattern := fmt.Sprintf("%s:*:*", id)

	fmt.Printf("fetching object by ID: %s\n", id)
	keys, err := r.client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, fmt.Errorf("object with id: %s not found", id)
	}

	if len(keys) > 1 {
		return nil, fmt.Errorf("more than one object ID's: %s", id)
	}

	valueString, err := r.client.WithContext(ctx).Get(keys[0]).Result()
	if err != nil {
		return nil, err
	}

	return objectExtraction(keys[0], valueString)
}

func (r *RedisDb) GetObjectByName(ctx context.Context, name string) ([]Object, error) {
	if name == "" {
		return nil, errors.New("empty object, name field missing")
	}
	pattern := fmt.Sprintf("*:%s:*", name)

	fmt.Printf("get objects by name %s\n", name)
	return r.listObjects(ctx, pattern)

}

func (r *RedisDb) ListObjects(ctx context.Context, kind string) ([]Object, error) {
	if kind == "" {
		return nil, errors.New("empty object, kind field missing")
	}

	return r.listObjects(ctx, fmt.Sprintf("*:*:%s", kind))

}

func (r *RedisDb) DeleteObject(ctx context.Context, id string) error {
	pattern := fmt.Sprintf("%s:*:*", id)
	fmt.Printf("delete objects by id %s\n", id)
	keys, err := r.client.WithContext(ctx).Keys(pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}
	if len(keys) > 1 {
		return fmt.Errorf("more than one object ID's: %s", id)
	}

	_, err = r.client.WithContext(ctx).Del(keys[0]).Result()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisConnection(ctx context.Context, address string, password string) (*RedisDb, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
	})
	_, err := rdb.Ping().Result()
	rdb.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	return &RedisDb{client: rdb}, nil

}

func (r *RedisDb) FlushAllRedisDb(ctx context.Context) error {
	_, err := r.client.WithContext(ctx).FlushAll().Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisDb) listObjects(ctx context.Context, pattern string) ([]Object, error) {
	var object Object
	var objects []Object

	keys, err := r.client.WithContext(ctx).Keys(pattern).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return objects, nil
	}

	for _, key := range keys {
		value, err := r.client.WithContext(ctx).Get(key).Result()
		if err != nil {
			return nil, err
		}
		object, err = objectExtraction(key, value)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	return objects, nil
}
