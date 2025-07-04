// --- Chat Logic ---

const socket = new WebSocket('ws://localhost:8080/ws');
socket.addEventListener('close', function(event) {
    alert("Room is full (max 4 players). Please try again later.");
});
// let playerName = prompt("Enter your name:") || "Player" + Math.floor(Math.random() * 1000);

const chatMessages = document.getElementById('chat-messages');
const chatInput = document.getElementById('chat-input');
const chatSend = document.getElementById('chat-send');

let playerName = "";
let playerColor = "#c8a060";

window.onload = function() {
    document.getElementById('player-setup').style.display = 'flex';
    // Optionally, disable chat input/button until player has joined
    chatInput.disabled = true;
    chatSend.disabled = true;
    document.getElementById('player-name').focus();
}

document.getElementById('player-join').onclick = function() {
    playerName = document.getElementById('player-name').value || "Player" + Math.floor(Math.random() * 1000);
    playerColor = document.getElementById('player-color').value || "#c8a060";
    document.getElementById('player-setup').style.display = 'none';
    chatInput.disabled = false;
    chatSend.disabled = false;
    chatInput.focus();
};

chatSend.onclick = sendCurrentChat;
chatInput.addEventListener('keydown', function(e) {
    if (e.key === 'Enter') sendCurrentChat();
});

function sendCurrentChat() {
    const text = chatInput.value.trim();
    if (text.length > 0) {
        sendChat(text);
        chatInput.value = '';
    }
}

function sendChat(text) {
    socket.send(JSON.stringify({ 
        type: "chat", 
        name: playerName,
        color: playerColor,
        text ,
        timestamp: Date.now()
    }));
}

socket.addEventListener('message', function (event) {
    let msg;
    try { msg = JSON.parse(event.data); } catch { return; }
    if (msg.type === "chat") {
        const div = document.createElement('div');
        const time = msg.timestamp
            ? new Date(msg.timestamp).toLocaleTimeString([], 
                {hour: '2-digit', minute:'2-digit', second:'2-digit'})
            : '';
        const nameColor = msg.color || "#c8a060"; // fallback color
        div.innerHTML = `<strong style="color:${nameColor};">${msg.name || 'Unknown'}</strong>            <span class="chat-time">[${time}]</span> : 
            <span class="chat-text">${msg.text}</span>`;
        chatMessages.appendChild(div);
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }
});


