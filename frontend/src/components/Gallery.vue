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
import { ref, onMounted } from 'vue';

const httpBaseUrl = import.meta.env.VITE_API_BASE
const images = ref([]);

onMounted(async () => {
    await getImages();
});

async function getImages() {
    const res = await fetch(`${httpBaseUrl}/image`, {
        method: "GET",
        headers: { "Content-Type": "application/json" }
    })
    const data = await res.json();

    for (let i = 0; i < data.length; ++i) {
        images.value.push({ id: i, url: `${httpBaseUrl}/image/${data[i]}` });
    }
}

const selected = ref(null);

function selectImage(img) {
  selected.value = img;
}

function deselect() {
  selected.value = null;
}
</script>
