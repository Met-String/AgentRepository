<script setup>
import { ref, onMounted } from 'vue'

const chatLog = ref(null)
const chatForm = ref(null)
const chatInput = ref(null)
const sendButton = ref(null)

onMounted(() => {
    // 为聊天栏新增一条聊天消息，并自动跳转到最底部。
    function appendMessage(text, type = 'user') {
        const entry = document.createElement('div')
        entry.className = `message ${type}`;  // 根据Type区分不同类型的消息，拓展性做的很好
        entry.textContent = text; // 我相信textContent不是div独有的
        chatLog.value.appendChild(entry); // 非常经典的appendChild
        chatLog.value.scrollTop = chatLog.value.scrollHeight;// 相当于有新消息时，直接自动拉到最低
    }

    const socket = new WebSocket('ws://qd2.mossfrp.cn:33331/ws'); // 新建一个socket对象，大概。

    socket.addEventListener('open', () => { // 添加一堆监听器之open
        appendMessage('已连接到服务器。', 'system');  // 哦，而且居然真的有system消息类型？
        sendButton.value.disabled = false; // 连接之后是可以发送消息的，我相信后面有机制让它天然关闭。
    });

    socket.addEventListener('message', (event) => { // 添加一堆监听器之message
        appendMessage(event.data, 'server'); // 嗯...对于非用户自己发的消息...似乎是核心代码
    });

    socket.addEventListener('close', () => {
        appendMessage('连接已关闭。', 'system');
        sendButton.value.disabled = true;
    });

    socket.addEventListener('error', () => {
        appendMessage('连接出现错误。', 'system');
    });

    chatForm.value.addEventListener('submit', (event) => { 
        event.preventDefault();
        const text = chatInput.value.value.trim();
        if (!text) {
            return;
        }
        if (socket.readyState === WebSocket.OPEN) {
            socket.send(text);
            chatInput.value.value = '';
            chatInput.value.focus();
        } else {
            appendMessage('当前未连接到服务器。', 'system');
        }
    });

    sendButton.value.disabled = socket.readyState !== WebSocket.OPEN;

    // ====================
    // [回车发送/回车+Shift换行]
    chatInput.value.addEventListener('keydown', function(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
        e.preventDefault();  // 阻止换行
        chatForm.value.requestSubmit();
    }
    });
    // ====================
})
</script>

<template>
    <section class="panel chat-panel">
        <div ref="chatLog" class="chat-log"></div>
        <form ref="chatForm" id="chatForm" class="input-row" autocomplete="off">
            <textarea ref="chatInput" style="resize:none; height:60px; width: 100%;" required></textarea>
            <!-- <input id="chatInput" type="text" placeholder="输入消息" required /> -->
            <button ref="sendButton" type="submit">发送</button>
        </form>
    </section>
    <section class="panel placeholder-panel">
        <span>右侧留空，可按需扩展。</span>
    </section>
</template>

<style scoped>
    * {
        box-sizing: border-box;
        margin: 0;
        padding: 0;
    }

    body {
        font-family: Arial, sans-serif;
        height: 100vh;
        display: flex;
        background-color: #f0f0f0;
    }

    .panel {
        flex: 1;
        display: flex;
        flex-direction: column;
    }

    .chat-panel {
        background-color: #ffffff;
        border-right: 1px solid #ddd;
    }

    .chat-log {
        flex: 1;
        overflow-y: auto;
        padding: 16px;
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    .message.user {
        color: black;
        font-size: 14px;
        word-break: break-word;
        white-space: pre-wrap; /* 保留换行和空格 */
    }

    .message.system {
        color: #666;
        font-style: italic;
        white-space: pre-wrap; /* 保留换行和空格 */
    }

    .input-row {
        display: flex;
        border-top: 1px solid #ddd;
        padding: 12px;
        gap: 8px;
    }

    .input-row input[type="text"] {
        flex: 1;
        padding: 8px 10px;
        border: 1px solid #bbb;
    }

    .input-row button {
        padding: 8px 16px;
        border: none;
        background-color: #1976d2;
        color: #fff;
        cursor: pointer;
    }

    .input-row button:disabled {
        background-color: #999;
        cursor: not-allowed;
    }

    .placeholder-panel {
        justify-content: center;
        align-items: center;
        color: #999;
    }
</style>

