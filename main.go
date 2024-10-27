package main

import "flag"

var (
	flagconf string 
)
func init(){
	flag.StringVar(&flagconf, "conf", "conf/config.yaml", "config path, eg: -conf conf/config.yaml")
}
func main(){
	flag.Parse()
	app := WireApp(flagconf)
	app.Run()
}

