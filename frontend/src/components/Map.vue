<script setup>
import { computed, onBeforeUnmount } from 'vue'
import { GoogleMap, AdvancedMarker, Polyline } from 'vue3-google-map'
import { useWebSocket } from '@/composables/useWebSocket'
import { useMission } from '@/composables/useMission'

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_KEY
const wsBaseUrl = import.meta.env.VITE_API_BASE_WS

const { dronePos, targetPos, pathPts, setDronePos, addTarget } = useMission()

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
    setDronePos(data.lat, data.lon)
})

useWebSocket(`${wsBaseUrl}/drone/mission/status`, (data) => {
    console.log("Mission Reached")
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
                <img src="/drone.png" style="height:25px;width:25px;transform:translate(0%,50%);" />
            </template>
        </AdvancedMarker>
        <Polyline :options="polyOpts" />
    </GoogleMap>
</template>

