package user

import (
	"todo-api-go/internal/application/usecase/user"
	"todo-api-go/internal/interface/api/response"

	userDTO "todo-api-go/internal/application/dto/user"
	errorHandler "todo-api-go/internal/interface/api/error"

	"github.com/gin-gonic/gin"
)

type RegisterUserHandler struct {
	registerUserUseCase *user.RegisterUserUseCase
}

func NewRegisterUserHandler(registerUserUseCase *user.RegisterUserUseCase) *RegisterUserHandler {
	return &RegisterUserHandler{
		registerUserUseCase: registerUserUseCase,
	}
}

func (h *RegisterUserHandler) Handle(c *gin.Context) {
	// リクエストの解析
	var input userDTO.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "リクエストが正しくありません", err)
		return
	}

	// ユースケースの実行
	output, err := h.registerUserUseCase.Execute(input)
	if err != nil {
		// エラーハンドリング
		switch {
		case errorHandler.IsDuplicateEmailError(err):
			response.Conflict(c, "このメールアドレスはすでに使用されています", err)
		case errorHandler.IsValidationError(err):
			response.BadRequest(c, "入力値が無効です", err)
		default:
			response.InternalServerError(c, "ユーザー登録に失敗しました", err)
		}
		return
	}

	response.Created(c, "ユーザー登録が完了しました", output)
}
