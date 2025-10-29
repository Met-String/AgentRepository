<template>
  <form class="input-row" @submit.prevent="handleSubmit">
    <textarea
      ref="textarea"
      v-model="text"
      class="chat-input"
      placeholder="输入消息"
      @keydown="handleKeydown"
      @input="autoResize"
      rows="1"
      required
    ></textarea>
    <button type="submit">发送</button>
  </form>
</template>

<script setup>
// 不用TS，就纯JS
import { ref, onMounted } from 'vue'

const text = ref('')
const textarea = ref(null)

// ✅ 自动调节高度（最多 10 行）
const autoResize = () => {
  const el = textarea.value
  if (!el) return
  el.style.height = 'auto' // 先清空，再计算实际高度
  const maxHeight = 10 * parseFloat(getComputedStyle(el).lineHeight)
  el.style.height = Math.min(el.scrollHeight, maxHeight) + 'px'
  el.style.overflowY = el.scrollHeight > maxHeight ? 'auto' : 'hidden'
}

// ✅ 提交逻辑
const handleSubmit = () => {
  const msg = text.value.trim()
  if (!msg) return
  // 👇 在这里发消息，比如 WebSocket
  console.log('发送消息:', msg)
  text.value = ''
  autoResize()
}

// ✅ 回车发送 / Shift+回车换行
const handleKeydown = (e) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    handleSubmit()
  }
}

// ✅ 挂载时初始化一次高度
onMounted(() => {
  autoResize()
})
</script>

<style scoped>
.input-row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
}

.chat-input {
  flex: 1;
  resize: none;
  overflow-y: hidden;
  height: auto;
  min-height: 1.5em;
  max-height: calc(1.5em * 10); /* 10行上限 */
  padding: 8px 10px;
  border: 1px solid #ccc;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.5em;
  outline: none;
  transition: border-color 0.2s;
}

.chat-input:focus {
  border-color: #5865f2; /* Discord蓝 */
}

button {
  padding: 8px 14px;
  background-color: #5865f2;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
