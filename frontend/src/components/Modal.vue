<template>
  <div id="myModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
    <div class="bg-white rounded-lg shadow-lg p-6 w-80 text-center">
      <h2 class="text-xl font-semibold mb-4">Mission reached</h2>
      <p class="mb-4">The drone has reached its destination</p>

      <img :src="imagePath"
           alt="Check mark"
           class="mx-auto mb-4 w-16 h-16">

      <button @click="emit('close')"
              class="mt-4 px-4 py-2 bg-blue-600 text-white rounded">
        Confirm
      </button>
    </div>
  </div>
</template>

<script setup>
const emit = defineEmits(['close'])
const props = defineProps({
  imagePath: {
    type: String,
    required: true
  }
})

</script>
