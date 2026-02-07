import { Renderer } from './components/renderer.js';


const renderer = new Renderer('sim-canvas');

const uiAlive = document.getElementById('stat-alive');
const uiFood = document.getElementById('stat-food');
const uiFps = document.getElementById('stat-fps');
const uiStatus = document.getElementById('connection-status');


const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
const host = window.location.host;
const wsUrl = `${protocol}://${host}/ws`;

let socket;
let lastFrameTime = performance.now();
let frameCount = 0;

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

        creatures.push({ id, x, y });
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
    }
}

connect();