<template>
  <div class="stock-monitor" @wheel="handleWheel">
    <div class="input-group" :class="{ 'show-input': showSearchWindow }">
      <input 
        v-model="stockCode" 
        placeholder="输入股票代码 (如: sh601606)" 
        @keyup.enter="startMonitoring"
      />
      <button @click="startMonitoring">{{ isMonitoring ? '停止监控' : '开始监控' }}</button>
    </div>
    
    <div v-if="stockData" class="stock-display">
      <div class="stock-header">
        <h2>{{ stockData.name }} ({{ stockData.code }})</h2>
        <div class="market-status" :class="marketStatusClass">
          {{ stockData.marketStatus || '未知' }}
        </div>
        <div class="time-display">最后更新: {{ lastUpdateTime }}</div>
      </div>
      
      <div class="price-container" :class="priceChangeClass">
        <div class="current-price">
          {{ formattedPrice }}
          <span class="change-indicator">
            <span v-if="priceChange > 0" class="up">▲</span>
            <span v-else-if="priceChange < 0" class="down">▼</span>
          </span>
        </div>
        
        <div class="change-info">
          <span v-if="priceChange > 0" class="up">
            +{{ formatNumber(stockData.change) }} ({{ formatNumber(stockData.changePct) }}%)
          </span>
          <span v-else-if="priceChange < 0" class="down">
            {{ formatNumber(stockData.change) }} ({{ formatNumber(stockData.changePct) }}%)
          </span>
          <span v-else class="neutral">
            {{ formatNumber(stockData.change) }} ({{ formatNumber(stockData.changePct) }}%)
          </span>
        </div>
      </div>
      
      <div class="details-toggle" @click="toggleDetails">
        <span>{{ showDetails ? '收起详情' : '查看详情' }}</span>
        <span class="toggle-icon">{{ showDetails ? '▲' : '▼' }}</span>
      </div>
      
      <transition name="slide">
        <div v-if="showDetails" class="stock-details">
          <div class="detail-row">
            <div class="detail-item">
              <div class="detail-label">开盘价</div>
              <div class="detail-value">{{ formatNumber(stockData.open) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">昨收价</div>
              <div class="detail-value">{{ formatNumber(stockData.prevClose) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">最高价</div>
              <div class="detail-value">{{ formatNumber(stockData.high) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">最低价</div>
              <div class="detail-value">{{ formatNumber(stockData.low) }}</div>
            </div>
          </div>
          
          <div class="detail-row">
            <div class="detail-item">
              <div class="detail-label">成交量</div>
              <div class="detail-value">{{ formatVolume(stockData.volume) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">成交额</div>
              <div class="detail-value">{{ formatAmount(stockData.amount) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">买入价</div>
              <div class="detail-value">{{ formatNumber(stockData.bid) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">卖出价</div>
              <div class="detail-value">{{ formatNumber(stockData.ask) }}</div>
            </div>
          </div>
          
          <div class="detail-row">
            <div class="detail-item">
              <div class="detail-label">买入量</div>
              <div class="detail-value">{{ formatVolume(stockData.bidVolume) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">卖出量</div>
              <div class="detail-value">{{ formatVolume(stockData.askVolume) }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">日期</div>
              <div class="detail-value">{{ stockData.date || '未知' }}</div>
            </div>
            <div class="detail-item">
              <div class="detail-label">时间</div>
              <div class="detail-value">{{ stockData.time || '未知' }}</div>
            </div>
          </div>
        </div>
      </transition>
    </div>
    
    <div v-else-if="isLoading" class="loading">加载中...</div>
    <div v-else class="placeholder">请输入股票代码开始监控</div>
    
    <div v-if="error" class="error">{{ error }}</div>
  </div>
</template>

<script setup>
  import { ref, computed, onBeforeUnmount } from 'vue';
  import { GetStockData } from '../../wailsjs/go/main/App';

  const stockCode = ref('sh601606');
  const stockData = ref(null);
  const lastPrice = ref(0);
  const priceChange = ref(0);
  const isLoading = ref(false);
  const error = ref('');
  const isMonitoring = ref(false);
  const timer = ref(null);
  const lastUpdateTime = ref('');
  const showDetails = ref(false); // 控制详细信息的显示/隐藏
  const showSearchWindow = ref(true); // 控制查询窗口的显示/隐藏

  // 格式化价格显示
  const formattedPrice = computed(() => {
    if (!stockData.value) return '--';
    return stockData.value.price.toFixed(2);
  });

  // 根据价格变化设置样式
  const priceChangeClass = computed(() => {
    if (priceChange.value > 0) return 'up';
    if (priceChange.value < 0) return 'down';
    return 'neutral';
  });

  // 根据市场状态设置样式
  const marketStatusClass = computed(() => {
    if (!stockData.value || !stockData.value.marketStatus) return '';
    
    if (stockData.value.marketStatus === '交易中') return 'status-trading';
    if (stockData.value.marketStatus === '已收盘') return 'status-closed';
    return 'status-unknown';
  });

  // 开始/停止监控
  const startMonitoring = async () => {
    if (isMonitoring.value) {
      stopMonitoring();
      return;
    }
    
    if (!stockCode.value) {
      error.value = '请输入股票代码';
      return;
    }
    
    error.value = '';
    isMonitoring.value = true;
    await fetchStockData();
    
    // 每5秒更新一次数据
    timer.value = setInterval(async () => {
      await fetchStockData();
    }, 3000);
  };

  // 停止监控
  const stopMonitoring = () => {
    isMonitoring.value = false;
    if (timer.value) {
      clearInterval(timer.value);
      timer.value = null;
    }
  };

  // 获取股票数据
  const fetchStockData = async () => {
    try {
      isLoading.value = true;
      const data = await GetStockData(stockCode.value);
      
      // 计算价格变化
      if (stockData.value && stockData.value.price) {
        priceChange.value = data.price - stockData.value.price;
      } else {
        priceChange.value = 0;
      }
      
      stockData.value = data;
      lastPrice.value = data.price;
      lastUpdateTime.value = new Date().toLocaleTimeString();
    } catch (err) {
      error.value = `获取数据失败: ${err.message || err}`;
    } finally {
      isLoading.value = false;
    }
  };

  // 格式化数字
  const formatNumber = (num) => {
    if (num === undefined || num === null) return '--';
    return num.toFixed(2);
  };

  // 格式化成交量
  const formatVolume = (volume) => {
    if (volume === undefined || volume === null) return '--';
    if (volume >= 100000000) {
      return (volume / 100000000).toFixed(2) + '亿';
    } else if (volume >= 10000) {
      return (volume / 10000).toFixed(2) + '万';
    }
    return volume;
  };

  // 格式化成交额
  const formatAmount = (amount) => {
    if (amount === undefined || amount === null) return '--';
    if (amount >= 10000) {
      return (amount / 10000).toFixed(2) + '亿';
    } else if (amount >= 1) {
      return amount.toFixed(2) + '万';
    }
    return amount;
  };
  
  // 切换详细信息的显示/隐藏
  const toggleDetails = () => {
    showDetails.value = !showDetails.value;
  };
  
  // 处理滚轮事件
  const handleWheel = (event) => {
    // 检测滚轮方向，deltaY < 0 表示向上滚动
    if (event.deltaY < 0) {
      showSearchWindow.value = true; // 向上滚动时显示查询窗口
    } else {
      showSearchWindow.value = false; // 向下滚动时隐藏查询窗口
    }
  };

  // 组件卸载前停止监控
  onBeforeUnmount(() => {
    stopMonitoring();
  });
</script>

<style scoped>
/* 使用全局CSS变量实现响应式设计 */

.stock-monitor {
  width: 100%;
  height: 100%;
  margin: 0 auto;
  padding: 0; /* 移除内边距，使内容紧贴窗口边缘 */
  font-family: 'Nunito', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  color: var(--text-color);
  font-size: var(--font-size-sm);
  box-sizing: border-box;
  overflow: auto;
}

.input-group {
  max-height: 0;
  overflow: hidden;
  transition: max-height 0.3s ease-in-out, opacity 0.3s ease-in-out, margin 0.3s ease-in-out;
  opacity: 0;
  margin-bottom: 0;
  display: flex;
  gap: var(--spacing-xs);
}

.input-group.show-input {
  max-height: 3em;
  opacity: 1;
  margin-bottom: var(--spacing-sm);
}

.input-group input {
  flex: 1;
  padding: var(--spacing-xs) var(--spacing-sm);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  font-size: var(--font-size-sm);
}

.input-group button {
  padding: var(--spacing-xs) var(--spacing-sm);
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: var(--border-radius);
  cursor: pointer;
  font-size: var(--font-size-sm);
  transition: background 0.3s;
}

.input-group button:hover {
  background: color-mix(in srgb, var(--primary-color), black 10%);
}

h1 {
  text-align: center;
  color: var(--secondary-color);
  margin-bottom: var(--spacing-lg);
}

input:focus {
  border-color: var(--primary-color);
}

.stock-display {
  background: white;
  border-radius: var(--border-radius);
  box-shadow: 0 var(--spacing-xs) var(--spacing-md) rgba(0, 0, 0, 0.1);
  padding: var(--spacing-md) 0 var(--spacing-md) 0; /* 完全移除左右内边距，只保留上下内边距 */
  margin-top: var(--spacing-sm);
  margin-left: 0; /* 确保股票显示区域紧贴左边 */
  margin-right: 0; /* 确保股票显示区域紧贴右边 */
  color: var(--secondary-color); /* 添加文本颜色，确保在白色背景上可见 */
  
}

.stock-header {
  display: flex;
  justify-content: flex-start; /* 将内容靠左对齐 */
  align-items: center;
  margin-bottom: var(--spacing-xs);
  margin-left: 0; /* 移除左外边距，使标题完全贴边 */
  border-bottom: 1px solid var(--border-color);
  padding-bottom: var(--spacing-xs);
  padding-left: 0; /* 移除左内边距，使内容紧贴左边 */
  padding-right: var(--spacing-xs); /* 保留少量右内边距 */
  flex-wrap: wrap;
  text-align: left !important; /* 强制左对齐 */
}

.stock-header h2 {
  margin: 0;
  padding: 0; /* 确保标题没有任何内边距 */
  font-size: var(--font-size-lg);
  flex: 1;
}

.market-status {
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--border-radius);
  font-size: var(--font-size-xs);
  font-weight: 500;
  margin-right: var(--spacing-xs);
}

.status-trading {
  background-color: var(--error-color);
  color: white;
}

.status-closed {
  background-color: #7f8c8d;
  color: white;
}

.status-unknown {
  background-color: var(--warning-color);
  color: white;
}

.time-display {
  font-size: var(--font-size-xs);
  color: #7f8c8d;
}

.price-container {
  text-align: center;
  padding: var(--spacing-md) 0;
  border-radius: var(--border-radius);
  transition: background-color 0.5s;
  margin-left: 0; /* 确保价格容器紧贴左边 */
}

.up {
  color: var(--error-color);
  background-color: rgba(231, 76, 60, 0.05);
}

.down {
  color: var(--success-color);
  background-color: rgba(39, 174, 96, 0.05);
}

.neutral {
  color: #7f8c8d;
}

.current-price {
  font-size: var(--font-size-xl);
  font-weight: bold;
  margin-bottom: var(--spacing-xs);
}

.change-indicator {
  font-size: var(--font-size-lg);
  vertical-align: middle;
}

.change-info {
  font-size: var(--font-size-md);
  font-weight: 500;
}

.stock-details {
  margin-top: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  padding-top: var(--spacing-md);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--spacing-sm);
}

.detail-item {
  flex: 1;
  text-align: center;
  padding: var(--spacing-sm);
  border-radius: var(--border-radius);
  background-color: var(--bg-light);
  margin: 0 var(--spacing-xs);
}

.detail-label {
  font-size: var(--font-size-xs);
  color: #7f8c8d;
  margin-bottom: var(--spacing-xs);
}

.detail-value {
  font-size: var(--font-size-sm);
  font-weight: 500;
}

.volume {
  margin-top: var(--spacing-md);
  text-align: center;
  font-size: var(--font-size-md);
  color: #7f8c8d;
}

.loading, .placeholder, .error {
  text-align: center;
  padding: var(--spacing-md);
  font-size: var(--font-size-sm);
  color: #7f8c8d;
}

.error {
  color: var(--error-color);
  background: rgba(231, 76, 60, 0.05);
  padding: var(--spacing-sm);
  border-radius: var(--border-radius);
  margin-top: var(--spacing-md);
  font-size: var(--font-size-sm);
}

.details-toggle {
  text-align: center;
  padding: var(--spacing-sm);
  margin-top: var(--spacing-sm);
  cursor: pointer;
  color: var(--primary-color);
  font-weight: 500;
  border-radius: var(--border-radius);
  transition: background-color 0.3s;
  font-size: var(--font-size-sm);
}

.details-toggle:hover {
  background-color: rgba(52, 152, 219, 0.1);
}

.toggle-icon {
  margin-left: var(--spacing-xs);
  font-size: var(--font-size-xs);
}

/* 过渡动画 */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
  max-height: 1000px;
  opacity: 1;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
  overflow: hidden;
}

@media (max-width: 768px) {
  .detail-row {
    flex-wrap: wrap;
  }
  
  .detail-item {
    flex-basis: calc(50% - var(--spacing-md));
    margin-bottom: var(--spacing-md);
  }
}

@media (max-width: 480px) {
  .detail-item {
    flex-basis: 100%;
  }
  
  .stock-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .market-status {
    margin-top: var(--spacing-sm);
    margin-bottom: var(--spacing-sm);
  }
}
</style>