![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2021

# Simplified Ui with Webview

User Interface (**UI**) is always a big problem in applications.You can use web development in most applications, but there may always be a need to create a **Desktop** application.

And you have to work it out to be **multiplatform** (Linux, Windows or MacOS). There are several frameworks for this, such as [**Qt**](https://www.qt.io/), which is multiplatform, but its can be an impediment.

And why not use Web technology as a frontend, keeping your business rules compiled within your Golang code? This is what the [**WebView**](https://github.com/webview/webview) component provides. A thin layer of web Ui that you can interface with your Golang code, in two ways: Web-> Golang and Golang-> Web.

Unfortunately, the project is poorly documented and lacks functional examples, but I tested it for you and created a very simple demo that solves 99.99% of your problems.

What do you need in a UI? Send messages and components to the user and get the result of events.

The [**source code**](../../portuguese/uidemo/main.go) of this project exemplifies this:

```
package main

import (
	"github.com/webview/webview"
)

var w webview.WebView

func acerta() {
	w.Eval("document.getElementById('txt01').value = '*****'")
}

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
	w.Bind("go_click", func(x string) string {
		acerta()
		return "Digitou: " + x
	})
	w.Init("window.alert('teste')")
	w.Navigate(page)
	w.Run()
}

```

To run it, install: 
```
go get github.com/webview/webview
```

On **Linux** systems you may have to install the GTK3 (or 4), GtkWebkit2 and Build Essentials packages.

The code is self-explanatory, but I will talk about some relevant methods:

- **New**: Instantiates the WebView;
- **Set..**: Change its properties (title, size, etc.);
- **Bind**: Associates a **Golang** object to an object in the **Javascript** global context. In this case, I am associating an anonymous function;
- **Eval**: Run Javascript in the web context. You can use it to change elements' properties;
- **Init**: Send **Javascript** code to be executed directly on the page;
- **Navigate**: Navigate to the URI. It can be an external file, an external website or, using the "data" protocol, a string generated inside the golang code, as in this example;

In this case, I send a page and a script to be executed in the interface. When entering a name and clicking the button, the fields are sent to the code **Golang**.

[**This article**] (https://medium.com/@master.rta/golang-create-a-web-view-app-for-any-platform-54917dea397) shows you how to package your app and distribute.