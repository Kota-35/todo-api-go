package repository

import (
	"context"
	"fmt"
	"strings"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
	valueobject "todo-api-go/internal/domain/valueobject"
	"todo-api-go/pkg/database"

	domainError "todo-api-go/internal/domain/error"

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
			if strings.Contains(err.Error(), "Unique constraint failed on the fields: (`email`)") {
				return domainError.NewDuplicateEmailError(user.Email().String())
			}
			return fmt.Errorf("[userRepository]ユーザーの作成に失敗しました: %w", err)
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
			return fmt.Errorf("[userRepository]ユーザーの更新に失敗しました: %w", err)
		}
	}

	return nil

}

func (r *userRepository) FindByEmail(email valueobject.Email) (*entity.User, error) {
	ctx := context.Background()

	user, err := r.client.User.FindUnique(
		db.User.Email.Equals(email.String()),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("[userRepository.FindByEmail]ユーザーの取得に失敗しました: %w", err)
	}

	reconstructUser, err := entity.ReconstructUser(
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("[[userRepository.FindByEmail]ドメインモデルの変換に失敗しました: %w", err)
	}

	return reconstructUser, nil
}
