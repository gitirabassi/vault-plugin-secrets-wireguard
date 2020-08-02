package main

import (
	"fmt"
	"io/ioutil"
	"path"
)

const corefile = `
. {
    bind %v
    forward . 8.8.8.8 9.9.9.9
    prometheus
    debug
    log
    errors
    ready
    cache 120
}
`

func renderCorefile(destinationFolder, gatewayIP string) error {
	body := fmt.Sprintf(corefile, gatewayIP)
	fullPath := path.Join(destinationFolder, "Corefile")
	return ioutil.WriteFile(fullPath, []byte(body), 0644)
}
