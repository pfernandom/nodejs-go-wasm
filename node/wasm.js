const util = require("util");
const fs = require("fs");
require("./wasm_exec");

function importWasm(callback = () => {}) {
  var source = fs.readFileSync("./app.wasm");

  const go = new globalThis.Go();
  var typedArray = new Uint8Array(source);

  globalThis.logCallback = callback;

  WebAssembly.instantiate(typedArray.buffer, go.importObject)
    .then((result) => {
      console.log(util.inspect(result, true, 0));
      return go.run(result.instance);
    })
    .catch((e) => {
      console.log(e);
    });
}

module.exports.importWasm = importWasm;
