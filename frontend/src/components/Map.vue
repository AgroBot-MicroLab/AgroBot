<script setup>
import { ref } from 'vue'
import { GoogleMap, AdvancedMarker } from 'vue3-google-map'

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_KEY
const pos = ref(null);
const dronePos = ref(null);
dronePos.value = { lat: 46.53834103516799, lng: 29.84049779990818 };

function onRightClick(e) {
    e.domEvent?.preventDefault?.()
    pos.value = { lat: e.latLng.lat(), lng: e.latLng.lng() }
}

</script>

<template>
    <GoogleMap
        :api-key="apiKey"
        map-id="main-map"
        :center="pos || { lat: 46.53834103516799, lng: 29.84049779990818 }"
        :zoom="18"
        map-type-id="satellite"
        style="width:100%;height:100vh"
        @rightclick="onRightClick"
    >
        <AdvancedMarker v-if="pos" :options="{ position: pos }" />

        <AdvancedMarker :options="{ position: dronePos }">
            <template #content>
                <img src="/drone.png" style="height: 25px; width: 25px;"/>
            </template>
        </AdvancedMarker>

    </GoogleMap>
</template>

