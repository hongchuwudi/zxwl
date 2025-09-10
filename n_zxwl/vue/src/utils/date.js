// 时间格式化工具
export const formatTime = (timestamp) => {
    if (!timestamp) return ''

    let date
    // 判断传入的是 Unix 时间戳数字还是 ISO 字符串
    if (typeof timestamp === 'number') {
        date = new Date(timestamp * 1000)  // Unix 时间戳（秒转毫秒）
    } else if (typeof timestamp === 'string') {
        date = new Date(timestamp)         // ISO 字符串
    } else {
        return ''
    }

    const now = new Date()
    const diff = now - date

    // 1分钟内
    if (diff < 60000) return '刚刚'

    // 1小时内
    if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`

    // 1天内
    if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`

    // 1周内
    if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`

    // 超过1周显示具体日期
    return date.toLocaleDateString()
}