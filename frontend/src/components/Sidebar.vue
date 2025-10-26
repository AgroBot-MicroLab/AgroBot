<script setup>
import {ref, onMounted} from "vue"
import { useMission } from '@/composables/useMission'
const { dronePos, targetPos, pathPts, clearPath } = useMission()

const missionActive = ref(false)

const httpBaseUrl = import.meta.env.VITE_API_BASE
async function startMission() {
    await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(pathPts.value)
    })
    missionActive.value = true 
}

async function stopMission() {
    await fetch(`${httpBaseUrl}/drone/mission`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" }
    })
    missionActive.value = false
    clearPath()
}


onMounted(() => {
  const ws = new WebSocket("ws://localhost:8080/drone/mission/status")

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data)
    console.log("Mission status update:", data)
    missionActive.value = data.status
  }
})

</script>

<template>
  <div class="sidebar w-[310px] h-[900px] bg-[#EBEBEB] p-6 flex flex-col gap-6 ">

    <div class="relative flex items-center gap-4 h-[92px]">
      <div class="absolute inset-0 bg-[#D5D5D5] rounded-md"></div>
      <img src="/logo_icon.png" alt="Logo" class="w-[47px] h-[61px] relative z-10" />
      <h3 class="font-kumbh text-[24px] leading-[30px] text-black relative z-10">
        SpectralCamera
      </h3>
    </div>

    <div class="relative flex flex-col gap-1 text-center">
      <h1 class="font-encode text-[24px] leading-[30px] text-black">
        Welcome to SpectralCamera System
      </h1>
      <div class="w-[247px] h-[5px] bg-[#06B000] mt-2 mx-auto"></div>
      <p class="font-encode text-[16px] leading-[20px] text-black">
        Your automated solution for vineyard monitoring
      </p>
    </div>

    <div class="flex flex-col gap-4 mt-14">
      <h2 class="font-konkhmer text-[27px] leading-[43px] text-[#5B5B5B] text-center">Quick Start Guide</h2>

<div class="flex flex-col items-center gap-6 mt-6">
 
  <div class="flex flex-col items-center text-center">
    <div class="flex items-center justify-center gap-2">
      <div class="w-[21px] h-[21px] bg-[#8CD05E] rounded-full flex items-center justify-center text-[14px] text-[#5B5B5B] font-konkhmer">
        1
      </div>
      <p class="font-konkhmer text-[20px] leading-[36px] text-[#5B5B5B]">Set Waypoints</p>
    </div>
    <p class="font-konkhmer text-[13px] leading-[23px] text-[#06B000] mt-1">
      Right-click on the map
    </p>
  </div>

  <div class="flex flex-col items-center text-center">
    <div class="flex items-center justify-center gap-2">
      <div class="w-[21px] h-[21px] bg-[#8CD05E] rounded-full flex items-center justify-center text-[14px] text-[#5B5B5B] font-konkhmer">
        2
      </div>
      <p class="font-konkhmer text-[20px] leading-[36px] text-[#5B5B5B]">Start Mission</p>
    </div>
    <p class="font-konkhmer text-[13px] leading-[23px] text-[#06B000] mt-1">
      Press the button in the control panel
    </p>
  </div>

  <div class="flex flex-col items-center text-center">
    <div class="flex items-center justify-center gap-2">
      <div class="w-[21px] h-[21px] bg-[#8CD05E] rounded-full flex items-center justify-center text-[14px] text-[#5B5B5B] font-konkhmer">
        3
      </div>
      <p class="font-konkhmer text-[20px] leading-[36px] text-[#5B5B5B]">Monitor Live Data</p>
    </div>
    <p class="font-konkhmer text-[13px] leading-[23px] text-[#06B000] mt-1 max-w-[260px]">
      Watch the rover's movement and image capture process here
    </p>
  </div>
</div>

    </div>

    <div class="mt-4">
      <button
        :class="missionActive ? 'bg-red-600' : 'bg-[#06B000]'"
        class="w-full h-[54px] rounded-[5px] text-white font-konkhmer text-[24px] leading-[43px] transition-colors duration-300"
        @click="missionActive ? stopMission() : startMission()"
      >
        {{ missionActive ? 'Stop Mission' : 'Start Mission' }}
      </button>
    </div>

     <p class="text-[#06B000] font-normal text-[16px] text-center">
    90% faster field inspection
  </p>
  </div>
</template>

<style scoped>
.font-encode {
  font-family: 'Encode Sans', sans-serif;
}

.font-konkhmer {
  font-family: 'Konkhmer Sleokchher', sans-serif;
}

.font-kumbh {
  font-family: 'Kumbh Sans', sans-serif;
}
</style>