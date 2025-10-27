package part

import (
	"github.com/Bladforceone/rocket/inventory/internal/repository"
	"github.com/Bladforceone/rocket/inventory/internal/service"
)

type serv struct {
	repo repository.PartRepository
}

func NewService(repo repository.PartRepository) service.PartService {
	return &serv{
		repo: repo,
	}
}
