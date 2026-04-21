package bot

import (
	baseModels "di/models/mongo"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	baseModels.BaseModel `bson:",inline"`

	OwnerID bson.ObjectID `json:"owner_id",bson:"owner_id"`

	Name string `json:"name" bson:"name"`
	Age  int    `json:"age" bson:"age"`
}

func (u *User) Validate() error {

	if u.Name == "" {
		return fmt.Errorf("Name is Required")
	}

	// if u.Age ==  {
	// 	return fmt.Errorf("Name is Required")
	// }

	return nil

}
