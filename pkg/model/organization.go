package model

type Organization struct {
	Model `bson:",inline"`
	Name  string `json:"name" bson:"name"`
	Owner string `json:"owner" bson:"owner"`
}
