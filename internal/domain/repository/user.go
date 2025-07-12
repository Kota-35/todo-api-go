package repository

import (
	"todo-api-go/internal/domain/entity"
)

// リポジトリインターフェース
/*
ドメインリポジトリ
- 抽象化: 具体的な実装を持たない
- ドメイン層: ビジネスロジックの観点から定義
- 依存関係の逆転: ドメインが外部技術に依存しない
- テスト容易性: モックを使ったテストが可能

イングラストラクチャリポジトリ
- 具体的実装: 実際のデータアクセス処理
- 技術詳細: Prisma, SQL, NoSQLなどの具体的な技術
- 外部依存: データベースやORMに依存
*/

type UserRepository interface {
	// 基本的なCRUD操作
	Save(user *entity.User) error
	// FindByID(id string) (*entity.User, error)
	// FindByEmail(email valueobject.Email) (*entity.User, error)
	// Delete(id string) error
}
