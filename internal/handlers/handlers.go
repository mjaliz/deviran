package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjaliz/deviran/internal/config"
	"github.com/mjaliz/deviran/internal/models"
	"net/http"
)

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) SignUp(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(user); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}
