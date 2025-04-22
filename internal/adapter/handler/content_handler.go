package handler

import (
	"blognewsportal/internal/adapter/handler/request"
	"blognewsportal/internal/adapter/handler/response"
	"blognewsportal/internal/core/domain/entity"
	"blognewsportal/internal/core/service"
	"blognewsportal/lib/conv"
	validatorLib "blognewsportal/lib/validator"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	EditContentByID(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error

	GetContentWithQuery(c *fiber.Ctx) error
	GetContentDetail(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

func (ch *contentHandler) GetContentDetail(c *fiber.Ctx) error {
	idParam := c.Params("contentID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetContentDetail - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetContentDetail - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	contentResponse := response.SuccessContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CreatedByID,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Categories fetched detail successfully"
	defaultSuccessResponse.Data = contentResponse

	return c.JSON(defaultSuccessResponse)
}

func (ch *contentHandler) GetContentWithQuery(c *fiber.Ctx) error {
	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[HANDLER] GetContentWithQuery - 1"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		if err != nil {
			code = "[HANDLER] GetContentWithQuery - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid limit number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	categoryID := 0
	if c.Query("categoryID") != "" {
		categoryID, err = conv.StringToInt(c.Query("categoryID"))
		if err != nil {
			code := "[HANDLER] GetContentWithQuery - 3"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid Category ID"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	reqEntity := entity.QueryString{
		Limit:      limit,
		Page:       page,
		OrderBy:    orderBy,
		OrderType:  orderType,
		Search:     search,
		Status:     "Publish",
		CategoryID: int64(categoryID),
	}

	results, totalData, totalPages, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] GetContentWithQuery - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	contentResponses := []response.SuccessContentResponse{}
	for _, content := range results {
		contentResponse := response.SuccessContentResponse{
			ID:           content.ID,
			Title:        content.Title,
			Excerpt:      content.Excerpt,
			Description:  content.Description,
			Image:        content.Image,
			Tags:         content.Tags,
			Status:       content.Status,
			CategoryID:   content.CategoryID,
			CreatedByID:  content.CreatedByID,
			CreatedAt:    content.CreatedAt.Format(time.RFC3339),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}
		contentResponses = append(contentResponses, contentResponse)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Contents fetched successfully"
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         page,
		PerPage:      limit,
		TotalPages:   int(totalPages),
	}
	defaultSuccessResponse.Data = contentResponses

	return c.JSON(defaultSuccessResponse)
}

func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	var req request.ContentRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] CreateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code = "[HANDLER] CreateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = ch.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content created successfully"
	defaultSuccessResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccessResponse)
}

func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] DeleteContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = ch.contentService.DeleteContent(c.Context(), id)
	if err != nil {
		code = "[HANDLER] DeleteCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content deleted successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (ch *contentHandler) EditContentByID(c *fiber.Ctx) error {
	var req request.ContentRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] EditContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] EditContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code = "[HANDLER] EditContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()
	}

	idParam := c.Params("contentID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] EditContentByID - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		ID:          id,
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = ch.contentService.EditContentByID(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] EditContentByID - 5"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Content updated successfully"
	defaultSuccessResponse.Data = nil

	return c.JSON(defaultSuccessResponse)
}

func (ch *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	contentResponse := response.SuccessContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CreatedByID,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Categories fetched detail successfully"
	defaultSuccessResponse.Data = contentResponse

	return c.JSON(defaultSuccessResponse)
}

func (ch *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetContents - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	reqEntity := entity.QueryString{
		Limit:      0,
		Page:       0,
		OrderBy:    "",
		OrderType:  "",
		Search:     "",
		CategoryID: 0,
	}

	results, totalData, totalPage, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] GetContents - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	contentResponses := []response.SuccessContentResponse{}
	for _, content := range results {
		contentResponse := response.SuccessContentResponse{
			ID:           content.ID,
			Title:        content.Title,
			Excerpt:      content.Excerpt,
			Description:  content.Description,
			Image:        content.Image,
			Tags:         content.Tags,
			Status:       content.Status,
			CategoryID:   content.CategoryID,
			CreatedByID:  content.CreatedByID,
			CreatedAt:    content.CreatedAt.Format(time.RFC3339),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}
		contentResponses = append(contentResponses, contentResponse)
	}

	defaultSuccessResponse.Meta.Status = true
	defaultSuccessResponse.Meta.Message = "Contents fetched successfully"
	defaultSuccessResponse.Pagination = &response.PaginationResponse{
		TotalRecords: int(totalData),
		Page:         0,
		PerPage:      0,
		TotalPages:   int(totalPage),
	}
	defaultSuccessResponse.Data = contentResponses

	return c.JSON(defaultSuccessResponse)
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{contentService: contentService}
}
