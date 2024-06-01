package user

import (
	"context"
	"fmt"

	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/entity"
	"csv-analyzer-api/internal/util"
	"csv-analyzer-api/internal/value"
)

type UserRepository interface {
	Create(context.Context, entity.User) error
	GetByID(context.Context, value.UserID) (*entity.User, error)
	GetByEmail(context.Context, string) (*entity.User, error)
	Delete(context.Context, value.UserID) error
	Update(context.Context, *entity.User) error
}

type UserService interface {
	Create(context.Context, *CreateArg) (*entity.User, error)
	GetByID(context.Context, value.UserID) (*entity.User, error)
	GetByEmail(context.Context, string) (*entity.User, error)
	Delete(context.Context, *DeleteArg) error
	Update(context.Context, *entity.User) error
}

type Service struct {
	config         *config.Configuration
	userRepository UserRepository
}

func NewService(
	config *config.Configuration,
	userRepository UserRepository,
) *Service {
	return &Service{
		config:         config,
		userRepository: userRepository,
	}
}

type CreateArg struct {
	Email    string
	Name     string
	Password string
	Role     value.Role
}

func (s *Service) Create(ctx context.Context, arg *CreateArg) (*entity.User, error) {
	hashedPassword, err := util.HashPassword(arg.Password)
	if err != nil {
		return nil, err
	}

	user := entity.User{
		Email:          arg.Email,
		Name:           arg.Name,
		HashedPassword: hashedPassword,
		Role:           arg.Role,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed user create: %w", err)
	}

	return &user, nil
}

func (s *Service) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed userRepository.Get: %w", err)
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id value.UserID) (*entity.User, error) {
	user, err := s.userRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed userRepository.Get: %w", err)
	}

	return user, nil
}

type DeleteArg struct {
	ID value.UserID
}

func (s *Service) Delete(ctx context.Context, arg *DeleteArg) error {
	err := s.userRepository.Delete(ctx, arg.ID)
	if err != nil {
		return fmt.Errorf("failed userRepository.Delete: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, arg *entity.User) error {
	//updatedAt := s.clock.Now().UTC()
	//user.UpdatedAt = &updatedAt

	err := s.userRepository.Update(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed userRepository.Update: %w", err)
	}

	return nil
}
