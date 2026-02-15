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

    render(state) {
        if (!state) return;

        this.ctx.fillStyle = 'rgba(11, 13, 20, 0.35)';
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);

        this.ctx.save();

        this.ctx.translate(this.offsetX, this.offsetY);
        this.ctx.scale(this.scaleFactor, this.scaleFactor);

        this.ctx.strokeStyle = '#333';
        this.ctx.lineWidth = 2;
        this.ctx.strokeRect(0, 0, this.worldWidth, this.worldHeight);

        // Draw Oasis zone
        this.ctx.fillStyle = 'rgba(0, 255, 100, 0.05)';
        const oW = this.worldWidth * 0.4;
        const oH = this.worldHeight * 0.4;
        this.ctx.fillRect(this.worldWidth * 0.3, this.worldHeight * 0.3, oW, oH);
        this.ctx.strokeStyle = 'rgba(0, 255, 100, 0.1)';
        this.ctx.strokeRect(this.worldWidth * 0.3, this.worldHeight * 0.3, oW, oH);

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