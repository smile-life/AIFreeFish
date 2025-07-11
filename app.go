package main

import (
	"context"

	"makedemo/internal/settings"
	"makedemo/internal/stock"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App 应用程序主结构体
type App struct {
	ctx             context.Context
	stockService    *stock.StockService
	settingsService *settings.Service
}

// NewApp 创建新的应用程序实例
func NewApp() *App {
	return &App{
		stockService:    stock.NewStockService(),
		settingsService: settings.NewService(),
	}
}

// startup 应用程序启动时的初始化
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化设置服务
	if err := a.settingsService.LoadSettings(); err != nil {
		runtime.LogErrorf(ctx, "加载设置失败: %v", err)
	}

	// 应用当前主题设置
	settings := a.settingsService.GetSettings()
	if err := a.UpdateTheme(settings.Theme); err != nil {
		runtime.LogErrorf(ctx, "应用主题失败: %v", err)
	}

	// 应用当前颜色设置
	if err := a.UpdateColors(
		settings.BackgroundColor,
		settings.TextColor,
		settings.AccentColor,
	); err != nil {
		runtime.LogErrorf(ctx, "应用颜色失败: %v", err)
	}
}

// GetStockData 获取股票数据，作为前端API的入口点
func (a *App) GetStockData(stockCode string) (*stock.StockData, error) {
	return a.stockService.GetStockData(stockCode)
}

// GetSettings 获取应用设置
func (a *App) GetSettings() settings.Settings {
	return a.settingsService.GetSettings()
}

// UpdateSettings 更新应用设置
func (a *App) UpdateSettings(settings settings.Settings) error {
	return a.settingsService.UpdateSettings(settings)
}

// UpdateTheme 更新主题
func (a *App) UpdateTheme(theme settings.Theme) error {
	return a.settingsService.UpdateTheme(theme)
}

// UpdateColors 更新颜色设置
func (a *App) UpdateColors(backgroundColor, textColor, accentColor string) error {
	return a.settingsService.UpdateColors(backgroundColor, textColor, accentColor)
}
