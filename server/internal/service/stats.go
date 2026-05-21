package service

import (
	"context"
	"math/big"
	"sort"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"

	"gorm.io/gorm"
)

// UsageReportReader 定义统计服务读取使用记录所需的最小能力。
type UsageReportReader interface {
	List(context.Context, *gorm.DB) ([]entity.UsageReport, error)
}

// StatsService 承载总览与趋势统计能力。
type StatsService struct {
	db     *gorm.DB
	reader UsageReportReader
}

// NewStatsService 创建统计服务。
func NewStatsService(db *gorm.DB) *StatsService {
	return &StatsService{
		db:     db,
		reader: repository.NewUsageReportRepository(),
	}
}

// Overview 聚合生成总览统计结果。
func (service *StatsService) Overview(ctx context.Context) (dto.StatsOverviewResponse, error) {
	reports, err := service.reader.List(ctx, service.db)
	if err != nil {
		return dto.StatsOverviewResponse{}, err
	}

	modelTokens := make(map[string]int64)
	clientCosts := make(map[string]*big.Rat)
	activeClients := make(map[string]struct{})
	totalCost := new(big.Rat)

	var totalTokens int64
	for _, report := range reports {
		tokenCount := report.InputTokens + report.OutputTokens + report.CacheReadTokens + report.CacheCreationTokens
		totalTokens += tokenCount
		modelTokens[report.Model] += tokenCount
		activeClients[report.ClientID] = struct{}{}

		reportCost := parseDecimal(report.TotalCostUSD)
		totalCost.Add(totalCost, reportCost)
		if _, exists := clientCosts[report.ClientID]; !exists {
			clientCosts[report.ClientID] = new(big.Rat)
		}
		clientCosts[report.ClientID].Add(clientCosts[report.ClientID], reportCost)
	}

	return dto.StatsOverviewResponse{
		TotalTokens:   totalTokens,
		TotalCostUSD:  totalCost.FloatString(10),
		TotalRequests: int64(len(reports)),
		ActiveClients: int64(len(activeClients)),
		TopModels:     buildTopModels(modelTokens),
		TopClients:    buildTopClients(clientCosts),
	}, nil
}

func buildTopModels(modelTokens map[string]int64) []dto.StatsModelRank {
	items := make([]dto.StatsModelRank, 0, len(modelTokens))
	for model, tokens := range modelTokens {
		items = append(items, dto.StatsModelRank{
			Model:  model,
			Tokens: tokens,
		})
	}

	sort.Slice(items, func(left int, right int) bool {
		if items[left].Tokens == items[right].Tokens {
			return items[left].Model < items[right].Model
		}
		return items[left].Tokens > items[right].Tokens
	})

	if len(items) > 5 {
		return items[:5]
	}
	return items
}

func buildTopClients(clientCosts map[string]*big.Rat) []dto.StatsClientRank {
	items := make([]dto.StatsClientRank, 0, len(clientCosts))
	for clientID, totalCost := range clientCosts {
		items = append(items, dto.StatsClientRank{
			ClientID:     clientID,
			TotalCostUSD: totalCost.FloatString(10),
		})
	}

	sort.Slice(items, func(left int, right int) bool {
		leftCost := parseDecimal(items[left].TotalCostUSD)
		rightCost := parseDecimal(items[right].TotalCostUSD)
		if leftCost.Cmp(rightCost) == 0 {
			return items[left].ClientID < items[right].ClientID
		}
		return leftCost.Cmp(rightCost) > 0
	})

	if len(items) > 5 {
		return items[:5]
	}
	return items
}

func parseDecimal(value string) *big.Rat {
	decimalValue := new(big.Rat)
	if _, ok := decimalValue.SetString(value); !ok {
		return new(big.Rat)
	}
	return decimalValue
}
