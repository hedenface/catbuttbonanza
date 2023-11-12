package main

import (
	"bytes"
	"fmt"
	"text/template"
)

func htmlTemplate(t string, vars interface{}) string {
	tmpl := template.New("template")
	tmpl, err := tmpl.Parse(t)
	if err != nil {
		return fmt.Sprintf("Error %v for tmpl.Parse(t)", err)
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, vars)
	if err != nil {
		return fmt.Sprintf("Error %v for tmpl.Execute()", err)
	}

	return b.String()
}
func htmlTemplateFormLogin(vars interface{}) string {
	return htmlTemplate(`
<form method="post" action="/login">
    <label>
        Username
        <input type="text" name="username">
    </label>
    <label>
        Password
        <input type="password" name="password">
    </label>
    <input type="submit" value="Login">
</form>
`, vars)
}

func htmlTemplatePage(vars interface{}) string {
	return htmlTemplate(`
<!doctype html>
<html lang="en">
  <head>
    <title>{{ .Title }}</title>
    <meta charset="utf-8">

    <link rel="icon" type="image/png" href="/favicon.png">

    <meta name="viewport" content="width=device-width, initial-scale=1">

    <!-- all for bootstrap -->
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-rbsA2VBKQhggwzxH7pPCaAqO46MgnOM80zW1RWuH61DGLwZJEdK2Kadq2F9CUG65" crossorigin="anonymous">
    <script src="https://code.jquery.com/jquery-3.3.1.min.js" integrity="sha256-FgpCb/KJQlLNfOu91ta32o/NMZxltwRo8QtmkMRdAu8=" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.11.6/dist/umd/popper.min.js" integrity="sha384-oBqDVmMz9ATKxIep9tiCxS/Z9fNfEXiDAYTujMAeBAsjFuCZSmKbSSUnQlmh/jp3" crossorigin="anonymous"></script>
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.3/dist/js/bootstrap.min.js" integrity="sha384-cuYeSxntonz0PPNlHhBs68uyIAVpIIOZZ5JqeqvYYIcEL727kskC66kF92t6Xl2V" crossorigin="anonymous"></script>

    <!-- all for a few lil icons -->
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.10.3/font/bootstrap-icons.css">
{{ .Head }}
  </head>
  <body>
{{ .Body }}
  </body>
</html>
`, vars)
}
