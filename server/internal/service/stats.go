package service

import (
	"context"
	"math/big"
	"sort"
	"time"

	"cc-status/server/internal/model/dto"
	"cc-status/server/internal/model/entity"
	"cc-status/server/internal/repository"

	"gorm.io/gorm"
)

// UsageReportReader 定义统计服务读取使用记录所需的最小能力。
type UsageReportReader interface {
	List(context.Context, *gorm.DB) ([]entity.UsageReport, error)
}

// StatsService 承载总览、趋势与仪表盘统计能力。
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

// Dashboard 返回仪表盘统计接口的稳定响应骨架。
func (service *StatsService) Dashboard(
	ctx context.Context,
	query dto.StatsDashboardQuery,
) (dto.StatsDashboardResponse, error) {
	if query.Interval != "hour" && query.Interval != "day" {
		return dto.StatsDashboardResponse{}, ValidationError{message: "interval 仅支持 hour 或 day"}
	}
	if query.StartAt > query.EndAt {
		return dto.StatsDashboardResponse{}, ValidationError{message: "start_at 必须小于等于 end_at"}
	}

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return dto.StatsDashboardResponse{}, err
	}

	reports, err := service.reader.List(ctx, service.db)
	if err != nil {
		return dto.StatsDashboardResponse{}, err
	}

	step := time.Hour
	if query.Interval == "day" {
		step = 24 * time.Hour
	}

	startAt := truncateTrendTime(time.Unix(query.StartAt, 0).In(location), query.Interval)
	endAt := truncateTrendTime(time.Unix(query.EndAt, 0).In(location), query.Interval)
	rangeEnd := endAt.Add(step - time.Nanosecond)

	type dashboardAggregate struct {
		inputTokens         int64
		outputTokens        int64
		cacheReadTokens     int64
		cacheCreationTokens int64
		totalRequests       int64
		totalCostUSD        *big.Rat
	}

	activeClients := make(map[string]struct{})
	buckets := make(map[time.Time]*dashboardAggregate)
	totalCost := new(big.Rat)

	var totalTokens int64
	var totalRequests int64

	for _, report := range reports {
		reportTime := time.Unix(report.CreatedAtUnix, 0).In(location)
		if reportTime.Before(startAt) || reportTime.After(rangeEnd) {
			continue
		}

		tokenCount := report.InputTokens + report.OutputTokens + report.CacheReadTokens + report.CacheCreationTokens
		totalTokens += tokenCount
		totalRequests++
		activeClients[report.ClientID] = struct{}{}
		totalCost.Add(totalCost, parseDecimal(report.TotalCostUSD))

		bucketTime := truncateTrendTime(reportTime, query.Interval)
		if _, exists := buckets[bucketTime]; !exists {
			buckets[bucketTime] = &dashboardAggregate{totalCostUSD: new(big.Rat)}
		}

		aggregate := buckets[bucketTime]
		aggregate.inputTokens += report.InputTokens
		aggregate.outputTokens += report.OutputTokens
		aggregate.cacheReadTokens += report.CacheReadTokens
		aggregate.cacheCreationTokens += report.CacheCreationTokens
		aggregate.totalRequests++
		aggregate.totalCostUSD.Add(aggregate.totalCostUSD, parseDecimal(report.TotalCostUSD))
	}

	trend := []dto.StatsDashboardTrendPoint{}
	if totalRequests > 0 {
		trend = make([]dto.StatsDashboardTrendPoint, 0)
		for current := startAt; !current.After(endAt); current = current.Add(step) {
			aggregate, exists := buckets[current]
			if !exists {
				trend = append(trend, dto.StatsDashboardTrendPoint{
					Bucket:       current.Format(time.RFC3339),
					TotalCostUSD: "0.0000000000",
				})
				continue
			}

			trend = append(trend, dto.StatsDashboardTrendPoint{
				Bucket:              current.Format(time.RFC3339),
				InputTokens:         aggregate.inputTokens,
				OutputTokens:        aggregate.outputTokens,
				CacheReadTokens:     aggregate.cacheReadTokens,
				CacheCreationTokens: aggregate.cacheCreationTokens,
				TotalRequests:       aggregate.totalRequests,
				TotalCostUSD:        aggregate.totalCostUSD.FloatString(10),
			})
		}
	}

	return dto.StatsDashboardResponse{
		Overview: dto.StatsDashboardOverview{
			TotalTokens:   totalTokens,
			TotalCostUSD:  totalCost.FloatString(10),
			TotalRequests: totalRequests,
			ActiveClients: int64(len(activeClients)),
		},
		Trend:      trend,
		TopModels:  []dto.StatsDashboardModelRank{},
		TopClients: []dto.StatsClientRank{},
		CacheAnalysis: dto.StatsDashboardCacheAnalysis{
			SavedCostUSD:         "0.0000000000",
			CacheReadCostUSD:     "0.0000000000",
			CacheCreationCostUSD: "0.0000000000",
		},
	}, nil
}

// Trend 按指定粒度返回基于业务时间的趋势统计结果，并补零缺失时间桶。
func (service *StatsService) Trend(
	ctx context.Context,
	query dto.StatsTrendQuery,
) ([]dto.StatsTrendPoint, error) {
	if query.Interval != "hour" && query.Interval != "day" {
		return nil, ValidationError{message: "interval 仅支持 hour 或 day"}
	}

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, err
	}

	reports, err := service.reader.List(ctx, service.db)
	if err != nil {
		return nil, err
	}
	if len(reports) == 0 {
		return []dto.StatsTrendPoint{}, nil
	}

	startAt, endAt := resolveTrendRange(location, reports, query)
	if startAt.After(endAt) {
		return []dto.StatsTrendPoint{}, nil
	}

	step := time.Hour
	if query.Interval == "day" {
		step = 24 * time.Hour
	}

	type aggregate struct {
		totalTokens   int64
		totalRequests int64
		totalCostUSD  *big.Rat
	}

	buckets := make(map[time.Time]*aggregate)
	for _, report := range reports {
		reportTime := time.Unix(report.CreatedAtUnix, 0).In(location)
		if reportTime.Before(startAt) || reportTime.After(endAt.Add(step-time.Nanosecond)) {
			continue
		}

		bucketTime := truncateTrendTime(reportTime, query.Interval)
		if _, exists := buckets[bucketTime]; !exists {
			buckets[bucketTime] = &aggregate{totalCostUSD: new(big.Rat)}
		}

		tokenCount := report.InputTokens + report.OutputTokens + report.CacheReadTokens + report.CacheCreationTokens
		buckets[bucketTime].totalTokens += tokenCount
		buckets[bucketTime].totalRequests++
		buckets[bucketTime].totalCostUSD.Add(buckets[bucketTime].totalCostUSD, parseDecimal(report.TotalCostUSD))
	}

	points := make([]dto.StatsTrendPoint, 0)
	for current := startAt; !current.After(endAt); current = current.Add(step) {
		aggregateValue, exists := buckets[current]
		if !exists {
			points = append(points, dto.StatsTrendPoint{
				Bucket:        current.Format(time.RFC3339),
				TotalTokens:   0,
				TotalRequests: 0,
				TotalCostUSD:  "0.0000000000",
			})
			continue
		}

		points = append(points, dto.StatsTrendPoint{
			Bucket:        current.Format(time.RFC3339),
			TotalTokens:   aggregateValue.totalTokens,
			TotalRequests: aggregateValue.totalRequests,
			TotalCostUSD:  aggregateValue.totalCostUSD.FloatString(10),
		})
	}

	return points, nil
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

	if len(items) > 10 {
		return items[:10]
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

	if len(items) > 10 {
		return items[:10]
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

func resolveTrendRange(
	location *time.Location,
	reports []entity.UsageReport,
	query dto.StatsTrendQuery,
) (time.Time, time.Time) {
	if query.StartAt > 0 && query.EndAt > 0 {
		startAt := truncateTrendTime(time.Unix(query.StartAt, 0).In(location), query.Interval)
		endAt := truncateTrendTime(time.Unix(query.EndAt, 0).In(location), query.Interval)
		return startAt, endAt
	}

	minTime := time.Unix(reports[0].CreatedAtUnix, 0).In(location)
	maxTime := minTime
	for _, report := range reports[1:] {
		reportTime := time.Unix(report.CreatedAtUnix, 0).In(location)
		if reportTime.Before(minTime) {
			minTime = reportTime
		}
		if reportTime.After(maxTime) {
			maxTime = reportTime
		}
	}

	return truncateTrendTime(minTime, query.Interval), truncateTrendTime(maxTime, query.Interval)
}

func truncateTrendTime(value time.Time, interval string) time.Time {
	if interval == "day" {
		return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
	}

	return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), 0, 0, 0, value.Location())
}
