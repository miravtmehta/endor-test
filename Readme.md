
# Object Store - Endor Labs


## Documentation

Redis is used to store Animal and Person Object. Key is combination of ULID, Name and Kind with deliminator.

Redis key struct 
```
"9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d" : "Tom" : "Animal"
"1b9d6bcd-bbfd-4b2d-9b5d-ab8dfbbd4bed" : "Jim" : "Person"
``` 

Redis Value struct
```
"name":"James","id":"01H03BVGJTBNNBN5F8ZVQ5VDET","last_name":"Bond","birthdate":"1960-01-07T00:00:00Z"
```

The primary reason for choose this structure in order
1. Searches are based on Name, ID, Kind ; easy to abstract information on redis
2. Redis support's pattern  / regex matching in elegant way
3. Scalable if type of objects increase, can partition data to redis-cluster nodes with help of consistent hashing
4. ULID provide lexicology sorting, helps in ordering and sorting information

## Getting Started

This Makefile provides a set of commands to help with the building, running, and cleaning of the endor project. The project includes creating go executables, setting up a Redis container, and running the executables for Mac and Linux. Below are the available commands:

```
make help        Show help for each of the Makefile recipes.
make clean       Removes go executables and  redis container
make build       Runs "clean" and "setup" to build infra and dependencies and create go executables
make setup       Pull latest Redis images and run as container @ port 6739
make run-mac     Run go executable for Mac
make run-linux   Run go executable for Linux
make run-tests   Run unit-tests cases
```

## Commands

### help

The `help` command shows help for each of the Makefile recipes.

### clean

The `clean` command removes the go executables and the Redis container.

### build

The `build` command runs `clean` + `setup`creates the go executables after setting up the Redis container in orderly fashion. 

### setup

The `setup` command pulls the latest Redis image and runs it as a container at port 6379.

### run-mac

The `run-mac` command runs the go executable for Mac that performs CRUD operations on object data present in main.go.

### run-linux

The `run-linux` command runs the go executable for Linux that performs CRUD operations on object data present in main.go.

### run-tests

The `run-tests` command runs unit-test cases for Objects interface

## Build local setup

Run `build` to install all necessary infrastructure and dependent packages and run `make run-linux` or `make run-mac` based on operating system

