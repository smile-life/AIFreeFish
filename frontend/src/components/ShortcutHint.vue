<template>
  <div class="shortcut-hint" v-if="visible">
    <div class="shortcut-item" v-for="(hint, index) in hints" :key="index">
      <span class="shortcut-key">{{ hint.key }}</span>
      <span class="shortcut-description">{{ hint.description }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  }
});

const hints = ref([
  { key: 'Ctrl + S', description: '保存设置' },
  { key: 'Esc', description: '关闭设置面板' }
]);
</script>

<style scoped>
.shortcut-hint {
  position: fixed;
  bottom: 20px;
  right: 20px;
  background-color: rgba(0, 0, 0, 0.7); /* 半透明黑色背景，在任何主题下都能清晰显示 */
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius);
  padding: var(--spacing-sm);
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
  z-index: 1001;
}

.shortcut-item {
  display: flex;
  align-items: center;
  margin-bottom: var(--spacing-xs);
}

.shortcut-item:last-child {
  margin-bottom: 0;
}

.shortcut-key {
  background-color: var(--secondary-color);
  color: var(--background-color);
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 0.8em;
  margin-right: var(--spacing-sm);
  font-family: monospace;
}

.shortcut-description {
  font-size: 0.9em;
  color: var(--text-color);
}
</style>