<template>
  <div class="map-container">
    <div id="map" ref="mapContainer"></div>
    <!-- 加载状态提示 -->
    <div v-if="loading" class="map-loading">
      <el-icon class="is-loading"><Loading /></el-icon>
      <span>地图加载中...</span>
    </div>

    <el-drawer
        v-model="drawerVisible"
        title="位置信息"
        direction="rtl"
        size="380px"
        class="info-drawer"
    >
      <div class="drawer-content">
        <el-divider />

        <div class="location-info">
          <div class="info-item">
            <span class="label">地址：</span>
            <span class="value">{{ props.address }}</span>
          </div>

          <div class="info-item" v-if="currentLocation.lng">
            <span class="label">经度：</span>
            <span class="value">{{ currentLocation.lng.toFixed(6) }}</span>
          </div>

          <div class="info-item" v-if="currentLocation.lat">
            <span class="label">纬度：</span>
            <span class="value">{{ currentLocation.lat.toFixed(6) }}</span>
          </div>

          <div class="info-item" v-if="locationDetails">
            <span class="label">详细信息：</span>
            <span class="value">{{ locationDetails }}</span>
          </div>
        </div>

        <el-divider />

        <div class="actions">
          <el-button type="primary" @click="centerMap">重新定位</el-button>
          <el-button @click="drawerVisible = false">关闭</el-button>
        </div>
      </div>
    </el-drawer>

    <div class="control-bar">
      <el-button @click="drawerVisible = !drawerVisible" icon="Info">
        位置信息
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Loading } from '@element-plus/icons-vue'
import { decryptAPIKey } from '@/utils/decryptAPIKey.js'
import axios from "axios";

// 定义组件属性
const props = defineProps({
  address: {
    type: String,
    required: true
  },
  zoom: {
    type: Number,
    default: 15
  },
  showMarker: {
    type: Boolean,
    default: true
  },
  showInfoWindow: {
    type: Boolean,
    default: true
  }
})

// 响应式数据
const drawerVisible = ref(false)
const mapContainer = ref(null)
const map = ref(null)
const currentLocation = ref({ lng: null, lat: null })
const locationDetails = ref('')
const geocoder = ref(null)
const marker = ref(null)
const infoWindow = ref(null)
const loading = ref(false)
const isComponentMounted = ref(true)
const bMapScriptLoaded = ref(false)
const baiduKey = ref('')

// 组件卸载时清理
onUnmounted(() => {
  isComponentMounted.value = false
  cleanupMap()
})

// 清理地图资源
const cleanupMap = () => {
  if (map.value) {
    try {
      map.value.clearOverlays()
      map.value.destroy()
    } catch (e) {
      console.warn('地图清理过程中出现警告:', e)
    }
    map.value = null
  }

  // 清理全局回调
  if (window._bMapCallback) {
    window._bMapCallback = null
  }
}

// 检查百度地图API是否已加载
const checkBMapLoaded = () => {
  return typeof BMap !== 'undefined' && typeof BMap.Map === 'function'
}

// 获取并解密API密钥
const getAndDecryptAPIKey = async () => {
  try {
    const res = await axios.get('/gapi/getBaiDuKey')
    baiduKey.value = decryptAPIKey(res.data.data)
    // console.log('百度地图API密钥:', baiduKey.value)
    return true
  } catch (error) {
    console.error('获取百度API密钥失败:', error)
    ElMessage.error('获取地图服务密钥失败')
    return false
  }
}

// 加载百度地图脚本
const loadBMapScript = () => {
  return new Promise((resolve, reject) => {
    if (checkBMapLoaded()) {
      resolve()
      return
    }

    // 检查密钥是否已获取
    if (!baiduKey.value) {
      reject(new Error('百度地图API密钥未获取'))
      return
    }

    // 使用唯一的回调函数名称
    const callbackName = '_bMapCallback_' + Date.now()

    window[callbackName] = function() {
      // 清理回调
      delete window[callbackName]

      // 给API一点时间完全初始化
      setTimeout(() => {
        if (checkBMapLoaded()) {
          resolve()
        } else {
          reject(new Error('百度地图API加载但未正确初始化'))
        }
      }, 100)
    }

    const script = document.createElement('script')
    script.src = `https://api.map.baidu.com/api?v=3.0&ak=${baiduKey.value}&callback=${callbackName}`

    script.onerror = (error) => {
      delete window[callbackName]
      reject(new Error('百度地图脚本加载失败: ' + error))
    }

    document.head.appendChild(script)

    // 设置超时时间
    setTimeout(() => {
      if (window[callbackName]) {
        delete window[callbackName]
        reject(new Error('百度地图加载超时'))
      }
    }, 20000)
  })
}

// 初始化地图
const initMap = async () => {
  if (!isComponentMounted.value) return

  loading.value = true

  try {
    // 如果API已加载，直接初始化
    if (checkBMapLoaded()) {
      bMapScriptLoaded.value = true
      setupMap()
      return
    }

    // 否则加载API
    await loadBMapScript()

    if (!isComponentMounted.value) return

    bMapScriptLoaded.value = true
    setupMap()
  } catch (error) {
    console.error('地图初始化失败:', error)
    if (isComponentMounted.value) {
      ElMessage.error('地图加载失败，请刷新页面重试')
    }
  } finally {
    if (isComponentMounted.value) {
      loading.value = false
    }
  }
}

// 设置地图
const setupMap = () => {
  if (!isComponentMounted.value || !mapContainer.value) return

  try {
    // 创建地图实例
    map.value = new BMap.Map(mapContainer.value)

    // 初始化地理编码器
    geocoder.value = new BMap.Geocoder()

    // 定位地址
    locateAddress()
  } catch (error) {
    console.error('地图设置失败:', error)
    if (isComponentMounted.value) {
      ElMessage.error('地图初始化失败')
    }
  }
}

// 地址定位
const locateAddress = () => {
  if (!geocoder.value || !isComponentMounted.value) return

  geocoder.value.getPoint(props.address, (point) => {
    if (!isComponentMounted.value) return

    if (point) {
      currentLocation.value = {
        lng: point.lng,
        lat: point.lat
      }

      // 设置地图中心点
      map.value.centerAndZoom(point, props.zoom)
      map.value.enableScrollWheelZoom(true)

      // 添加控件
      map.value.addControl(new BMap.NavigationControl())
      map.value.addControl(new BMap.ScaleControl())

      // 添加标记
      if (props.showMarker) {
        addMarker(point)
      }

      // 获取地址详情
      getLocationDetails(point)
    } else {
      if (isComponentMounted.value) {
        ElMessage.error('地址解析失败，请检查地址是否正确')
      }
    }
  }, '全国')
}

// 添加标记点
const addMarker = (point) => {
  if (!map.value || !isComponentMounted.value) return

  // 清除现有标记
  if (marker.value) {
    map.value.removeOverlay(marker.value)
  }

  // 创建新标记
  marker.value = new BMap.Marker(point)
  map.value.addOverlay(marker.value)

  // 添加点击事件
  if (props.showInfoWindow) {
    marker.value.addEventListener('click', () => {
      if (isComponentMounted.value) {
        showInfoWindow()
      }
    })
  }
}

// 显示信息窗口
const showInfoWindow = () => {
  if (!map.value || !currentLocation.value.lng) return

  const content = `
    <div style="margin:10px;">
      <h4>位置信息</h4>
      <p>${props.address}</p>
      <p>坐标: ${currentLocation.value.lng.toFixed(6)}, ${currentLocation.value.lat.toFixed(6)}</p>
    </div>
  `

  const infoWindow = new BMap.InfoWindow(content)
  map.value.openInfoWindow(infoWindow, currentLocation.value)
}

// 获取位置详情
const getLocationDetails = (point) => {
  geocoder.value.getLocation(point, (result) => {
    if (!isComponentMounted.value) return

    if (result && result.addressComponents) {
      const comp = result.addressComponents
      locationDetails.value = `${comp.province}${comp.city}${comp.district}${comp.street}${comp.streetNumber}`
    }
  })
}

// 重新定位到中心
const centerMap = () => {
  if (map.value && currentLocation.value.lng && isComponentMounted.value) {
    map.value.panTo(new BMap.Point(currentLocation.value.lng, currentLocation.value.lat))
  }
}

// 监听地址变化
watch(() => props.address, (newAddress) => {
  if (newAddress && map.value && isComponentMounted.value && bMapScriptLoaded.value) {
    locateAddress()
  }
})

// 组件挂载时初始化地图
onMounted(async () => {
  isComponentMounted.value = true

  // 先获取密钥
  const keyLoaded = await getAndDecryptAPIKey()
  if (!keyLoaded) {
    loading.value = false
    return
  }

  nextTick(() => {
    // 延迟初始化，确保DOM完全渲染
    setTimeout(() => {
      if (isComponentMounted.value && mapContainer.value) {
        initMap()
      }
    }, 100)
  })
})
</script>

<style scoped>
.map-container {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 400px;
}

#map {
  width: 100%;
  height: 100%;
}

.map-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.9);
  padding: 12px 20px;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.control-bar {
  position: absolute;
  top: 16px;
  right: 16px;
  z-index: 100;
}

.info-drawer {
  z-index: 1000;
}

.drawer-content {
  padding: 20px;
}

.drawer-content h3 {
  margin-bottom: 16px;
  color: #1f2f3d;
}

.location-info {
  margin: 20px 0;
}

.info-item {
  display: flex;
  margin-bottom: 12px;
  line-height: 1.6;
}

.info-item .label {
  font-weight: 500;
  min-width: 70px;
  color: #606266;
}

.info-item .value {
  flex: 1;
  color: #303133;
  word-break: break-all;
}

.actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>