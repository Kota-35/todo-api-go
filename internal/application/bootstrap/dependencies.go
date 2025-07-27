package bootstrap

import (
	"todo-api-go/internal/application/usecase/auth"
	"todo-api-go/internal/application/usecase/project"
	"todo-api-go/internal/application/usecase/team"
	"todo-api-go/internal/application/usecase/user"
	"todo-api-go/internal/config"
	"todo-api-go/internal/domain/repository"
	"todo-api-go/internal/domain/security"
	"todo-api-go/internal/domain/valueobject"
	persistenceRepo "todo-api-go/internal/infrastructure/persistence/repository"
	sessionHandler "todo-api-go/internal/interface/api/handler/session"
	teamHandler "todo-api-go/internal/interface/api/handler/team"
	projectHandler "todo-api-go/internal/interface/api/handler/team/project"
	userHandler "todo-api-go/internal/interface/api/handler/user"
	"todo-api-go/internal/interface/api/middleware"
)

// Repositories は全てのリポジトリを集約
type Repositories struct {
	User    repository.UserRepository
	Session repository.SessionRepository
	Team    repository.TeamRepository
	Project repository.ProjectRepository
}

// Security は認証関連のコンポーネントを集約
type Security struct {
	JWTGenerator          security.JWTGenerator
	RefreshTokenGenerator security.RefreshTokenGenerator
	Pepper                valueobject.Pepper
	AuthMiddleware        *middleware.AuthMiddleware
}

// UseCases は全てのユースケースを集約
type UseCases struct {
	// User
	RegisterUser *user.RegisterUserUseCase

	// Auth
	AuthenticateUser  *auth.AuthenticateUserUseCase
	GetCurrentSession *auth.GetCurrentSessionUseCase
	RefreshSession    *auth.RefreshSessionUseCase
	Logout            *auth.LogoutUseCase

	// Team
	CreateTeam  *team.CreateTeamUseCase
	FindMyTeams *team.FindMyTeamsUseCase

	// Project
	GetProjects   *project.GetProjectsUseCase
	CreateProject *project.CreateProjectUseCase
}

// Handlers は全てのハンドラーを集約
type Handlers struct {
	User    *userHandler.UserHandler
	Session *sessionHandler.SessionHandler
	Team    *teamHandler.TeamHandler
}

// Application は全ての依存関係を集約
type Application struct {
	Config       *config.Config
	Repositories *Repositories
	Security     *Security
	UseCases     *UseCases
	Handlers     *Handlers
}

// NewRepositories はリポジトリ層を初期化
func NewRepositories() *Repositories {
	return &Repositories{
		User:    persistenceRepo.NewUserRepository(),
		Session: persistenceRepo.NewSessionRepository(),
		Team:    persistenceRepo.NewTeamRepository(),
		Project: persistenceRepo.NewProjectRepository(),
	}
}

// NewSecurity は認証関連コンポーネントを初期化
func NewSecurity(cfg config.Config, repos *Repositories) *Security {
	jwtSecret := cfg.JwtSecret
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // 開発用デフォルト値
	}

	jwtGenerator := security.NewJWTGenerator(jwtSecret)
	refreshTokenGenerator := security.NewRefreshTokenGenerator()
	pepper := valueobject.Pepper([]byte("default-pepper-key"))

	authMiddleware := middleware.NewAuthMiddleware(
		jwtGenerator,
		repos.Session,
	)

	return &Security{
		JWTGenerator:          jwtGenerator,
		RefreshTokenGenerator: refreshTokenGenerator,
		Pepper:                pepper,
		AuthMiddleware:        authMiddleware,
	}
}

// NewUseCases はユースケース層を初期化
func NewUseCases(repos *Repositories, security *Security) *UseCases {
	return &UseCases{
		// User
		RegisterUser: user.NewRegisterUserUseCase(repos.User),

		// Auth
		AuthenticateUser: auth.NewAuthenticateUserUseCase(
			repos.User,
			repos.Session,
			security.JWTGenerator,
			security.RefreshTokenGenerator,
			security.Pepper,
		),
		GetCurrentSession: auth.NewGetCurrentSessionUseCase(
			security.JWTGenerator,
			repos.User,
			repos.Session,
		),
		RefreshSession: auth.NewRefreshSessionUseCase(
			repos.Session,
			repos.User,
			security.JWTGenerator,
		),
		Logout: auth.NewLogoutUserCase(
			security.Pepper,
			repos.Session,
		),

		// Team
		CreateTeam:  team.NewCreateTeamUseCase(repos.Team),
		FindMyTeams: team.NewFindMyTeamsUseCase(repos.Team),

		// Project
		GetProjects:   project.NewGetProjectsUseCase(repos.Project),
		CreateProject: project.NewCreateProjectUseCase(repos.Project),
	}
}

// NewHandlers はハンドラー層を初期化
func NewHandlers(useCases *UseCases, security *Security) *Handlers {
	// User handlers
	registerUserHandler := userHandler.NewRegisterUserHandler(
		useCases.RegisterUser,
	)
	userHandlerInstance := userHandler.NewUserHandler(registerUserHandler)

	// Session handlers
	loginUserHandler := sessionHandler.NewLoginUserHandler(
		useCases.AuthenticateUser,
	)
	getCurrentSessionHandler := sessionHandler.NewGetCurrentSessionHandler(
		useCases.GetCurrentSession,
	)
	sessionRefreshHandler := sessionHandler.NewSessionRefreshHandler(
		&security.Pepper,
		useCases.RefreshSession,
	)
	logoutUserHandler := sessionHandler.NewLogoutUserHandler(
		security.JWTGenerator,
		useCases.Logout,
	)
	sessionHandlerInstance := sessionHandler.NewSessionHandler(
		loginUserHandler,
		getCurrentSessionHandler,
		sessionRefreshHandler,
		logoutUserHandler,
		security.AuthMiddleware,
	)

	// Team handlers
	teamCreateHandler := teamHandler.NewCreateTeamHandler(useCases.CreateTeam)
	findMyTeamsHandler := teamHandler.NewFindMyTeamsHandler(
		useCases.FindMyTeams,
	)

	// Project handlers
	listProjectsHandler := projectHandler.NewListProjectsHandler(
		useCases.GetProjects,
	)
	createProjectHandler := projectHandler.NewCraeteProjectHandler(
		useCases.CreateProject,
	)
	projectHandlerInstance := projectHandler.NewProjectHandler(
		listProjectsHandler,
		createProjectHandler,
	)

	teamHandlerInstance := teamHandler.NewTeamHandler(
		teamCreateHandler,
		security.AuthMiddleware,
		findMyTeamsHandler,
		projectHandlerInstance,
	)

	return &Handlers{
		User:    userHandlerInstance,
		Session: sessionHandlerInstance,
		Team:    teamHandlerInstance,
	}
}

// NewApplication はアプリケーション全体を初期化
func NewApplication() *Application {
	cfg := config.LoadEnv()
	repos := NewRepositories()
	security := NewSecurity(cfg, repos)
	useCases := NewUseCases(repos, security)
	handlers := NewHandlers(useCases, security)

	return &Application{
		Config:       &cfg,
		Repositories: repos,
		Security:     security,
		UseCases:     useCases,
		Handlers:     handlers,
	}
}
