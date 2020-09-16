package users

import (
	"log"
	"time"

	database "github.com/avaldevilap/go-auth/internal/pkg/db/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
}

func (user *User) Create() {
	hashedPassword, err := HashPassword(user.Password)
	user.Password = hashedPassword

	_, err = database.Collection.InsertOne(database.Ctx, user)
	if err != nil {
		log.Fatal(err)
	}
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
