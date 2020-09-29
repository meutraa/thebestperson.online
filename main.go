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
