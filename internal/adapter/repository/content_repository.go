package repository

import (
	"blognewsportal/internal/core/domain/entity"
	"blognewsportal/internal/core/domain/model"
	"context"
	"strings"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ContentRepository interface {
	GetContents(ctx context.Context) ([]entity.ContentEntity, error)
	GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error)
	CreateContent(ctx context.Context, req entity.ContentEntity) error
	EditContentByID(ctx context.Context, req entity.ContentEntity) error
	DeleteContent(ctx context.Context, id int64) error
}

type contentRepository struct {
	db *gorm.DB
}

func (c *contentRepository) CreateContent(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Create(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] CreateContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentRepository) DeleteContent(ctx context.Context, id int64) error {
	err = c.db.Where("id = ?", id).Delete(&model.Content{}).Error
	if err != nil {
		code = "[REPOSITORY] DeleteContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentRepository) EditContentByID(ctx context.Context, req entity.ContentEntity) error {
	tags := strings.Join(req.Tags, ",")
	modelContent := model.Content{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: req.CreatedByID,
	}

	err = c.db.Where("id = ?", req.ID).Updates(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] EditContent - 1"
		log.Errorw(code, err)
		return err
	}

	return nil
}

func (c *contentRepository) GetContentByID(ctx context.Context, id int64) (*entity.ContentEntity, error) {
	var modelContent model.Content

	err = c.db.Where("id = ?", id).Preload(clause.Associations).First(&modelContent).Error
	if err != nil {
		code = "[REPOSITORY] GetContentByID - 1"
		log.Errorw(code, err)
		return nil, err
	}

	tags := strings.Split(modelContent.Tags, ",")
	resp := entity.ContentEntity{
		ID:          modelContent.ID,
		Title:       modelContent.Title,
		Excerpt:     modelContent.Excerpt,
		Description: modelContent.Description,
		Tags:        tags,
		Status:      modelContent.Status,
		CategoryID:  modelContent.CategoryID,
		CreatedByID: modelContent.CreatedByID,
		CreatedAt:   modelContent.CreatedAt,
		Category: entity.CategoryEntity{
			ID:    modelContent.CategoryID,
			Title: modelContent.Category.Title,
			Slug:  modelContent.Category.Slug,
		},
		User: entity.UserEntity{
			ID:   modelContent.User.ID,
			Name: modelContent.User.Name,
		},
	}

	return &resp, nil
}

func (c *contentRepository) GetContents(ctx context.Context) ([]entity.ContentEntity, error) {
	var modelContents []model.Content

	err = c.db.Order("created_at DESC").Preload(clause.Associations).Find(&modelContents).Error
	if err != nil {
		code = "[REPOSITORY] GetContents - 1"
		log.Errorw(code, err)
		return nil, err
	}

	resps := []entity.ContentEntity{}
	for _, val := range modelContents {
		tags := strings.Split(val.Tags, ",")
		resp := entity.ContentEntity{
			ID:          val.ID,
			Title:       val.Title,
			Excerpt:     val.Excerpt,
			Description: val.Description,
			Tags:        tags,
			Status:      val.Status,
			CategoryID:  val.CategoryID,
			CreatedByID: val.CreatedByID,
			CreatedAt:   val.CreatedAt,
			Category: entity.CategoryEntity{
				ID:    val.CategoryID,
				Title: val.Category.Title,
				Slug:  val.Category.Slug,
			},
			User: entity.UserEntity{
				ID:   val.User.ID,
				Name: val.User.Name,
			},
		}

		resps = append(resps, resp)
	}

	return resps, nil
}

func NewContentRepository(db *gorm.DB) ContentRepository {
	return &contentRepository{
		db: db,
	}
}
