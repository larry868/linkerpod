
spaInitWebAssembly();

/*
* Web Assembly
*/

async function spaInitWebAssembly() {
    if (!spaCanLoadWebAssembly()) {
        document.getElementById("spa-wasm-status").innerText = "unable to load the web assembly code";
        return;
    }

    try {
        const goWasm = new Go()

        WebAssembly.instantiateStreaming(fetch("spa.wasm"), goWasm.importObject)
            .then((result) => {
                goWasm.run(result.instance)

            })
    } catch (err) {
        document.getElementById("spa-wasm-status").innerText = "loading wasm failed:" + err;
        console.error("loading wasm failed: ", err);
    }
}

function spaCanLoadWebAssembly() {
    return !/bot|googlebot|crawler|spider|robot|crawling/i.test(
        navigator.userAgent
    );
}
