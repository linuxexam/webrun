// websocket for admin
const wsAdmin = new WebSocket("/ws-admin");
wsAdmin.addEventListener("open", ()=>{
    console.log("ws-admin connected");
});
wsAdmin.addEventListener("close", ()=>{
    console.log("ws-admin closed.");
});
wsAdmin.addEventListener("error", ()=>{
    console.log("ws-admin error.");
});
wsAdmin.addEventListener("message", (e)=>{
    console.log("ws-admin message: " + e.data.toString())
})

// websocket for term data
const wsData = new WebSocket("/ws-term");
wsData.addEventListener("open", ()=>{
    console.log("WebSocket connected");
});
wsData.addEventListener("close", ()=>{
    console.log("WebSocket closed.");
    term.write('\r\n\x1B[31mDisconnected!\r\nRefresh this page to reconnect.\x1B[0m')
});
wsData.addEventListener("error", ()=>{
    console.log("WebSocket error.");
});
wsData.addEventListener("message", (e)=>{
    console.log("WebSocket message: " + e.data.toString())
})

const termContainer = document.getElementById("terminal");
const term = new Terminal();

term.onResize(({cols, rows})=>{
    console.log(`${cols}, ${rows}`);
    if(wsAdmin.readyState == WebSocket.OPEN){
        wsAdmin.send(`{ "type":"winsize", "cols": ${cols}, "rows": ${rows} }`);
    }
});

term.open(termContainer);

const attachAddon = new AttachAddon.AttachAddon(wsData);
term.loadAddon(attachAddon);

const fitAddon = new FitAddon.FitAddon();
term.loadAddon(fitAddon);

fitAddon.fit();

window.addEventListener("resize", ()=>{
    fitAddon.fit();
});