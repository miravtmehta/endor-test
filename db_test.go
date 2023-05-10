package main

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"math/rand"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(prefix string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return fmt.Sprintf("%s-%s", prefix, string(b))
}

type testRedis struct {
	client          *RedisDb
	t               *testing.T
	requestObjects  []Object
	responseObjects []Object
}

func (r *testRedis) givenObjects(n int) *testRedis {
	for i := 0; i < n; i++ {
		r.requestObjects = append(r.requestObjects, newObject())
	}
	return r
}

func newObject() Object {
	return &Animal{Name: RandStringBytes("dog", 4), Type: "tamed", OwnerID: "john"}

}

func TestRedisDb(t *testing.T) {

	ctx := context.Background()
	redisDb, err := NewRedisConnection(ctx, "127.0.0.1:6379", "")
	if err != nil {
		t.Fatal(err)

	}

	testcases := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "testRedisGetByID",
			run:  testRedisGetByID(redisDb),
		},
		{
			name: "testRedisGetByName",
			run:  testRedisGetByName(redisDb),
		},
		{
			name: "testRedisListObjects",
			run:  testRedisListObjects(redisDb),
		},
		{
			name: "testRedisDeleteObject",
			run:  testRedisDeleteObject(redisDb),
		},
	}

	t.Run("test-redis", func(t *testing.T) {
		for _, testcase := range testcases {
			tc := testcase
			t.Run(tc.name, func(t *testing.T) {
				err = redisDb.FlushAllRedisDb(ctx)
				if err != nil {
					t.Fatal(err)
				}
				tc.run(t)
			})
		}
	})

}

func (r *testRedis) storeObjects() *testRedis {
	r.responseObjects = nil
	for _, object := range r.requestObjects {
		err := r.client.Store(context.TODO(), object)
		if err != nil {
			r.t.Fatal(err)
		}
	}

	return r
}

func (r *testRedis) getObjectsByID() *testRedis {
	r.responseObjects = nil

	for _, object := range r.requestObjects {
		response, err := r.client.GetObjectByID(context.TODO(), object.GetID())
		if err != nil {
			r.t.Fatal(err)
		}

		r.responseObjects = append(r.responseObjects, response)
	}

	return r
}

func (r *testRedis) getObjectsByName() *testRedis {
	r.responseObjects = nil

	for _, object := range r.requestObjects {
		response, err := r.client.GetObjectByName(context.TODO(), object.GetName())
		if err != nil {
			r.t.Fatal(err)
		}

		r.responseObjects = append(r.responseObjects, response...)
	}

	return r
}

func (r *testRedis) listObjects() *testRedis {
	r.responseObjects = nil

	response, err := r.client.ListObjects(context.TODO(), r.requestObjects[0].GetKind())
	if err != nil {
		r.t.Fatal(err)
	}
	r.responseObjects = append(r.responseObjects, response...)
	return r

}

func (r *testRedis) deleteObject() *testRedis {
	for _, object := range r.requestObjects {
		err := r.client.DeleteObject(context.TODO(), object.GetID())
		if err != nil {
			r.t.Fatal(err)
		}
	}
	return r

}

func (r *testRedis) thenAssert() {
	comparator := func(a, b Object) bool {
		return a.GetID() < b.GetID()
	}

	if diff := cmp.Diff(r.requestObjects, r.responseObjects, cmpopts.SortSlices(comparator)); diff != "" {
		r.t.Fatal(diff)
	}
}

func (r *testRedis) thenAssertNotFound() {

	for _, object := range r.requestObjects {
		response, err := r.client.GetObjectByID(context.TODO(), object.GetID())
		if err == nil {
			r.t.Fatalf("unexpected object %+v in DB", response)
		}

	}
}

func testRedisGetByID(db *RedisDb) func(*testing.T) {
	return func(t *testing.T) {
		r := &testRedis{client: db, t: t}
		r.givenObjects(1).storeObjects().getObjectsByID().thenAssert()

	}
}

func testRedisGetByName(db *RedisDb) func(*testing.T) {
	return func(t *testing.T) {
		r := &testRedis{client: db, t: t}
		r.givenObjects(1).storeObjects().getObjectsByName().thenAssert()

	}
}

func testRedisListObjects(db *RedisDb) func(*testing.T) {
	return func(t *testing.T) {
		r := &testRedis{client: db, t: t}
		r.givenObjects(5).storeObjects().listObjects().thenAssert()

	}
}

func testRedisDeleteObject(db *RedisDb) func(*testing.T) {
	return func(t *testing.T) {
		r := &testRedis{client: db, t: t}
		r.givenObjects(1).storeObjects().getObjectsByID().thenAssert()
		r.deleteObject().thenAssert()

	}
}
