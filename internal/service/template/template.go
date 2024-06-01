package template

import (
	"context"
	"fmt"

	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/entity"
	"csv-analyzer-api/internal/value"
)

type TemplateRepository interface {
	Create(context.Context, *entity.Template) (*entity.Template, error)
	GetByID(context.Context, value.TemplateID) (*entity.Template, error)
	Delete(context.Context, value.TemplateID) error
	Update(context.Context, *entity.Template) error
}

type TemplateService interface {
	Create(context.Context, *entity.Template) (*entity.Template, error)
	GetByID(context.Context, value.TemplateID) (*entity.Template, error)
	Delete(context.Context, *DeleteArg) error
	Update(context.Context, *entity.Template) error
}

type Service struct {
	config             *config.Configuration
	templateRepository TemplateRepository
}

func NewService(
	config *config.Configuration,
	templateRepository TemplateRepository,
) *Service {
	return &Service{
		config:             config,
		templateRepository: templateRepository,
	}
}

func (s *Service) Create(ctx context.Context, template *entity.Template) (*entity.Template, error) {
	template, err := s.templateRepository.Create(ctx, template)
	if err != nil {
		return nil, fmt.Errorf("failed template create: %w", err)
	}

	return template, nil
}

func (s *Service) GetByID(ctx context.Context, id value.TemplateID) (*entity.Template, error) {
	template, err := s.templateRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed templateRepository.Get: %w", err)
	}

	return template, nil
}

type DeleteArg struct {
	ID value.TemplateID
}

func (s *Service) Delete(ctx context.Context, arg *DeleteArg) error {
	err := s.templateRepository.Delete(ctx, arg.ID)
	if err != nil {
		return fmt.Errorf("failed templateRepository.Delete: %w", err)
	}

	return nil
}

func (s *Service) Update(ctx context.Context, arg *entity.Template) error {
	//updatedAt := s.clock.Now().UTC()
	//template.UpdatedAt = &updatedAt

	err := s.templateRepository.Update(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed templateRepository.Update: %w", err)
	}

	return nil
}
