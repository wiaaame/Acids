package repository

import (
	"errors"

	"acids/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepo - репозиторий пользователей (все операции через ORM).
type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// FindByID - получить пользователя по ID.
func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	err := r.db.
		Preload("Acids").
		Preload("Likes").
		First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByUsername - поиск по логину (для аутентификации).
func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.
		Where("username = ?", username).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByEmail - поиск по email (для аутентификации / регистрации).
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.
		Where("email = ?", email).
		First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ExistsByUsername - проверка занятости логина.
func (r *UserRepo) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("username = ?", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ExistsByEmail - проверка занятости email.
func (r *UserRepo) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Create - регистрация нового пользователя (пароль хешируется bcrypt).
func (r *UserRepo) Create(username, email, password string) (*models.User, error) {
	if exists, err := r.ExistsByUsername(username); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.New("username already taken")
	}
	if exists, err := r.ExistsByEmail(email); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.New("email already registered")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}
	if err := r.db.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// Authenticate - проверка логина и пароля.
func (r *UserRepo) Authenticate(username, password string) (*models.User, error) {
	u, err := r.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}

// UpdateEmail - сменить email пользователя.
func (r *UserRepo) UpdateEmail(id uint, email string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("email", email).Error
}

// UpdatePassword - сменить пароль (хеш пересчитывается).
func (r *UserRepo) UpdatePassword(id uint, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return r.db.Model(&models.User{}).
		Where("id = ?", id).
		Update("password_hash", string(hash)).Error
}

// ListAll - все пользователи (для админки / отладки).
func (r *UserRepo) ListAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("id ASC").Find(&users).Error
	return users, err
}