import { Renderer } from './components/render.js';


const renderer = new Renderer('sim-canvas');

const uiAlive = document.getElementById('stat-alive');
const uiFood = document.getElementById('stat-food');
const uiFps = document.getElementById('stat-fps');
const uiUptime = document.getElementById('stat-uptime');
const uiStatus = document.getElementById('connection-status');

const startTime = Date.now();


const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
const host = window.location.host;
const wsUrl = `${protocol}://${host}/ws`;

let socket;
let lastFrameTime = performance.now();
let frameCount = 0;

// Fetch static map data
fetch('/api/map')
    .then(res => res.json())
    .then(data => {
        console.log("Map data loaded:", data);
        renderer.setMap(data);
    })
    .catch(err => console.error("Failed to load map:", err));

function connect() {
    console.log(`Connecting to ${wsUrl}...`);
    socket = new WebSocket(wsUrl);
    socket.binaryType = "arraybuffer";

    socket.onopen = () => {
        console.log("WebSocket Connected");
        uiStatus.classList.add('connected');
    };

    socket.onclose = () => {
        console.log("WebSocket Disconnected");
        uiStatus.classList.remove('connected');
        setTimeout(connect, 2000);
    };

    socket.onerror = (err) => {
        console.error("WebSocket Error:", err);
        socket.close();
    };

    socket.onmessage = (event) => {
        const buffer = event.data;
        const state = parseWorldState(buffer);

        renderer.render(state);

        updateStats(state);
    };
}


function parseWorldState(buffer) {
    const view = new DataView(buffer);
    let offset = 0;

    const creaturesCount = view.getUint16(offset, true); // true = Little Endian
    offset += 2;

    const creatures = [];

    for (let i = 0; i < creaturesCount; i++) {
        const id = view.getUint16(offset, true);
        offset += 2;

        const x = view.getFloat32(offset, true);
        offset += 4;
        const y = view.getFloat32(offset, true);
        offset += 4;

        const isCarnivore = view.getUint8(offset) === 1;
        offset += 1;

        const size = view.getFloat32(offset, true);
        offset += 4;

        creatures.push({ id, x, y, isCarnivore, size });
    }

    // --- 2. Food ---
    let food = [];

    if (offset < buffer.byteLength) {
        const foodCount = view.getUint16(offset, true);
        offset += 2;

        for (let i = 0; i < foodCount; i++) {
            const x = view.getFloat32(offset, true);
            offset += 4;
            const y = view.getFloat32(offset, true);
            offset += 4;

            food.push({ x, y });
        }
    }

    return { creatures, food };
}

function updateStats(state) {
    if (state.creatures) uiAlive.innerText = state.creatures.length;
    if (state.food) uiFood.innerText = state.food.length;

    frameCount++;
    const now = performance.now();
    if (now - lastFrameTime >= 1000) {
        uiFps.innerText = frameCount;
        frameCount = 0;
        lastFrameTime = now;

        // Update uptime
        const diff = Math.floor((Date.now() - startTime) / 1000);
        const hrs = Math.floor(diff / 3600).toString().padStart(2, '0');
        const mins = Math.floor((diff % 3600) / 60).toString().padStart(2, '0');
        const secs = (diff % 60).toString().padStart(2, '0');
        uiUptime.innerText = `${hrs}:${mins}:${secs}`;
    }
}

connect();