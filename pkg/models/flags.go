package models

// GetFlags returns all flags.
// func (db *MongoDatabase) GetFlags(t Tenant) ([]Flag, error) {
// 	var flags []Flag

// 	err := db.C("flags").Find(bson.M{"tenant": t.ID}).All(&flags)
// 	if err != nil {
// 		return flags, err
// 	}
// 	return flags, nil
// }
