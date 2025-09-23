
const API = import.meta.env.VITE_API_BASE || "http://localhost:8080";

async function asJson(res) {
    if (!res.ok) throw new Error(await res.text());
    return res.json();
}

export async function fetchMissions() {
    const res = await fetch(`${API}/missions`);
    return asJson(res);
}

export async function saveMission(waypoints = []) {
    const res = await fetch(`${API}/mission`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ waypoints }),
    });
    return asJson(res); // { mission_id }
}

export async function fetchMissionPoints(id) {
    const res = await fetch(`${API}/missions/${id}/points`);
    return asJson(res); // [{lat, lon}, ...]
}

export async function deleteMission(id) {
    const res = await fetch(`${API}/missions/${id}`, { method: "DELETE" });
    if (!res.ok && res.status !== 204) throw new Error(await res.text());
}
