package substko

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

func getBody(subdomain string, https bool, timeout int) (body []byte) {
	var url string

	if https {
		url = fmt.Sprintf("https://%s/", subdomain)
	} else {
		url = fmt.Sprintf("http://%s/", subdomain)
	}

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.Add("Connection", "close")

	resp := fasthttp.AcquireResponse()

	client := &fasthttp.Client{
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client.DoTimeout(req, resp, time.Duration(timeout)*time.Second)

	return resp.Body()
}
