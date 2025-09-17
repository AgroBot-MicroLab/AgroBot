<script setup>
import Map from './components/Map.vue'
import Sidebar from './components/Sidebar.vue'


</script>

<template>
    <div class="layout">
        <div class="map">
            <Map />
        </div>

        <div class="sidebar">
            <Sidebar />
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
