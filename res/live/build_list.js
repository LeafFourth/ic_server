(function (g, factory) {
  if (typeof exports === 'object' && typeof module !== 'undefined') {
    console.log("only support amd require!");
  } else {
    console.log("AMD require!");
    factory(g);
  }
}(window, function (g) { 
  'use strict';

  g.f = function () {
    console.log("f()");
  }
}));