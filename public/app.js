const statusEl = document.getElementById("status");
const form = document.getElementById("message-form");
const eventInput = document.getElementById("event");
const input = document.getElementById("message");
const log = document.getElementById("log");
const button = form.querySelector("button");

const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
const socket = new WebSocket(`${protocol}//${window.location.host}/ws`);

function writeLog(prefix, message) {
  const entry = document.createElement("div");
  entry.className = "log-entry";
  entry.textContent = `[${new Date().toLocaleTimeString()}] ${prefix}: ${message}`;
  log.appendChild(entry);
  log.scrollTop = log.scrollHeight;
}

socket.addEventListener("open", () => {
  statusEl.textContent = "Connected";
  statusEl.className = "status status-on";
  button.disabled = false;
  writeLog("system", "connected to WebSocket");
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
  statusEl.textContent = "Disconnected";
  statusEl.className = "status status-off";
  button.disabled = true;
  writeLog("system", "connection closed");
});

socket.addEventListener("error", () => {
  writeLog("system", "WebSocket error");
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
