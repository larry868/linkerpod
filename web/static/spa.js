
spaInitWebAssembly();

/*
* Web Assembly
*/

async function spaInitWebAssembly() {
    ews = document.getElementById("spa-wasm-status");

    if (!spaCanLoadWebAssembly()) {
        msg = "unable to load the web assembly code with this useragant";
        if (ews !== null) {
            ews.innerText = msg;
        }
        console.error(msg);
        return;
    }

    const goWasm = new Go()

    WebAssembly.instantiateStreaming(fetch("spa.wasm"), goWasm.importObject)
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

function spaCanLoadWebAssembly() {
    return !/bot|googlebot|crawler|spider|robot|crawling/i.test(
        navigator.userAgent
    );
}
