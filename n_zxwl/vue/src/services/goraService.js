import AgoraRTC from 'agora-rtc-sdk-ng'

class AgoraService {
    constructor() {
        this.client = null
        this.localAudioTrack = null
        this.localVideoTrack = null
        this.isJoined = false
        this.remoteUsers = new Map()
    }

    // 初始化Agora
    async initAgora(appId) {
        this.client = AgoraRTC.createClient({ mode: "rtc", codec: "vp8" })

        // 监听远程用户事件
        this.client.on("user-published", async (user, mediaType) => {
            await this.client.subscribe(user, mediaType)

            if (mediaType === "video") {
                const remoteVideoTrack = user.videoTrack
                this.remoteUsers.set(user.uid, remoteVideoTrack)
                this.onRemoteVideoPublished(user.uid, remoteVideoTrack)
            }

            if (mediaType === "audio") {
                const remoteAudioTrack = user.audioTrack
                remoteAudioTrack.play()
            }
        })

        this.client.on("user-unpublished", (user, mediaType) => {
            if (mediaType === "video") {
                this.remoteUsers.delete(user.uid)
                this.onRemoteVideoUnpublished(user.uid)
            }
        })
    }

    // 加入频道
    async joinChannel(appId, channel, token, uid) {
        try {
            await this.client.join(appId, channel, token || null, uid || null)
            this.isJoined = true

            // 创建本地音视频轨道
            this.localAudioTrack = await AgoraRTC.createMicrophoneAudioTrack()
            this.localVideoTrack = await AgoraRTC.createCameraVideoTrack()

            // 发布本地轨道
            await this.client.publish([this.localAudioTrack, this.localVideoTrack])

            return {
                localAudioTrack: this.localAudioTrack,
                localVideoTrack: this.localVideoTrack
            }
        } catch (error) {
            console.error("加入频道失败:", error)
            throw error
        }
    }

    // 离开频道
    async leaveChannel() {
        try {
            // 关闭本地轨道
            if (this.localAudioTrack) {
                this.localAudioTrack.close()
                this.localAudioTrack = null
            }
            if (this.localVideoTrack) {
                this.localVideoTrack.close()
                this.localVideoTrack = null
            }

            // 离开频道
            await this.client.leave()
            this.isJoined = false
            this.remoteUsers.clear()

        } catch (error) {
            console.error("离开频道失败:", error)
        }
    }

    // 切换麦克风
    async toggleMicrophone(enabled) {
        if (this.localAudioTrack) {
            await this.localAudioTrack.setEnabled(enabled)
        }
        return enabled
    }

    // 切换摄像头
    async toggleCamera(enabled) {
        if (this.localVideoTrack) {
            await this.localVideoTrack.setEnabled(enabled)
        }
        return enabled
    }

    // 切换摄像头设备
    async switchCamera(deviceId) {
        if (this.localVideoTrack) {
            await this.localVideoTrack.setDevice(deviceId)
        }
    }

    // 获取设备列表
    async getDevices() {
        try {
            const devices = await AgoraRTC.getDevices()
            return {
                cameras: devices.filter(device => device.kind === 'videoinput'),
                microphones: devices.filter(device => device.kind === 'audioinput'),
                speakers: devices.filter(device => device.kind === 'audiooutput')
            }
        } catch (error) {
            console.error("获取设备列表失败:", error)
            return { cameras: [], microphones: [], speakers: [] }
        }
    }

    // 事件回调（需要在组件中重写）
    onRemoteVideoPublished(uid, videoTrack) {
        console.log("远程视频发布:", uid)
    }

    onRemoteVideoUnpublished(uid) {
        console.log("远程视频取消发布:", uid)
    }
}

export default new AgoraService()