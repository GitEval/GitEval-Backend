package main

import "flag"

var (
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "conf/config.yaml", "config path, eg: -conf conf/config.yaml")
}
func main() {
	flag.Parse()
	app, clean := WireApp(flagconf)
	//启动清理map的协程
	go clean()
	app.Run()
}
