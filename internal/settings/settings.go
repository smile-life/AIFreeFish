package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Theme 主题类型
type Theme string

const (
	// LightTheme 亮色主题
	LightTheme Theme = "light"
	// DarkTheme 暗色主题
	DarkTheme Theme = "dark"
	// CustomTheme 自定义主题
	CustomTheme Theme = "custom"
)

// Settings 应用设置结构体
type Settings struct {
	Theme           Theme  `json:"theme"`
	BackgroundColor string `json:"backgroundColor"`
	TextColor       string `json:"textColor"`
	AccentColor     string `json:"accentColor"`
}

// Service 设置服务
type Service struct {
	settings     Settings
	settingsPath string
}

// NewService 创建新的设置服务
func NewService() *Service {
	// 获取用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}

	// 创建应用配置目录
	appConfigDir := filepath.Join(configDir, "makedemo")
	if _, err := os.Stat(appConfigDir); os.IsNotExist(err) {
		os.MkdirAll(appConfigDir, 0755)
	}

	settingsPath := filepath.Join(appConfigDir, "settings.json")

	// 创建服务实例
	service := &Service{
		settingsPath: settingsPath,
		settings: Settings{
			Theme:           LightTheme,
			BackgroundColor: "#ffffff",
			TextColor:       "#000000",
			AccentColor:     "#3b82f6",
		},
	}

	// 加载设置
	service.LoadSettings()

	return service
}

// LoadSettings 从文件加载设置
func (s *Service) LoadSettings() error {
	// 检查设置文件是否存在
	if _, err := os.Stat(s.settingsPath); os.IsNotExist(err) {
		log.Println("设置文件不存在，创建默认设置")
		// 如果不存在，保存默认设置
		if err := s.SaveSettings(); err != nil {
			log.Printf("保存默认设置失败: %v", err)
			return fmt.Errorf("保存默认设置失败: %w", err)
		}
		log.Println("默认设置已创建")
		return nil
	}

	// 读取设置文件
	data, err := ioutil.ReadFile(s.settingsPath)
	if err != nil {
		log.Printf("读取设置文件失败: %v", err)
		return fmt.Errorf("读取设置文件失败: %w", err)
	}

	// 解析JSON
	err = json.Unmarshal(data, &s.settings)
	if err != nil {
		log.Printf("解析设置JSON失败: %v", err)
		return fmt.Errorf("解析设置JSON失败: %w", err)
	}

	log.Println("设置加载成功")
	return nil
}

// SaveSettings 保存设置到文件
func (s *Service) SaveSettings() error {
	// 确保目录存在
	dir := filepath.Dir(s.settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("创建设置目录失败: %v", err)
		return fmt.Errorf("创建设置目录失败: %w", err)
	}

	// 将设置转换为JSON
	data, err := json.MarshalIndent(s.settings, "", "  ")
	if err != nil {
		log.Printf("序列化设置失败: %v", err)
		return fmt.Errorf("序列化设置失败: %w", err)
	}

	// 写入文件
	err = ioutil.WriteFile(s.settingsPath, data, 0644)
	if err != nil {
		log.Printf("写入设置文件失败: %v", err)
		return fmt.Errorf("写入设置文件失败: %w", err)
	}

	log.Println("设置保存成功")
	return nil
}

// GetSettings 获取当前设置
func (s *Service) GetSettings() Settings {
	return s.settings
}

// UpdateSettings 更新设置
func (s *Service) UpdateSettings(settings Settings) error {
	s.settings = settings
	return s.SaveSettings()
}

// UpdateTheme 更新主题
func (s *Service) UpdateTheme(theme Theme) error {
	s.settings.Theme = theme
	return s.SaveSettings()
}

// UpdateColors 更新颜色设置
func (s *Service) UpdateColors(backgroundColor, textColor, accentColor string) error {
	s.settings.BackgroundColor = backgroundColor
	s.settings.TextColor = textColor
	s.settings.AccentColor = accentColor
	return s.SaveSettings()
}
