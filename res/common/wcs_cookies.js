'use strict';

(function (g, factory) {
  if (typeof exports === 'object' && typeof module !== 'undefined') {
    console.log("only support amd require!");
  } else {
    console.log("AMD require!");
    factory(g);
  }
}(window, function (g) { 
  function getCookie(c_name)
  {
    if (document.cookie.length>0)
    {
      let c_start;
      let c_end;
      c_start=document.cookie.indexOf(c_name + "=")
      if (c_start!=-1)
      { 
        c_start=c_start + c_name.length+1 
        c_end=document.cookie.indexOf(";",c_start)
        if (c_end==-1) c_end=document.cookie.length
          return unescape(document.cookie.substring(c_start,c_end))
      } 
    }
    return ""
  }
  
  function setCookie(c_name,value,expiredays)
  {
    var exdate=new Date()
    exdate.setDate(exdate.getDate()+expiredays)
    document.cookie=c_name+ "=" +escape(value)+
      ((expiredays==null) ? "" : ";expires="+exdate.toGMTString()) + 
      ";path=/"
  }
  
  
  g.setCookie = setCookie;
  g.getCookie = getCookie;
  
  if ( typeof define === "function" && define.amd ) {
    define("wcs_cookies", [], function() {
      return { setCookie: setCookie, getCookie: getCookie };
    })
  }
}));