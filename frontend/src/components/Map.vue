<script setup>
import { computed, onBeforeUnmount, ref } from 'vue'
import { GoogleMap, AdvancedMarker, Polyline } from 'vue3-google-map'
import { useWebSocket } from '@/composables/useWebSocket'
import { useMission } from '@/composables/useMission'
import Modal from '@/components/Modal.vue'

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_KEY
const wsBaseUrl = import.meta.env.VITE_API_BASE_WS
const apiBaseUrl = import.meta.env.VITE_API_BASE

const { dronePos, targetPos, pathPts, setDronePos, addTarget, clearPath, currentMissionId } = useMission()
const arrived = ref(false)
const image = ref('')

function onRightClick(e) {
  e.domEvent?.preventDefault?.()
  addTarget(e.latLng.lat(), e.latLng.lng())
}

const polyOpts = computed(() => ({
  path: pathPts.value.slice(),
  geodesic: true,
  strokeColor: '#FF0000',
  strokeOpacity: 1,
  strokeWeight: 2
}))

const { close: closePosWs } = useWebSocket(`${wsBaseUrl}/drone/position`, (data) => {
  setDronePos(data.lat, data.lon, data.yaw)
})

const { close: closeMissionWs } = useWebSocket(`${wsBaseUrl}/drone/mission/status`, (event) => {
  switch (event.type) {
    case 'waypoint_passed':
      if (event.data.is_last) {
        arrived.value = true
        clearPath()
        // currentMissionId.value = event.data.missionId
      }
      break
    case 'photo_received':
      image.value = event.data.path
      break
    default:
      break
  }
})
async function handleMissionConfirm(payload) {
  arrived.value = false

  if (currentMissionId.value == null) {
    console.warn('Нет currentMissionId, не могу сохранить имя миссии')
    return
  }

  try {
    const res = await fetch(`${apiBaseUrl}/mission/${currentMissionId.value}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: payload.name,
        description: payload.description,
      }),
    })

    if (!res.ok) {
      const text = await res.text()
      console.error('Update mission failed:', res.status, text)
      return
    }

    window.dispatchEvent(new Event('mission-updated'))
  } catch (e) {
    console.error('PATCH /mission failed:', e)
  }
}




onBeforeUnmount(() => {
  closePosWs()
  closeMissionWs()
})
</script>

<template>
  <GoogleMap
      :api-key="apiKey"
      map-id="main-map"
      :center="{ lat: 47.061657183060966, lng: 28.867524495508608 }"
      :zoom="18"
      map-type-id="satellite"
      style="width: 100%; height: 100vh"
      @rightclick="onRightClick"
  >
    <AdvancedMarker v-if="targetPos" :options="{ position: targetPos }" />

    <AdvancedMarker v-if="dronePos" :options="{ position: dronePos }">
      <template #content>
        <img
            src="/drone.png"
            style="height: 50px; width: 50px"
            :style="{ transform: `translate(0%,50%) rotate(${dronePos.yaw + 180}deg)` }"
        />
      </template>
    </AdvancedMarker>

    <Polyline :options="polyOpts" />
  </GoogleMap>

  <Modal
      v-if="arrived"
      :imagePath="`${apiBaseUrl}/image/${image}`"
      @close="arrived = false"
      @confirm="handleMissionConfirm"
  />
</template>

