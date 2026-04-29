package users

import (
	baseModels "bot-engine/models/mongo"
	"fmt"
)

type User struct {
	baseModels.BaseModel `bson:",inline"`

	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}

func (u *User) Validate() error {

	if u.Name == "" {
		return fmt.Errorf("Name is Required")
	}

	if u.Age > 100 || u.Age < 10 {
		return fmt.Errorf("Invalid Age")
	}

	return nil

}
