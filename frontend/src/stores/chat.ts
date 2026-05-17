import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

interface Message {
  id: string
  session_id: string
  role: 'user' | 'assistant'
  content: string
  created_at: string
}

interface Session {
  id: string
  title: string
  provider: string
  model: string
  created_at: string
  updated_at: string
  messages: Message[]
}

export const useChatStore = defineStore('chat', () => {
  const sessions = ref<Session[]>([])
  const currentSessionId = ref<string | null>(null)

  const currentSession = computed(() => {
    return sessions.value.find(s => s.id === currentSessionId.value)
  })

  const api = axios.create({
    baseURL: '/api/v1',
  })

  async function fetchSessions() {
    try {
      const response = await api.get('/chat/sessions')
      sessions.value = response.data || []
      if (sessions.value.length > 0 && !currentSessionId.value) {
        currentSessionId.value = sessions.value[0].id
        await loadSessionMessages(sessions.value[0].id)
      }
    } catch (error) {
      console.error('Failed to fetch sessions:', error)
    }
  }

  async function loadSessionMessages(sessionId: string) {
    try {
      const response = await api.get(`/chat/sessions/${sessionId}`)
      const session = sessions.value.find(s => s.id === sessionId)
      if (session) {
        session.messages = response.data.messages || []
      }
    } catch (error) {
      console.error('Failed to load session messages:', error)
    }
  }

  async function createSession(title: string) {
    try {
      const response = await api.post('/chat/sessions', {
        title,
        provider: 'ollama',
        model: 'llama2',
      })
      const newSession: Session = {
        ...response.data,
        messages: [],
      }
      sessions.value.unshift(newSession)
      currentSessionId.value = newSession.id
    } catch (error) {
      console.error('Failed to create session:', error)
    }
  }

  async function deleteSession(sessionId: string) {
    try {
      await api.delete(`/chat/sessions/${sessionId}`)
      sessions.value = sessions.value.filter(s => s.id !== sessionId)
      if (currentSessionId.value === sessionId) {
        currentSessionId.value = sessions.value[0]?.id || null
      }
    } catch (error) {
      console.error('Failed to delete session:', error)
    }
  }

  async function sendMessage(content: string) {
    if (!currentSessionId.value) return

    try {
      const userMessage: Message = {
        id: Date.now().toString(),
        session_id: currentSessionId.value,
        role: 'user',
        content,
        created_at: new Date().toISOString(),
      }
      const session = currentSession.value
      if (session) {
        session.messages.push(userMessage)
      }

      const response = await api.post('/chat/messages', {
        session_id: currentSessionId.value,
        content,
      })

      const assistantMessage: Message = {
        ...response.data,
        role: 'assistant',
      }
      if (session) {
        session.messages.push(assistantMessage)
      }
    } catch (error) {
      console.error('Failed to send message:', error)
    }
  }

  function setCurrentSession(sessionId: string) {
    currentSessionId.value = sessionId
    loadSessionMessages(sessionId)
  }

  return {
    sessions,
    currentSessionId,
    currentSession,
    fetchSessions,
    createSession,
    deleteSession,
    sendMessage,
    setCurrentSession,
  }
})
