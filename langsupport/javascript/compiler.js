(function(module, exports) {
  exports._compile = function(src) {
    return babel.transform(src, { presets: ["es2015"] }).code;
  }
})