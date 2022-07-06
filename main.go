package main

import servers "gorski.mateusz/webcalc/server"

func main() {
	go servers.StartServer()
	servers.StartHealthcheck()
}
