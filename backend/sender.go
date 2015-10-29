package backend

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/jmcvetta/napping"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
)

//Response represent a response from the LB
type Response struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func init() {
	InitCredentials()
}

//Request represent a request to LB
func Request(method int, u string, body interface{}) (*Response, error) {
	InitSession()
	var (
		err error
		r   *napping.Response
	)
	if method == common.GET {
		glog.Infof("Contacting host: %s", u)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		glog.Infof("Contacting host: %s json data %s", u, string(data))
	}
	switch method {
	case common.GET:
		{
			r, err = sess.Get(u, nil, body, nil)
		}
	case common.POST:
		{
			r, err = sess.Post(u, body, nil, nil)
		}
	case common.PUT:
		{
			r, err = sess.Put(u, body, nil, nil)
		}
	case common.PATCH:
		{
			r, err = sess.Patch(u, body, nil, nil)
		}
	case common.DELETE:
		{
			r, err = sess.Delete(u, nil, nil, nil)
		}
	}
	return &Response{r.Status(), r.RawText()}, err
}
