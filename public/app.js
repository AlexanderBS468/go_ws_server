const statusEl = document.getElementById("status");
const connectionForm = document.getElementById("connection-form");
const channelInput = document.getElementById("channel");
const connectButton = document.getElementById("connect");
const form = document.getElementById("message-form");
const leadForm = document.getElementById("lead-form");
const leadInput = document.getElementById("lead");
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
    sendPayload({
      event: "login",
      data: {
        hash: Math.random().toString(36).slice(2),
        id: Math.floor(Math.random() * 1000),
        fullname: "Browser User",
      },
    });
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

function sendPayload(payload) {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    return false;
  }

  socket.send(JSON.stringify(payload));
  writeLog("sent", JSON.stringify(payload, null, 2));
  return true;
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

  if (sendPayload(payload)) {
    input.select();
  }
});

leadForm.addEventListener("click", (event) => {
  if (!event.target.matches("button[data-event]")) {
    return;
  }

  const leadID = Number.parseInt(leadInput.value, 10);
  if (!leadID) {
    return;
  }

  sendPayload({
    event: event.target.dataset.event,
    data: {
      leadId: leadID,
    },
  });
});
