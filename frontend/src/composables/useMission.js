import { reactive, toRefs } from 'vue'

const state = reactive({
    dronePos: null,
    targetPos: null,
    pathPts: []
})

function setDronePos(lat, lng, yaw) {
    state.dronePos = { lat, lng, yaw }
    if (state.pathPts.length === 0) state.pathPts.push({ lat, lng, yaw })
    else {
        state.pathPts[0].lat = lat
        state.pathPts[0].lng = lng
        state.pathPts[0].yaw = yaw
    }
}

function addTarget(lat, lng) {
    state.targetPos = { lat, lng }
    state.pathPts.push({ lat, lng })
}

function clearPath() {
    state.targetPos = null
    state.pathPts.splice(0)
    state.pathPts.push({ lat: state.dronePos.lat, lng: state.dronePos.lng })
}

export function useMission() {
    return {
        ...toRefs(state),
        setDronePos,
        addTarget,
        clearPath,
    }
}

