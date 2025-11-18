<script setup>
import Map from './components/Map.vue'
import Sidebar from './components/Sidebar.vue'
import Gallery from './components/Gallery.vue'
import { ref } from 'vue'

const showGallery = ref(false)

function openGallery() {
  showGallery.value = true
}

function closeGallery() {
  showGallery.value = false
}
</script>

<template>
  <div class="layout">
    <div class="map">
      <Map />
    </div>

    <div class="sidebar">
      <Sidebar />
      <button
        @click="openGallery"
        class="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:opacity-80 transition"
      >
        Open Gallery
      </button>
    </div>

    <!-- Галерея открывается поверх layout -->
    <div v-if="showGallery" class="fixed inset-0 z-50">
      <Gallery />
      <!-- Кнопка закрытия галереи -->
      <button
        @click="closeGallery"
        class="absolute top-4 right-4 px-4 py-2 bg-red-600 text-white rounded hover:opacity-80 transition"
      >
        Close
      </button>
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
  padding: 10px;
}
</style>
