![](../../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2018

# Packages

**Go** programmers organize their code into a single **workspace**. A folder, usually within your home folder, pointed to by the **GOPATH** environment variable. I use a folder like this:

```
echo $GOPATH
/home/cleuton/go
```

You can update the environment variable within: ```/etc/bash.bashrc``` e ```/etc/profile.d```.

**GOPATH's** structure is like that: 

```
go 
+-src: projects grouped by repository and package
+-bin: commands (Go executables)
```

If you download a Go package, it will be installed into a subfolder under **src**: 

```
go get -v golang.org/x/text
```

Analizing **GOPATH** we can see the source code under **src**: 

**/home/cleuton/go/src/golang.org/x/text/...**

This is the source code of the downloaded package. And we'll see a folder **pkg** which contains the binaries: 

**/home/cleuton/go/pkg/linux_amd64/golang.org/x/text.a**

As you may have noticed, packages are grouped by the hostname of the repository they are in. In this case, the **text** package is under **golang.org**. There is no version control and no central package repository, as in other languages.

If you want to create a public **Go** package for others to use, just post it on **Github** and provide your address: github.com/**username**/**repository name**.

## Executable Packages and Libraries

Everything in **Go** must be in one package. For example, let's create a package inside $ GOPATH with the functions we use to calculate the quadratic equation (see the [**bhaskara file**](./code/bhaskara.go)).

Inside $GOPATH/src create a folder called "bhaskarautils" and copy this file there. You will have a structure like this: **$GOPATH/src/bhaskarautils/bhaskara.go**.

Now run the command below:

```
go install bhaskarautils
```

And you will see that a file **$GOPATH/pkg/linux_amd64/bhaskarautils.a** has been created. Now we can run the [**testbhaskara.go code**](./code/test/testbhaskara.go) without problems.

```
import (
    ...
	"bhaskarautils"
)
...
	deltaV := bhaskarautils.Delta(coefA, coefB, coefC)
	x1, x2 := bhaskarautils.Roots(deltaV, coefA, coefB)
```

I had to rename the functions **Delta()** and **Roots()** to uppercase, so the package will export them. Otherwise, they will be considered private.

The package **bhaskarautils** is a library as it contains no files with the **func main()**, so a library file (extension .a) will be generated.

If we compile an executable, we will generate a native operating system executable file:

```
go build testbhaskara.go
```

And we can run it as we do with any program.

## Repositories

You should create your project in a repository such as Github and import it into your workspace with ```go get <package>```.

## Is this the end?

Of course not. It is just the end of this basic course. I know there is a lot more to learn, but I believe you now have the indispensable tools to learn on your own. Good luck and follow [**golang.network**](http://golang.network).