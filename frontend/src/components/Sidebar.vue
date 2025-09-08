<script setup>
const props = defineProps({
    point: String
})
</script>

<template>
    <p>{{ point }}</p>
    <button>Go</button>
</template>
