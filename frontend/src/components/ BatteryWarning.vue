<template>
  <div v-if="show" class="battery-modal">
    <div class="battery-modal-content">
      <h3>Rover battery is low!</h3>
      <p>Charging is required to continue operation.</p>
      <button @click="chargeBattery">Charge Battery</button>
    </div>
  </div>
</template>

<script setup>
import { defineProps, defineEmits } from "vue"

const props = defineProps({
  show: Boolean
})

const emit = defineEmits(["charging"])

function chargeBattery() {
  emit("charging")
  alert("Charging started!")
}
</script>

<style scoped>
.battery-modal {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}
.battery-modal-content {
  background: #ebebeb;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
}
.battery-modal-content h3 {
  color: #c73234;
  margin-bottom: 10px;
}
.battery-modal-content p {
  color: #c73234;
  margin-bottom: 15px;
}
.battery-modal-content button {
  background: #4CAF50;
  color: white;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
}
.battery-modal-content button:hover {
  background: #3a9b42;
}
</style>
