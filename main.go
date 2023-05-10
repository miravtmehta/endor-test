package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	AnimalKind = "Animal"
	PersonKind = "Person"
)

func main() {

	ctx := context.Background()
	redisDb, err := NewRedisConnection(ctx, "127.0.0.1:6379", "")
	if err != nil {
		log.Fatalf("---- Redis connection error ---- %s", err.Error())
	}

	err = redisDb.FlushAllRedisDb(ctx)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("*******************************\n")
	fmt.Printf("[INFO] Redis DB flushed out, Good to go.... \n")
	fmt.Printf("*******************************\n")

	objects := []Object{
		&Animal{Name: "Dog", Type: "tamed", OwnerID: "john"},
		&Animal{Name: "Wolf", Type: "un-tamed", OwnerID: "nature"},
		&Animal{Name: "Wolf", Type: "un-tamed", OwnerID: "mountains"},
		&Person{Name: "Rey", LastName: "Mysterious", Birthday: time.DateOnly, BirthDate: time.Date(2006, 01, 02, 0, 0, 0, 0, time.UTC)},
		&Person{Name: "James", LastName: "Bond", BirthDate: time.Date(1960, 01, 7, 0, 0, 0, 0, time.UTC)},
	}

	/*
		Here starts the storage of Animal and person Objects
	*/

	for _, obj := range objects {
		err = redisDb.Store(ctx, obj)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Get object by its Name | GetObjectByName

	for _, obj := range objects {
		objDB, err := redisDb.GetObjectByName(ctx, obj.GetName())
		if err != nil {
			log.Fatalf(err.Error())
		}
		animalData, _ := json.Marshal(objDB)
		fmt.Println(string(animalData))
	}

	// Get object by its ID  | GetObjectByID
	for _, obj := range objects {
		ObjDb, err := redisDb.GetObjectByID(ctx, obj.GetID())
		if err != nil {
			log.Fatalf(err.Error())
		}
		animalData, _ := json.Marshal(ObjDb)
		fmt.Println(string(animalData))
	}

	// List object by its kind  | ListObjects

	fmt.Printf("objects of kind : %s \n", AnimalKind)
	animalKindObjects, _ := redisDb.ListObjects(ctx, AnimalKind)
	for _, value := range animalKindObjects {
		data, _ := json.Marshal(value)
		fmt.Println(string(data))
	}

	fmt.Printf("objects of kind : %s \n", PersonKind)
	personKindObjects, _ := redisDb.ListObjects(ctx, PersonKind)
	for _, value := range personKindObjects {
		data, _ := json.Marshal(value)
		fmt.Println(string(data))
	}

	fmt.Printf("Total count of person object: %d \n", len(personKindObjects))

	// Delete Object
	err = redisDb.DeleteObject(ctx, objects[3].GetID())
	if err != nil {
		log.Fatalf(err.Error())
	}
	zip, _ := redisDb.ListObjects(ctx, PersonKind)

	for _, value := range zip {
		data, _ := json.Marshal(value)
		fmt.Println(string(data))
	}

	fmt.Printf("[POST deletion] Total count of person object: %d \n", len(zip))

}
