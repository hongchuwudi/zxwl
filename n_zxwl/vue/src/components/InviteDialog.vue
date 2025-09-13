<template>
  <a-modal
      v-model:open="open"
      title="邀请他人"
      :footer="null"
      width="400px"
  >
    <div class="invite-dialog">
      <a-input
          v-model:value="inviteLink"
          placeholder="生成邀请链接..."
          readonly
          class="link-input"
      >
        <template #addonAfter>
          <a-button type="text" @click="copyLink">
            <CopyOutlined />
          </a-button>
        </template>
      </a-input>

      <div class="divider">或</div>

      <div class="invite-options">
        <h4>通过以下方式邀请</h4>
        <div class="option-buttons">
          <a-button>
            <template #icon>
              <WechatOutlined />
            </template>
            微信
          </a-button>
          <a-button>
            <template #icon>
              <QqOutlined />
            </template>
            QQ
          </a-button>
          <a-button>
            <template #icon>
              <MessageOutlined />
            </template>
            短信
          </a-button>
        </div>
      </div>

      <div class="action-buttons">
        <a-button type="primary" @click="handleOk">完成</a-button>
      </div>
    </div>
  </a-modal>
</template>

<script setup>
import { ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { CopyOutlined, WechatOutlined, QqOutlined, MessageOutlined } from '@ant-design/icons-vue'

const props = defineProps({
  open: {
    type: Boolean,
    default: false
  }
})
const open = ref(false)

const emit = defineEmits(['update:visible', 'invite'])

const visible = ref(false)
const inviteLink = ref('https://meet.example.com/room/123456')

watch(() => props.open, (val) => {
  open.value = val
})

watch(open, (val) => {
  emit('update:open', val)
})

const copyLink = () => {
  navigator.clipboard.writeText(inviteLink.value)
      .then(() => {
        message.success('链接已复制到剪贴板')
      })
      .catch(err => {
        console.error('复制失败:', err)
        message.error('复制失败')
      })
}

const handleOk = () => {
  emit('invite', { link: inviteLink.value })
  visible.value = false
}
</script>

<style scoped>
.invite-dialog {
  padding: 10px 0;
}

.link-input {
  margin-bottom: 16px;
}

.divider {
  text-align: center;
  margin: 16px 0;
  position: relative;
  color: #8c8c8c;
}

.divider:before,
.divider:after {
  content: '';
  position: absolute;
  top: 50%;
  width: 40%;
  height: 1px;
  background-color: #e8e8e8;
}

.divider:before {
  left: 0;
}

.divider:after {
  right: 0;
}

.invite-options h4 {
  margin-bottom: 12px;
  color: #262626;
}

.option-buttons {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.option-buttons button {
  flex: 1;
  margin: 0 4px;
}

.action-buttons {
  text-align: right;
}
</style>