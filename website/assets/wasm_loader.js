
const start = Date.now();
initWebAssembly();

/*
* Web Assembly
*/

async function initWebAssembly() {
    ews = document.getElementById("webapp");

    if (!canLoadWebAssembly()) {
        if (ews !== null) {
            ews.innerText = "Unable to load the App in this browser.";
        }
        console.error("unable to load the web assembly code with this useragent");
        return;
    }

    const goWasm = new Go()

    WebAssembly.instantiateStreaming(fetch("webapp.wasm"), goWasm.importObject)
        .then((result) => {
            goWasm.run(result.instance);
            const end = Date.now();
            console.log("webapp.wasm code loaded in " + (end - start) + "ms");
        })
        .catch((err) => {
            if (ews !== null) {
                ews.innerText = "Error loading the App.";
            }
            console.error("loading wasm failed:" + err);
        })
}

function canLoadWebAssembly() {
    return !/bot|googlebot|crawler|spider|robot|crawling/i.test(
        navigator.userAgent
    );
}

function ickError(msg) { console.error(msg) }

function ickWarn(msg) { console.warn(msg) }
