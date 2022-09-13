package user

import "context"

type Storage interface {
	Create(ctx context.Context, user User) (string, error)
	FindById(ctx context.Context, id string) (User, error)
	FindByNumber(ctx context.Context, id string) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}
