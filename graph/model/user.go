package model

import (
	"log"
	"time"

	database "github.com/avaldevilap/go-auth/internal/pkg/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	// ID        string     `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
}

func (user *User) Create() {
	hashedPassword, err := HashPassword(user.Password)
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = hashedPassword

	_, err = database.Collection.InsertOne(database.Ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func (user *User) Authenticate() bool {
	filter := bson.D{primitive.E{Key: "_id", Value: user.Email}}
	var dbUser *User
	err := database.Collection.FindOne(database.Ctx, filter).Decode(&dbUser)
	if err != nil {
		return false
	}

	return CheckPasswordHash(user.Password, dbUser.Password)
}

func GetAll() ([]User, error) {
	filter := bson.D{{}}
	return filterUsers(filter)
}

// HashPassword hashes given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPassword hash compares raw password with it's hashed values
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func filterUsers(filter interface{}) ([]User, error) {
	var users []User

	cur, err := database.Collection.Find(database.Ctx, filter)
	if err != nil {
		return users, err
	}

	for cur.Next(database.Ctx) {
		var u User
		err := cur.Decode(&u)
		if err != nil {
			return users, err
		}

		users = append(users, u)
	}

	if err := cur.Err(); err != nil {
		return users, err
	}

	cur.Close(database.Ctx)

	if len(users) == 0 {
		return users, mongo.ErrNoDocuments
	}
	log.Print(users)

	return users, nil
}
