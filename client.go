
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"crypto/tls"
)

func main(){

	url := "https://192.168.100.3"
	req, _ := http.NewRequest("GET", url, nil)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	res, _ := client.Do(req)
	
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)

	sadeas;dfkldiskx
	asdfasdfasdfsadfs
	fmt.Println(string(body))
}