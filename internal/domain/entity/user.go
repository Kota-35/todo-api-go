package entity

import (
	"errors"
	"fmt"
	"time"
	domainErr "todo-api-go/internal/domain/error"

	valueobject "todo-api-go/internal/domain/valueobject"
)

type User struct {
	id           string
	email        valueobject.Email
	username     string
	passwordHash string
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
}

func NewUser(email, username, password string) (*User, error) {
	// Validation
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, domainErr.NewInvalidUserDataError(
			fmt.Sprintf("メールアドレスに問題があります: %s", err.Error()),
		)
	}

	passwordVO, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, domainErr.NewInvalidUserDataError(
			fmt.Sprintf("パスワードに問題があります: %s", err.Error()),
		)
	}

	if len(username) < 3 {
		return nil, domainErr.NewInvalidUserDataError(
			fmt.Sprintf("ユーザー名に問題があります: %s", " ユーザー名は3文字以上である必要があります"),
		)
	}

	/*
		新規作成時はEntityではIDを空文字にする, リポジトリで設定
		理由は以下:

		シンプルさ: エンティティの責務が明確
		柔軟性: 異なるデータストアでも対応可能
		テストしやすさ: IDの生成タイミングを制御可能
		Prismaとの親和性: Prismaの自動生成機能を活用
	*/
	return &User{
		id:           "",
		email:        emailVO,
		username:     username,
		passwordHash: string(passwordVO.Hash()),
		isActive:     true,
		createdAt:    time.Now(),
		updatedAt:    time.Now(),
	}, nil

}

// 復元用コンストラクタ(DB取得時)
func ReconstructUser(id string, email string, username string, passwordHash string, isActive bool, createdAt time.Time, updatedAt time.Time) (*User, error) {
	emailVO, err := valueobject.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           id,
		email:        emailVO,
		username:     username,
		passwordHash: passwordHash,
		isActive:     isActive,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}, nil
}

// パスワードの変更
func (u *User) ChangePassword(currentPassword, newPassword string) error {
	// 1. アクティブ状態の確認
	if !u.isActive {
		return domainErr.NewInvalidUserDataError("非アクティブなユーザーはパスワードを変更できません")
	}

	// 2. 現在のパスワードを検証
	currentPasswordVO := valueobject.RestorePassword(u.passwordHash)
	if !currentPasswordVO.Verify(currentPassword) {
		return domainErr.NewInvalidUserDataError("パスワードが正しくありません")
	}

	// 3. 新しいパスワードの作成
	newPasswordVO, err := valueobject.NewPassword(newPassword)
	if err != nil {
		return domainErr.NewInvalidUserDataError(
			fmt.Sprintf("新しいパスワードに問題があります: %s", err),
		)
	}

	// 4. 同じパスワードかチェック
	if currentPasswordVO.Equals(newPasswordVO) {
		return domainErr.NewInvalidUserDataError(
			"新しいパスワードは現在のパスワードと異なる必要があります.",
		)
	}

	u.passwordHash = newPasswordVO.Hash()
	u.updatedAt = time.Now()
	return nil
}

func (u *User) IsNew() bool {
	return u.id == ""
}

func (u *User) SetID(id string) error {
	if u.id != "" {
		return errors.New("IDはすでに設定されています")
	}
	if id == "" {
		return errors.New("IDは空にできません")
	}

	u.id = id
	return nil
}

func (u *User) ID() string {
	return u.id
}

func (u *User) Email() valueobject.Email {
	return u.email
}

func (u *User) Username() string {
	return u.username
}

func (u *User) PasswordHash() string {
	return u.passwordHash
}

func (u *User) IsActive() bool {
	return u.isActive
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// パスワードがあっているか検証
func (u *User) VerifyPassword(plainPassword string) bool {
	collectPassword := valueobject.RestorePassword(u.passwordHash)

	return collectPassword.Verify(plainPassword)
}
