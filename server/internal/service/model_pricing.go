package service

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"

	"gorm.io/gorm"
)

// ModelPricingReader 定义定价管理读取所需的最小能力。
type ModelPricingReader interface {
	List(context.Context, *gorm.DB) ([]entity.ModelPricing, error)
	GetByID(context.Context, *gorm.DB, uint) (entity.ModelPricing, error)
}

// ModelPricingWriter 定义定价管理写入所需的最小能力。
type ModelPricingWriter interface {
	Create(context.Context, *gorm.DB, *entity.ModelPricing) error
	Update(context.Context, *gorm.DB, *entity.ModelPricing) error
	HasPlaceholder(context.Context, *gorm.DB, uint) (bool, error)
}

// ConflictError 用于把资源冲突映射成 409。
type ConflictError struct {
	message string
}

func (conflictError ConflictError) Error() string {
	return conflictError.message
}

// NotFoundError 用于把资源不存在映射成 404。
type NotFoundError struct {
	message string
}

func (notFoundError NotFoundError) Error() string {
	return notFoundError.message
}

// ModelPricingService 承载模型定价管理能力。
type ModelPricingService struct {
	db     *gorm.DB
	reader ModelPricingReader
	writer ModelPricingWriter
}

// NewModelPricingService 创建模型定价管理服务。
func NewModelPricingService(db *gorm.DB) *ModelPricingService {
	return &ModelPricingService{
		db:     db,
		reader: repository.NewModelPricingRepository(),
		writer: repository.NewModelPricingRepository(),
	}
}

// List 返回全部模型定价记录。
func (service *ModelPricingService) List(ctx context.Context) ([]dto.ModelPricingResponse, error) {
	pricings, err := service.reader.List(ctx, service.db)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ModelPricingResponse, 0, len(pricings))
	for _, pricing := range pricings {
		result = append(result, dto.ModelPricingResponse{
			ID:                          pricing.ID,
			ModelID:                     pricing.ModelID,
			DisplayName:                 pricing.DisplayName,
			InputCostPerMillion:         pricing.InputCostPerMillion,
			OutputCostPerMillion:        pricing.OutputCostPerMillion,
			CacheReadCostPerMillion:     pricing.CacheReadCostPerMillion,
			CacheCreationCostPerMillion: pricing.CacheCreationCostPerMillion,
			IsPlaceholder:               pricing.IsPlaceholder,
			CreatedAt:                   pricing.CreatedAt,
			UpdatedAt:                   pricing.UpdatedAt,
		})
	}

	return result, nil
}

// Create 创建一条新的模型定价记录。
func (service *ModelPricingService) Create(
	ctx context.Context,
	request dto.ModelPricingUpsertRequest,
) (dto.ModelPricingResponse, error) {
	pricing, err := normalizeModelPricingRequest(request)
	if err != nil {
		return dto.ModelPricingResponse{}, err
	}

	if pricing.IsPlaceholder {
		exists, err := service.writer.HasPlaceholder(ctx, service.db, 0)
		if err != nil {
			return dto.ModelPricingResponse{}, err
		}
		if exists {
			return dto.ModelPricingResponse{}, ConflictError{message: "placeholder 默认定价已存在"}
		}
	}

	if err := service.writer.Create(ctx, service.db, &pricing); err != nil {
		return dto.ModelPricingResponse{}, err
	}

	return dto.ModelPricingResponse{
		ID:                          pricing.ID,
		ModelID:                     pricing.ModelID,
		DisplayName:                 pricing.DisplayName,
		InputCostPerMillion:         pricing.InputCostPerMillion,
		OutputCostPerMillion:        pricing.OutputCostPerMillion,
		CacheReadCostPerMillion:     pricing.CacheReadCostPerMillion,
		CacheCreationCostPerMillion: pricing.CacheCreationCostPerMillion,
		IsPlaceholder:               pricing.IsPlaceholder,
		CreatedAt:                   pricing.CreatedAt,
		UpdatedAt:                   pricing.UpdatedAt,
	}, nil
}

// Update 按全量更新语义覆盖已有模型定价。
func (service *ModelPricingService) Update(
	ctx context.Context,
	id uint,
	request dto.ModelPricingUpsertRequest,
) (dto.ModelPricingResponse, error) {
	existing, err := service.reader.GetByID(ctx, service.db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto.ModelPricingResponse{}, NotFoundError{message: "模型定价不存在"}
		}
		return dto.ModelPricingResponse{}, err
	}

	pricing, err := normalizeModelPricingRequest(request)
	if err != nil {
		return dto.ModelPricingResponse{}, err
	}

	if pricing.IsPlaceholder {
		exists, err := service.writer.HasPlaceholder(ctx, service.db, id)
		if err != nil {
			return dto.ModelPricingResponse{}, err
		}
		if exists {
			return dto.ModelPricingResponse{}, ConflictError{message: "placeholder 默认定价已存在"}
		}
	}

	existing.ModelID = pricing.ModelID
	existing.DisplayName = pricing.DisplayName
	existing.InputCostPerMillion = pricing.InputCostPerMillion
	existing.OutputCostPerMillion = pricing.OutputCostPerMillion
	existing.CacheReadCostPerMillion = pricing.CacheReadCostPerMillion
	existing.CacheCreationCostPerMillion = pricing.CacheCreationCostPerMillion
	existing.IsPlaceholder = pricing.IsPlaceholder

	if err := service.writer.Update(ctx, service.db, &existing); err != nil {
		return dto.ModelPricingResponse{}, err
	}

	return dto.ModelPricingResponse{
		ID:                          existing.ID,
		ModelID:                     existing.ModelID,
		DisplayName:                 existing.DisplayName,
		InputCostPerMillion:         existing.InputCostPerMillion,
		OutputCostPerMillion:        existing.OutputCostPerMillion,
		CacheReadCostPerMillion:     existing.CacheReadCostPerMillion,
		CacheCreationCostPerMillion: existing.CacheCreationCostPerMillion,
		IsPlaceholder:               existing.IsPlaceholder,
		CreatedAt:                   existing.CreatedAt,
		UpdatedAt:                   existing.UpdatedAt,
	}, nil
}

func normalizeModelPricingRequest(request dto.ModelPricingUpsertRequest) (entity.ModelPricing, error) {
	modelID := strings.ToLower(strings.TrimSpace(request.ModelID))
	if modelID == "" {
		return entity.ModelPricing{}, ValidationError{message: "model_id 不能为空"}
	}
	if len(modelID) > 128 {
		return entity.ModelPricing{}, ValidationError{message: "model_id 长度不能超过 128"}
	}

	for fieldName, value := range map[string]string{
		"input_cost_per_million":          request.InputCostPerMillion,
		"output_cost_per_million":         request.OutputCostPerMillion,
		"cache_read_cost_per_million":     request.CacheReadCostPerMillion,
		"cache_creation_cost_per_million": request.CacheCreationCostPerMillion,
	} {
		if err := validateDecimalField(fieldName, value); err != nil {
			return entity.ModelPricing{}, err
		}
	}

	if request.IsPlaceholder == nil {
		return entity.ModelPricing{}, ValidationError{message: "is_placeholder 不能为空"}
	}

	return entity.ModelPricing{
		ModelID:                     modelID,
		DisplayName:                 strings.TrimSpace(request.DisplayName),
		InputCostPerMillion:         request.InputCostPerMillion,
		OutputCostPerMillion:        request.OutputCostPerMillion,
		CacheReadCostPerMillion:     request.CacheReadCostPerMillion,
		CacheCreationCostPerMillion: request.CacheCreationCostPerMillion,
		IsPlaceholder:               *request.IsPlaceholder,
	}, nil
}

func validateDecimalField(fieldName string, value string) error {
	decimalValue := new(big.Rat)
	if _, ok := decimalValue.SetString(strings.TrimSpace(value)); !ok {
		return ValidationError{message: fieldName + " 必须是合法数字"}
	}
	if decimalValue.Sign() < 0 {
		return ValidationError{message: fieldName + " 不能为负数"}
	}
	return nil
}

// IsConflictError 判断错误是否属于资源冲突。
func IsConflictError(err error) bool {
	var conflictError ConflictError
	return errors.As(err, &conflictError)
}

// IsNotFoundError 判断错误是否属于资源不存在。
func IsNotFoundError(err error) bool {
	var notFoundError NotFoundError
	return errors.As(err, &notFoundError)
}
