<script setup>
import { ref, onMounted } from "vue";
import { useMission } from "@/composables/useMission";
import { useMissionStore } from "../stores/mission";
import { fetchMissionPoints } from "../services/mission";

const { pathPts, clearPath } = useMission();
const missionStore = useMissionStore();

const missionActive = ref(false);
const API = import.meta.env.VITE_API_BASE || "http://localhost:8080";

async function startMission() {
  await fetch(`${API}/drone/mission`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(pathPts.value),
  });
  missionActive.value = true;
}

async function stopMission() {
  await fetch(`${API}/drone/mission`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
  });
  missionActive.value = false;
  clearPath();
}

async function runMission(id) {
  const points = await fetchMissionPoints(id);
  await fetch(`${API}/drone/mission`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(points),
  });
}

onMounted(() => missionStore.loadMissions());
</script>

<template>
  <div class="flex gap-6 mt-6 ml-6">
    <div>
      <button
          v-show="!missionActive"
          class="bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-blue-600 hover:to-indigo-700 transition-all duration-500"
          @click="startMission"
      >
        Start Mission
      </button>
      <button
          v-show="missionActive"
          class="bg-gradient-to-r from-red-500 to-red-700 text-white font-semibold py-2 px-4 rounded-lg shadow-md hover:from-red-600 hover:to-red-800 transition-all duration-500"
          @click="stopMission"
      >
        Stop Mission
      </button>
    </div>

    <aside class="w-72 min-h-screen p-4 bg-gradient-to-b from-gray-100 to-gray-300 shadow-lg">
      <h3 class="font-bold mb-2 text-gray-800">Saved Missions</h3>
      <ul class="space-y-2">
        <<li v-for="m in missionStore.missions" :key="m.id" class="flex items-center gap-2">
        <button class="px-2 py-1 bg-emerald-600 text-white rounded" @click="runMission(m.id)">Run</button>
        <button class="px-2 py-1 bg-red-600 text-white rounded" @click="missionStore.removeMission(m.id)">Delete</button>
        <span>#{{ m.id }} — {{ new Date(m.created_at).toLocaleString() }}</span>
      </li>
      </ul>
    </aside>


  </div>
</template>
