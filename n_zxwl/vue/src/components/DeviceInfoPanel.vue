<template>
  <div class="device-info-panel">
    <div class="info-section">
      <h3><laptop-outlined /> 设备信息</h3>
      <div class="info-item">
        <span class="label">设备类型:</span>
        <span class="value">{{ deviceInfo.device }}</span>
      </div>
      <div class="info-item">
        <span class="label">操作系统:</span>
        <span class="value">{{ deviceInfo.os.name }} {{ deviceInfo.os.version }}</span>
      </div>
      <div class="info-item">
        <span class="label">浏览器:</span>
        <span class="value">{{ deviceInfo.browser.name }} {{ deviceInfo.browser.version }}</span>
      </div>
    </div>

    <div class="info-section">
      <h3><desktop-outlined /> 屏幕信息</h3>
      <div class="info-item">
        <span class="label">分辨率:</span>
        <span class="value">{{ deviceInfo.screen.width }} × {{ deviceInfo.screen.height }}</span>
      </div>
      <div class="info-item">
        <span class="label">色彩深度:</span>
        <span class="value">{{ deviceInfo.screen.colorDepth }} 位</span>
      </div>
      <div class="info-item">
        <span class="label">方向:</span>
        <span class="value">{{ deviceInfo.screen.orientation }}</span>
      </div>
    </div>

    <div class="info-section">
      <h3><wifi-outlined /> 网络信息</h3>
      <div class="info-item">
        <span class="label">网络类型:</span>
        <span class="value">{{ deviceInfo.network.effectiveType }}</span>
      </div>
      <div class="info-item">
        <span class="label">下行速度:</span>
        <span class="value">{{ deviceInfo.network.downlink }} Mbps</span>
      </div>
      <div class="info-item">
        <span class="label">往返延迟:</span>
        <span class="value">{{ deviceInfo.network.rtt }} ms</span>
      </div>
    </div>

    <div class="info-section">
      <h3><global-outlined /> 其他信息</h3>
      <div class="info-item">
        <span class="label">语言:</span>
        <span class="value">{{ deviceInfo.language }}</span>
      </div>
      <div class="info-item">
        <span class="label">时区:</span>
        <span class="value">{{ deviceInfo.timezone }}</span>
      </div>
      <div class="info-item">
        <span class="label">检测时间:</span>
        <span class="value">{{ formatDate(deviceInfo.timestamp) }}</span>
      </div>
    </div>

    <div class="visualization">
      <div class="radar">
        <div class="radar-sweep"></div>
      </div>
    </div>
  </div>
</template>

<script>
import {
  LaptopOutlined,
  DesktopOutlined,
  WifiOutlined,
  GlobalOutlined
} from '@ant-design/icons-vue'

export default {
  name: 'DeviceInfoPanel',
  components: {
    LaptopOutlined,
    DesktopOutlined,
    WifiOutlined,
    GlobalOutlined
  },
  props: {
    deviceInfo: {
      type: Object,
      required: true
    }
  },
  methods: {
    formatDate(dateString) {
      return new Date(dateString).toLocaleString('zh-CN')
    }
  }
}
</script>

<style scoped>
.device-info-panel {
  position: relative;
  padding: 15px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border-radius: 8px;
  color: #e6e6e6;
  overflow: hidden;
}

.device-info-panel::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg,
  transparent,
  rgba(0, 183, 255, 0.5),
  transparent);
  animation: scanline 3s linear infinite;
}

@keyframes scanline {
  0% { transform: translateY(0); }
  100% { transform: translateY(100%); }
}

.info-section {
  margin-bottom: 20px;
  position: relative;
  z-index: 1;
}

.info-section h3 {
  color: #00b7ff;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  font-size: 16px;
}

.info-section h3 i {
  margin-right: 8px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.info-item:last-child {
  border-bottom: none;
}

.label {
  font-weight: 500;
  color: #a0a0a0;
}

.value {
  color: #ffffff;
  font-weight: 500;
}

.visualization {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}

.radar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  border: 2px solid #00b7ff;
  position: relative;
  overflow: hidden;
  background: rgba(0, 183, 255, 0.05);
}

.radar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 50%;
  border: 2px solid rgba(0, 183, 255, 0.2);
  transform: scale(0.8);
}

.radar::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 50%;
  border: 2px solid rgba(0, 183, 255, 0.1);
  transform: scale(0.6);
}

.radar-sweep {
  position: absolute;
  top: 0;
  left: 50%;
  width: 50%;
  height: 100%;
  transform-origin: left;
  background: linear-gradient(90deg,
  transparent,
  rgba(0, 183, 255, 0.3),
  transparent);
  animation: sweep 4s infinite linear;
}

@keyframes sweep {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 响应式设计 */
@media (max-width: 768px) {
  .info-item {
    flex-direction: column;
  }

  .value {
    margin-top: 4px;
  }
}
</style>