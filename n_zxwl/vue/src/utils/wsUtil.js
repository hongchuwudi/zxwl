// 全局实例
let globalSocket = null;

class NotificationWebSocket {
    constructor(userId, userEmail) {
        this.ws = null                         // WebSocket实例
        this.userId = userId                   // 用户ID
        this.userEmail = userEmail             // 用户邮箱
        this.reconnectAttempts = 0             // 重连尝试次数
        this.maxReconnectAttempts = 5         // 最大重连次数
        this.heartbeatInterval = null          // 心跳间隔定时器
        this.init()
    }
    // 初始连接
    init = () => {
        try {
            const wsUrl = `/gapi/ws/user?userID=${this.userId}` // 直接使用相对路径，Vite 会在开发时代理，生产环境使用绝对路径
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
            const fullWsUrl = wsUrl.startsWith('ws') ? wsUrl : `${protocol}${window.location.host}${wsUrl}`
            this.ws = new WebSocket(fullWsUrl)
            this.setupEventListeners();
        } catch (error) {
            console.error('WebSocket初始化失败:', error)
        }
    }

    // 事件监听
    setupEventListeners = () => {
        this.ws.onopen = () => {
            console.log('✅ WebSocket连接已建立')
            this.reconnectAttempts = 0
            this.startHeartbeat()
        }

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                this.handleMessage(message);
            } catch (error) {
                console.error('消息解析失败:', error);
            }
        };

        this.ws.onclose = (event) => {
            console.log('❌ WebSocket连接已关闭', event.code, event.reason);
            this.stopHeartbeat();

            // 自动重连（如果不是正常关闭）
            if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
                this.reconnect();
            }
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket错误:', error);
        };
    }

    // 处理消息
    handleMessage = (message) => {
        console.log('📨 收到消息:', message);

        const handlers = {
            friend_status: this.handleFriendStatus,
            friend_online: this.handleFriendOnline,
            friend_offline: this.handleFriendOffline,
            new_message: this.handleNewMessage,
            message_read: this.handleMessageRead,
            typing_status: this.handleTypingStatus,
            new_request: this.handleNewRequest,
            pong: () => console.log('🏓 收到心跳响应')
        };
        const handler = handlers[message.type] || (() => console.log('未知消息类型:', message.type));
        handler(message.data);
    }
    // 添加新的处理方法
    handleFriendStatus = (data) => this.dispatchEvent('friendStatus', data);
    handleFriendOnline = (data) => this.dispatchEvent('friendOnline', data);
    handleFriendOffline = (data) => this.dispatchEvent('friendOffline', data);
    handleNewMessage = (data) => this.dispatchEvent('newMessage', data);
    handleMessageRead = (data) => this.dispatchEvent('messageRead', data);
    handleTypingStatus = (data) => this.dispatchEvent('typingStatus', data);
    handleNewRequest = (data) => this.dispatchEvent('newRequest', data);
    // 分发自定义事件的方法
    dispatchEvent = (eventName, data) => {
        const event = new CustomEvent(eventName, { detail: data });
        window.dispatchEvent(event);
    }

    send = (message) => {
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        }
    }

    startHeartbeat = () => {
        this.heartbeatInterval = setInterval(() => {
            if (this.ws.readyState === WebSocket.OPEN) {
                this.send({ type: 'ping' });
            }
        }, 25000); // 25秒一次心跳
    }

    stopHeartbeat = () => {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
    }

    reconnect = () => {
        this.reconnectAttempts++;
        const delay = Math.min(3000 * this.reconnectAttempts, 15000); // 最大15秒

        console.log(`♻️ ${delay/1000}秒后尝试第${this.reconnectAttempts}次重连...`);

        setTimeout(() => {
            this.init();
        }, delay);
    }

    close = () => {
        this.stopHeartbeat();
        this.ws?.close(1000, '正常关闭');
    }

    showNotification = (title, body) => {
        if (!('Notification' in window)) return;

        if (Notification.permission === 'granted') {
            new Notification(title, { body, icon: '/favicon.ico' });
        } else if (Notification.permission !== 'denied') {
            Notification.requestPermission().then(permission => {
                if (permission === 'granted') {
                    new Notification(title, { body, icon: '/favicon.ico' });
                }
            });
        }
    }
}

// 创建 WebSocket 连接
export const initWebSocket = (userId, userEmail) => {
    globalSocket?.close();
    globalSocket = new NotificationWebSocket(userId, userEmail);
    return globalSocket;
}

// 获取当前 WebSocket 实例
export const getWebSocket = () => globalSocket;

// 关闭 WebSocket 连接
export const closeWebSocket = () => {
    globalSocket?.close();
    globalSocket = null;
}

export default NotificationWebSocket;