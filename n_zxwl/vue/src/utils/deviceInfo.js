// utils/deviceInfo.js

/**
 * 获取设备信息工具类
 * 返回设备类型、操作系统、浏览器、屏幕信息等
 */
class DeviceUtils {
    /**
     * 获取完整的设备信息
     * @returns {Object} 设备信息对象
     */
    static getDeviceInfo() {
        return {
            device: this.getDeviceType(),
            os: this.getOSInfo(),
            browser: this.getBrowserInfo(),
            screen: this.getScreenInfo(),
            language: this.getLanguage(),
            timezone: this.getTimezone(),
            network: this.getNetworkInfo(),
            timestamp: new Date().toISOString(),
            userAgent: navigator.userAgent
        }
    }

    /**
     * 获取设备类型
     * @returns {string} mobile | tablet | desktop
     */
    static getDeviceType() {
        const ua = navigator.userAgent
        const isMobile = /iPhone|iPad|iPod|Android|webOS|BlackBerry|Windows Phone/i.test(ua)
        const isTablet = /iPad|Android|Tablet|PlayBook|Silk/i.test(ua) && !/Mobile/i.test(ua)

        if (isMobile) return 'mobile'
        if (isTablet) return 'tablet'
        return 'desktop'
    }

    /**
     * 获取操作系统信息
     * @returns {Object} 操作系统信息
     */
    static getOSInfo() {
        const ua = navigator.userAgent
        let os = 'unknown'
        let version = 'unknown'

        if (/Windows NT 10.0/i.test(ua)) {
            os = 'Windows'
            version = '10'
        } else if (/Windows NT 6.3/i.test(ua)) {
            os = 'Windows'
            version = '8.1'
        } else if (/Windows NT 6.2/i.test(ua)) {
            os = 'Windows'
            version = '8'
        } else if (/Windows NT 6.1/i.test(ua)) {
            os = 'Windows'
            version = '7'
        } else if (/Macintosh|Mac OS X/i.test(ua)) {
            os = 'macOS'
            const match = /Mac OS X (\d+[._]\d+)/i.exec(ua)
            if (match) version = match[1].replace('_', '.')
        } else if (/Android/i.test(ua)) {
            os = 'Android'
            const match = /Android (\d+[._]\d+)/i.exec(ua)
            if (match) version = match[1]
        } else if (/iOS|iPhone|iPad|iPod/i.test(ua)) {
            os = 'iOS'
            const match = /OS (\d+[._]\d+)/i.exec(ua)
            if (match) version = match[1].replace('_', '.')
        } else if (/Linux/i.test(ua)) {
            os = 'Linux'
        }

        return { name: os, version }
    }

    /**
     * 获取浏览器信息
     * @returns {Object} 浏览器信息
     */
    static getBrowserInfo() {
        const ua = navigator.userAgent
        let browser = 'unknown'
        let version = 'unknown'

        if (/Edg\/|Edge\/|EdgA\/|EdgiOS\//i.test(ua)) {
            browser = 'Edge'
            const match = /(Edg|Edge|EdgA|EdgiOS)\/(\d+)/i.exec(ua)
            if (match) version = match[2]
        } else if (/Chrome|CriOS/i.test(ua) && !/Edg/i.test(ua)) {
            browser = 'Chrome'
            const match = /(Chrome|CriOS)\/(\d+)/i.exec(ua)
            if (match) version = match[2]
        } else if (/Firefox|FxiOS/i.test(ua)) {
            browser = 'Firefox'
            const match = /(Firefox|FxiOS)\/(\d+)/i.exec(ua)
            if (match) version = match[2]
        } else if (/Safari/i.test(ua) && !/Chrome/i.test(ua)) {
            browser = 'Safari'
            const match = /Version\/(\d+)/i.exec(ua)
            if (match) version = match[1]
        } else if (/MSIE|Trident/i.test(ua)) {
            browser = 'IE'
            if (/MSIE (\d+)/i.test(ua)) {
                version = /MSIE (\d+)/i.exec(ua)[1]
            } else if (/Trident\/.*rv:(\d+)/i.test(ua)) {
                version = /Trident\/.*rv:(\d+)/i.exec(ua)[1]
            }
        }

        return { name: browser, version }
    }

    /**
     * 获取屏幕信息
     * @returns {Object} 屏幕信息
     */
    static getScreenInfo() {
        return {
            width: window.screen.width,
            height: window.screen.height,
            availWidth: window.screen.availWidth,
            availHeight: window.screen.availHeight,
            colorDepth: window.screen.colorDepth,
            pixelDepth: window.screen.pixelDepth,
            orientation: window.screen.orientation ? window.screen.orientation.type : 'unknown'
        }
    }

    /**
     * 获取语言信息
     * @returns {string} 浏览器语言
     */
    static getLanguage() {
        return navigator.language || navigator.userLanguage || navigator.browserLanguage || 'unknown'
    }

    /**
     * 获取时区信息
     * @returns {string} 时区
     */
    static getTimezone() {
        try {
            return Intl.DateTimeFormat().resolvedOptions().timeZone
        } catch (e) {
            return 'unknown'
        }
    }

    /**
     * 获取网络信息
     * @returns {Object} 网络信息
     */
    static getNetworkInfo() {
        const connection = navigator.connection || navigator.mozConnection || navigator.webkitConnection

        if (connection) {
            return {
                effectiveType: connection.effectiveType || 'unknown',
                downlink: connection.downlink || 'unknown',
                rtt: connection.rtt || 'unknown',
                saveData: connection.saveData || false,
                type: connection.type || 'unknown'
            }
        }

        return { effectiveType: 'unknown', downlink: 'unknown', rtt: 'unknown', saveData: false, type: 'unknown' }
    }

    /**
     * 获取设备信息并转换为JSON字符串
     * @returns {string} JSON字符串
     */
    static getDeviceInfoJSON() {
        return JSON.stringify(this.getDeviceInfo(), null, 2)
    }

    /**
     * 检查是否支持触摸
     * @returns {boolean} 是否支持触摸
     */
    static isTouchSupported() {
        return 'ontouchstart' in window || navigator.maxTouchPoints > 0
    }

    /**
     * 检查是否在线
     * @returns {boolean} 是否在线
     */
    static isOnline() {
        return navigator.onLine
    }

    /**
     * 获取设备内存信息（仅限支持的设备）
     * @returns {Object} 内存信息
     */
    static getMemoryInfo() {
        if (navigator.deviceMemory) {
            return {
                deviceMemory: navigator.deviceMemory + ' GB',
                totalJSHeapSize: performance.memory ? performance.memory.totalJSHeapSize : 'not supported',
                usedJSHeapSize: performance.memory ? performance.memory.usedJSHeapSize : 'not supported'
            }
        }
        return { deviceMemory: 'not supported' }
    }
}

export default DeviceUtils