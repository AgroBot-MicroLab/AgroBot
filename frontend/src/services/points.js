const API_URL = import.meta.env.VITE_API_BASE


export async function getCurrentPoint() {
    const res = await fetch(`${API_URL}/points`);
    if (!res.ok) throw new Error('Request failed');
    return res.json();
}

function createPoint() {


}
