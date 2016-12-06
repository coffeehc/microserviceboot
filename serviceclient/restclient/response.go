package restclient

import (
	"net/http"

	"encoding/json"
	"github.com/coffeehc/microserviceboot/base"
	"github.com/coreos/etcd/Godeps/_workspace/src/github.com/gogo/protobuf/proto"
	"io"
	"io/ioutil"
	"net/url"
)

type Response interface {
	GetStatusCode() int
	GetBody() io.ReadCloser
	DecodeBody(decoder ResponseBodyDecoder, target interface{}) base.Error
}

func buildResponse(res *http.Response) Response {
	return &_Response{
		response: res,
	}
}

type _Response struct {
	response *http.Response
}

func (this *_Response) GetStatusCode() int {
	return this.response.StatusCode
}

func (this *_Response) GetBody() io.ReadCloser {
	return this.response.Body
}
func (this *_Response) DecodeBody(decoder ResponseBodyDecoder, target interface{}) base.Error {
	return decoder(this.GetBody(), target)
}

type ResponseBodyDecoder func(body io.ReadCloser, target interface{}) base.Error

func ResponseFormBodyDecoder(body io.ReadCloser, target interface{}) base.Error {
	defer body.Close()
	if vs, ok := target.(url.Values); ok {
		data, err := ioutil.ReadAll(body)
		if err != nil {
			return base.NewErrorWrapper(err)
		}
		values, err1 := url.ParseQuery(string(data))
		if err1 != nil {
			return base.NewErrorWrapper(err1)
		}
		for k, vss := range values {
			for _, v := range vss {
				vs.Add(k, v)
			}
		}
		return nil
	}
	return base.NewError(-1, "target type is not url.Value")
}

func ResponsePBBodyDecoder(body io.ReadCloser, target interface{}) base.Error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return base.NewErrorWrapper(err)
	}
	if message, ok := target.(proto.Message); ok {
		err = proto.Unmarshal(data, message)
		if err != nil {
			return base.NewErrorWrapper(err)
		}
		return nil
	}
	return base.NewError(-1, "target type is not proto.Message")
}

func ResponseJsonBodyDecoder(body io.ReadCloser, target interface{}) base.Error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return base.NewErrorWrapper(err)
	}
	err = json.Unmarshal(data, target)
	if err != nil {
		return base.NewErrorWrapper(err)
	}
	return nil
}