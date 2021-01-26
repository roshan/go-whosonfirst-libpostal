package http

import (
	"encoding/json"
	"github.com/whosonfirst/go-sanitize"
	"io/ioutil"
	gohttp "net/http"
	"strings"
)

func GetAddress(r *gohttp.Request) (string, *Error) {

	_, err := IsValidMethod(r, []string{"GET"})

	if err != nil {
		return "", err
	}

	query, err := ParseRequest(r)

	if err != nil {
		return "", err
	}

	q, err := ParseQuery(query)

	if err != nil {
		return "", err
	}

	return q, nil
}

func GetAddressList(r *gohttp.Request) ([]string, *Error) {

	_, err := IsValidMethod(r, []string{"POST"})

	if err != nil {
		return make([]string, 0), err
	}

	queries, err := ParseMultiRequest(r)

	if err != nil {
		return make([]string, 0), err
	}

	var addresses = make([]string, len(queries))
	for i, q := range queries {
		address, err := ParseQuery(&q)
		if err != nil {
			addresses[i] = ""
		} else {
			addresses[i] = address
		}
	}

	return addresses, nil
}

func ParseMultiRequest(r *gohttp.Request) ([]Query, *Error) {
	q := make([]Query, 0)

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, NewError(400, err.Error())
	}

	err = json.Unmarshal(bodyBytes, &q)
	if err != nil {
		return nil, NewError(400, err.Error())
	}

	return q, nil
}

func ParseRequest(r *gohttp.Request) (*Query, *Error) {

	query := r.URL.Query()

	address := strings.Trim(query.Get("address"), " ")

	if address == "" {
		return nil, NewError(gohttp.StatusBadRequest, "E_INSUFFICIENT_QUERY")
	}

	q := Query{
		Address: address,
	}

	return &q, nil
}

func ParseQuery(query *Query) (string, *Error) {

	opts := sanitize.DefaultOptions()
	q, err := sanitize.SanitizeString(query.Address, opts)

	if err != nil {
		return "", NewError(gohttp.StatusBadRequest, err.Error())
	}

	if q == "" {
		return "", NewError(gohttp.StatusBadRequest, "E_INVALID_QUERY")
	}

	return q, nil
}

func IsValidMethod(r *gohttp.Request, allowed []string) (bool, *Error) {

	method := r.Method
	ok := false

	for _, this := range allowed {

		if this == method {
			ok = true
			break
		}
	}

	if !ok {
		return false, NewError(gohttp.StatusMethodNotAllowed, "")
	}

	return true, nil
}

func WriteResponse(w gohttp.ResponseWriter, rsp interface{}) {

	rsp_encoded, err := json.Marshal(rsp)

	if err != nil {
		gohttp.Error(w, err.Error(), gohttp.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(rsp_encoded)
}
