<script setup>
import {ref, onMounted} from "vue"
import { useMission } from '@/composables/useMission'
import Mission from '@/components/Mission.vue'
const { dronePos, targetPos, pathPts, clearPath } = useMission()

const missionActive = ref(false)

const httpBaseUrl = import.meta.env.VITE_API_BASE
async function startMission() {
    const res = await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(pathPts.value)
    })

    const data = await res.json();
    missionActive.value = true
}

async function stopMission() {
    await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" }
    })
    missionActive.value = false
    clearPath()
}


onMounted(() => {
  const ws = new WebSocket("ws://localhost:8080/drone/mission/status")

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    missionActive.value = data.status
  }
})

</script>

<template>
  <div style="margin-left: 30px;">
    <div class="flex gap-2 mt-4">
      <button
        v-show="!missionActive"
        class="bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-blue-600 hover:to-indigo-700 transition-all duration-500"
        @click="startMission()"
      >
        Start Mission
      </button>
<button
        v-show="missionActive"
        class="bg-gradient-to-r from-red-500 to-red-700 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-red-600 hover:to-red-800 transition-all duration-500"
        @click="stopMission()"
      >
        Stop Mission
      </button>
    </div>
    <Mission/>
  </div>
</template>
