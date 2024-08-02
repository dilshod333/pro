package pkg

type User struct {
	Username     string  `bson:"username" json:"username"`
	Email        string  `bson:"email" json:"email"`
	PasswordHash string  `bson:"password_hash" json:"password_hash"`
	Profile      Profile `bson:"profile" json:"profile"`
}

type Profile struct {
	Name    string `bson:"name" json:"name"`
	Address string `bson:"address" json:"address"`
}
