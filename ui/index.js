const termContainer = document.getElementById("terminal");
const term = new Terminal();
term.open(termContainer);

const fitAddon = new FitAddon.FitAddon();
term.loadAddon(fitAddon);

fitAddon.fit();
window.addEventListener("resize", ()=>{
    fitAddon.fit();
});

term.write('Hello from \x1B[1;3;31mxterm.js\x1B0m $');

const socket = new WebSocket("ws://localhost:8080/ws");
socket.addEventListener("open", ()=>{
    console.log("WebSocket connected");
});
socket.addEventListener("close", ()=>{
    console.log("WebSocket closed.");
});
socket.addEventListener("error", ()=>{
    console.log("WebSocket error.");
});

const attachAddon = new AttachAddon.AttachAddon(socket);
term.loadAddon(attachAddon);
