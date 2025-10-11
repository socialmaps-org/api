package templates

import "text/template"

var Index = template.Must(template.ParseFiles("web/template/index.gohtml", "web/template/base.gohtml"))
var Login = template.Must(template.ParseFiles("web/template/login.gohtml", "web/template/base.gohtml"))
var Logout = template.Must(template.ParseFiles("web/template/logout.gohtml", "web/template/base.gohtml"))
