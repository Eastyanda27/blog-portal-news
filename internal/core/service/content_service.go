package service

import (
	"blognewsportal/internal/adapter/repository"
	"blognewsportal/internal/core/domain/entity"
	"context"

	"github.com/gofiber/fiber/v2/log"
)

type ContentService interface {
	GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	EditContentByID(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentService struct {
	contentRepository repository.ContentRepository
}

func (c *contentService) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepository.CreateContent(ctx, req)
	if err != nil {
		code = "[SERVICE] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentService) DeleteContent(ctx context.Context, id int64) error {
	err = c.contentRepository.DeleteContent(ctx, id)
	if err != nil {
		code = "[SERVICE] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentService) EditContentByID(ctx context.Context, req entity.ContentEntity) error {
	err = c.contentRepository.EditContentByID(ctx, req)
	if err != nil {
		code = "[SERVICE] EditContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentService) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	result, err := c.contentRepository.GetContentByID(ctx, id)
	if err != nil {
		code = "[SERVICE] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func (c *contentService) GetContents(ctx context.Context, query entity.QueryString) ([]entity.ContentEntity, error) {
	result, err := c.contentRepository.GetContents(ctx, query)
	if err != nil {
		code = "[SERVICE] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	return result, nil
}

func NewContentService(contentRepo repository.ContentRepository) ContentService {
	return &contentService{contentRepository: contentRepo}
}
