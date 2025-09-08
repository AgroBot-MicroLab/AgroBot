<script setup>
import { ref, onMounted } from 'vue'
import { GoogleMap, AdvancedMarker } from 'vue3-google-map'
import { getCurrentPoint } from '../services/points.js'

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_KEY
const targetPos = ref(null);
const dronePos = ref(null);
dronePos.value = { lat: 46.53834103516799, lng: 29.84049779990818 };

onMounted(async () => {
    const [point] = await getCurrentPoint();
    targetPos.value = { lat: null, lng: null };
    targetPos.value.lat = point.latitude;
    targetPos.value.lng = point.longitude;
});

function onRightClick(e) {
    e.domEvent?.preventDefault?.()
    pos.value = { lat: e.latLng.lat(), lng: e.latLng.lng() }
    emit('update:point', pos.value)
}

const props = defineProps({
    point: Object
})

const emit = defineEmits(['update:point']);
</script>

<template>
    <GoogleMap
        :api-key="apiKey"
        map-id="main-map"
        :center="targetPos || { lat: 46.53834103516799, lng: 29.84049779990818 }"
        :zoom="18"
        map-type-id="satellite"
        style="width:100%; height:100vh"
        @rightclick="onRightClick"
    >
        <AdvancedMarker v-if="targetPos" :options="{ position: targetPos }" />
        <AdvancedMarker :options="{ position: dronePos }">
            <template #content>
                <img src="/drone.png" style="height: 25px; width: 25px;"/>
            </template>
        </AdvancedMarker>
    </GoogleMap>
</template>

