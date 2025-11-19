<template>
    <div class="mt-4 flex items-center justify-between" v-show="!deleted">
      <h1 class="text-xl font-semibold">{{ mission.id }}</h1>

      <div class="flex gap-4">
        <button
          class="bg-emerald-500 hover:bg-emerald-600 text-white font-semibold
                 py-2 px-5 rounded-xl shadow-md transition-all duration-300"
          @click="previewMission"
        >
          Select
        </button>

        <button
          class="bg-red-500 hover:bg-red-600 text-white font-semibold
                 py-2 px-5 rounded-xl shadow-md transition-all duration-300"
          @click="deleteMission"
        >
          Delete
        </button>
      </div>
    </div>
</template>

<script setup>
import { useMission } from '@/composables/useMission'
import { ref } from 'vue'

const deleted = ref(false);
const httpBaseUrl = import.meta.env.VITE_API_BASE
const props = defineProps({
  mission: {
    type: Object,
    required: true
  }
})

const { setMission } = useMission();

function previewMission() {
    setMission(props.mission.waypoints);
}

async function deleteMission() {
    await fetch(`${httpBaseUrl}/mission/${props.mission.id}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" }
    })
    deleted.value = true;
}

</script>

