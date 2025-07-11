<template>
  <div class="settings-panel" :class="{ 'show': isVisible }">
    <div class="settings-header">
      <h2>{{ currentTab === 'theme' ? '主题设置' : '自定义颜色' }}</h2>
      <button class="close-btn" @click="close">×</button>
    </div>
    
    <div class="settings-tabs">
      <button 
        :class="{ 'active': currentTab === 'theme' }" 
        @click="currentTab = 'theme'"
      >
        主题设置
      </button>
      <button 
        :class="{ 'active': currentTab === 'colors' }" 
        @click="currentTab = 'colors'"
      >
        自定义颜色
      </button>
    </div>
    
    <div class="settings-content">
      <!-- 主题设置 -->
      <div v-if="currentTab === 'theme'" class="theme-settings">
        <div class="theme-options">
          <div 
            class="theme-option" 
            :class="{ 'selected': settings.theme === 'light' }"
            @click="selectTheme('light')"
          >
            <div class="theme-preview light-theme"></div>
            <div class="theme-name">亮色主题</div>
          </div>
          
          <div 
            class="theme-option" 
            :class="{ 'selected': settings.theme === 'dark' }"
            @click="selectTheme('dark')"
          >
            <div class="theme-preview dark-theme"></div>
            <div class="theme-name">暗色主题</div>
          </div>
          
          <div 
            class="theme-option" 
            :class="{ 'selected': settings.theme === 'custom' }"
            @click="selectTheme('custom')"
          >
            <div class="theme-preview custom-theme" :style="customThemeStyle"></div>
            <div class="theme-name">自定义主题</div>
          </div>
        </div>
      </div>
      
      <!-- 自定义颜色设置 -->
      <div v-if="currentTab === 'colors'" class="color-settings">
        <div class="color-picker">
          <label>背景颜色</label>
          <div class="color-input-group">
            <input 
              type="color" 
              v-model="settings.backgroundColor" 
              @change="updateColors"
            />
            <input 
              type="text" 
              v-model="settings.backgroundColor" 
              @change="updateColors"
            />
          </div>
        </div>
        
        <div class="color-picker">
          <label>文本颜色</label>
          <div class="color-input-group">
            <input 
              type="color" 
              v-model="settings.textColor" 
              @change="updateColors"
            />
            <input 
              type="text" 
              v-model="settings.textColor" 
              @change="updateColors"
            />
          </div>
        </div>
        
        <div class="color-picker">
          <label>强调颜色</label>
          <div class="color-input-group">
            <input 
              type="color" 
              v-model="settings.accentColor" 
              @change="updateColors"
            />
            <input 
              type="text" 
              v-model="settings.accentColor" 
              @change="updateColors"
            />
          </div>
        </div>
        
        <div class="color-preview" :style="customThemeStyle">
          <div class="preview-text">预览效果</div>
          <button class="preview-button">按钮示例</button>
        </div>
      </div>
    </div>
    
    <div class="settings-footer">
      <button class="save-btn" @click="saveSettings">保存设置</button>
      <button class="cancel-btn" @click="close">取消</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { GetSettings, UpdateSettings, UpdateTheme, UpdateColors } from '../../wailsjs/go/main/App';

// 组件状态
const isVisible = ref(false);
const currentTab = ref('theme');
const settings = ref({
  theme: 'light',
  backgroundColor: '#ffffff',
  textColor: '#000000',
  accentColor: '#3498db'
});

// 自定义主题样式
const customThemeStyle = computed(() => {
  return {
    backgroundColor: settings.value.backgroundColor,
    color: settings.value.textColor,
    '--accent-color': settings.value.accentColor
  };
});

// 加载设置
const loadSettings = async () => {
  try {
    const savedSettings = await GetSettings();
    settings.value = savedSettings;
    applyTheme(savedSettings.theme);
  } catch (error) {
    console.error('加载设置失败:', error);
  }
};

// 保存设置
const saveSettings = async () => {
  try {
    await UpdateSettings(settings.value);
    applyTheme(settings.value.theme);
    close();
  } catch (error) {
    console.error('保存设置失败:', error);
  }
};

// 选择主题
const selectTheme = async (theme) => {
  settings.value.theme = theme;
  applyTheme(theme);
  
  try {
    await UpdateTheme(theme);
  } catch (error) {
    console.error('更新主题失败:', error);
  }
};

// 更新颜色
const updateColors = async () => {
  if (settings.value.theme !== 'custom') {
    settings.value.theme = 'custom';
  }
  
  try {
    await UpdateColors(
      settings.value.backgroundColor,
      settings.value.textColor,
      settings.value.accentColor
    );
    applyTheme('custom');
  } catch (error) {
    console.error('更新颜色失败:', error);
  }
};

// 应用主题
const applyTheme = (theme) => {
  const root = document.documentElement;
  
  if (theme === 'light') {
    root.style.setProperty('--background-color', '#ffffff');
    root.style.setProperty('--text-color', '#2c3e50');
    root.style.setProperty('--primary-color', '#3498db');
    root.style.setProperty('--secondary-color', '#2c3e50');
    root.style.setProperty('--border-color', '#ddd');
  } else if (theme === 'dark') {
    root.style.setProperty('--background-color', '#1a1a1a');
    root.style.setProperty('--text-color', '#ecf0f1');
    root.style.setProperty('--primary-color', '#3498db');
    root.style.setProperty('--secondary-color', '#ecf0f1');
    root.style.setProperty('--border-color', '#444');
  } else if (theme === 'custom') {
    root.style.setProperty('--background-color', settings.value.backgroundColor);
    root.style.setProperty('--text-color', settings.value.textColor);
    root.style.setProperty('--primary-color', settings.value.accentColor);
    root.style.setProperty('--secondary-color', settings.value.textColor);
  }
};

// 关闭设置面板
const close = () => {
  isVisible.value = false;
};

// 打开设置面板
const open = (tab = 'theme') => {
  currentTab.value = tab;
  isVisible.value = true;
};

// 监听事件
onMounted(() => {
  loadSettings();
  
  // 监听打开主题设置事件
  window.runtime.EventsOn('open-theme-settings', () => {
    console.log('open-theme-settings');
    open('theme');
  });
  
  // 监听打开颜色设置事件
  window.runtime.EventsOn('open-color-settings', () => {
    open('colors');
  });
});

// 导出方法供父组件调用
defineExpose({
  open,
  close
});
</script>

<style scoped>
.settings-panel {
  position: fixed;
  top: 0;
  right: -400px;
  width: 400px;
  height: 100vh;
  background-color: var(--background-color);
  box-shadow: -2px 0 10px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  transition: right 0.3s ease;
  display: flex;
  flex-direction: column;
  color: var(--text-color);
}

.settings-panel.show {
  right: 0;
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--border-color);
}

.settings-header h2 {
  margin: 0;
  font-size: var(--font-size-lg);
}

.close-btn {
  background: none;
  border: none;
  font-size: var(--font-size-xl);
  cursor: pointer;
  color: var(--text-color);
  padding: 0;
}

.settings-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
}

.settings-tabs button {
  flex: 1;
  background: none;
  border: none;
  padding: var(--spacing-sm);
  cursor: pointer;
  color: var(--text-color);
  font-weight: 500;
  transition: all 0.3s;
}

.settings-tabs button.active {
  color: var(--primary-color);
  border-bottom: 2px solid var(--primary-color);
}

.settings-content {
  flex: 1;
  padding: var(--spacing-md);
  overflow-y: auto;
}

.theme-options {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  justify-content: center;
}

.theme-option {
  width: 120px;
  cursor: pointer;
  border-radius: var(--border-radius);
  padding: var(--spacing-sm);
  transition: all 0.3s;
  border: 2px solid transparent;
}

.theme-option.selected {
  border-color: var(--primary-color);
}

.theme-preview {
  height: 80px;
  border-radius: var(--border-radius);
  margin-bottom: var(--spacing-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.theme-preview::after {
  content: '';
  position: absolute;
  width: 60%;
  height: 20%;
  bottom: 10px;
  left: 20%;
  border-radius: var(--border-radius);
}

.light-theme {
  background-color: #ffffff;
  color: #2c3e50;
}

.light-theme::after {
  background-color: #3498db;
}

.dark-theme {
  background-color: #1a1a1a;
  color: #ecf0f1;
}

.dark-theme::after {
  background-color: #3498db;
}

.custom-theme::after {
  background-color: var(--accent-color, #3498db);
}

.theme-name {
  text-align: center;
  font-size: var(--font-size-sm);
}

.color-picker {
  margin-bottom: var(--spacing-md);
}

.color-picker label {
  display: block;
  margin-bottom: var(--spacing-xs);
  font-weight: 500;
}

.color-input-group {
  display: flex;
  gap: var(--spacing-sm);
}

.color-input-group input[type="color"] {
  width: 50px;
  height: 40px;
  padding: 0;
  border: 1px solid var(--border-color);
}

.color-input-group input[type="text"] {
  flex: 1;
}

.color-preview {
  margin-top: var(--spacing-lg);
  padding: var(--spacing-md);
  border-radius: var(--border-radius);
  text-align: center;
}

.preview-text {
  margin-bottom: var(--spacing-md);
  font-size: var(--font-size-md);
}

.preview-button {
  background-color: var(--accent-color, var(--primary-color));
  color: white;
  border: none;
  padding: var(--spacing-sm) var(--spacing-md);
  border-radius: var(--border-radius);
  cursor: pointer;
}

.settings-footer {
  padding: var(--spacing-md);
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
}

.save-btn {
  background-color: var(--primary-color);
  color: white;
}

.cancel-btn {
  background-color: transparent;
  color: var(--text-color);
  border: 1px solid var(--border-color);
}

@media (max-width: 480px) {
  .settings-panel {
    width: 100%;
    right: -100%;
  }
  
  .theme-options {
    flex-direction: column;
    align-items: center;
  }
  
  .theme-option {
    width: 80%;
  }
}
</style>