(function (g, factory) {
  if (typeof exports === 'object' && typeof module !== 'undefined') {
    console.log("only support amd require!");
  } else {
    console.log("AMD require!");
    factory(g);
  }
}(window, function (g) { 
  'use strict';

  let con = $("#list");
  for (let i = 0; i < 10; ++i) {
    con.append($("<span></span>").text("line " + i));
    con.append($("<br />"));
  }
  
}));