// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"html/template"
	"log"
	"time"
	"bytes"
	"encoding/base64"
	"net/http"
)

var body = "Me"
const page string = `<!DOCTYPE html>
<html>
<head>
<title>{{.}}</title>
</head>
<style>
  html {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translateY(-50%) translateX(-50%);
    text-align: center;
  }
</style>
<div style="font-size:48px">is</div>
<form 
  action="save">
  <input 
    name=b 
    onfocus="var t=this.value; this.value=''; this.value=t"
    autofocus 
    style="border:none;text-align:center;font-size:72px" 
    value="{{.}}">
  </form>
</html>`
var tmpl = template.Must(template.New("page").Parse(page))

//const favicon = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 36 36"><path fill="#FFAC33" d="M5.123 5h6C12.227 5 13 4.896 13 6V4c0-1.104-.773-2-1.877-2h-8c-2 0-3.583 2.125-3 5 0 0 1.791 9.375 1.917 9.958C2.373 18.5 4.164 20 6.081 20h6.958c1.105 0-.039-1.896-.039-3v-2c0 1.104-.773 2-1.877 2h-4c-1.104 0-1.833-1.042-2-2S3.539 7.667 3.539 7.667C3.206 5.75 4.018 5 5.123 5zm25.812 0h-6C23.831 5 22 4.896 22 6V4c0-1.104 1.831-2 2.935-2h8c2 0 3.584 2.125 3 5 0 0-1.633 9.419-1.771 10-.354 1.5-2.042 3-4 3h-7.146C21.914 20 22 18.104 22 17v-2c0 1.104 1.831 2 2.935 2h4c1.104 0 1.834-1.042 2-2s1.584-7.333 1.584-7.333C32.851 5.75 32.04 5 30.935 5zM20.832 22c0-6.958-2.709 0-2.709 0s-3-6.958-3 0-3.291 10-3.291 10h12.292c-.001 0-3.292-3.042-3.292-10z"/><path fill="#FFCC4D" d="M29.123 6.577c0 6.775-6.77 18.192-11 18.192-4.231 0-11-11.417-11-18.192 0-5.195 1-6.319 3-6.319 1.374 0 6.025-.027 8-.027l7-.001c2.917-.001 4 .684 4 6.347z"/><path fill="#C1694F" d="M27 33c0 1.104.227 2-.877 2h-16C9.018 35 9 34.104 9 33v-1c0-1.104 1.164-2 2.206-2h13.917c1.042 0 1.877.896 1.877 2v1z"/><path fill="#C1694F" d="M29 34.625c0 .76.165 1.375-1.252 1.375H8.498C7.206 36 7 35.385 7 34.625v-.25C7 33.615 7.738 33 8.498 33h19.25c.759 0 1.252.615 1.252 1.375v.25z"/></svg>`

const favicon = `AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABPasHGT2nB/09pwf9PacH/T2nB/09pwf9PacH/T2nB/09pwf9PacDUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAATWvCMk9pwf1PacH/T2nB/09pwf9PacH/T2nB/09pwf9PacH9T2vDNwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABOasFvTW/IvkZ/1v9Gf9b/Rn/W/0aA1v5Nb8e8TmnCdQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA0rf/GM6z//zOs//8zrP+/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANKz/ezOs//8zrP//Mqz/dQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD26/31Ixv//Scj//z+8/4EAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE3M/yhNy//yTcz//03M//9MzP/4Tc3/OAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADKr/2EzrP/VM6z/3TOs/95Jx//2Tcz//03M//9NzP//Tcz//0rJ//o1rf/fM6z/3TOs/9gzq/90AAAAAC6u/xYzrP/7M63/zDOr/3dCv/+qTcz//03M//9NzP//Tcz//03M//9NzP//RcL/uDOr/3czrf/JM6z//zGq/yozrf9LM6z//zOt/0EA//8BTcz/1E3M//9NzP//Tcz//03M//9NzP//Tcz//03M/+lVxv8JNav/OjOs//8zrP9fM63/fDOs//wuov8LTc3/TE3M//9NzP//Tcz//03M//9NzP//Tcz//03M//9NzP//Ts3/ZiS2/wczrP/4M6z/izSs/60zrP/PAAAAAEzM/5pNzP//Tcz//03M//9NzP//Tcz//03M//9NzP//Tcz//03M/7QAAAAAM63/yTOs/7gzrP/eM63/mwAAAABNy//GTcz//03M//9NzP//Tcz//03M//9NzP//Tcz//03M//9NzP/gAAAAADOr/5UzrP/mM6z/8zKt/8AzrP83TMr/2U3M//9NzP//Tcz//03M//9NzP//Tcz//03M//9NzP//Tcz/7jSq/zYzq/+9M6z/+TOs/5wzrP//M6z//0bD//9NzP//Tcz//03M//9NzP//Tcz//03M//9NzP//Tcz//0nH//8zrP//M6z//zKs/6IAAAAAMbH/GjWt/yJJx/9pTcv/3k3M/91NzP/dTcz/3U3M/91NzP/dTcz/3U3K/99HyP9zNa3/Ii+q/xsAAAAA4AcAAOAHAADwDwAA/D8AAPw/AAD8PwAA+B8AAIABAAAAAAAAAAAAAAAAAAAgBAAAIAQAAAAAAAAAAAAAgAEAAA==`

func main() {
	fc, _ := base64.StdEncoding.DecodeString(favicon)
	var time time.Time

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, body); nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		body = r.FormValue("b")
		http.Redirect(w, r, "/", http.StatusFound)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "favicon.ico", time, bytes.NewReader(fc))
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
