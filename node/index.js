const { importWasm } = require("./wasm.js");

importWasm((log) => {
  console.log({ log });
});
