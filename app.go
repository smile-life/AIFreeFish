package main

import (
	"context"

	"makedemo/internal/stock"
)

// App 应用程序主结构体
type App struct {
	ctx          context.Context
	stockService *stock.StockService
}

// NewApp 创建新的应用程序实例
func NewApp() *App {
	return &App{
		stockService: stock.NewStockService(),
	}
}

// startup 应用程序启动时的初始化
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetStockData 获取股票数据，作为前端API的入口点
func (a *App) GetStockData(stockCode string) (*stock.StockData, error) {
	return a.stockService.GetStockData(stockCode)
}
