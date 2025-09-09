import { defineStore } from 'pinia'

export const useMission = defineStore('mission', {
    state: () => ({ point: null as null | { lat: number; lng: number } }),
})
