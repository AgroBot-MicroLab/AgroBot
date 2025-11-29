<template>
  <div
      class="mt-4 border rounded-xl px-4 py-3 bg-white shadow-sm"
      v-show="!deleted"
  >
    <div class="flex items-center justify-between">
      <button
          class="flex items-center gap-2"
          @click="expanded = !expanded"
      >
        <span
            class="inline-block transition-transform"
            :class="expanded ? 'rotate-90' : ''"
        >
          ▶
        </span>

        <h1 class="text-lg font-semibold">
          {{ missionTitle }}
        </h1>
      </button>

      <div class="flex gap-4">
        <button
            class="bg-emerald-500 hover:bg-emerald-600 text-white font-semibold
                 py-2 px-5 rounded-xl shadow-md transition-all duration-300"
            @click="previewMission"
        >
          Select
        </button>

        <button
            class="bg-red-500 hover:bg-red-600 text-white font-semibold
                 py-2 px-5 rounded-xl shadow-md transition-all duration-300"
            @click="deleteMission"
        >
          Delete
        </button>
      </div>
    </div>

    <div v-if="expanded" class="mt-3 text-sm text-gray-700 space-y-1">
      <p v-if="mission.description">
        <span class="font-medium">Description:</span>
        {{ mission.description }}
      </p>

      <p v-if="mission.savedAt || mission.createdAt">
        <span class="font-medium">Saved at:</span>
        {{ formatDate(mission.savedAt || mission.createdAt) }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useMission } from '@/composables/useMission'

const deleted = ref(false)
const expanded = ref(false)
const emit = defineEmits(['updated'])

const httpBaseUrl = import.meta.env.VITE_API_BASE

const props = defineProps({
  mission: {
    type: Object,
    required: true
  }
})

const { setMission } = useMission()

const missionTitle = computed(() =>
    props.mission.name?.trim()
        ? props.mission.name
        : `Mission #${props.mission.id}`
)

function previewMission () {
  setMission(props.mission.waypoints)
}

async function deleteMission () {
  await fetch(`${httpBaseUrl}/mission/${props.mission.id}`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' }
  })
  deleted.value = true
  emit('updated')
}

function formatDate (value) {
  if (!value) return ''
  return new Date(value).toLocaleString()
}
</script>
