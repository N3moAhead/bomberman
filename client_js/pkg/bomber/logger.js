const util = require('util');

const logPrefix = "[BOMBER] ";
const colorReset = "\x1b[0m";
const colorRed = "\x1b[31m";
const colorGreen = "\x1b[32m";
const colorBlue = "\x1b[34m";
const colorCyan = "\x1b[36m";

function green(format, ...args) {
    const message = util.format(format, ...args);
    return `${colorGreen}${message}${colorReset}`;
}

function blue(format, ...args) {
    const message = util.format(format, ...args);
    return `${colorBlue}${message}${colorReset}`;
}

function red(format, ...args) {
    const message = util.format(format, ...args);
    return `${colorRed}${message}${colorReset}`;
}

function success(format, ...args) {
    const message = util.format(format, ...args);
    console.log(`${logPrefix}${colorGreen}${message}${colorReset}`);
}

function error(format, ...args) {
    const message = util.format(format, ...args);
    console.log(`${logPrefix}${colorRed}${message}${colorReset}`);
}

function info(format, ...args) {
    const message = util.format(format, ...args);
    console.log(`${logPrefix}${colorBlue}${message}${colorReset}`);
}

function debug(format, ...args) {
    const message = util.format(format, ...args);
    console.log(`${logPrefix}${colorCyan}${message}${colorReset}`);
}

module.exports = {
    green,
    blue,
    red,
    success,
    error,
    info,
    debug,
};
