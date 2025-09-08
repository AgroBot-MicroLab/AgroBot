<script setup>
import Map from './components/Map.vue'
import Sidebar from './components/Sidebar.vue'
import { ref } from 'vue'

const point = ref(null);

</script>

<template>
    <div class="layout">
        <div class="map">
            <Map v-model:point="point" />
        </div>

        <div class="sidebar">
            <Sidebar :point="point" />
        </div>
    </div>
</template>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
}

.map {
  flex: 0 0 80%;
}

.sidebar {
  flex: 0 0 20%;
  border-left: 1px solid #ddd;
}
</style>
