<script setup>
import { ref } from 'vue';
import Stock from './components/Stock.vue';
import Settings from './components/Settings.vue';
</script>

<template>
  <div class="app-container">
    <Stock />
    <Settings ref="settingsComponent" />
  </div>
</template>

<style>
.app-container {
  width: 100%;
  height: 100%;
  position: relative;
}
</style>