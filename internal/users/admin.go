package users

import (
	"context"
	"errors"
	"fmt"
	"goschool/internal/auth"
	db "goschool/internal/database"
	"goschool/internal/models"

	"github.com/jackc/pgx/v5"
)

func CreateAdmin(ctx context.Context, config *models.Config, d *db.Database) error {
	_, err := d.GetUserByEmail(ctx, config.ADMIN_EMAIL)
	if err == nil {
		return nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check admin: %w", err)
	}

	if config.ADMIN_EMAIL == "" || config.ADMIN_PASSWORD == "" || config.ADMIN_NAME == "" {
		return nil
	}

	AdminPassword, err := auth.HashPassword(config.ADMIN_PASSWORD)
	if err != nil {
		return err
	}
	err = d.AddUser(ctx, &models.User{Email: config.ADMIN_EMAIL, Name: config.ADMIN_NAME, Password: AdminPassword, Role: models.RoleAdmin})
	if err != nil {
		return err
	}
	return nil
}
