'use strict';

(function (g, factory) {
  if (typeof exports === 'object' && typeof module !== 'undefined') {
    console.log("only support amd require!");
  } else {
    console.log("AMD require!");
    factory(g);
  }
}(window, function (g) { 
  

  const depsModules = [
      // "../common/jQuery/jquery-3.3.1.min.js",
      // "../common/vue/vue.js",
      // "../common/wcs_cookies.js"
  ];

  define(depsModules, function(){
    let tk = getCookie("token");
    console.log(tk);
    $.post("livelist", {token: tk}).done(replySuccess).fail(replyFail);
    return {};
  })
  
  


  function replySuccess(data, status)
  {
    const data2 = eval(data);
    console.log(data2[1]);
    const list = $("#list");
    for (let i = 0; i < data2.length; ++i) {
      list.append($("<button type='button'></input>").text(data2[i].Rname)).append($("<br />"));
    }

    
  }
  function replyFail(o, status, errT)
  {
    const list = $("#list");
    list.html($("<span></span>").text("error:" + o.responseText));

  }
  
}));