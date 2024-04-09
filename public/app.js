const statusEl = document.getElementById("status");
const connectionForm = document.getElementById("connection-form");
const channelInput = document.getElementById("channel");
const connectButton = document.getElementById("connect");
const form = document.getElementById("message-form");
const eventInput = document.getElementById("event");
const input = document.getElementById("message");
const log = document.getElementById("log");
const button = form.querySelector("button");

const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
let socket;
button.disabled = true;

function writeLog(prefix, message) {
  const entry = document.createElement("div");
  entry.className = "log-entry";
  entry.textContent = `[${new Date().toLocaleTimeString()}] ${prefix}: ${message}`;
  log.appendChild(entry);
  log.scrollTop = log.scrollHeight;
}

function setConnected(channel) {
  statusEl.textContent = `Connected: ${channel}`;
  statusEl.className = "status status-on";
  button.disabled = false;
  connectButton.textContent = "Reconnect";
}

function setDisconnected() {
  statusEl.textContent = "Disconnected";
  statusEl.className = "status status-off";
  button.disabled = true;
}

function connect(channel) {
  if (socket) {
    socket.close();
  }

  socket = new WebSocket(`${protocol}//${window.location.host}/ws/${encodeURIComponent(channel)}`);

  socket.addEventListener("open", () => {
    setConnected(channel);
    writeLog("system", `connected to channel ${channel}`);
  });

  socket.addEventListener("message", (event) => {
    try {
      const payload = JSON.parse(event.data);
      writeLog("received", JSON.stringify(payload, null, 2));
    } catch (error) {
      writeLog("received", event.data);
    }
  });

  socket.addEventListener("close", () => {
    setDisconnected();
    writeLog("system", "connection closed");
  });

  socket.addEventListener("error", () => {
    writeLog("system", "WebSocket error");
  });
}

connectionForm.addEventListener("submit", (event) => {
  event.preventDefault();

  const channel = channelInput.value.trim() || "system";
  connect(channel);
});

form.addEventListener("submit", (event) => {
  event.preventDefault();

  const message = input.value.trim();
  if (!message || socket.readyState !== WebSocket.OPEN) {
    return;
  }

  const payload = {
    event: eventInput.value.trim() || "message",
    data: message,
  };

  socket.send(JSON.stringify(payload));
  writeLog("sent", JSON.stringify(payload, null, 2));
  input.select();
});
