<script setup>
import { ref, onMounted, onBeforeUnmount, computed} from 'vue'
import { GoogleMap, AdvancedMarker, Polyline } from 'vue3-google-map'
import { getCurrentPoint } from '@/services/points.js'
import { useWebSocket } from '@/composables/useWebSocket'

const apiKey = import.meta.env.VITE_GOOGLE_MAPS_KEY
const wsBaseUrl = import.meta.env.VITE_API_BASE_WS
const httpBaseUrl = import.meta.env.VITE_API_BASE
const dronePos = ref(null);
const targetPos = ref(null);

const pathPts = ref([]);

async function onRightClick(e) {
    e.domEvent?.preventDefault?.()
    targetPos.value = { lat: e.latLng.lat(), lng: e.latLng.lng() }

    pathPts.value.push({
        lat: e.latLng.lat(),
        lng: e.latLng.lng()
    })

    emit('update:point', targetPos.value)
}

const polyOpts = computed(() => {
    return {
        path: pathPts.value.slice(),
        geodesic: true,
        strokeColor: '#FF0000',
        strokeOpacity: 1,
        strokeWeight: 2,
    }
})

const props = defineProps({
    point: Object
})

const {send, close} = useWebSocket(`${wsBaseUrl}/drone/position`, (data) => {
    dronePos.value = { lat: null, lng: null };
    dronePos.value.lat = data.lat;
    dronePos.value.lng = data.lon;

    pathPts.value[0] = { lat: null, lng: null };
    pathPts.value[0].lat = data.lat;
    pathPts.value[0].lng = data.lon;
});

onBeforeUnmount(() => {
    close();
});

const emit = defineEmits(['update:point']);

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
        <AdvancedMarker :options="{ position: dronePos }">
            <template #content>
                <img src="/drone.png" style="height: 25px; width: 25px;"/>
            </template>
        </AdvancedMarker>

        <Polyline
            :options="polyOpts"
        />
    </GoogleMap>
</template>

