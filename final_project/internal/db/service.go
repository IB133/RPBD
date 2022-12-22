package db

import (
	"context"

	"github.com/IB133/RPBD/final_project/internal/models"
)

type Service struct {
	db *Store
}

func (s *Service) Registration(ctx context.Context, u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.HashPassword(); err != nil {
		return err
	}
	if err := s.db.addUSer(ctx, u); err != nil {
		return err
	}
	return nil
}
