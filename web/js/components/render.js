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

        const scaleX = parent.clientWidth / this.worldWidth;
        const scaleY = parent.clientHeight / this.worldHeight;
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

        this.ctx.shadowBlur = 10;
        this.ctx.shadowColor = '#ff0055';
        this.ctx.fillStyle = '#ff0055';

        if (state.food) {
            for (let i = 0; i < state.food.length; i++) {
                const f = state.food[i];
                this.ctx.beginPath();
                this.ctx.arc(f.x, f.y, 4, 0, Math.PI * 2);
                this.ctx.fill();
            }
        }

        this.ctx.shadowColor = '#00ff9d';
        this.ctx.fillStyle = '#00ff9d';

        if (state.creatures) {
            for (let i = 0; i < state.creatures.length; i++) {
                const c = state.creatures[i];
                this.ctx.beginPath();
                this.ctx.arc(c.x, c.y, 5, 0, Math.PI * 2);
                this.ctx.fill();
            }
        }

        this.ctx.restore();
    }
}