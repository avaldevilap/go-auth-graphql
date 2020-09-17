package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/avaldevilap/go-auth/graph/generated"
	"github.com/avaldevilap/go-auth/graph/model"
	"github.com/avaldevilap/go-auth/pkg/jwt"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user model.User
	user.Email = input.Email
	user.Password = input.Password
	user.Create()
	token, err := jwt.GenerateToken(user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user model.User
	user.Email = input.Email
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		return "", &model.WrongEmailOrPassword{}
	}
	token, err := jwt.GenerateToken(user.Email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	email, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(email)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var resultUsers []*model.User
	var dbUsers []model.User
	dbUsers, _ = model.GetAll()
	for _, user := range dbUsers {
		resultUsers = append(resultUsers, &model.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email, Password: user.Password})
	}
	return resultUsers, nil
	// panic(fmt.Errorf("not implemented"))
}

func (r *userResolver) ID(ctx context.Context, obj *model.User) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
