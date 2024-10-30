const termContainer = document.getElementById("terminal");
const term = new Terminal();
term.open(termContainer);

const fitAddon = new FitAddon.FitAddon();
term.loadAddon(fitAddon);

fitAddon.fit();
window.addEventListener("resize", ()=>{
    fitAddon.fit();
});

//term.write('Hello from \x1B[1;3;31mxterm.js\x1B[0m $');

const socket = new WebSocket("/ws");
socket.addEventListener("open", ()=>{
    console.log("WebSocket connected");
});
socket.addEventListener("close", ()=>{
    console.log("WebSocket closed.");
    term.write('\r\n\x1B[31mDisconnected!\r\nRefresh this page to reconnect.\x1B[0m')
});
socket.addEventListener("error", ()=>{
    console.log("WebSocket error.");
});
socket.addEventListener("message", (e)=>{
    console.log("WebSocket message: " + e.data.toString())
})

const attachAddon = new AttachAddon.AttachAddon(socket);
term.loadAddon(attachAddon);
