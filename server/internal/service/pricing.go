package service

import (
	"context"
	"fmt"
	"math/big"

	"cc-status/server/internal/model/entity"

	"gorm.io/gorm"
)

var tokensPerMillion = big.NewRat(1_000_000, 1)

// ModelPricingMatcher 定义同步计费所需的最小定价查询能力。
type ModelPricingMatcher interface {
	FindMatch(context.Context, *gorm.DB, string) (entity.ModelPricing, string, error)
}

// PricingService 负责命中定价规则并计算费用字段。
type PricingService struct {
	matcher ModelPricingMatcher
}

// NewPricingService 创建定价服务。
func NewPricingService(matcher ModelPricingMatcher) *PricingService {
	return &PricingService{matcher: matcher}
}

// ApplyToReport 根据模型定价补全费用字段与 pricing_source。
func (service *PricingService) ApplyToReport(
	ctx context.Context,
	tx *gorm.DB,
	report entity.UsageReport,
) (entity.UsageReport, error) {
	pricing, pricingSource, err := service.matcher.FindMatch(ctx, tx, report.Model)
	if err != nil {
		return entity.UsageReport{}, fmt.Errorf("find model pricing: %w", err)
	}

	inputCost, err := calculateUSD(report.InputTokens, pricing.InputCostPerMillion)
	if err != nil {
		return entity.UsageReport{}, err
	}
	outputCost, err := calculateUSD(report.OutputTokens, pricing.OutputCostPerMillion)
	if err != nil {
		return entity.UsageReport{}, err
	}
	cacheReadCost, err := calculateUSD(report.CacheReadTokens, pricing.CacheReadCostPerMillion)
	if err != nil {
		return entity.UsageReport{}, err
	}
	cacheCreationCost, err := calculateUSD(report.CacheCreationTokens, pricing.CacheCreationCostPerMillion)
	if err != nil {
		return entity.UsageReport{}, err
	}

	totalCost := new(big.Rat).Add(inputCost, outputCost)
	totalCost.Add(totalCost, cacheReadCost)
	totalCost.Add(totalCost, cacheCreationCost)

	report.InputCostUSD = inputCost.FloatString(10)
	report.OutputCostUSD = outputCost.FloatString(10)
	report.CacheReadCostUSD = cacheReadCost.FloatString(10)
	report.CacheCreationCostUSD = cacheCreationCost.FloatString(10)
	report.TotalCostUSD = totalCost.FloatString(10)
	report.PricingSource = pricingSource

	return report, nil
}

func calculateUSD(tokens int64, pricePerMillion string) (*big.Rat, error) {
	price := new(big.Rat)
	if _, ok := price.SetString(pricePerMillion); !ok {
		return nil, fmt.Errorf("invalid decimal value %q", pricePerMillion)
	}

	cost := new(big.Rat).Mul(big.NewRat(tokens, 1), price)
	cost.Quo(cost, tokensPerMillion)
	return cost, nil
}
