package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/silvioubaldino/ilia-users/internal/adapter/client/wallet"
	"github.com/silvioubaldino/ilia-users/internal/adapter/http/handler"
	postgresrepo "github.com/silvioubaldino/ilia-users/internal/adapter/repository/postgres"
	"github.com/silvioubaldino/ilia-users/internal/infrastructure/config"
	"github.com/silvioubaldino/ilia-users/internal/usecase"
	"gorm.io/gorm"
)

func SetupUser(db *gorm.DB, cfg *config.Config, public gin.IRouter, auth gin.IRouter) {
	repo := postgresrepo.NewUserRepository(db)
	_ = wallet.NewClient(cfg.WalletBaseURL, cfg.JWTInternalSecret)

	createUC := usecase.NewCreateUser(repo)
	listUC := usecase.NewListUsers(repo)
	getUC := usecase.NewGetUser(repo)
	updateUC := usecase.NewUpdateUser(repo)
	deleteUC := usecase.NewDeleteUser(repo)
	authenticateUC := usecase.NewAuthenticateUser(repo, cfg.JWTSecret)

	userHandler := handler.NewUserHandler(createUC, listUC, getUC, updateUC, deleteUC)
	authHandler := handler.NewAuthHandler(authenticateUC)

	public.POST("/users", userHandler.Create)
	public.POST("/auth", authHandler.Login)

	auth.GET("/users", userHandler.List)
	auth.GET("/users/:id", userHandler.Get)
	auth.PATCH("/users/:id", userHandler.Update)
	auth.DELETE("/users/:id", userHandler.Delete)
}
