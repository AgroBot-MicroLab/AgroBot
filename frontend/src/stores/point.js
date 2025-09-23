import { defineStore } from "pinia";

export const usePointStore = defineStore("point", {
    state: () => ({
        waypoints: /** @type {{lat:number, lon:number}[]} */ ([]),
    }),
    actions: {
        setWaypoints(list) { this.waypoints = list || []; },
        addPoint(p) { this.waypoints.push(p); },
        clear() { this.waypoints = []; },
    },
});
