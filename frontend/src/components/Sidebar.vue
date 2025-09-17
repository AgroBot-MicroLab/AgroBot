<script setup>
import { useMission } from '@/composables/useMission'
const { dronePos, targetPos, pathPts, clearPath } = useMission()

const httpBaseUrl = import.meta.env.VITE_API_BASE
async function startMission() {
    await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(pathPts.value)
    })
}

</script>

<template>
  <div style="margin-left: 30px;">
    <div class="flex gap-2 mt-4">
      <button
        class="bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-blue-600 hover:to-indigo-700 transition-all duration-500"
        @click="startMission()"
      >
        Start Mission
      </button>

      <button
        class="bg-gray-200 text-gray-800 py-2 px-4 rounded-lg hover:bg-gray-300 transition-colors duration-500"
        @click="clearPath()"
      >
        Clear
      </button>
    </div>
  </div>
</template>
