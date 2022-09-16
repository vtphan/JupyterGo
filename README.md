# JupyterGo

Guideline on how to create a Jupyter Lab extension that talks to a Go server.

----

## Server

+ main.go -- the main server
+ config.json -- configuration file with information about IP address, port, etc.


### Server controllers

+ echo - simply echos the "content" variable.

### How to run the server

```
go run main.go -c config.json
```

----

## Client

A button in Jupyter Lab, when clicked, sends the content of the current cell to the server.