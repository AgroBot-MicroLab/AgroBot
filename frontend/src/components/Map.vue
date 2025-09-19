<script setup>
import { computed, onBeforeUnmount } from 'vue'
import { GoogleMap, AdvancedMarker, Polyline } from 'vue3-google-map'
import { useWebSocket } from '@/composables/useWebSocket'
import { useMission } from '@/composables/useMission'
import { ref } from 'vue'
import Modal from './Modal.vue'

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_KEY
const wsBaseUrl = import.meta.env.VITE_API_BASE_WS

const { dronePos, targetPos, pathPts, setDronePos, addTarget, clearPath } = useMission()

function onRightClick(e) {
    e.domEvent?.preventDefault?.()
    addTarget(e.latLng.lat(), e.latLng.lng())
}

const polyOpts = computed(() => ({
    path: pathPts.value.slice(),
    geodesic: true,
    strokeColor: '#FF0000',
    strokeOpacity: 1,
    strokeWeight: 2,
}))

useWebSocket(`${wsBaseUrl}/drone/position`, (data) => {
    setDronePos(data.lat, data.lon, data.yaw)
    console.log("Drone position update:", data)
})

const arrived = ref(false)
useWebSocket(`${wsBaseUrl}/drone/mission/status`, (data) => {
    if (data.is_last){
        arrived.value = true
        console.log("Mission reached",data)
        clearPath()
    }
})

onBeforeUnmount(close)
</script>

<template>
    <GoogleMap
        :api-key="apiKey"
        map-id="main-map"
        :center="{ lat: -35.363163, lng: 149.1652221 }"
        :zoom="18"
        map-type-id="satellite"
        style="width:100%; height:100vh"
        @rightclick="onRightClick"
    >
        <AdvancedMarker v-if="targetPos" :options="{ position: targetPos }" />
        <AdvancedMarker v-if="dronePos" :options="{ position: dronePos }">
            <template #content>
                <img 
                    src="/drone.png" 
                    style="height:40px;width:40px"
                    :style="{ transform: `translate(0%,50%) rotate(${dronePos.yaw+180}deg)` }" 
                />
            </template>
        </AdvancedMarker>
        <Polyline :options="polyOpts" />
    </GoogleMap>
    <Modal v-if="arrived" @close="arrived = false" />
</template>

