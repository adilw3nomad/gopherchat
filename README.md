# gopherchat
A simple chat application which uses TCP sockets and go channels. I made this to learn more about Go channels and the `net` package

To run it, first build the binary using `go build`. 

Then start the server by running this command: 

```
gopherchat --mode server
```

After that, you can add as many clients as you want by running 

```
gopherchat --mode client
```

