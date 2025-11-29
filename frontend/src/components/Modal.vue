<template>
  <div
      id="myModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
  >
    <div class="bg-white rounded-lg shadow-lg p-6 w-96 text-center">
      <h2 class="text-xl font-semibold mb-2">Mission reached</h2>
      <p class="mb-4">The drone has reached its destination</p>

      <img
          v-if="imagePath"
          :src="imagePath"
          alt="Check mark"
          class="mx-auto mb-4 max-h-40 object-contain"
      />

      <div class="text-left space-y-3 mb-2">
        <div>
          <label class="block text-sm font-medium mb-1">
            Mission name
          </label>
          <input
              v-model="name"
              type="text"
              class="w-full border rounded px-3 py-2 text-sm focus:outline-none focus:ring focus:ring-blue-300"
              placeholder="North Field – Morning Sweep"
          />
          <p v-if="error" class="text-xs text-red-500 mt-1">
            Please enter a mission name.
          </p>
        </div>

        <div>
          <label class="block text-sm font-medium mb-1">
            Description (optional)
          </label>
          <textarea
              v-model="description"
              rows="3"
              class="w-full border rounded px-3 py-2 text-sm resize-none focus:outline-none focus:ring focus:ring-blue-300"
              placeholder="Test spectral camera, check irrigation in sector B2…"
          />
        </div>
      </div>

      <div class="flex justify-end gap-3 mt-2">
        <button
            @click="onCancel"
            class="px-4 py-2 text-sm border rounded hover:bg-gray-100"
        >
          Cancel
        </button>
        <button
            @click="onConfirm"
            class="px-4 py-2 text-sm bg-blue-600 text-white rounded hover:bg-blue-700"
        >
          Confirm
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const emit = defineEmits(['close', 'confirm'])

defineProps({
  imagePath: {
    type: String,
    required: true
  }
})

const name = ref('')
const description = ref('')
const error = ref(false)

function reset () {
  name.value = ''
  description.value = ''
  error.value = false
}

function onCancel () {
  reset()
  emit('close')
}

function onConfirm () {
  if (!name.value.trim()) {
    error.value = true
    return
  }

  emit('confirm', {
    name: name.value.trim(),
    description: description.value.trim()
  })

  reset()
}
</script>
