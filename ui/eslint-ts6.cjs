const Module = require("node:module");
const originalRequire = Module.prototype.require;
Module.prototype.require = function (id) {
  if (id === "typescript") {
    return require("typescript6");
  }
  return originalRequire.apply(this, arguments);
};
