<template>
  <div class="app-container">
    <div class="chat-wrapper">
      <!-- Sidebar -->
      <div class="sidebar">
        <div class="sidebar-header">
          <h1>💬 智能助手</h1>
          <button class="btn-new-chat" @click="createNewChat">+</button>
        </div>
        
        <div class="sessions-list">
          <div
            v-for="session in chatStore.sessions"
            :key="session.id"
            class="session-item"
            :class="{ active: chatStore.currentSessionId === session.id }"
            @click="selectSession(session.id)"
          >
            <div class="session-title">{{ session.title }}</div>
            <button class="btn-delete" @click.stop="deleteSession(session.id)">×</button>
          </div>
        </div>
      </div>

      <!-- Main Chat Area -->
      <div class="chat-area">
        <div class="chat-container">
          <!-- Empty State -->
          <div v-if="!chatStore.currentSession" class="empty-state">
            <div class="empty-icon">🤖</div>
            <p>选择一个会话或创建新的会话</p>
          </div>

          <!-- Messages -->
          <div v-else class="messages-container" ref="messagesContainer">
            <div
              v-for="message in chatStore.currentSession.messages"
              :key="message.id"
              class="message"
              :class="message.role"
            >
              <div class="message-avatar">{{ message.role === 'user' ? '👤' : '🤖' }}</div>
              <div class="message-content"><p>{{ message.content }}</p></div>
            </div>

            <!-- Loading indicator -->
            <div v-if="isLoading" class="message assistant loading">
              <div class="message-avatar">🤖</div>
              <div class="message-content">
                <div class="typing-indicator">
                  <span></span><span></span><span></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Input Area -->
        <div v-if="chatStore.currentSession" class="input-area">
          <div class="input-controls">
            <button class="btn-voice" :class="{ recording: isRecording }" @click="toggleRecording">
              {{ isRecording ? '🔴 停止' : '🎤 开始' }}
            </button>
            <input
              v-model="inputText"
              type="text"
              class="text-input"
              placeholder="输入您的问题..."
              @keyup.enter="sendMessage"
            />
            <button class="btn-send" :disabled="!inputText.trim()" @click="sendMessage">发送</button>
          </div>
          <div v-if="isRecording" class="recording-indicator">🎙️ 录音中... {{ recordingTime }}s</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue'
import { useChatStore } from '@/stores/chat'

const chatStore = useChatStore()
const inputText = ref('')
const isRecording = ref(false)
const isLoading = ref(false)
const recordingTime = ref(0)
const messagesContainer = ref<HTMLElement>()

onMounted(() => {
  chatStore.fetchSessions()
})

const selectSession = (sessionId: string) => {
  chatStore.setCurrentSession(sessionId)
}

const createNewChat = async () => {
  const title = `对话 ${new Date().toLocaleTimeString()}`
  await chatStore.createSession(title)
}

const deleteSession = async (sessionId: string) => {
  if (confirm('确定要删除这个会话吗？')) {
    await chatStore.deleteSession(sessionId)
  }
}

const toggleRecording = () => {
  isRecording.value = !isRecording.value
  if (isRecording.value) {
    recordingTime.value = 0
  }
}

const sendMessage = async () => {
  if (!inputText.value.trim()) return

  const message = inputText.value
  inputText.value = ''
  isLoading.value = true

  try {
    await chatStore.sendMessage(message)
  } catch (error) {
    console.error('Failed to send message:', error)
  } finally {
    isLoading.value = false
  }
}

watch(
  () => chatStore.currentSession?.messages,
  () => {
    nextTick(() => {
      if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
      }
    })
  },
  { deep: true }
)
</script>

<style scoped>
:root {
  --primary-color: #4a90e2;
  --primary-light: #f0f5ff;
  --text-primary: #1f2937;
  --text-secondary: #6b7280;
  --border-color: #e5e7eb;
  --bg-color: #ffffff;
  --message-user: #e7f3ff;
  --message-assistant: #f3f4f6;
}

.app-container {
  display: flex;
  height: 100vh;
  background-color: var(--bg-color);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

.chat-wrapper {
  display: flex;
  width: 100%;
  height: 100%;
}

.sidebar {
  width: 280px;
  background-color: #f8f9fa;
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header h1 {
  font-size: 18px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.btn-new-chat {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  background-color: var(--bg-color);
  cursor: pointer;
  font-size: 18px;
  transition: all 0.2s;
}

.btn-new-chat:hover {
  background-color: var(--primary-light);
  border-color: var(--primary-color);
}

.sessions-list {
  flex: 1;
  padding: 8px;
  overflow-y: auto;
}

.session-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background-color: transparent;
}

.session-item:hover {
  background-color: var(--primary-light);
}

.session-item.active {
  background-color: var(--primary-light);
  border: 1px solid var(--primary-color);
}

.session-title {
  flex: 1;
  font-size: 14px;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.btn-delete {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 18px;
  padding: 0;
  width: 24px;
  height: 24px;
  opacity: 0;
  transition: all 0.2s;
}

.session-item:hover .btn-delete {
  opacity: 1;
}

.btn-delete:hover {
  background-color: rgba(234, 88, 12, 0.1);
  color: #ea580c;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-color);
}

.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.empty-state {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.messages-container {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.message {
  display: flex;
  gap: 12px;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message.user {
  justify-content: flex-end;
}

.message.assistant {
  justify-content: flex-start;
}

.message-avatar {
  font-size: 28px;
  flex-shrink: 0;
}

.message-content {
  max-width: 60%;
  padding: 12px 16px;
  border-radius: 12px;
  line-height: 1.5;
  word-break: break-word;
}

.message.user .message-content {
  background-color: var(--primary-color);
  color: white;
}

.message.assistant .message-content {
  background-color: var(--message-assistant);
  color: var(--text-primary);
}

.message-content p {
  margin: 0;
  font-size: 14px;
}

.typing-indicator {
  display: flex;
  gap: 4px;
}

.typing-indicator span {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: var(--text-secondary);
  animation: typing 1.4s infinite;
}

.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing {
  0%, 60%, 100% {
    opacity: 0.5;
    transform: translateY(0);
  }
  30% {
    opacity: 1;
    transform: translateY(-10px);
  }
}

.input-area {
  border-top: 1px solid var(--border-color);
  padding: 16px;
  background-color: var(--bg-color);
}

.input-controls {
  display: flex;
  gap: 12px;
}

.btn-voice {
  padding: 10px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background-color: var(--bg-color);
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
  white-space: nowrap;
}

.btn-voice:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.btn-voice.recording {
  background-color: #fee2e2;
  border-color: #dc2626;
  color: #dc2626;
}

.text-input {
  flex: 1;
  padding: 10px 16px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  font-size: 14px;
  outline: none;
  transition: all 0.2s;
}

.text-input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--primary-light);
}

.btn-send {
  padding: 10px 24px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-send:hover:not(:disabled) {
  background-color: #357abd;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(74, 144, 226, 0.4);
}

.btn-send:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.recording-indicator {
  margin-top: 12px;
  padding: 8px 12px;
  background-color: #fef2f2;
  border-left: 3px solid #dc2626;
  border-radius: 4px;
  color: #991b1b;
  font-size: 12px;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

@media (max-width: 640px) {
  .sidebar {
    width: 100%;
    height: auto;
    max-height: 150px;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
    flex-direction: row;
  }

  .message-content {
    max-width: 90%;
  }
}
</style>
