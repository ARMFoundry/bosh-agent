package client

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	davconf "github.com/cloudfoundry/bosh-agent/davcli/config"
	boshhttp "github.com/cloudfoundry/bosh-agent/http"
)

type Client interface {
	Get(path string) (content io.ReadCloser, err error)
	Put(path string, content io.ReadCloser, contentLength int64) (err error)
}

func NewClient(config davconf.Config, httpClient boshhttp.Client) (c Client) {
	return client{
		config:     config,
		httpClient: httpClient,
	}
}

type client struct {
	config     davconf.Config
	httpClient boshhttp.Client
}

func (c client) Get(path string) (content io.ReadCloser, err error) {
	req, err := c.createReq("GET", path, nil)
	if err != nil {
		return
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		body := ""
		if resp.Body != nil {
			buf := make([]byte, 1024)
			n, _ := resp.Body.Read(buf)
			body = string(buf[0:n])
		}
		err = fmt.Errorf("Getting dav blob %s: Wrong response code: %d; body: %s", path, resp.StatusCode, body)
		return
	}

	content = resp.Body
	return
}

func (c client) Put(path string, content io.ReadCloser, contentLength int64) (err error) {
	req, err := c.createReq("PUT", path, content)
	if err != nil {
		return
	}
	defer content.Close()
	req.ContentLength = contentLength

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 201 && resp.StatusCode != 204 {
		body := ""
		if resp.Body != nil {
			buf := make([]byte, 1024)
			n, _ := resp.Body.Read(buf)
			body = string(buf[0:n])
		}
		err = fmt.Errorf("Putting dav blob %s: Wrong response code: %d; body: %s", path, resp.StatusCode, body)
		return
	}

	return
}

func (c client) createReq(method, blobID string, body io.Reader) (req *http.Request, err error) {
	blobURL, err := url.Parse(c.config.Endpoint)
	if err != nil {
		return
	}

	digester := sha1.New()
	digester.Write([]byte(blobID))
	blobPrefix := fmt.Sprintf("%02x", digester.Sum(nil)[0])

	newPath := path.Join(blobURL.Path, blobPrefix, blobID)
	if !strings.HasPrefix(newPath, "/") {
		newPath = "/" + newPath
	}

	blobURL.Path = newPath

	req, err = http.NewRequest(method, blobURL.String(), body)
	if err != nil {
		return
	}

	req.SetBasicAuth(c.config.User, c.config.Password)
	return
}
