package main

type textFormatter struct{}

func (tf *textFormatter) format(body []byte) ([]byte, error) {
	return body, nil
}
