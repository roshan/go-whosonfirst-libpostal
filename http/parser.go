package http

import (
	"github.com/openvenues/gopostal/parser"
	"github.com/whosonfirst/go-whosonfirst-log"
	gohttp "net/http"
	"time"
)

func ParserHandler(logger *log.WOFLogger) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		address, err := GetAddress(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), err.Code)
			return
		}

		t1 := time.Now()

		defer func() {
			t2 := time.Since(t1)
			logger.Status("parse '%s' %v", address, t2)
		}()

		parsed := postal.ParseAddress(address)

		query := req.URL.Query()

		if query.Get("format") == "keys" {
			WriteResponse(rsp, FormatParsed(parsed))
		} else {
			WriteResponse(rsp, parsed)
		}
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}

func MultiParserHandler(logger *log.WOFLogger) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		addresses, err := GetAddressList(req)

		if err != nil {
			gohttp.Error(rsp, err.Error(), err.Code)
			return
		}

		t1 := time.Now()

		defer func() {
			t2 := time.Since(t1)
			logger.Status("parse '%s' %v", addresses, t2)
		}()

		parsed_addresses := make([][]postal.ParsedComponent, len(addresses))

		for i,a := range addresses {
			if a == "" {
				parsed_addresses[i] = make([]postal.ParsedComponent, 0)
			} else {
				parsed_addresses[i] = postal.ParseAddress(a)
			}
		}

		formatted_addresses := make([]map[string][]string, len(parsed_addresses))

		for i,p := range parsed_addresses {
			if len(p) == 0 {
				formatted_addresses[i] = make(map[string][]string)
			} else {
				formatted_addresses[i] = FormatParsed(p)
			}
		}

		WriteResponse(rsp, formatted_addresses)
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}

func FormatParsed(parsed []postal.ParsedComponent) map[string][]string {

	rsp := make(map[string][]string)

	for _, component := range parsed {

		key := component.Label
		value := component.Value

		possible, ok := rsp[key]

		if !ok {
			possible = make([]string, 0)
		}

		possible = append(possible, value)
		rsp[key] = possible
	}

	return rsp
}
