<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'

const canvasRef = ref(null)
let raf = 0
let balls = []
let ctx = null
let w = 0
let h = 0

function createBall() {
  return {
    x: Math.random() * w,
    y: Math.random() * h,
    r: Math.random() * 1.8 + 0.6,
    vx: (Math.random() - 0.5) * 0.45,
    vy: (Math.random() - 0.5) * 0.45
  }
}

function resize() {
  const canvas = canvasRef.value
  if (!canvas) return
  const dpr = Math.min(window.devicePixelRatio || 1, 2)
  w = window.innerWidth
  h = window.innerHeight
  canvas.width = w * dpr
  canvas.height = h * dpr
  canvas.style.width = w + 'px'
  canvas.style.height = h + 'px'
  ctx = canvas.getContext('2d')
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0)
  const count = Math.min(70, Math.floor((w * h) / 18000))
  balls = Array.from({ length: count }, createBall)
}

function draw() {
  if (!ctx) return
  ctx.clearRect(0, 0, w, h)
  ctx.fillStyle = 'rgba(23, 114, 246, 0.28)'
  ctx.strokeStyle = 'rgba(23, 114, 246, 0.12)'
  ctx.lineWidth = 0.8

  for (let i = 0; i < balls.length; i++) {
    const a = balls[i]
    a.x += a.vx
    a.y += a.vy
    if (a.x < 0 || a.x > w) a.vx *= -1
    if (a.y < 0 || a.y > h) a.vy *= -1

    ctx.beginPath()
    ctx.arc(a.x, a.y, a.r, 0, Math.PI * 2)
    ctx.fill()

    for (let j = i + 1; j < balls.length; j++) {
      const b = balls[j]
      const dx = a.x - b.x
      const dy = a.y - b.y
      const dist = Math.sqrt(dx * dx + dy * dy)
      if (dist < 140) {
        ctx.globalAlpha = 1 - dist / 140
        ctx.beginPath()
        ctx.moveTo(a.x, a.y)
        ctx.lineTo(b.x, b.y)
        ctx.stroke()
        ctx.globalAlpha = 1
      }
    }
  }
  raf = requestAnimationFrame(draw)
}

onMounted(() => {
  resize()
  draw()
  window.addEventListener('resize', resize)
})

onBeforeUnmount(() => {
  cancelAnimationFrame(raf)
  window.removeEventListener('resize', resize)
})
</script>

<template>
  <canvas ref="canvasRef" class="particles" aria-hidden="true"></canvas>
</template>

<style scoped>
.particles {
  position: fixed;
  inset: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
}
</style>
