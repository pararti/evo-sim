export class Renderer {
    constructor(canvasId) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext('2d', { alpha: false }); // alpha: false ускоряет рендер

        this.worldWidth = 800;
        this.worldHeight = 600;

        this.resize();
        window.addEventListener('resize', () => this.resize());
    }

    resize() {
        const dpr = window.devicePixelRatio || 1;
        const parent = this.canvas.parentElement;


        this.canvas.width = parent.clientWidth * dpr;
        this.canvas.height = parent.clientHeight * dpr;

        this.canvas.style.width = `${parent.clientWidth}px`;
        this.canvas.style.height = `${parent.clientHeight}px`;


        this.ctx.scale(dpr, dpr);

        const padding = 40; // Total padding (20px each side)
        const availableW = parent.clientWidth - padding;
        const availableH = parent.clientHeight - padding;

        const scaleX = availableW / this.worldWidth;
        const scaleY = availableH / this.worldHeight;
        this.scaleFactor = Math.min(scaleX, scaleY);

        this.offsetX = (parent.clientWidth - (this.worldWidth * this.scaleFactor)) / 2;
        this.offsetY = (parent.clientHeight - (this.worldHeight * this.scaleFactor)) / 2;
    }

    setMap(mapData) {
        this.terrainCanvas = document.createElement('canvas');
        this.terrainCanvas.width = this.worldWidth;
        this.terrainCanvas.height = this.worldHeight;
        const tCtx = this.terrainCanvas.getContext('2d');

        const w = mapData.Width;
        const h = mapData.Height;
        const scale = mapData.Scale;
        
        // Decode Base64 string to Uint8Array if necessary
        let cells = mapData.Cells;
        if (typeof cells === 'string') {
            const binaryString = atob(cells);
            const len = binaryString.length;
            const bytes = new Uint8Array(len);
            for (let i = 0; i < len; i++) {
                bytes[i] = binaryString.charCodeAt(i);
            }
            cells = bytes;
        }

        for (let y = 0; y < h; y++) {
            for (let x = 0; x < w; x++) {
                const type = cells[y * w + x];
                const px = x * scale;
                const py = y * scale;

                // 0: Water, 1: Sand, 2: Grass
                if (type === 0) {
                    tCtx.fillStyle = '#1a3c6e'; // Deep Water
                } else if (type === 1) {
                    tCtx.fillStyle = '#e6c288'; // Sand
                } else {
                    tCtx.fillStyle = '#2d6e32'; // Grass
                }
                
                tCtx.fillRect(px, py, scale, scale);
            }
        }
    }

    render(state) {
        if (!state) return;

        this.ctx.fillStyle = '#111'; // Fallback bg
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

        this.ctx.save();

        this.ctx.translate(this.offsetX, this.offsetY);
        this.ctx.scale(this.scaleFactor, this.scaleFactor);

        // Draw Terrain (Cached)
        if (this.terrainCanvas) {
            this.ctx.drawImage(this.terrainCanvas, 0, 0);
        } else {
            // Draw border if no map yet
            this.ctx.strokeStyle = '#333';
            this.ctx.lineWidth = 2;
            this.ctx.strokeRect(0, 0, this.worldWidth, this.worldHeight);
        }

        // Render Food
        this.ctx.shadowBlur = 8;
        this.ctx.shadowColor = '#ffe100';
        this.ctx.fillStyle = '#ffe100';

        if (state.food) {
            for (let i = 0; i < state.food.length; i++) {
                const f = state.food[i];
                this.ctx.beginPath();
                this.ctx.arc(f.x, f.y, 3, 0, Math.PI * 2);
                this.ctx.fill();
            }
        }

        // Render Creatures
        if (state.creatures) {
            for (let i = 0; i < state.creatures.length; i++) {
                const c = state.creatures[i];
                
                if (c.isCarnivore) {
                    this.ctx.shadowColor = '#ff4d4d';
                    this.ctx.fillStyle = '#ff4d4d';
                } else {
                    this.ctx.shadowColor = '#00ff9d';
                    this.ctx.fillStyle = '#00ff9d';
                }

                const radius = 5 * (c.size || 1.0);
                this.ctx.beginPath();
                this.ctx.arc(c.x, c.y, radius, 0, Math.PI * 2);
                this.ctx.fill();
                
                // Optional: border for carnivores
                if (c.isCarnivore) {
                    this.ctx.strokeStyle = '#fff';
                    this.ctx.lineWidth = 1;
                    this.ctx.stroke();
                }
            }
        }

        this.ctx.restore();
    }
}