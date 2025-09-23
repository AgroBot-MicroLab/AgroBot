<template>
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white rounded-xl shadow-2xl p-6 w-96 text-center">
      <h2 class="text-xl font-semibold mb-3">Mission reached</h2>
      <p class="mb-4">The drone has reached its destination</p>

      <img
          src="https://www.iconpacks.net/icons/2/free-check-mark-icon-3280-thumb.png"
          alt="Check"
          class="mx-auto mb-4 w-16 h-16"
      />

      <div class="flex gap-3 justify-center">
        <button
            class="px-4 py-2 rounded bg-emerald-600 text-white hover:bg-emerald-700"
            @click="onSave"
        >
          Save mission
        </button>
        <button
            class="px-4 py-2 rounded bg-blue-600 text-white hover:bg-blue-700"
            @click="emit('close')"
        >
          Confirm
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useMission } from "@/composables/useMission";
import { useMissionStore } from "@/stores/mission";

const emit = defineEmits(["close", "saved"]);

const { pathPts, clearPath } = useMission();
const missionStore = useMissionStore();

async function onSave() {
  await missionStore.createMission(pathPts.value);
  clearPath();
  emit("saved");
  emit("close");
}
</script>
