<template>
  <div id="myModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
    <div class="bg-[#EBEBEB] rounded-lg shadow-lg p-6 w-[797px] h-[477px] text-center relative">

      <h2 class="font-konkhmer text-[24px] leading-[30px] text-[#000000] mb-4">
        Mission reached
      </h2>

      <p class="font-encode text-[16px] leading-[20px] text-[#000000] mb-6">
        The drone has reached its destination
      </p>

      <img src="/field.png"
           alt="Check mark"
           class="mx-auto mb-2 w-[500px] h-[260px]">

      <button @click="emit('close')"
              class="w-[183px] h-[54px] bg-[#06B000] rounded-[5px] text-white font-konkhmer text-[24px] leading-[43px] mt-2">
        Confirm
      </button>

      <p class="font-encode text-[16px] leading-[20px] text-[#06B000] mt-2">
        90% faster field inspection
      </p>
    </div>
  </div>
</template>

<script setup>
const emit = defineEmits(['close'])
</script>
