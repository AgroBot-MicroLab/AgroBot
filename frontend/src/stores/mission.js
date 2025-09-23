import { defineStore } from "pinia";
import { fetchMissions, saveMission } from "../services/mission";

export const useMissionStore = defineStore("mission", {
    state: () => ({ missions: [], loading: false }),
    actions: {
        async loadMissions() {
            this.loading = true;
            try { this.missions = await fetchMissions(); }
            finally { this.loading = false; }
        },
        async createMission(waypoints) {
            await saveMission(waypoints);
            await this.loadMissions();
        },
        async removeMission(id) {
            await deleteMission(id);
            await this.loadMissions();
        },
    },
});
