<script setup>
import {ref, onMounted} from "vue"
import { useMission } from '@/composables/useMission'
import Gallery from '@/components/Gallery.vue'
import Mission from '@/components/Mission.vue'
const { dronePos, targetPos, pathPts, clearPath } = useMission()

const missionActive = ref(false)
const missionsList = ref([])
const showGallery = ref(false)
const httpBaseUrl = import.meta.env.VITE_API_BASE

function openGallery() {
  showGallery.value = true
}

function closeGallery() {
  showGallery.value = false
}

async function startMission() {
    const res = await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(pathPts.value)
    })

    const data = await res.json();
    missionActive.value = true
    missionsList.value = await getMissions();
}

async function stopMission() {
    await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" }
    })
    missionActive.value = false
    clearPath()
}

async function getMissions() {
    const res = await fetch(`${httpBaseUrl}/mission`, {
        method: "GET",
        headers: { "Content-Type": "application/json" }
    })

    const data = await res.json();
    return data;
}

onMounted(async () => {
  const ws = new WebSocket("ws://localhost:8080/drone/mission/status")
  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    missionActive.value = data.status
  }

  const missions = await getMissions();
  missionsList.value = missions;
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
        class="g-gradient-to-r from-red-500 to-red-700 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-red-600 hover:to-red-800 transition-all duration-500"
        @click="stopMission()"
      >
        Stop Mission
      </button>

      <button
        class="bg-gradient-to-r from-orange-500 to-amber-600 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-orange-600 hover:to-amber-700 transition-all duration-500"
        @click="clearPath()"
      >
        Clear
      </button>

      <button
        @click="openGallery"
        class="bg-gradient-to-r from-teal-500 to-cyan-600 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-teal-600 hover:to-cyan-700 transition-all duration-500"
      >
        Open Gallery
      </button>

      <div v-if="showGallery" class="fixed inset-0 z-50">
        <Gallery />
        <button
          @click="closeGallery"
          class="absolute top-4 right-4 px-4 py-2 bg-red-600 text-white rounded hover:opacity-80 transition"
        >
          Close
        </button>
      </div>
    </div>

    <div v-for="mission in missionsList">
        <Mission :mission="mission" />
    </div>

  </div>
</template>
