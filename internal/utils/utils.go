package utils

import "github.com/mjaliz/deviran/internal/config"

type Repository struct {
	App *config.AppConfig
}

var Repo *Repository

func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

func NewUtils(r *Repository) {
	Repo = r
}
