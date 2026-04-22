package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/silvioubaldino/ilia-users/internal/domain"
	"github.com/silvioubaldino/ilia-users/pkg/apperrors"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, apperrors.ErrConflict
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperrors.ErrNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperrors.ErrNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, updates domain.User) (domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, apperrors.ErrNotFound
		}
		return domain.User{}, err
	}
	if err := r.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		if isUniqueViolation(err) {
			return domain.User{}, apperrors.ErrConflict
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&domain.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}

func isUniqueViolation(err error) bool {
	return err != nil && errors.Is(err, gorm.ErrDuplicatedKey)
}
