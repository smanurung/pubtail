package main

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/linkedin/goavro"
)

type avroFormatter struct{}

func (af *avroFormatter) format(body []byte) ([]byte, error) {
	r := bytes.NewReader(body)
	br := bufio.NewReader(r)

	ocfr, err := goavro.NewOCFReader(br)
	if err != nil {
		return nil, err
	}

	var datum interface{}
	for ocfr.Scan() {
		datum, err = ocfr.Read()
		if err != nil {
			return nil, err
		}
	}

	return []byte(fmt.Sprint(datum)), nil
}
