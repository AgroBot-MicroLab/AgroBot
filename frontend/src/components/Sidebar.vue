<script setup>
import { ref, onMounted, onBeforeUnmount} from 'vue'
import Gallery from '@/components/Gallery.vue'
import Mission from '@/components/Mission.vue'

import { useMission } from '@/composables/useMission'

const { dronePos, targetPos, pathPts, clearPath, setCurrentMissionId } = useMission()

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
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(pathPts.value)
  })

  const text = await res.text()
  console.log('startMission status =', res.status)
  console.log('startMission body =', text)

  if (!res.ok) {
    console.error('Mission start failed:', text)
    return
  }

  let data
  try {
    data = JSON.parse(text)
  } catch (e) {
    console.error('Не смог распарсить mission_id:', e)
    return
  }

  setCurrentMissionId(data.mission_id)

  missionActive.value = true
  missionsList.value = await getMissions()
}



async function stopMission() {
  await fetch(`${httpBaseUrl}/drone/mission`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' }
  })
  missionActive.value = false
  clearPath()
}

async function getMissions() {
  const res = await fetch(`${httpBaseUrl}/mission`, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' }
  })

  const data = await res.json()
  return data
}

async function reloadMissions() {
  missionsList.value = await getMissions()
}


onMounted(async () => {
  await reloadMissions()
  window.addEventListener('mission-updated', reloadMissions)
})

onBeforeUnmount(() => {
  window.removeEventListener('mission-updated', reloadMissions)
})
</script>

<template>
  <div class="sidebar-inner">
    <div class="flex flex-col gap-2 mt-4">
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
    </div>

    <div v-if="showGallery" class="fixed inset-0 z-50">
      <Gallery />
      <button
          @click="closeGallery"
          class="absolute top-4 right-4 px-4 py-2 bg-red-600 text-white rounded hover:opacity-80 transition"
      >
        Close
      </button>
    </div>

    <div class="mt-6 space-y-2">
      <Mission
          v-for="mission in missionsList"
          :key="mission.id"
          :mission="mission"
          @updated="missionsList = missionsList.filter(m => m.id !== mission.id)"
      />
    </div>
  </div>
</template>

<style scoped>
.sidebar-inner {
  display: flex;
  flex-direction: column;
}
</style>
