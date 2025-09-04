<script setup>
import { GoogleMap, AdvancedMarker } from 'vue3-google-map'

const center = { lat: 37.7749, lng: -122.4194 }
const apiKey = import.meta.env.VITE_GOOGLE_MAPS_KEY;

</script>

<template>
  <GoogleMap
    :api-key="apiKey"
    :center="center"
    :zoom="12"
    map-id="main-map"
    style="width: 100%; height: 100vh"
  >
    <AdvancedMarker :options="{ position: center }" />
  </GoogleMap>
</template>

