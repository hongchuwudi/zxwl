<template>
  <div class="snake-animation-container">
    <div class="grid-container">
      <div
          v-for="(cell, index) in grid"
          :key="index"
          :class="['grid-cell', cell.type, { 'rain-light': cell.rainLight }]"
          :style="{ backgroundColor: getCellColor(cell) }"
      ></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import {timestamp} from "@vueuse/core";

// 网格配置
const rows = 13
const cols = 31
const grid = ref([])
const snake = ref([])
const food = ref({ row: 0, col: 0 })
const direction = ref('right')
const animationInterval = ref(null)
const mode = ref('rain-show') // 'rain-show', 'pause', 'snake-game'
const rainShowInterval = ref(null)
const rainShowTimeout = ref(null)
const showSpeed = ref(70)
const rainDrops = ref([])
const vibrantColors = [
  '#FF5252', // 鲜艳红色
  '#FF4081', // 粉红色
  '#E040FB', // 紫色
  '#7C4DFF', // 深紫色
  '#536DFE', // 蓝色
  '#448AFF', // 亮蓝色
  '#40C4FF', // 天蓝色
  '#18FFFF', // 青色
  '#64FFDA', // 绿青色
  '#69F0AE', // 亮绿色
  '#B2FF59', // 黄绿色
  '#EEFF41', // 黄色
  '#FFFF00', // 纯黄色
  '#FFD740', // 橙色
  '#FFAB40', // 橙红色
  '#FF6E40'  // 红橙色
]
const snakeColors = ref({
  head: '#1890ff',
  body: '#69c0ff'
})

const snakeColorThemes = [
  // 基础色系
  { head: '#1890ff', body: '#69c0ff' }, // 蓝色系
  { head: '#FF5252', body: '#FF8A80' }, // 红色系
  { head: '#69F0AE', body: '#B9F6CA' }, // 绿色系
  { head: '#FFD740', body: '#FFE57F' }, // 黄色系
  { head: '#E040FB', body: '#EA80FC' }, // 紫色系
  { head: '#FF6E40', body: '#FF9E80' }, // 橙色系
  // 新增色系
  { head: '#00BFA5', body: '#1DE9B6' }, // 青绿色系
  { head: '#FF4081', body: '#FF80AB' }, // 粉红色系
  { head: '#7C4DFF', body: '#B388FF' }, // 深紫色系
  { head: '#448AFF', body: '#82B1FF' }, // 亮蓝色系
  { head: '#FFC107', body: '#FFD54F' }, // 金黄色系
  // 霓虹色系
  { head: '#00E5FF', body: '#18FFFF' }, // 霓虹青色
  { head: '#76FF03', body: '#B2FF59' }, // 霓虹绿色
  { head: '#FF1744', body: '#FF5252' }, // 霓虹红色
  { head: '#F50057', body: '#FF4081' }, // 霓虹粉红
  { head: '#651FFF', body: '#7C4DFF' }, // 霓虹紫色
  // 柔和色系
  { head: '#64B5F6', body: '#90CAF9' }, // 柔和蓝色
  { head: '#4DB6AC', body: '#80CBC4' }, // 柔和绿色
  { head: '#FFB74D', body: '#FFCC80' }, // 柔和橙色
  { head: '#9575CD', body: '#B39DDB' }, // 柔和紫色
  { head: '#F06292', body: '#F48FB1' }, // 柔和粉色
  // 深色系
  { head: '#1565C0', body: '#1976D2' }, // 深蓝色
  { head: '#C62828', body: '#D32F2F' }, // 深红色
  { head: '#2E7D32', body: '#388E3C' }, // 深绿色
  { head: '#F9A825', body: '#FBC02D' }, // 深黄色
  { head: '#6A1B9A', body: '#7B1FA2' }, // 深紫色
  // 特殊效果色系
  { head: '#00B0FF', body: '#40C4FF' }, // 电光蓝
  { head: '#FF3D00', body: '#FF6E40' }, // 熔岩橙
  { head: '#00C853', body: '#00E676' }, // 荧光绿
  { head: '#AA00FF', body: '#D500F9' }, // 紫罗兰
  { head: '#FFAB00', body: '#FFD600' }, // 琥珀黄
  // 渐变色系
  { head: '#FF4081', body: '#E040FB' }, // 粉紫渐变
  { head: '#00BFA5', body: '#00E5FF' }, // 青蓝渐变
  { head: '#FF6D00', body: '#FFAB00' }, // 橙黄渐变
  { head: '#304FFE', body: '#3D5AFE' }, // 深蓝渐变
  { head: '#C51162', body: '#F50057' }  // 玫红渐变
]
let currentColorIndex = 0

// 初始化网格
const initializeGrid = () => {
  const newGrid = []
  for (let row = 0; row < rows; row++) {
    for (let col = 0; col < cols; col++) {
      newGrid.push({ row, col, type: 'empty', rainLight: false })
    }
  }
  grid.value = newGrid
}

// 更新网格状态
const updateGrid = () => {
  // 重置所有单元格
  grid.value.forEach(cell => {
    cell.type = 'empty'
  })

  // 设置食物
  const foodIndex = food.value.row * cols + food.value.col
  if (grid.value[foodIndex]) {
    grid.value[foodIndex].type = 'food'
  }

  // 设置蛇身
  snake.value.forEach((segment, index) => {
    const segmentIndex = segment.row * cols + segment.col
    if (grid.value[segmentIndex]) {
      grid.value[segmentIndex].type = index === 0 ? 'head' : 'body'
    }
  })
}

// 自动寻路算法 - 使用A*算法寻找最短路径
const findPathToFood = () => {
  const head = snake.value[0]
  const queue = [{ ...head, path: [] }]
  const visited = new Set()
  visited.add(`${head.row},${head.col}`)

  while (queue.length > 0) {
    const current = queue.shift()

    // 如果找到食物，返回路径
    if (current.row === food.value.row && current.col === food.value.col) {
      return current.path
    }

    // 检查四个方向
    const directions = [
      { row: -1, col: 0, dir: 'up' },
      { row: 1, col: 0, dir: 'down' },
      { row: 0, col: -1, dir: 'left' },
      { row: 0, col: 1, dir: 'right' }
    ]

    for (const dir of directions) {
      const newRow = (current.row + dir.row + rows) % rows
      const newCol = (current.col + dir.col + cols) % cols
      const key = `${newRow},${newCol}`

      // 检查是否已访问或是否是蛇身（除了尾部，因为移动后尾部会空出来）
      if (!visited.has(key) && !isSnakeBody(newRow, newCol, true)) {
        visited.add(key)
        queue.push({
          row: newRow,
          col: newCol,
          path: [...current.path, dir.dir]
        })
      }
    }
  }

  // 如果没有找到路径，返回安全移动方向
  return getSafeDirection()
}

// 检查是否是蛇身（排除尾部）
const isSnakeBody = (row, col, excludeTail = false) => {
  const endIndex = excludeTail ? snake.value.length - 1 : snake.value.length
  return snake.value.slice(0, endIndex).some(segment =>
      segment.row === row && segment.col === col
  )
}

// 获取安全移动方向（避免撞到自己）
const getSafeDirection = () => {
  const head = snake.value[0]
  const directions = [
    { row: -1, col: 0, dir: 'up' },
    { row: 1, col: 0, dir: 'down' },
    { row: 0, col: -1, dir: 'left' },
    { row: 0, col: 1, dir: 'right' }
  ]

  // 随机尝试各个方向，找到安全的移动方向
  const safeDirections = directions.filter(dir => {
    const newRow = (head.row + dir.row + rows) % rows
    const newCol = (head.col + dir.col + cols) % cols
    return !isSnakeBody(newRow, newCol, true)
  })

  if (safeDirections.length > 0) {
    // 优先保持当前方向
    const currentDir = safeDirections.find(d => d.dir === direction.value)
    if (currentDir) {
      return [direction.value]
    }
    // 否则返回随机安全方向
    return [safeDirections[Math.floor(Math.random() * safeDirections.length)].dir]
  }

  // 如果没有安全方向，返回当前方向（游戏结束）
  return [direction.value]
}

// 移动蛇
const moveSnake = () => {
  // 获取路径
  const path = findPathToFood()
  if (path && path.length > 0) {
    direction.value = path[0]
  }

  // 获取蛇头
  const head = {...snake.value[0]}

  // 根据方向移动蛇头
  switch (direction.value) {
    case 'up': head.row = (head.row - 1 + rows) % rows; break
    case 'down': head.row = (head.row + 1) % rows; break
    case 'left': head.col = (head.col - 1 + cols) % cols; break
    case 'right': head.col = (head.col + 1) % cols; break
  }

  // 检查是否吃到食物
  const ateFood = head.row === food.value.row && head.col === food.value.col

  // 移动蛇
  snake.value.unshift(head)
  if (!ateFood) {
    snake.value.pop()
  } else {
    // 吃到食物时变换颜色
    changeSnakeColor()
    // 生成新食物
    generateFood()
  }

  // 更新网格
  updateGrid()

  // 头部方向指示
  updateHeadDirection()
}

// 更新头部方向指示
const updateHeadDirection = () => {
  const headIndex = snake.value[0].row * cols + snake.value[0].col
  if (grid.value[headIndex]) {
    grid.value[headIndex].direction = direction.value
  }
}

// 生成食物
const generateFood = () => {
  let newFood
  do {
    newFood = {
      row: Math.floor(Math.random() * rows),
      col: Math.floor(Math.random() * cols)
    }
    // 确保食物不会出现在蛇身上
  } while (snake.value.some(segment => segment.row === newFood.row && segment.col === newFood.col))

  food.value = newFood
}

// 下雨灯效 - 使用setInterval控制速度
// 下雨灯效 - 使用对象跟踪每个雨滴
const startRainShow = () => {
  let startTime = Date.now()
  const duration = 3000 // 2秒总时长
  mode.value = 'rain-show'
  rainDrops.value = [] // 清空雨滴数组

  // 清除之前的interval（如果有）
  if (rainShowInterval.value) clearInterval(rainShowInterval.value)

  // 使用setInterval精确控制速度
  rainShowInterval.value = setInterval(() => {
    const elapsed = Date.now() - startTime
    const progress = elapsed / duration

    if (elapsed >= duration) {
      clearInterval(rainShowInterval.value)
      rainShowInterval.value = null
      mode.value = 'pause'
      rainDrops.value = []

      // 清除所有雨滴效果
      grid.value.forEach(cell => {
        cell.rainLight = false
        cell.rainColor = null
      })

      // 0.1秒后开始贪吃蛇游戏
      rainShowTimeout.value = setTimeout(() => {
        mode.value = 'snake-game'
        initializeSnake()
        generateFood()
        updateGrid()
        startAnimation()
      }, 100)
      return
    }

    // 清除上一帧的所有雨滴效果
    grid.value.forEach(cell => {
      cell.rainLight = false
      cell.rainColor = null
    })

    // 控制密度：随时间增加
    const intensity = 6 + progress * 8

    // 随机生成新雨滴
    if (Math.random() < intensity * 0.3) {
      createNewRainDrop()
    }

    // 更新所有现有雨滴的位置
    updateRainDrops()

  }, 30) // 70ms的更新间隔
}

// 创建新雨滴
const createNewRainDrop = () => {
  const edgeType = Math.floor(Math.random() * 4) // 0:上, 1:下, 2:左, 3:右
  const color = vibrantColors[Math.floor(Math.random() * vibrantColors.length)]

  let rainDrop = {
    id: Date.now() + Math.random(),
    color: color,
    progress: 0,
    speed: 0.03 + Math.random() * 0.02,
    length: 3 + Math.floor(Math.random() * 3)
  }

  switch (edgeType) {
    case 0: // 上边缘
      rainDrop = {
        ...rainDrop,
        startCol: Math.floor(Math.random() * cols),
        direction: 'down',
        row: 0
      }
      break
    case 1: // 下边缘
      rainDrop = {
        ...rainDrop,
        startCol: Math.floor(Math.random() * cols),
        direction: 'up',
        row: rows - 1
      }
      break
    case 2: // 左边缘
      rainDrop = {
        ...rainDrop,
        startRow: Math.floor(Math.random() * rows),
        direction: 'right',
        col: 0
      }
      break
    case 3: // 右边缘
      rainDrop = {
        ...rainDrop,
        startRow: Math.floor(Math.random() * rows),
        direction: 'left',
        col: cols - 1
      }
      break
  }

  rainDrops.value.push(rainDrop)
}

// 更新所有雨滴位置
const updateRainDrops = () => {
  const updatedDrops = []

  rainDrops.value.forEach(drop => {
    drop.progress += drop.speed

    if (drop.progress > 1) {
      return // 雨滴已经到达终点，移除
    }

    updatedDrops.push(drop)
    renderRainDrop(drop)
  })

  rainDrops.value = updatedDrops
}

// 渲染单个雨滴
const renderRainDrop = (drop) => {
  switch (drop.direction) {
    case 'down':
      for (let i = 0; i < drop.length; i++) {
        const row = Math.floor(drop.row + drop.progress * (rows - 1) - i)
        if (row >= 0 && row < rows) {
          const cellIndex = row * cols + drop.startCol
          if (grid.value[cellIndex]) {
            grid.value[cellIndex].rainLight = true
            grid.value[cellIndex].rainColor = drop.color
            grid.value[cellIndex].rainOpacity = 1 - (i * 0.3)
          }
        }
      }
      break

    case 'up':
      for (let i = 0; i < drop.length; i++) {
        const row = Math.floor(drop.row - drop.progress * (rows - 1) + i)
        if (row >= 0 && row < rows) {
          const cellIndex = row * cols + drop.startCol
          if (grid.value[cellIndex]) {
            grid.value[cellIndex].rainLight = true
            grid.value[cellIndex].rainColor = drop.color
            grid.value[cellIndex].rainOpacity = 1 - (i * 0.3)
          }
        }
      }
      break

    case 'right':
      for (let i = 0; i < drop.length; i++) {
        const col = Math.floor(drop.col + drop.progress * (cols - 1) - i)
        if (col >= 0 && col < cols) {
          const cellIndex = drop.startRow * cols + col
          if (grid.value[cellIndex]) {
            grid.value[cellIndex].rainLight = true
            grid.value[cellIndex].rainColor = drop.color
            grid.value[cellIndex].rainOpacity = 1 - (i * 0.3)
          }
        }
      }
      break

    case 'left':
      for (let i = 0; i < drop.length; i++) {
        const col = Math.floor(drop.col - drop.progress * (cols - 1) + i)
        if (col >= 0 && col < cols) {
          const cellIndex = drop.startRow * cols + col
          if (grid.value[cellIndex]) {
            grid.value[cellIndex].rainLight = true
            grid.value[cellIndex].rainColor = drop.color
            grid.value[cellIndex].rainOpacity = 1 - (i * 0.3)
          }
        }
      }
      break
  }
}

// 变换蛇的颜色（吃到食物时调用）-- 随机调用
const changeSnakeColor = () => {
  // 随机选择下一个颜色（确保不与当前颜色相同）
  let newIndex
  do {
    newIndex = Math.floor(Math.random() * snakeColorThemes.length)
  } while (newIndex === currentColorIndex && snakeColorThemes.length > 1)

  currentColorIndex = newIndex
  snakeColors.value = snakeColorThemes[currentColorIndex]
}

// 变换蛇的颜色（顺序切换版本） -- 顺序调用
// const changeSnakeColor = () => {
//   currentColorIndex = (currentColorIndex + 1) % snakeColorThemes.length
//   snakeColors.value = snakeColorThemes[currentColorIndex]
// }

// 获取单元格颜色和样式
const getCellColor = (cell) => {
  if (cell.rainLight && cell.rainColor) {
    const opacity = cell.rainOpacity !== undefined ? cell.rainOpacity : 0.8
    return cell.rainColor.replace(')', `, ${opacity})`).replace('hsl', 'hsla')
  }

  switch (cell.type) {
    case 'head':
      return snakeColors.value.head
    case 'body':
      return snakeColors.value.body
    case 'food':
      return '#4fb80d' // 食物改为更醒目的红色
    default:
      return '#f0f0f0'
  }
}
// 初始化蛇
const initializeSnake = () => {
  // 初始长度为10，水平放置在中间位置
  const startCol = Math.floor(cols / 2) - 5
  const startRow = Math.floor(rows / 2)
  snake.value = []

  for (let i = 0; i < 10; i++) {
    snake.value.push({ row: startRow, col: startCol + i })
  }

  direction.value = 'right'
  // 重置为初始颜色
  currentColorIndex = 0
  snakeColors.value = snakeColorThemes[0]
}

// 启动动画
const startAnimation = () => animationInterval.value = setInterval(moveSnake, showSpeed.value)

// 停止动画
const stopAnimation = () => {
  if (animationInterval.value) {
    clearInterval(animationInterval.value)
    animationInterval.value = null
  }
}

// 重新开始动画循环
const restartAnimation = () => {
  stopAnimation()
  // 删除 stopColorChange() 调用
  currentColorIndex = 0 // 重置颜色索引
  snakeColors.value = snakeColorThemes[0] // 重置为初始颜色
  mode.value = 'rain-show'
  startRainShow()
}

// 监听模式变化，游戏结束后重新开始
onMounted(() => {
  initializeGrid()
  startRainShow()
})

onUnmounted(() => {
  stopAnimation()
  if (rainShowInterval.value) {
    clearInterval(rainShowInterval.value)
    rainShowInterval.value = null
  }
  if (rainShowTimeout.value) clearTimeout(rainShowTimeout.value)
})
</script>

<style scoped>
.snake-animation-container {
  padding: 10px;
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.grid-container {
  display: grid;
  grid-template-columns: repeat(31, 1fr);
  grid-template-rows: repeat(13, 1fr);
  gap: 1px;
  width: 100%;
  aspect-ratio: 31 / 13;
}

.grid-cell {
  border-radius: 2px;
  transition: all 0.2s ease;
}

.grid-cell.head {
  box-shadow: 0 0 0 2px #000, 0 0 5px 3px rgba(255, 255, 255, 0.7);
  z-index: 2;
  transform: scale(1.05);
}

.grid-cell.body {
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.3);
  z-index: 1;
}

.grid-cell.food {
  box-shadow: 0 0 0 2px #000, 0 0 8px 4px rgba(85, 255, 0, 0.5);
  animation: pulse 1s infinite alternate;
  z-index: 3;
}

/* 添加头部方向指示的CSS */
.grid-cell.head::after {
  content: '';
  position: absolute;
  width: 30%;
  height: 30%;
  background: #fa4545;
  border-radius: 50%;
}

.grid-cell.head[data-direction="right"]::after {
  right: 10%;
  top: 50%;
  transform: translateY(-50%);
}

.grid-cell.head[data-direction="left"]::after {
  left: 10%;
  top: 50%;
  transform: translateY(-50%);
}

.grid-cell.head[data-direction="up"]::after {
  top: 10%;
  left: 50%;
  transform: translateX(-50%);
}

.grid-cell.head[data-direction="down"]::after {
  bottom: 10%;
  left: 50%;
  transform: translateX(-50%);
}

@keyframes pulse {
  0% {
    transform: scale(1);
    opacity: 1;
  }
  100% {
    transform: scale(1.1);
    opacity: 0.8;
  }
}
@keyframes raindrop {
  0% {
    opacity: 0.3;
  }
  50% {
    opacity: 0.8;
  }
  100% {
    opacity: 0.3;
  }
}
</style>