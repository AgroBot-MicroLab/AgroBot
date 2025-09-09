<script setup>

const httpBaseUrl = import.meta.env.VITE_API_BASE
const props = defineProps({
    point: {
        type: Object,
        default: null
    }
})

async function onGo() {
    const res = await fetch(`${httpBaseUrl}/drone/goto`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            lat: props.point.lat,
            lon: props.point.lng,
        })
    })
}

</script>

<template>
    <p>{{ point }}</p>
    <button @click="onGo()">Go</button>
</template>
