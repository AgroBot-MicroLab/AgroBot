<template>
  <div class="min-h-screen bg-white p-8">
    <div class="max-w-6xl mx-auto">
      <h1 class="text-3xl font-light mb-8 text-gray-900">Gallery</h1>

      <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div
          v-for="img in images"
          :key="img.id"
          @click="selectImage(img)"
          class="aspect-square overflow-hidden cursor-pointer group relative"
        >
          <img
            :src="img.url"
            :alt="img.title"
            class="w-full h-full object-cover transition-all duration-300 group-hover:scale-110 group-hover:brightness-75"
          />
          <div class="absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-30 transition-all duration-300 flex items-center justify-center">
            <span class="text-white text-lg font-light opacity-0 group-hover:opacity-100 transition-opacity duration-300">
              {{ img.title }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div
      v-if="selected"
      @click="deselect"
      class="fixed inset-0 bg-black bg-opacity-90 flex items-center justify-center p-4 z-50"
    >
      <div class="relative max-w-4xl max-h-full" @click.stop>
        <img
          :src="selected.url"
          :alt="selected.title"
          class="max-w-full max-h-screen object-contain"
        />
        <button
          @click="deselect"
          class="absolute top-4 right-4 text-white text-3xl font-light hover:opacity-70"
        >
          ×
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const images = ref([
  { id: 1, url: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=400'},
  { id: 2, url: 'https://images.unsplash.com/photo-1469474968028-56623f02e42e?w=400'},
  { id: 3, url: 'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=400'},
  { id: 4, url: 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=400' },
  { id: 5, url: 'https://images.unsplash.com/photo-1426604966848-d7adac402bff?w=400'},
  { id: 6, url: 'https://images.unsplash.com/photo-1472214103451-9374bd1c798e?w=400'}
]);

const selected = ref(null);

function selectImage(img) {
  selected.value = img;
}

function deselect() {
  selected.value = null;
}
</script>