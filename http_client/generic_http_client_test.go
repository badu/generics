package http_client

import (
	"bytes"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	ContentTypeHeader = "Content-Type"
	JSONContentType   = "application/json"
	ServerPath        = "/data"
)

var payload = `
{
	"data": [
	  {
	    "userId": 1,
	    "id": 1,
	    "title": "Gaudeamus igitur",
	    "body": "Mi-a intrat o musca-n c*r"
	  },
	  {
	    "userId": 1,
	    "id": 2,
	    "title": "Iuvenes dum sumus",
	    "body": "Si te rog pe dumneata, sa bagi nasul dupa ea"
	  }
	]
}
			`

type Data struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	ID      int    `json:"id"`
	OwnerID int    `json:"userId"`
}

type ListResponse struct {
	Data []Data `json:"data"`
}

func TestJSONGenericHTTPClient(t *testing.T) {
	var (
		url         string
		req         *http.Request
		body        []byte
		contentType string
	)

	serverHandler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		body, _ = ioutil.ReadAll(req.Body)
		_, err := writer.Write([]byte(payload))
		assert.NoError(t, err)
	})

	server := httptest.NewServer(serverHandler)
	defer server.Close()

	var response *ReqRespErr

	client := NewClient()

	testInterceptors := Interceptor(func(request *http.Request) error {
		url = request.URL.Path
		req = request
		contentType = req.Header.Get(ContentTypeHeader)
		return nil
	})
	client.AddMiddlewares(&testInterceptors)

	response = client.Options(server.URL + ServerPath)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, ServerPath, url)

	response = client.Head(server.URL + ServerPath)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, ServerPath, url)

	response = client.Get(server.URL + ServerPath)
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, ServerPath, url)

	response = client.Delete(server.URL + ServerPath + "/1")
	assert.Equal(t, ServerPath+"/1", url)
	assert.Equal(t, nil, response.Err)

	response = client.Post(server.URL+ServerPath, JSONContentType, bytes.NewReader([]byte(`{"userId":0,"id":5,"title":"Пу́тин — хуйло́","body":""}`)))
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, JSONContentType, contentType)
	assert.Equal(t, `{"userId":0,"id":5,"title":"Пу́тин — хуйло́","body":""}`, string(body))

	contentType = ""
	response = client.Put(server.URL+ServerPath, JSONContentType, bytes.NewReader([]byte(`{"userId":0,"id":4,"title":"Пу́тін — хуйло́","body":""}`)))
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, JSONContentType, contentType)
	assert.Equal(t, `{"userId":0,"id":4,"title":"Пу́тін — хуйло́","body":""}`, string(body))

	contentType = ""
	response = client.Patch(server.URL+ServerPath, JSONContentType, bytes.NewReader([]byte(`{"userId":0,"id":3,"title":"Пу́цін хуйло́","body":""}`)))
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, JSONContentType, contentType)
	assert.Equal(t, `{"userId":0,"id":3,"title":"Пу́цін хуйло́","body":""}`, string(body))

	client.ForgetMiddlewares(&testInterceptors)
	contentType = ""
	response = client.Patch(server.URL+ServerPath, JSONContentType, bytes.NewReader([]byte(`{"userId":0,"id":3,"title":"cc","body":""}`)))
	assert.Equal(t, "", contentType)

	api := NewGenericHTTPClient(server.URL)
	api.GetSimpleHTTP().AddMiddlewares(&testInterceptors)

	var resp *APIResponse[ListResponse]

	postsGet := GET[ListResponse](api, ServerPath)
	resp = postsGet(nil, &ListResponse{}).Callback()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, "/"+ServerPath, url)
	assert.Equal(t, 2, len(resp.Response.Data))

	postsGetOne := GET[ListResponse](api, ServerPath+"/{id}")
	resp = postsGetOne(PathParam{"id": 1}, &ListResponse{}).Callback()
	assert.Equal(t, "/"+ServerPath+"/1", url)
	assert.Equal(t, nil, response.Err)

	postsDeleteOne := DELETE[ListResponse](api, ServerPath+"/{id}")
	resp = postsDeleteOne(PathParam{"id": 1}, &ListResponse{}).Callback()
	assert.Equal(t, "/"+ServerPath+"/1", url)
	assert.Equal(t, nil, response.Err)

	contentType = ""
	postsPost := POST[Data, ListResponse](api, ServerPath)
	resp = postsPost(nil, Data{ID: 5, Title: "Пу́цін хуйло́", Body: "Пу́цін хуйло́"}, &ListResponse{}).Callback()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, JSONContentType, contentType)
	assert.Equal(t, `{"id":5,"userId":0,"title":"Пу́цін хуйло́","body":"Пу́цін хуйло́"}`, string(body))

	contentType = ""
	postsPut := PUT[Data, ListResponse](api, ServerPath)
	resp = postsPut(nil, Data{ID: 4, Title: "Пу́цін хуйло́", Body: "Пу́цін хуйло́"}, &ListResponse{}).Callback()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, JSONContentType, contentType)
	assert.Equal(t, `{"id":4,"userId":0,"title":"Пу́цін хуйло́","body":"Пу́цін хуйло́"}`, string(body))

	contentType = ""
	postsPatch := PATCH[Data, ListResponse](api, ServerPath)
	resp = postsPatch(nil, Data{ID: 3, Title: "Пу́цін хуйло́", Body: "Пу́цін хуйло́"}, &ListResponse{}).Callback()
	assert.Equal(t, nil, response.Err)
	assert.Equal(t, JSONContentType, contentType)
	assert.Equal(t, `{"id":3,"userId":0,"title":"Пу́цін хуйло́","body":"Пу́цін хуйло́"}`, string(body))

	api.GetSimpleHTTP().ClearMiddlewares()
	contentType = ""
	resp = postsPatch(nil, Data{ID: 3, Title: "Пу́цін хуйло́"}, &ListResponse{}).Callback()
	assert.Equal(t, "", contentType)
}

func TestMultipartGenericHTTPClient(t *testing.T) {
	var (
		req         *http.Request
		mpartReader *multipart.Reader
		params      map[string]string
		form        *multipart.Form
		sentValues  map[string][]string
		requestBody []byte
		contentType string
	)

	serverHandler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		requestBody, _ = ioutil.ReadAll(req.Body)
		_, err := writer.Write([]byte(`{}`))
		assert.NoError(t, err)
	})

	server := httptest.NewServer(serverHandler)
	defer server.Close()

	var resp *APIResponse[ListResponse]

	interceptorForTest := Interceptor(func(request *http.Request) error {
		req = request
		contentType = req.Header.Get(ContentTypeHeader)
		return nil
	})

	client := NewGenericHTTPClient(server.URL)
	client.GetSimpleHTTP().AddMiddlewares(&interceptorForTest)

	whereNow, _ := os.Getwd()
	filePath := path.Join(whereNow, "generic_http_client_test.go")

	contentType = ""
	postsPost := POSTMultipart[ListResponse](client, ServerPath)
	sentValues = map[string][]string{"userId": {"badu"}, "id": {"5"}, "title": {"Пу́цін хуйло́"}, "body": {"Пу́цін хуйло́"}}
	sentFiles := map[string][]string{"file": {filePath}}
	resp = postsPost(nil, &MultipartForm{Value: sentValues, File: sentFiles}, &ListResponse{}).Callback()

	assert.Equal(t, nil, resp.Err)

	_, params, _ = mime.ParseMediaType(contentType)
	mpartReader = multipart.NewReader(bytes.NewReader(requestBody), params["boundary"])
	form, _ = mpartReader.ReadForm(1024)

	assert.Equal(t, sentValues, form.Value)
	assert.Equal(t, 1, len(form.File["file"]))

	contentType = ""
	postsPut := PUTMultipart[ListResponse](client, ServerPath)
	sentValues = map[string][]string{"userId": {"badu"}, "id": {"4"}, "title": {"Пу́цін хуйло́"}, "body": {"Пу́цін хуйло́"}}
	resp = postsPut(nil, &MultipartForm{Value: sentValues}, &ListResponse{}).Callback()

	assert.Equal(t, nil, resp.Err)

	_, params, _ = mime.ParseMediaType(contentType)
	mpartReader = multipart.NewReader(bytes.NewReader(requestBody), params["boundary"])
	form, _ = mpartReader.ReadForm(1024)

	assert.Equal(t, sentValues, form.Value)
	assert.Equal(t, 0, len(form.File["file"]))

	contentType = ""
	postsPatch := PATCHMultipart[ListResponse](client, ServerPath)
	sentValues = map[string][]string{"userId": {"badu"}, "id": {"3"}, "title": {"Пу́цін хуйло́"}, "body": {"Пу́цін хуйло́"}}
	resp = postsPatch(nil, &MultipartForm{Value: sentValues}, &ListResponse{}).Callback()

	assert.Equal(t, nil, resp.Err)

	_, params, _ = mime.ParseMediaType(contentType)
	mpartReader = multipart.NewReader(bytes.NewReader(requestBody), params["boundary"])
	form, _ = mpartReader.ReadForm(1024)

	assert.Equal(t, sentValues, form.Value)
	assert.Equal(t, 0, len(form.File["file"]))
}
