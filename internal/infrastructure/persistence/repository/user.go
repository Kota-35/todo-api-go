package repository

import (
	"context"
	"fmt"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/pkg/database"

	"todo-api-go/prisma/db"
)

type userRepository struct {
	client *db.PrismaClient
}

func NewUserRepository() repository.UserRepository {
	return &userRepository{
		client: database.PrismaClient,
	}
}

func (r *userRepository) Save(user *entity.User) error {
	ctx := context.Background()

	if user.IsNew() {
		// 新規作成
		createdUser, err := r.client.User.CreateOne(
			db.User.Email.Set(user.Email().String()),
			db.User.Username.Set(user.Username()),
			db.User.PasswordHash.Set(user.PasswordHash()),
			db.User.IsActive.Set(user.IsActive()),
		).Exec(ctx)

		if err != nil {
			return fmt.Errorf("ユーザーの作成に失敗しました: %w", err)
		}

		return user.SetID(createdUser.ID)
	} else {
		// 更新
		_, err := r.client.User.FindUnique(
			db.User.ID.Equals(user.ID()),
		).Update(
			db.User.Username.Set(user.Username()),
			db.User.PasswordHash.Set(user.PasswordHash()),
			db.User.IsActive.Set(user.IsActive()),
			db.User.UpdatedAt.Set(user.UpdatedAt()),
		).Exec(ctx)

		if err != nil {
			return fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
		}
	}

	return nil

}
