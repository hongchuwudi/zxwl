<template>
  <div class="chat-container" :class="{ 'collapsed': isCollapsed }">
    <!-- 控制栏 -->
    <div class="control-bar">
      <!-- 左侧控制区 -->
      <div class="left-controls">
        <div class="back-button" @click="handleBack">
          <el-icon :size="24" class="back-icon">
            <Back />
          </el-icon>
        </div>

        <!-- 非收起状态下显示模型选择器 -->
        <div v-if="!isCollapsed" class="model-selector">
          <label>选择模型：</label>
          <select v-model="selectedModel" class="modern-select">
            <option v-for="model in availableModels" :key="model" :value="model">
              {{ model }}
            </option>
          </select>
          <span class="status-indicator ready"></span>
        </div>
      </div>

      <!-- 中央标题 - 非收起状态下显示 -->
      <div v-if="!isCollapsed" class="app-title">
        <h1 class="title-text">智选AI</h1>
        <div class="title-underline"></div>
      </div>

      <!-- 右侧控制区 -->
      <div class="right-controls">
        <!-- 清空按钮 - 添加悬浮提示 -->
        <el-tooltip content="清空对话" placement="bottom">
          <button class="clear-btn modern-btn" @click="clearHistory">
            <svg viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor" d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
            </svg>
          </button>
        </el-tooltip>

        <!-- 用户菜单 -->
        <div class="user-menu" v-click-outside="closeUserMenu">
          <button class="user-btn modern-btn" @click="toggleUserMenu">
            <svg viewBox="0 0 24 24" width="20" height="20">
              <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
            </svg>
          </button>
          <transition name="fade-slide">
            <div v-if="showUserMenu" class="user-dropdown modern-dropdown">
              <div class="dropdown-item" @click="viewProfile">
                <svg viewBox="0 0 24 24" width="16" height="16">
                  <path fill="currentColor" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
                </svg>
                <span>个人资料</span>
              </div>
              <div class="dropdown-item" @click="handleLogout">
                <svg viewBox="0 0 24 24" width="16" height="16">
                  <path fill="currentColor" d="M17 7l-1.41 1.41L18.17 11H8v2h10.17l-2.58 2.59L17 17l5-5zM4 5h8V3H4c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h8v-2H4V5z"/>
                </svg>
                <span>退出登录</span>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </div>

    <!-- 文件查看抽屉 -->
    <el-drawer v-model="showFileDrawer" title="文件内容" direction="rtl" size="50%">
      <div class="file-content">
        <pre>{{ selectedFileContent }}</pre>
      </div>
    </el-drawer>

    <!-- 聊天消息区域 -->
    <div ref="messagesContainer" class="chat-messages">
      <div
          v-for="(message, index) in messages"
          :key="index"
          :class="['message', message.role]"
      >
        <div class="message-header">
          <div class="role-avatar">
            <svg v-if="message.role === 'user'" viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" width="16" height="16">
              <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 18c-4.41 0-8-3.59-8-8s3.59-8 8-8 8 3.59 8 8-3.59 8-8 8zm-1-13h2v6h-2V7zm0 8h2v2h-2v-2z"/>
            </svg>
          </div>
          <span class="role-badge">{{ roleNames[message.role] }}</span>
          <span class="time">{{ message.timestamp }}</span>
        </div>
        <div
            class="message-content"
            v-html="formatContent(message.content)"
        ></div>
      </div>
      <div v-if="isLoading" class="thinking-indicator">
        <div class="typing-animation">
          <div class="typing-dot"></div>
          <div class="typing-dot"></div>
          <div class="typing-dot"></div>
        </div>
        <span class="thinking-text">AI回答中...</span>
      </div>
    </div>

    <!-- 输入区域 -->
    <div class="input-area">
      <!-- 文件列表区域 - 非收起状态下显示 -->
      <div v-if="!isCollapsed && uploadedFiles.length > 0" class="file-list">
        <div class="file-item" v-for="(file, index) in uploadedFiles" :key="index" @click="openFileDrawer(file)">
          <span class="file-name">{{ file.name }}</span>
          <el-icon @click.stop="removeFile(index)"><Close /></el-icon>
        </div>
      </div>

      <div class="input-wrapper">
        <!-- 文本输入框 - 非收起状态下显示 -->
        <textarea
            v-if="!isCollapsed"
            v-model="inputText"
            @keydown.enter.exact.prevent="sendMessage"
            @keydown.shift.enter.exact.prevent="inputText += '\n'"
            placeholder="输入消息 (Shift + Enter 换行)..."
            ref="textArea"
            class="modern-textarea"
            :disabled="isLoading"
            @input="adjustTextareaHeight"
            rows="1"
        ></textarea>

        <div class="input-actions">
          <!-- 上传按钮 - 非收起状态下显示 -->
          <el-upload
              v-if="!isCollapsed"
              class="upload-btn"
              action="#"
              :show-file-list="false"
              :before-upload="beforeUpload"
          >
            <el-button circle class="action-btn">
              <el-icon><Upload /></el-icon>
            </el-button>
          </el-upload>
          <!-- 发送按钮 - 非收起状态下显示 -->
          <el-button
              v-if="!isCollapsed"
              class="send-btn action-btn"
              type="warning"
              circle
              :disabled="isLoading || !inputText.trim()"
              @click="sendMessage"
          >
            <el-icon><Top /></el-icon>
          </el-button>
          <!-- 展开/收起按钮 -->
          <el-button
              class="toggle-btn action-btn"
              :class="{ 'collapsed': isCollapsed }"
              circle
              @click="isCollapsed = !isCollapsed"
          >
            <el-icon v-if="isCollapsed"><Expand /></el-icon>
            <el-icon v-else><Fold /></el-icon>
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick, computed, watch } from 'vue'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'
import { Back, Top, Upload,Close } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import axios from "axios";
import { ElMessage } from 'element-plus'
import { marked } from 'marked';
import {deleteUser} from "@/vuex/userStorage.js";
import {useUserStore} from "@/utils/auth.js";
import { message } from 'ant-design-vue';


// 用户信息
const { userName, userEmail, getUser, checkLoginStatus } = useUserStore();
const router = useRouter()

// 前文大模型
const QWEN_API = 'xapi/api/v1/services/aigc/text-generation/generation'
let API_KEY = ''  // 生产环境建议通过后端代理

// 响应式状态
const selectedModel = ref('qwen-plus')
const availableModels = ref(['qwen-max', 'qwen-plus', 'qwen-turbo'])
const messages = ref([])
const inputText = ref('')
const isLoading = ref(false)
const messagesContainer = ref(null)
const showUserMenu = ref(false)
const textArea = ref(null)

// 控制是否收起
const isCollapsed = ref(false)

// 文件状态
const uploadedFiles = ref([])
const showFileDrawer = ref(false)
const selectedFileContent = ref('')

// 打开文件抽屉
const openFileDrawer = (file) => {
  selectedFileContent.value = file.content
  showFileDrawer.value = true
}

// 添加 click-outside 指令
const vClickOutside = {
  beforeMount(el, binding) {
    el.clickOutsideEvent = function(event) {
      if (!(el === event.target || el.contains(event.target))) {
        binding.value(event);
      }
    };
    document.addEventListener('click', el.clickOutsideEvent);
  },
  unmounted(el) {
    document.removeEventListener('click', el.clickOutsideEvent);
  }
}

// 移除文件
const removeFile = (index) => {
  uploadedFiles.value.splice(index, 1)
  message.info('文件已移除')
}

// 配置marked
marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value;
    }
    return hljs.highlightAuto(code).value;
  },
  breaks: true,
  gfm: true
});

// 用户菜单
const toggleUserMenu = () => showUserMenu.value = !showUserMenu.value
const closeUserMenu = () => showUserMenu.value = false
const viewProfile = () => {
  console.log('查看资料');
  router.push('/profile')
  closeUserMenu()
}
const handleLogout = () => {
  confirm('确定退出？')
  closeUserMenu()
  deleteUser()
  router.push('/');
}

// 计算属性
// const modelInfo = computed(() => `当前模型: ${selectedModel.value}`)
const roleNames = { user: '您', assistant: '为广大高考生出言献计的伟大AI小助手', system: '系统' }

// 内容格式化
// 修改formatContent方法
const formatContent = (content) => {
  // 如果是空内容，返回空字符串
  if (!content) return '';

  try {
    // 使用marked解析Markdown
    return marked.parse(content);
  } catch (e) {
    console.error('Markdown解析错误:', e);
    // 如果解析失败，返回原始内容
    return content.replace(/\n/g, '<br>');
  }
};

// 调整文本区域高度
const adjustTextareaHeight = () => {
  nextTick(() => {
    if (textArea.value) {
      textArea.value.style.height = 'auto'
      textArea.value.style.height = Math.min(textArea.value.scrollHeight, 150) + 'px'
    }
  })
}

// 滚动处理
const scrollToBottom = () => {
  nextTick(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTo({
        top: messagesContainer.value.scrollHeight,
        behavior: 'smooth'
      })
    }
  })
}

// sendMessage SSE解析流 (调用AI的核心代码 文件调用+流式输出+非记忆模式混合输入)
const sendMessage = async () => {
  if (!inputText.value.trim() || isLoading.value) return;

  // 构建发送的消息数组（保持原有逻辑）
  const messagesToSend = [];
  const recentMessages = messages.value.slice(-20);
  recentMessages.forEach(msg => {
    if (msg.role !== 'system' && !msg.content.includes('📎 附')) {
      messagesToSend.push({
        role: msg.role,
        content: msg.content
      });
    }
  });

  messagesToSend.push({
    role: "user",
    content: inputText.value
  });

  if (uploadedFiles.value.length > 0) {
    uploadedFiles.value.forEach(file => {
      const maxLength = 1500;
      const truncatedContent = file.content.length > maxLength
          ? file.content.substring(0, maxLength) + '...（内容过长已截断）'
          : file.content;

      messagesToSend.push({
        role: "user",
        content: `【文件: ${file.name}】\n${truncatedContent}`
      });
    });
  }

  // 创建界面显示的用户消息
  const fileHint = uploadedFiles.value.length > 0 ? `\n\n📎 附 ${uploadedFiles.value.length} 个文件` : '';
  const userMessage = {
    role: 'user',
    content: inputText.value + fileHint,
    timestamp: new Date().toLocaleTimeString()
  };

  messages.value.push(userMessage);
  inputText.value = '';

  if (textArea.value) textArea.value.style.height = 'auto';

  isLoading.value = true;
  scrollToBottom();

  const aiMessageIndex = messages.value.length;
  messages.value.push({
    role: 'assistant',
    content: '',
    timestamp: new Date().toLocaleTimeString()
  });

  try {
    // 改为调用后端API（保持SSE流式传输）
    const response = await fetch("/gapi/aiChat", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        model: selectedModel.value,
        messages: messagesToSend,
        stream: true,
        temperature: 0.7,
        max_tokens: 5000
      })
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    let buffer = '';

    // 清空初始内容
    messages.value[aiMessageIndex].content = '';

    while (true) {
      const { done, value } = await reader.read();
      if (done) {
        isLoading.value = false;
        break;
      }

      buffer += decoder.decode(value, { stream: true });

      // 保持原有的SSE解析逻辑
      let lineIndex;
      while ((lineIndex = buffer.indexOf('\n')) !== -1) {
        const line = buffer.slice(0, lineIndex).trim();
        buffer = buffer.slice(lineIndex + 1);

        if (!line || line.startsWith(':')) continue;

        if (line.startsWith('data:')) {
          const data = line.slice(5).trim();

          if (data === '[DONE]') {
            isLoading.value = false;
            continue;
          }

          try {
            const parsed = JSON.parse(data);
            let content = '';

            // 提取内容（保持原有逻辑）
            if (parsed.output?.text) {
              content = parsed.output.text;
            } else if (parsed.output?.choices) {
              content = parsed.output.choices[0]?.message?.content || '';
            } else if (parsed.choices) {
              content = parsed.choices[0]?.delta?.content || '';
            } else if (parsed.text) {
              content = parsed.text;
            }

            if (content) {
              messages.value[aiMessageIndex].content = content;
              scrollToBottom();
            }
          } catch (e) {
            // 不是JSON，直接作为文本处理
            if (data && data !== '[DONE]') {
              messages.value[aiMessageIndex].content = data;
              scrollToBottom();
            }
          }
        }
      }
    }
  } catch (error) {
    console.error('请求错误:', error);
    isLoading.value = false;
    messages.value[aiMessageIndex].content = `错误: ${error.message}`;
    scrollToBottom();
  } finally {
    // 发送完成后清空文件列表
    uploadedFiles.value = [];
  }
}

// 文件上传
const beforeUpload = async (file) => {
  const isText = file.type.includes('text/') ||
      file.type === 'application/pdf' ||
      file.name.endsWith('.txt') ||
      file.name.endsWith('.pdf')  ||
      file.name.endsWith('.json') ||
      file.name.endsWith('.md')

  if (!isText) {
    ElMessage.error('只能上传文本文件,PDF,JSON和markdown文件')
    return false
  }

  try {
    // 读取文件内容
    const content = await readFileContent(file)
    uploadedFiles.value.push({
      name: file.name,
      type: file.type,
      content: content,
      size: file.size,
      uploadTime: new Date().toLocaleTimeString()
    })

    ElMessage.success(`文件 ${file.name} 上传成功`)
  } catch (error) {
    ElMessage.error('文件读取失败')
    console.error('File reading error:', error)
  }

  return false // 阻止自动上传
}

// 新增文件读取方法
const readFileContent = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => resolve(e.target.result)
    reader.onerror = reject

    if (file.type === 'application/pdf') {
      // PDF文件特殊处理（这里需要PDF.js等库，简化处理）
      resolve(`[PDF文件: ${file.name}]`)
    } else {
      reader.readAsText(file)
    }
  })
}

// 其他功能
const clearHistory = () => {
  if (messages.value.length > 0) {
    messages.value = []
  }
}
const handleBack = () => router.back()


onMounted(async () => {
  // 1.检查用户登录状态
  getUser()
  if (!checkLoginStatus()) {
    ElMessage.error('请先登录！')
    router.push('/login')
  }

  // 提交用户操作日志
  const logData = {
    email: localStorage.getItem('userEmail'),
    date: new Date().toISOString().slice(0, 19).replace('T', ' '),
    operation: "进行AI对话"
  }
  await axios.post("gapi/log", logData)

  // 监听输入框变化以调整高度
  if (textArea.value) {
    textArea.value.addEventListener('input', adjustTextareaHeight)
  }
})
</script>

<style scoped>
/* 基础样式重置与优化 */
.chat-container {
  width: 100vw;
  height: 100vh;
  margin: 0;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4e8f0 100%);
  color: #2d3748;
  font-family: 'Inter', 'SF Pro Display', -apple-system, BlinkMacSystemFont, sans-serif;
  overflow: hidden;
}

/* 控制栏样式 */
.control-bar {
  position: relative;
  padding: 1rem 2rem;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(226, 232, 240, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.app-title {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
}

.title-text {
  font-size: 1.8rem;
  font-weight: 700;
  background: linear-gradient(90deg, #6a11cb 0%, #2575fc 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 1px;
  margin: 0;
  padding-bottom: 8px;
}

.title-underline {
  position: absolute;
  bottom: -2px;
  left: 50%;
  transform: translateX(-50%);
  width: 80%;
  height: 3px;
  background: linear-gradient(
      90deg,
      rgba(106, 17, 203, 0) 0%,
      rgba(37, 117, 252, 0.6) 50%,
      rgba(106, 17, 203, 0) 100%
  );
  border-radius: 2px;
  animation: underlineFlow 3s infinite linear;
}

@keyframes underlineFlow {
  0% { background-position-x: -100%; }
  100% { background-position-x: 200%; }
}

.left-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
  z-index: 1;
}

.right-controls {
  display: flex;
  align-items: center;
  gap: 1rem;
  z-index: 1;
}

.model-selector {
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.modern-select {
  padding: 0.6rem 1rem;
  background: white;
  border: 1px solid rgba(203, 213, 225, 0.6);
  color: #475569;
  border-radius: 10px;
  font-size: 0.9rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 200px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.modern-select:hover {
  border-color: #94a3b8;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
}

.modern-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  position: relative;
}

.status-indicator.ready {
  background: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.3);
}

.status-indicator.ready::after {
  content: '';
  position: absolute;
  top: -3px;
  left: -3px;
  right: -3px;
  bottom: -3px;
  border-radius: 50%;
  background: rgba(16, 185, 129, 0.3);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { opacity: 0.8; transform: scale(1); }
  50% { opacity: 0; transform: scale(1.8); }
  100% { opacity: 0; transform: scale(1.8); }
}

/* 按钮样式 */
.modern-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.6rem 1rem;
  border: none;
  border-radius: 10px;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
  background: white;
  color: #475569;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.modern-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.modern-btn.primary {
  background: linear-gradient(90deg, #6a11cb 0%, #2575fc 100%);
  color: white;
}

.modern-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.clear-btn {
  padding: 0.6rem;
  background: white;
  color: #64748b;
  border-radius: 50% !important; /* 添加这个 */
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-menu {
  position: relative;
  margin-right: 0.5rem;
}

.user-btn {
  padding: 0.6rem;
  background: transparent;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.3s;
  color: #64748b;
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.user-btn:hover {
  background: #f8fafc;
  color: #3b82f6;
  transform: rotate(15deg);
}

/* 文件列表样式 */
.file-list {
  padding: 0.5rem 2rem;
  background: rgba(255, 255, 255, 0.9);
  border-bottom: 1px solid rgba(226, 232, 240, 0.6);
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.85rem;

}

.file-item:hover {
  border-color: #3b82f6;
  box-shadow: 0 2px 8px rgba(59, 130, 246, 0.2);
}

.file-name {
  color: #475569;
}

.file-content {
  padding: 1rem;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: 'Fira Code', monospace;
  font-size: 0.9rem;
  line-height: 1.5;
}

/* 下拉菜单 */
.modern-dropdown {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
  min-width: 140px;
  z-index: 100;
  overflow: hidden;
  border: 1px solid rgba(226, 232, 240, 0.8);
}

.dropdown-item {
  padding: 0.75rem 1rem;
  cursor: pointer;
  transition: all 0.2s;
  color: #475569;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.9rem;
}

.dropdown-item:hover {
  background: rgba(59, 130, 246, 0.08);
  color: #3b82f6;
}

/* 消息区域 */
.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 1.5rem 2rem;
  background: rgba(248, 250, 252, 0.6);
  scroll-behavior: smooth;
}

.message {
  margin: 1rem 0;
  padding: 1.25rem;
  border-radius: 18px;
  max-width: 80%;
  min-width: 120px;
  width: fit-content;
  animation: messageAppear 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: transform 0.3s, box-shadow 0.3s;
  line-height: 1.6;
  position: relative;
  overflow: hidden;
}

/* Markdown内容样式 */
.message-content ::v-deep h1,
.message-content ::v-deep h2,
.message-content ::v-deep h3,
.message-content ::v-deep h4 {
  margin-top: 1.2em;
  margin-bottom: 0.8em;
  font-weight: 600;
  line-height: 1.3;
}

.message-content ::v-deep h1 {
  font-size: 1.5em;
  border-bottom: 1px solid rgba(226, 232, 240, 0.8);
  padding-bottom: 0.3em;
}

.message-content ::v-deep h2 {
  font-size: 1.3em;
  border-bottom: 1px solid rgba(226, 232, 240, 0.6);
  padding-bottom: 0.3em;
}

.message-content ::v-deep p {
  margin-bottom: 1em;
  line-height: 1.6;
}

.message-content ::v-deep ul,
.message-content ::v-deep ol {
  margin-bottom: 1em;
  padding-left: 2em;
}

.message-content ::v-deep li {
  margin-bottom: 0.5em;
}

.message-content ::v-deep blockquote {
  border-left: 4px solid #e2e8f0;
  padding-left: 1em;
  margin: 1em 0;
  color: #64748b;
  font-style: italic;
}

.message-content ::v-deep table {
  border-collapse: collapse;
  width: 100%;
  margin: 1em 0;
}

.message-content ::v-deep th,
.message-content ::v-deep td {
  border: 1px solid #e2e8f0;
  padding: 0.5em;
}

.message-content ::v-deep th {
  background-color: #f8fafc;
  font-weight: 600;
}

.message-content ::v-deep tr:nth-child(even) {
  background-color: #f8fafc;
}

.message-content ::v-deep a {
  color: #3b82f6;
  text-decoration: none;
}

.message-content ::v-deep a:hover {
  text-decoration: underline;
}

.message-content ::v-deep code {
  font-family: 'Fira Code', monospace;
  background: rgba(59, 130, 246, 0.08);
  padding: 0.2em 0.4em;
  border-radius: 4px;
  color: #1e40af;
  border: 1px solid rgba(59, 130, 246, 0.1);
}

.message-content ::v-deep pre {
  background: rgba(59, 130, 246, 0.05);
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1rem 0;
  border: 1px solid rgba(59, 130, 246, 0.1);
}

.message-content ::v-deep pre code {
  background: none;
  padding: 0;
  border: none;
  color: inherit;
  border-radius: 0;
}

.message::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, rgba(255,255,255,0) 0%, rgba(255,255,255,0.4) 50%, rgba(255,255,255,0) 100%);
}

.message:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
}

/* 用户消息背景色调整 */
.message.user {
  background: #e6f7ff;
  color: #1f2937;
  margin-left: auto;
  border-bottom-right-radius: 6px;
  border: 1px solid #bae7ff;
}

.message.assistant {
  background: #fafafa;
  color: #1f2937;
  border: 1px solid #e8e8e8;
  margin-right: auto;
  border-bottom-left-radius: 6px;
}

.message-header {
  display: flex;
  align-items: center;
  margin-bottom: 0.75rem;
  font-size: 0.85rem;
  gap: 0.5rem;
}

.role-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
}

.message.assistant .role-avatar {
  background: rgba(59, 130, 246, 0.1);
  color: #1d4ed8;
}

.message.user .role-avatar {
  background: rgba(255, 255, 255, 0.3);
  color: white;
}

.role-badge {
  padding: 0.2rem 0.6rem;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.15);
  font-weight: 500;
  font-size: 0.75rem;
}

.message.assistant .role-badge {
  background: rgba(59, 130, 246, 0.1);
  color: #1d4ed8;
}

.time {
  margin-left: auto;
  opacity: 0.7;
  font-size: 0.75rem;
}

.message-content {
  font-size: 0.95rem;
  word-break: break-word;
}

.message-content ::v-deep code {
  font-family: 'Fira Code', monospace;
  background: rgba(59, 130, 246, 0.08);
  padding: 0.2em 0.4em;
  border-radius: 4px;
  color: #1e40af;
  border: 1px solid rgba(59, 130, 246, 0.1);
}

.message-content ::v-deep pre {
  background: rgba(59, 130, 246, 0.05);
  padding: 1rem;
  border-radius: 8px;
  overflow-x: auto;
  margin: 1rem 0;
  border: 1px solid rgba(59, 130, 246, 0.1);
}

/* 思考指示器 */
.thinking-indicator {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 18px;
  margin: 1rem 0;
  width: fit-content;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(226, 232, 240, 0.8);
}

.typing-animation {
  display: flex;
  align-items: center;
  gap: 4px;
}

.typing-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: linear-gradient(90deg, #6a11cb 0%, #2575fc 100%);
  animation: typingAnimation 1.4s infinite ease-in-out;
}

.typing-dot:nth-child(1) { animation-delay: 0s; }
.typing-dot:nth-child(2) { animation-delay: 0.2s; }
.typing-dot:nth-child(3) { animation-delay: 0.4s; }

@keyframes typingAnimation {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-5px); }
}

.thinking-text {
  font-size: 0.9rem;
  color: #64748b;
}

/* 输入区域 */
.input-area {
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(226, 232, 240, 0.8);
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.05);
}

.input-wrapper {
  display: flex;
  align-items: flex-end;
  gap: 0.75rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  position: relative;
}

.modern-textarea {
  flex: 1;
  padding: 1rem 1.25rem;
  padding-right: 100px; /* 为按钮留出空间 */
  background: white;
  border: 1px solid rgba(203, 213, 225, 0.6);
  color: #475569;
  border-radius: 16px;
  min-height: 125px;
  max-height: 150px;
  font-size: 0.95rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  line-height: 1.6;
  resize: none;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  overflow-y: auto;
}

.modern-textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15), 0 4px 10px rgba(0, 0, 0, 0.05);
}

.modern-textarea:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* 输入操作按钮 */
.input-actions {
  position: absolute;
  right: 12px;
  bottom: 12px;
  display: flex;
  gap: 8px;
  z-index: 5;
  outline: none;
}

.action-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.upload-btn {
  background: white;
  color: #64748b;
  border: 1px solid #ffffff;
}

.upload-btn:hover {
  background: #f8fafc;
  color: #3b82f6;
}

.send-btn {
  background: linear-gradient(135deg, #2571af 0%, #2575fc 100%);
  color: white;
}

.send-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #5a0db9 0%, #1c64d8 100%);
}

.send-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.hint-text {
  text-align: center;
  width: 100%;
  margin-top: 0.75rem;
  color: #64748b;
  font-size: 0.85rem;
  opacity: 0.8;
}

/* 返回按钮 */
.back-button {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 50%;
  padding: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(226, 232, 240, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-icon {
  color: #64748b;
  transition: color 0.3s ease;
}

.back-button:hover {
  transform: scale(1.05) translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.back-button:hover .back-icon {
  color: #6a11cb;
}

/* 动画效果 */
@keyframes messageAppear {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

/* 滚动条样式 */
.chat-messages::-webkit-scrollbar {
  width: 8px;
}

.chat-messages::-webkit-scrollbar-track {
  background: rgba(241, 245, 249, 0.5);
  border-radius: 4px;
}

.chat-messages::-webkit-scrollbar-thumb {
  background: rgba(203, 213, 225, 0.7);
  border-radius: 4px;
}

.chat-messages::-webkit-scrollbar-thumb:hover {
  background: rgba(148, 163, 184, 0.9);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .control-bar {
    padding: 0.75rem 1rem;
    flex-wrap: wrap;
  }

  .app-title {
    position: static;
    transform: none;
    order: -1;
    width: 100%;
    margin-bottom: 0.5rem;
  }

  .left-controls, .right-controls {
    width: 50%;
  }

  .right-controls {
    justify-content: flex-end;
  }

  .message {
    max-width: 90%;
  }

  .input-wrapper {
    flex-direction: column;
  }

  .back-button {
    margin-right: 0.5rem;
  }
}

/* 收起状态下的消息区域 */
.chat-container.collapsed .chat-messages {
  height: calc(100vh - 80px); /* 调整高度 */
}

/* 展开/收起按钮样式 */
.toggle-btn {
  background: white;
  color: #64748b;
  border: 1px solid #e2e8f0;
}

.toggle-btn:hover {
  background: #f8fafc;
  color: #3b82f6;
}

.toggle-btn.collapsed {
  background: linear-gradient(135deg, #2571af 0%, #2575fc 100%);
  color: white;
}

/* 移除折叠按钮的左边距 */
.input-actions ::v-deep .toggle-btn.el-button {
  margin-left: 0 !important;
}
.input-actions ::v-deep .el-button + .toggle-btn.el-button {
  margin-left: 0 !important;
}
</style>