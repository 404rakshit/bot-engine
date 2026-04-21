package db

import "fmt"

type Database struct {
	dbString string
}

func NewDatabase(dbString string) *Database {
	fmt.Printf("dbString: %s", dbString)
	return &Database{dbString: dbString}
}
