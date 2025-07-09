package stock

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese" // GBK编码支持
	"golang.org/x/text/transform"                  // 编码转换
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// StockService 股票服务结构体
type StockService struct {
	clients map[string]*http.Client
	mu      sync.Mutex
}

// StockData 股票数据结构体
type StockData struct {
	Name         string    `json:"name"`         // 股票名称
	Code         string    `json:"code"`         // 股票代码
	Price        float64   `json:"price"`        // 当前价格
	Change       float64   `json:"change"`       // 涨跌额
	ChangePct    float64   `json:"changePct"`    // 涨跌幅(%)
	Volume       int64     `json:"volume"`       // 成交量(手)
	Amount       float64   `json:"amount"`       // 成交额(万元)
	Open         float64   `json:"open"`         // 开盘价
	PrevClose    float64   `json:"prevClose"`    // 昨收价
	High         float64   `json:"high"`         // 最高价
	Low          float64   `json:"low"`          // 最低价
	Bid          float64   `json:"bid"`          // 买入价
	Ask          float64   `json:"ask"`          // 卖出价
	BidVolume    int64     `json:"bidVolume"`    // 买入量(手)
	AskVolume    int64     `json:"askVolume"`    // 卖出量(手)
	MarketStatus string    `json:"marketStatus"` // 市场状态
	UpdateTime   time.Time `json:"updateTime"`   // 数据更新时间
	Date         string    `json:"date"`         // 交易日期
	Time         string    `json:"time"`         // 交易时间
}

// NewStockService 创建新的股票服务实例
func NewStockService() *StockService {
	return &StockService{
		clients: make(map[string]*http.Client),
	}
}

// GetStockData 获取股票数据
func (s *StockService) GetStockData(stockCode string) (*StockData, error) {
	url := fmt.Sprintf("https://qt.gtimg.cn/q=%s", stockCode)

	// 为每个股票代码创建独立的HTTP客户端
	s.mu.Lock()
	client, exists := s.clients[stockCode]
	if !exists {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
		s.clients[stockCode] = client
	}
	s.mu.Unlock()

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP错误: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	decoder := simplifiedchinese.GBK.NewDecoder()
	utf8Body, _, err := transform.Bytes(decoder, body)
	if err != nil {
		return nil, fmt.Errorf("编码转换失败: %v", err)
	}

	return parseStockData(string(utf8Body), stockCode)
}

// parseStockData 解析股票数据
func parseStockData(data, stockCode string) (*StockData, error) {
	// 格式: v_sh601606="1~长城军工~601606~28.71~28.69~28.28~652127~316651~335311~28.70~52~..."

	// 提取有效数据部分
	start := strings.Index(data, "\"")
	end := strings.LastIndex(data, "\"")
	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("无效的数据格式")
	}

	content := data[start+1 : end]
	fields := strings.Split(content, "~")
	if len(fields) < 40 {
		return nil, fmt.Errorf("数据字段不足")
	}

	// 解析关键字段
	name := fields[1]

	// 解析价格相关字段
	price, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return nil, fmt.Errorf("解析价格失败: %v", err)
	}

	// 解析开盘价
	open, _ := strconv.ParseFloat(fields[5], 64)

	// 解析昨收价
	prevClose, _ := strconv.ParseFloat(fields[4], 64)

	// 解析最高价和最低价
	high, _ := strconv.ParseFloat(fields[41], 64)
	low, _ := strconv.ParseFloat(fields[42], 64)

	// 如果最高价和最低价解析失败，尝试其他可能的字段
	if high == 0 && len(fields) > 33 {
		high, _ = strconv.ParseFloat(fields[33], 64)
	}
	if low == 0 && len(fields) > 34 {
		low, _ = strconv.ParseFloat(fields[34], 64)
	}

	// 解析成交量(手)
	volume, _ := strconv.ParseInt(fields[6], 10, 64)

	// 解析成交额(万元)
	amount, _ := strconv.ParseFloat(fields[37], 64)

	// 解析涨跌额和涨跌幅
	change, _ := strconv.ParseFloat(fields[31], 64)
	changePct, _ := strconv.ParseFloat(fields[32], 64)

	// 解析买入价和卖出价
	bid, _ := strconv.ParseFloat(fields[9], 64)
	ask, _ := strconv.ParseFloat(fields[19], 64)

	// 解析买入量和卖出量
	bidVolume, _ := strconv.ParseInt(fields[10], 10, 64)
	askVolume, _ := strconv.ParseInt(fields[20], 10, 64)

	// 解析日期和时间
	date := ""
	timeStr := ""
	if len(fields) > 30 {
		date = fields[30]
	}
	// 时间可能在不同位置，尝试获取
	if len(fields) > 31 && fields[31] != "" && !strings.Contains(fields[31], ".") {
		timeStr = fields[31]
	} else if len(fields) > 29 && fields[29] != "" {
		timeStr = fields[29]
	}

	// 确定市场状态
	marketStatus := "未知"
	if len(fields) > 40 {
		// 根据API返回的状态码确定市场状态
		statusCode := fields[40]
		switch statusCode {
		case "1":
			marketStatus = "交易中"
		case "2":
			marketStatus = "已收盘"
		case "3":
			marketStatus = "开盘前"
		case "4":
			marketStatus = "盘中休息"
		case "5":
			marketStatus = "已暂停"
		case "6":
			marketStatus = "已熔断"
		case "7":
			marketStatus = "未开盘"
		}
	} else if strings.Contains(data, "\"市场已收盘\"") {
		marketStatus = "已收盘"
	} else if price > 0 {
		marketStatus = "交易中"
	} else {
		marketStatus = "未开市"
	}

	return &StockData{
		Name:         name,
		Code:         stockCode,
		Price:        price,
		Change:       change,
		ChangePct:    changePct,
		Volume:       volume,
		Amount:       amount,
		Open:         open,
		PrevClose:    prevClose,
		High:         high,
		Low:          low,
		Bid:          bid,
		Ask:          ask,
		BidVolume:    bidVolume,
		AskVolume:    askVolume,
		MarketStatus: marketStatus,
		UpdateTime:   time.Now(),
		Date:         date,
		Time:         timeStr,
	}, nil
}
