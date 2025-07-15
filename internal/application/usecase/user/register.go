package user

import (
	"fmt"
	userDTO "todo-api-go/internal/application/dto/user"
	"todo-api-go/internal/domain/entity"
	"todo-api-go/internal/domain/repository"
)

type RegisterUserUseCase struct {
	userRepo repository.UserRepository
}

func NewRegisterUserUseCase(
	userRepo repository.UserRepository,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo: userRepo,
	}
}

func (uc *RegisterUserUseCase) Execute(
	input userDTO.RegisterUserInput,
) (*userDTO.RegisterUserOutput, error) {

	// ドメインエンティティの作成
	user, err := entity.NewUser(input.Email, input.Username, input.Password)
	if err != nil {
		return nil, fmt.Errorf("[RegisterUserUseCase]ユーザーのドメインエンティティの作成に失敗しました: %w", err)
	}

	// 永続化
	if err := uc.userRepo.Save(user); err != nil {
		return nil, fmt.Errorf("[RegisterUserUseCase]ユーザーの永続化に失敗しました: %w", err)
	}

	// 出力DTOの作成
	return &userDTO.RegisterUserOutput{
		UserID:   user.ID(),
		Email:    user.Email().String(),
		Username: user.Username(),
		Message:  "ユーザー登録が完了しました",
	}, nil
}
