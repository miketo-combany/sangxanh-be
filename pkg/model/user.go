package model

type User struct {
	Model         `bson:",inline"`
	Email         string         `json:"email" bson:"email"`
	Name          string         `json:"name" bson:"name"`
	GivenName     string         `json:"given_name" bson:"given_name"`
	FamilyName    string         `json:"family_name" bson:"family_name"`
	Avatar        string         `json:"avatar" bson:"avatar"`
	Metadata      map[string]any `json:"metadata" bson:"metadata"`
	Organizations []Organization `json:"organizations" bson:"organizations"`
}

func (User) CollectionName() string {
	return "users"
}
