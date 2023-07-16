
initWebAssembly();

/*
* Web Assembly
*/

async function initWebAssembly() {
    ews = document.getElementById("wasm-status");

    if (!canLoadWebAssembly()) {
        msg = "unable to load the web assembly code with this useragent";
        if (ews !== null) {
            ews.innerText = msg;
        }
        console.error(msg);
        return;
    }

    const goWasm = new Go()

    WebAssembly.instantiateStreaming(fetch("webapp.wasm"), goWasm.importObject)
        .then((result) => {
            goWasm.run(result.instance)
        })
        .catch((err) => {
            msg = "loading wasm failed:" + err
            if (ews !== null) {
                ews.innerText = msg;
            }
            console.error(msg);
        })
}

function canLoadWebAssembly() {
    return !/bot|googlebot|crawler|spider|robot|crawling/i.test(
        navigator.userAgent
    );
}
