# Sloth

Sloth is a simple interpeter written in Golang. This is hoby project for learning golang. I converted this to WebAssembly :) program, now i can run it on browser.

## Syntax

This is a simple syntax for sloth. currently, sloth works with following syntax in repl.

Variable declaration in sloth is like this:

    var name = value

1. Addition of two numbers

```shell
>>> var a = 10
>>> var b = 20
>>> a + b
30
```

Note: Sloth support all basic arithmetic operations except `>=` and `<=`.

2. Function

```shell
>>> var add = fun(a,b) { return a + b};
>>> add(1,2)
3
```

3. If else statement

```shell
>>> var isten = fun(a) { if(a == 10) { return "yes" } else { return "no"} }
>>> isten(10)
yes
>>> isten(11)
no
```

4. String concatenation

```shell
>>> var say = "hello" + " " + "world"
>>> say
hello world
>>> var say = concat("hello", "world")
>>> say
hello world
```
5. Built-in functions

- `len(a)`: returns length of array
- `concat(a,b)`: concatenates two strings

6. Array

```shell
>>> var a = [1,2 + 2,3]
>>> a[0]
1
>>> a[1]
4
```

## Build

1. WASM build

    $ cd wasm
    $ $ENV:GOOS=js 
    $ $ENV:GOARCH=wasm 
    $ go build -o sloth.wasm

# Btw

I wrote this project based on Thorsten Ball's [Writing an Interpreter in Go](https://interpreterbook.com/).