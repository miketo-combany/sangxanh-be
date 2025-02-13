package query

import "go.mongodb.org/mongo-driver/bson"

type Q bson.M

func Query() Q {
	return make(Q)
}

func (q Q) Like(name string, value string) Q {
	if value == "" {
		return q
	}
	q[name] = bson.M{"$regex": value, "$options": "i"}
	return q
}

func (q Q) Eq(name string, value interface{}) Q {
	if value == "" {
		return q
	}
	q[name] = value
	return q
}

func (q Q) BSON() bson.M {
	return bson.M(q)
}
