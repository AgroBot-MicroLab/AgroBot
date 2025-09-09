import { reactive, toRefs } from 'vue'

const state = reactive({
    dronePos: null,
    targetPos: null,
    pathPts: []
})

function setDronePos(lat, lng) {
    state.dronePos = { lat, lng }
    if (state.pathPts.length === 0) state.pathPts.push({ lat, lng })
}
function addTarget(lat, lng) {
    state.targetPos = { lat, lng }
    state.pathPts.push({ lat, lng })
}
function clearPath() {
    state.targetPos = null
    state.pathPts.splice(0)
    state.pathPts.push({ lat: dronePos.lat, lng: dronePos.lng })
}

export function useMission() {
    return {
        ...toRefs(state),
        setDronePos,
        addTarget,
        clearPath,
    }
}

