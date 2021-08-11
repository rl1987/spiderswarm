package spiderswarm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPAction(t *testing.T) {
	baseURL := "https://httpbin.org/post"
	method := "POST"

	httpAction := NewHTTPAction(baseURL, method, false)

	assert.NotNil(t, httpAction)
	assert.False(t, httpAction.AbstractAction.ExpectMany)
	assert.Equal(t, httpAction.BaseURL, baseURL)
	assert.Equal(t, httpAction.Method, method)
	assert.Equal(t, len(httpAction.AbstractAction.AllowedInputNames), 5)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[0], HTTPActionInputBaseURL)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[1], HTTPActionInputURLParams)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[2], HTTPActionInputHeaders)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[3], HTTPActionInputCookies)
	assert.Equal(t, httpAction.AbstractAction.AllowedInputNames[4], HTTPActionInputBody)
	assert.Equal(t, len(httpAction.AbstractAction.AllowedOutputNames), 4)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[0], HTTPActionOutputBody)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[1], HTTPActionOutputHeaders)
	assert.Equal(t, httpAction.AbstractAction.AllowedOutputNames[2], HTTPActionOutputStatusCode)

}

func TestNewHTTPActionFromTemplate(t *testing.T) {
	baseURL := "https://github.com/rl1987/spiderswarm"
	method := "HEAD"
	canFail := false

	constructorParams := map[string]interface{}{
		"baseURL": baseURL,
		"method":  method,
		"canFail": canFail,
	}

	actionTempl := &ActionTemplate{
		StructName:        "HTTPAction",
		ConstructorParams: constructorParams,
	}

	action := NewHTTPActionFromTemplate(actionTempl)

	assert.NotNil(t, action)
	assert.Equal(t, baseURL, action.BaseURL)
	assert.Equal(t, method, action.Method)
	assert.Equal(t, canFail, action.CanFail)
}

func TestHTTPActionRunGET(t *testing.T) {
	testHeaders := http.Header{
		"User-Agent": []string{"spiderswarm"},
		"Accept":     []string{"text/plain"},
	}

	testParams := map[string][]string{
		"a": []string{"1"},
		"b": []string{"2"},
	}

	expectedBody := []byte("Test Payload")

	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			u := req.URL
			m, _ := url.ParseQuery(u.RawQuery)
			assert.Equal(t, 2, len(m))
			assert.Equal(t, "1", m["a"][0])
			assert.Equal(t, "2", m["b"][0])

			assert.Equal(t, "spiderswarm", req.Header["User-Agent"][0])
			assert.Equal(t, "text/plain", req.Header["Accept"][0])

			res.Header()["Server"] = []string{"TestServer"}
			res.WriteHeader(200)
			res.Write(expectedBody)
		}))

	defer testServer.Close()

	httpAction := NewHTTPAction(testServer.URL, http.MethodGet, false)

	headersIn := NewDataPipe()
	err := headersIn.Add(testHeaders)
	assert.Nil(t, err)

	err = httpAction.AddInput(HTTPActionInputHeaders, headersIn)
	assert.Nil(t, err)

	paramsIn := NewDataPipe()
	paramsIn.Add(testParams)

	err = httpAction.AddInput(HTTPActionInputURLParams, paramsIn)
	assert.Nil(t, err)

	headersOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputHeaders, headersOut)
	assert.Nil(t, err)

	bodyOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputBody, bodyOut)
	assert.Nil(t, err)

	statusOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputStatusCode, statusOut)
	assert.Nil(t, err)

	err = httpAction.Run()
	assert.Nil(t, err)

	gotBody, ok1 := bodyOut.Remove().([]byte)
	assert.True(t, ok1)
	assert.Equal(t, expectedBody, gotBody)

	gotHeaders, ok2 := headersOut.Remove().(http.Header)
	assert.True(t, ok2)
	assert.True(t, len(gotHeaders) > 1)
	assert.Equal(t, "TestServer", gotHeaders["Server"][0])

	gotStatus, ok3 := statusOut.Remove().(int)
	assert.True(t, ok3)
	assert.Equal(t, 200, gotStatus)

}

func TestHTTPActionRunHEAD(t *testing.T) {
	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodHead, req.Method)

			res.WriteHeader(200)
		}))

	defer testServer.Close()

	httpAction := NewHTTPAction(testServer.URL, http.MethodHead, false)

	err := httpAction.Run()
	assert.Nil(t, err)
}

func TestHTTPActionRunPOST(t *testing.T) {
	expectedBody := []byte("Test Payload")

	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodPost, req.Method)

			body, err := io.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.Equal(t, expectedBody, body)

			res.WriteHeader(201)
		}))

	defer testServer.Close()

	bodyIn := NewDataPipe()
	bodyIn.Add(expectedBody)

	httpAction := NewHTTPAction(testServer.URL, http.MethodPost, false)

	err := httpAction.AddInput(HTTPActionInputBody, bodyIn)
	assert.Nil(t, err)

	err = httpAction.Run()
	assert.Nil(t, err)
}

func TestHTTPActionRunCookies(t *testing.T) {
	cookieName := "SessionID"
	cookieValue1 := "DFB32DF6-ABB5-4877-9B6A-D7F8D9791880"
	cookieValue2 := "2FED9EDE-E3DD-40EA-8F53-A9FBD8959F49"

	testServer := httptest.NewServer(
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, http.MethodGet, req.Method)

			foundExpectedCookie := false

			for _, cookie := range req.Cookies() {
				if cookie.Name == cookieName && cookie.Value == cookieValue1 {
					foundExpectedCookie = true
					break
				}
			}

			assert.True(t, foundExpectedCookie)

			newCookie := &http.Cookie{
				Name:  cookieName,
				Value: cookieValue2,
			}

			http.SetCookie(res, newCookie)

			res.WriteHeader(200)
		}))

	defer testServer.Close()

	cookieDict := map[string]string{cookieName: cookieValue1}

	cookiesIn := NewDataPipe()
	cookiesIn.Add(cookieDict)

	httpAction := NewHTTPAction(testServer.URL, http.MethodGet, false)
	err := httpAction.AddInput(HTTPActionInputCookies, cookiesIn)
	assert.Nil(t, err)

	cookiesOut := NewDataPipe()
	err = httpAction.AddOutput(HTTPActionOutputCookies, cookiesOut)
	assert.Nil(t, err)

	err = httpAction.Run()
	assert.Nil(t, err)

	var gotCookies map[string]string
	var ok bool

	gotCookies, ok = cookiesOut.Remove().(map[string]string)
	assert.True(t, ok)

	foundExpectedCookie := false

	for name, value := range gotCookies {
		if name == cookieName && value == cookieValue2 {
			foundExpectedCookie = true
		}
	}

	assert.True(t, foundExpectedCookie)
}

func TestAddInput(t *testing.T) {
	baseURL := "https://httpbin.org/post"
	method := "POST"

	httpAction := NewHTTPAction(baseURL, method, false)

	dp := NewDataPipe()

	err := httpAction.AddInput("bad_name", dp)
	assert.NotNil(t, err)

	err = httpAction.AddInput(HTTPActionInputURLParams, dp)
	assert.Nil(t, err)
	assert.Equal(t, httpAction.AbstractAction.Inputs[HTTPActionInputURLParams], dp)

}

func TestAddOutput(t *testing.T) {
	baseURL := "https://httpbin.org/post"
	method := "POST"

	httpAction := NewHTTPAction(baseURL, method, false)

	dp := NewDataPipe()

	err := httpAction.AddOutput("bad_name", dp)
	assert.NotNil(t, err)

	err = httpAction.AddOutput(HTTPActionOutputBody, dp)
	assert.Nil(t, err)
	assert.Equal(t, dp, httpAction.AbstractAction.Outputs[HTTPActionOutputBody][0])
}