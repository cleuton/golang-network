package main

import (
	"github.com/webview/webview"
)

var w webview.WebView

func main() {
	debug := true
	page := `data:text/html,
	<!DOCTYPE html>	
	<html>
		<body>
			<label for="txt01">Nome</label>
			<input type="text" id="txt01">
			<button onclick="tratar()">Enviar</button>
			<script>
				function tratar() {
					go_click(document.getElementById("txt01").value)
						.then(function (campo) {
							window.alert(campo)
						})
				}
			</script>
		</body>
	</html>
	`

	w = webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)
	w.Bind("go_click", func(x string) string { return "Digitou: " + x })
	w.Init("window.alert('teste')")
	w.Navigate(page)
	w.Run()
}
