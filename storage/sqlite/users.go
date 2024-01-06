package sqlite

import (
	"fmt"

	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/gorm"
)

func (s *SqliteStorage) CreateUser(user *models.User) error {
	return s.DB.Create(user).Error
}

func (s *SqliteStorage) GetUserById(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("User with id: %d", id))
		}
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewObjectNotFoundError(fmt.Sprintf("User with username: %s", username))
		}
		return nil, err
	}

	return &user, nil
}

func (s *SqliteStorage) GetAllUsersByUsernameFuzzy(username string) ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (s *SqliteStorage) GetUsersByUsernameFuzzy(username string, offset, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username)).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return users, err
	}

	return users, nil
}

func (s *SqliteStorage) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *SqliteStorage) GetUsers(offset, limit int) ([]*models.User, error) {
	var users []*models.User
	if err := s.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *SqliteStorage) GetUsersCount() (int, error) {
	var count int64
	if err := s.DB.Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *SqliteStorage) UpdateUser(user *models.User) error {
	if _, err := s.GetUserById(user.ID); err != nil {
		return err
	}

	return s.DB.Save(user).Error
}

func (s *SqliteStorage) DeleteUserById(id uint) error {
	if id == 1 {
		return errors.NewMainAdminDeletionError("Cannot delete main admin")
	}

	if _, err := s.GetUserById(id); err != nil {
		return err
	}

	return s.DB.Delete(&models.User{}, id).Error
}
