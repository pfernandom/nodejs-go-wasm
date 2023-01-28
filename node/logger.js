const fs = require("fs");
const proc = require("process");
const util = require("util");

var logFile = fs.createWriteStream("test.log", { flags: "a" });

function log() {
  // console.log(...arguments);
  logFile.write(util.format.apply(null, arguments) + "\n");
}

function randomElement(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

function randomTime() {
  const time = [1000, 100, 3000, 0];
  return randomElement(time);
}

function randomLevel() {
  const levels = ["WARN", "INFO", "ERROR"];
  return randomElement(levels);
}

let i = 0;
function startLogger() {
  log(`${randomLevel()} Log number ${i++}`);

  setTimeout(() => {
    startLogger();
  }, randomTime());
}

startLogger();
