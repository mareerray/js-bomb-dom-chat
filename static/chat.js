// --- Chat Logic ---

const socket = new WebSocket('ws://localhost:8080/ws');

const chatMessages = document.getElementById('chat-messages');
const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');

chatSend.onclick = sendCurrentChat;
chatInput.addEventListener('keydown', function(e) {
    if (e.key === 'Enter') sendCurrentChat();
});

function sendCurrentChat() {
    const text = chatInput.value.trim();
    if (text) {
        sendChat(text);
        chatInput.value = '';
    }
}

function sendChat(text) {
    socket.send(JSON.stringify({ type: "chat", text }));
}

socket.addEventListener('message', function (event) {
    let msg;
    try { msg = JSON.parse(event.data); } catch { return; }
    if (msg.type === "chat") {
        const div = document.createElement('div');
        div.textContent = msg.text;
        chatMessages.appendChild(div);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
});
