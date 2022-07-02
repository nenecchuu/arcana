package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
)

func ParseJsonStructToBytesBuffer(req interface{}) (*bytes.Buffer, error) {
	var (
		body []byte
		err  error
	)

	body, err = json.Marshal(req)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return nil, err
	}

	return bytes.NewBuffer(body), nil
}

func ParseResponseBodyToJson(res *http.Response, v interface{}) error {
	var (
		err         error
		resBodyByte []byte
	)

	resBodyByte, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Err(err).Msg(err.Error())
		return nil
	}

	err = json.Unmarshal(resBodyByte, v)

	if err != nil {
		log.Err(err).Msg(err.Error())
		return err
	}

	return nil
}
