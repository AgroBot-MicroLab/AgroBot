<script setup>
import { useMission } from '@/composables/useMission'
const { dronePos, targetPos, pathPts, clearPath } = useMission()

const httpBaseUrl = import.meta.env.VITE_API_BASE
async function startMission() {
    await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(pathPts.value)
    })
}

</script>

<template>
    {{ targetPos }}
    <button @click="startMission()">Start Mission</button>
    <button @click="clearPath()">Clear</button>
</template>
