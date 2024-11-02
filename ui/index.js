const T = {
    term: null,
    wsAdmin: null,
    wsTerm: null,
    rows: 24,
    cols: 80,
};

// websocket for admin
function initWsAdmin() {
    T.wsAdmin = new WebSocket("/ws-admin");
    T.wsAdmin.addEventListener("open", () => {
        console.log("ws-admin connected");
        resizeTerminal();
    });
    T.wsAdmin.addEventListener("close", () => {
        console.log("ws-admin closed.");
    });
    T.wsAdmin.addEventListener("error", () => {
        console.log("ws-admin error.");
    });
    T.wsAdmin.addEventListener("message", (e) => {
        console.log("ws-admin message: " + e.data.toString())
    })
}

function initWsTerm() {
    // websocket for term data
    T.wsTerm = new WebSocket("/ws-term");
    T.wsTerm.addEventListener("open", () => {
        console.log("WebSocket connected");
        // connect terminal io with websocket  
        const attachAddon = new AttachAddon.AttachAddon(T.wsTerm);
        T.term.loadAddon(attachAddon);

        // admin is connected after term created on the server side
        initWsAdmin();
    });

    T.wsTerm.addEventListener("close", () => {
        console.log("WebSocket closed.");
        T.term.write('\r\n\x1B[31mDisconnected!\r\nRefresh this page to reconnect.\x1B[0m')
    });
    T.wsTerm.addEventListener("error", () => {
        console.log("WebSocket error.");
    });
}

// create terminal
function init() {
    const termContainer = document.getElementById("terminal");
    T.term = new Terminal();
    T.term.open(termContainer);

    // inform the server side terminal size changes
    T.term.onResize(({ cols, rows }) => {
        console.log(`new winsize: ${cols}, ${rows}`);
        T.cols = cols;
        T.rows = rows;
        resizeTerminal();
    });

    // keep the terminal fit its container
    const fitAddon = new FitAddon.FitAddon();
    T.term.loadAddon(fitAddon);
    fitAddon.fit();

    window.addEventListener("resize", () => {
        fitAddon.fit();
    });

    initWsTerm();
}

function resizeTerminal() {
    if(T.wsAdmin == null){
        return;
    }
    if (T.wsAdmin.readyState != WebSocket.OPEN) {
        return;
    }
    T.wsAdmin.send(JSON.stringify(
        {
            type: "winsize",
            cols: T.cols,
            rows: T.rows,
        }
    ));
}

init();