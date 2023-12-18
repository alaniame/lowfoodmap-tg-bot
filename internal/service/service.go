package service

import "lowfoodmap-tg-bot/internal/repository"

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProduct(name string) (*repository.Product, error) {
	return s.repo.GetProduct(name)
}
