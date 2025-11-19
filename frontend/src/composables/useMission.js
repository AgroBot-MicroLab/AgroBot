import { reactive, toRefs } from 'vue'

const state = reactive({
    dronePos: null,
    targetPos: null,
    pathPts: []
})

function setMission(pts) {
    clearPath();

    for (const point of pts) {
        state.pathPts.push({ lat: point.lat, lng: point.lon })
    }
}

function setDronePos(lat, lng, yaw) {
    state.dronePos = { lat, lng, yaw }
    if (state.pathPts.length === 0) state.pathPts.push({ lat, lng, yaw })
    else {
        state.pathPts[0].lat = lat
        state.pathPts[0].lng = lng
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
        setMission,
    }
}

