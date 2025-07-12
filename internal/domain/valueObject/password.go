package valueobject

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashedValue string
}

// プレーンパスワードから生成
func NewPassword(plainPassword string) (Password, error) {
	if err := validatePassword(plainPassword); err != nil {
		return Password{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}

	return Password{hashedValue: string(passwordHash)}, nil
}

// ハッシュ済みパスワードから復元(DBから取得時)
func RestorePassword(hashedPassword string) Password {
	return Password{hashedValue: hashedPassword}
}

// パスワード検証
func (p Password) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
	return err == nil
}

// ハッシュ値の取得
func (p Password) Hash() string {
	return p.hashedValue
}

func (p Password) Equals(newPasswordVO Password) bool {
	return p.Hash() == newPasswordVO.Hash()
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("パスワードは8文字以上である必要があります")
	}

	if len(password) > 128 {
		return errors.New("パスワードは128文字以下である必要があります")
	}

	// TODO: 正規表現で以下の条件を実装
	// .regex(/[a-zA-Z]/, { message: '少なくとも1つの英文字を含めてください' })
	// .regex(/[0-9]/, { message: '少なくとも1つの数字を含めてください' })
	// .regex(/[^a-zA-Z0-9]/, {
	//  message: '少なくとも1つの記号(特殊文字)を含めてください',
	// })

	return nil
}
