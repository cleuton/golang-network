![](../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2021

# Ui simplificada com webview

Interface de Usuário (**UI**) é sempre um grande problema nas aplicações. É claro que você pode utilizar desenvolvimento web na maior parte das aplicações, mas sempre pode haver necessidade de criar uma aplicação **Desktop**.

E você tem que trabalhar isso **multiplataforma** (Linux, Windows ou MacOS). Há vários frameworks para isto, como o [**Qt**](https://www.qt.io/), que é multiplataforma, mas a complexidade de uso pode ser um impedimento. 

E por que não utilizar a tecnologia Web como frontend, mantendo suas regras de negócio compiladas dentro do seu código Golang? É isso que o componente [**WebView**](https://github.com/webview/webview) proporciona. Uma camada fina de Ui web que você pode interfacear com seu código Golang, em duas vias: Web->Golang e Golang->Web.

Infelizmente, o projeto é pouco documentado e carece de exemplos funcionais, mas eu testei para você e criei uma demonstração bem simples, que resolve 99,99% dos seus problemas. 

O que você precisa em uma UI? Enviar mensagens e componentes para o usuário e pegar o resultado de eventos. 

O [**código-fonte**](./main.go) deste projeto exemplifica isto: 

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

Para executá-lo, é preciso instalar o pacote: 
```
go get github.com/webview/webview
```

No caso do **Linux** você tem que ter os pacotes GTK3 (ou 4), GtkWebkit2 e também o Build Essentials. 

O código é auto explicativo, mas vou falar de alguns métodos relevantes: 

- **New**: Instancia a WebView;
- **Set..**: Muda as propriedades (título, tamanho etc);
- **Bind**: Associa um objeto **Golang** a um objeto no contexto global **Javascript**. Neste caso, estou associando uma função anônima;
- **Eval**: Executa código Javascript no contexto da web. Serve para alterar propriedades;
- **Init**: Envia código **Javascript** para ser executado diretamente na página;
- **Navigate**: Navega para a URI. Pode ser um arquivo externo, um site externo ou, utilizando o protocolo "data", um string gerado dentro do seu próprio código, como neste exemplo;

Neste caso, eu envio uma página e um script para serem executados na interface. Ao digitar um nome e clicar no botão, os campos são enviados para o código **Golang**. 

[**Este artigo**](https://medium.com/@master.rta/golang-create-a-web-view-app-for-any-platform-54917dea397) mostra como empacotar sua app e distribuir. 

