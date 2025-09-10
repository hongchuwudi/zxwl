export class ParticleEngine {
    constructor(canvas) {
        this.canvas = canvas
        this.ctx = canvas.getContext('2d')
        this.particles = []
        this.mouse = { x: null, y: null, radius: 100 }
        this.animationFrameId = null

        // 可配置参数
        this.config = {
            particleCount: 150,
            particleColor: 'rgba(255, 255, 255, 0.8)',
            lineColor: 'rgba(255, 255, 255, 0.3)',
            maxLineDistance: 100,
            baseSpeed: 0.5,
            mouseForce: 0.1
        }
    }

    init() {
        this._setupCanvas()
        this._createParticles()
        this._bindEvents()
        this.start()
    }

    _setupCanvas() {
        this.canvas.width = window.innerWidth
        this.canvas.height = window.innerHeight
    }

    _createParticles() {
        for (let i = 0; i < this.config.particleCount; i++) {
            this.particles.push({
                x: Math.random() * this.canvas.width,
                y: Math.random() * this.canvas.height,
                vx: (Math.random() - 0.5) * this.config.baseSpeed,
                vy: (Math.random() - 0.5) * this.config.baseSpeed,
                radius: Math.random() * 1.5 + 0.5,
                baseX: Math.random() * this.canvas.width,
                baseY: Math.random() * this.canvas.height
            })
        }
    }

    _bindEvents() {
        window.addEventListener('resize', this._handleResize.bind(this))
        window.addEventListener('mousemove', this._handleMouseMove.bind(this))
    }

    _handleResize() {
        this.canvas.width = window.innerWidth
        this.canvas.height = window.innerHeight
    }

    _handleMouseMove(event) {
        this.mouse.x = event.clientX
        this.mouse.y = event.clientY
    }

    start() {
        this._animate()
    }

    _animate() {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height)
        this._updateParticles()
        this._drawLines()
        this.animationFrameId = requestAnimationFrame(this._animate.bind(this))
    }

    _updateParticles() {
        this.particles.forEach(particle => {
            // 基础运动
            particle.x += particle.vx
            particle.y += particle.vy

            // 边界反弹
            if (particle.x < 0 || particle.x > this.canvas.width) particle.vx *= -1
            if (particle.y < 0 || particle.y > this.canvas.height) particle.vy *= -1

            // 鼠标互动
            if (this.mouse.x && this.mouse.y) {
                const dx = this.mouse.x - particle.x
                const dy = this.mouse.y - particle.y
                const distance = Math.sqrt(dx * dx + dy * dy)

                if (distance < this.mouse.radius) {
                    const force = (this.mouse.radius - distance) / this.mouse.radius
                    particle.vx -= (dx / distance) * force * this.config.mouseForce
                    particle.vy -= (dy / distance) * force * this.config.mouseForce
                }
            }

            // 回归基础位置
            particle.x += (particle.baseX - particle.x) * 0.01
            particle.y += (particle.baseY - particle.y) * 0.01
        })
    }

    _drawLines() {
        this.particles.forEach(p1 => {
            this.particles.forEach(p2 => {
                const dx = p1.x - p2.x
                const dy = p1.y - p2.y
                const distance = Math.sqrt(dx * dx + dy * dy)

                if (distance < this.config.maxLineDistance) {
                    this.ctx.beginPath()
                    this.ctx.strokeStyle = this.config.lineColor
                    this.ctx.lineWidth = 1 - (distance / this.config.maxLineDistance)
                    this.ctx.moveTo(p1.x, p1.y)
                    this.ctx.lineTo(p2.x, p2.y)
                    this.ctx.stroke()
                }
            })
        })

        this.particles.forEach(particle => {
            this.ctx.beginPath()
            this.ctx.arc(particle.x, particle.y, particle.radius, 0, Math.PI * 2)
            this.ctx.fillStyle = this.config.particleColor
            this.ctx.fill()
        })
    }

    destroy() {
        cancelAnimationFrame(this.animationFrameId)
        window.removeEventListener('resize', this._handleResize)
        window.removeEventListener('mousemove', this._handleMouseMove)
    }
}