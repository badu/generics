package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type GenericCallbackWrapper[R any] struct {
	Callback func() R
}

func NewCallbackWrapper[R any](callback func() R) *GenericCallbackWrapper[R] {
	return &GenericCallbackWrapper[R]{Callback: callback}
}

type GenericSlice[T comparable] []T

func GenericSliceFromSlice[T comparable](source []T) *GenericSlice[T] {
	result := GenericSlice[T](source)
	return &result
}

func (s *GenericSlice[T]) ToArray() []T {
	return Copy(*s)
}

func (s *GenericSlice[E]) RemoveItem(input ...E) *GenericSlice[E] {
	if len(input) > 0 {
		result := GenericSlice[E](Sub(*s, input))
		return &result
	}

	return s
}

func (s *GenericSlice[E]) Append(item ...E) *GenericSlice[E] {
	return s.Concat(item)
}

func (s *GenericSlice[E]) Len() int {
	return len(*s)
}

func (s *GenericSlice[E]) Concat(slices ...[]E) *GenericSlice[E] {
	if len(slices) == 0 {
		return s
	}

	return GenericSliceFromSlice(Add(s.ToArray(), slices...))
}

const (
	DefaultTimeoutMillisecond = 30 * time.Second // default timeout in milliseconds
)

type Interceptor func(req *http.Request) error

type ClientWMiddlewares struct {
	currentTransport   http.RoundTripper
	previousTransport  http.RoundTripper
	client             *http.Client
	middleWares        GenericSlice[*Interceptor]
	TimeoutMillisecond int64
}

type ReqRespErr struct {
	Req  *http.Request
	Resp *http.Response
	Err  error
}

func NewClient() *ClientWMiddlewares {
	return NewClientWithMiddlewares(&http.Client{})
}

func NewClientWithMiddlewares(client *http.Client, mwares ...*Interceptor) *ClientWMiddlewares {
	result := &ClientWMiddlewares{
		client:      client,
		middleWares: GenericSlice[*Interceptor](mwares),
	}
	result.SetHTTPClient(client)
	return result
}

func (s *ClientWMiddlewares) ForgetMiddlewares(mwares ...*Interceptor) {
	for _, interceptor := range mwares {
		s.middleWares = *s.middleWares.RemoveItem(interceptor)
	}
}

func (s *ClientWMiddlewares) AddMiddlewares(mwares ...*Interceptor) {
	for _, interceptor := range mwares {
		s.middleWares = *s.middleWares.Append(interceptor)
	}
}

func (s *ClientWMiddlewares) ClearMiddlewares() {
	s.middleWares = GenericSlice[*Interceptor]{}
}

func (s *ClientWMiddlewares) HTTPClient() *http.Client {
	return s.client
}

func (s *ClientWMiddlewares) SetHTTPClient(client *http.Client) {
	if client.Transport == nil { // fix non-initialized
		client.Transport = http.DefaultTransport
	}

	if client.Transport != s.previousTransport { // avoid setting up again next time
		s.currentTransport = client.Transport  // remember previous one
		client.Transport = s                   // setup for self
		s.previousTransport = client.Transport // avoid setting up again next time
	}

	s.client = client
}

// RoundTrip interface implementation
func (s *ClientWMiddlewares) RoundTrip(req *http.Request) (*http.Response, error) {
	return s.withMiddlewares(req, 0)
}

func (s *ClientWMiddlewares) withMiddlewares(req *http.Request, index int) (*http.Response, error) {
	if index >= s.middleWares.Len() && s.currentTransport != nil {
		return s.currentTransport.RoundTrip(req)
	}

	err := (*s.middleWares[index])(req)
	if err != nil {
		return nil, err
	}

	return s.withMiddlewares(req, index+1)
}

func (s *ClientWMiddlewares) ContextTimeout() (context.Context, context.CancelFunc) {
	if s.TimeoutMillisecond > 0 {
		return context.WithTimeout(context.Background(), time.Duration(s.TimeoutMillisecond))
	}

	return context.WithTimeout(context.Background(), DefaultTimeoutMillisecond)
}

// Get HTTP Method Get
func (s *ClientWMiddlewares) Get(targetURL string) *ReqRespErr {
	ctx, cancel := s.ContextTimeout()
	defer cancel()

	return s.DoNewRequest(ctx, nil, http.MethodGet, targetURL)
}

// Head HTTP Method Head
func (s *ClientWMiddlewares) Head(targetURL string) *ReqRespErr {
	ctx, cancel := s.ContextTimeout()
	defer cancel()

	return s.DoNewRequest(ctx, nil, http.MethodHead, targetURL)
}

// Options HTTP Method Options
func (s *ClientWMiddlewares) Options(targetURL string) *ReqRespErr {
	ctx, cancel := s.ContextTimeout()
	defer cancel()

	return s.DoNewRequest(ctx, nil, http.MethodOptions, targetURL)
}

// Delete HTTP Method Delete
func (s *ClientWMiddlewares) Delete(targetURL string) *ReqRespErr {
	ctx, cancel := s.ContextTimeout()
	defer cancel()

	return s.DoNewRequest(ctx, nil, http.MethodDelete, targetURL)
}

// Post HTTP Method Post
func (s *ClientWMiddlewares) Post(targetURL, contentType string, body io.Reader) *ReqRespErr {
	ctx, cancel := s.ContextTimeout()
	defer cancel()

	return s.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPost, targetURL, body, contentType)
}

// Put HTTP Method Put
func (s *ClientWMiddlewares) Put(targetURL, contentType string, body io.Reader) *ReqRespErr {
	ctx, cancel := s.ContextTimeout()
	defer cancel()

	return s.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPut, targetURL, body, contentType)
}

// Patch HTTP Method Patch
func (s *ClientWMiddlewares) Patch(targetURL, contentType string, body io.Reader) *ReqRespErr {
	ctx, cancel := s.ContextTimeout()
	defer cancel()

	return s.DoNewRequestWithBodyOptions(ctx, nil, http.MethodPatch, targetURL, body, contentType)
}

// DoNewRequest
func (s *ClientWMiddlewares) DoNewRequest(ctx context.Context, header http.Header, method string, targetURL string) *ReqRespErr {
	request, newRequestErr := http.NewRequestWithContext(ctx, method, targetURL, nil)
	if newRequestErr != nil {
		return &ReqRespErr{
			Req: request,
			Err: newRequestErr,
		}
	}

	if header != nil {
		request.Header = header
	}

	return s.DoRequest(request)
}

// DoNewRequestWithBodyOptions
func (s *ClientWMiddlewares) DoNewRequestWithBodyOptions(ctx context.Context, header http.Header, method string, targetURL string, body io.Reader, contentType string) *ReqRespErr {
	request, newRequestErr := http.NewRequestWithContext(ctx, method, targetURL, body)
	if newRequestErr != nil {
		return &ReqRespErr{
			Req: request,
			Err: newRequestErr,
		}
	}

	if header != nil {
		request.Header = header
	}

	if len(contentType) > 0 {
		request.Header.Add("Content-Type", contentType)
	}

	return s.DoRequest(request)
}

// DoRequest
func (s *ClientWMiddlewares) DoRequest(req *http.Request) *ReqRespErr {
	response, err := s.client.Do(req)

	return &ReqRespErr{
		Req:  req,
		Resp: response,
		Err:  err,
	}
}

// SimpleAPI

// PathParam Path params for API usages
type PathParam map[string]any

// MultipartForm Path params for API usages
type MultipartForm struct {
	Value map[string][]string
	// File The absolute paths of files
	File map[string][]string
}

// APIResponse Response with Error & Type
type APIResponse[R any] struct {
	ReqRespErr
	Response *R
}

// WithoutBody API without request body options
type WithoutBody[R any] func(pathParam PathParam, target *R) *GenericCallbackWrapper[*APIResponse[R]]

// WithBody API with request body options
type WithBody[T any, R any] func(pathParam PathParam, body T, target *R) *GenericCallbackWrapper[*APIResponse[R]]

// WithMultipart API with request body options
type WithMultipart[R any] func(pathParam PathParam, body *MultipartForm, target *R) *GenericCallbackWrapper[*APIResponse[R]]

// BodySerializer Serialize the body (for put/post/patch etc)
type BodySerializer func(body any) (io.Reader, error)

// BodyDeserializer Deserialize the body (for response)
type BodyDeserializer func(body []byte, target any) (any, error)

// MultipartSerializer Serialize the multipart body (for put/post/patch etc)
type MultipartSerializer func(body *MultipartForm) (io.Reader, string, error)

// JSONBodyDeserializer Default JSON Body deserializer
func JSONBodyDeserializer(body []byte, target any) (any, error) {
	err := json.Unmarshal(body, target)
	return target, err
}

// JSONBodySerializer Default JSON Body serializer
func JSONBodySerializer(body any) (io.Reader, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(jsonBytes), err
}

// GeneralMultipartSerializer Default Multipart Body serializer
func GeneralMultipartSerializer(form *MultipartForm) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for fieldName, values := range form.Value {
		for _, value := range values {
			writeFieldErr := writer.WriteField(fieldName, value)
			if writeFieldErr != nil {
				return nil, "", writeFieldErr
			}
		}
	}

	var fileCloseErr error
	for fieldName, filePaths := range form.File {
		for _, filePath := range filePaths {
			part, createFormFileErr := writer.CreateFormFile(fieldName, filepath.Base(filePath))
			if createFormFileErr != nil {
				return nil, "", createFormFileErr
			}
			file, openFileErr := os.Open(filePath)
			if openFileErr != nil {
				return nil, "", openFileErr
			}

			defer func(file *os.File) {
				err := file.Close()
				if err != nil {
					fileCloseErr = err
				}
			}(file)

			_, err := io.Copy(part, file)
			if err != nil {
				return nil, "", err
			}
		}
	}

	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	if fileCloseErr != nil {
		return nil, "", fileCloseErr
	}

	return body, writer.FormDataContentType(), nil
}

type GenericHTTPClient struct {
	genericClient     *ClientWMiddlewares
	DefaultHeader     http.Header
	Serializer        MultipartSerializer
	SerializerForJSON BodySerializer
	Deserializer      BodyDeserializer
	BaseURL           string
}

func (w *GenericHTTPClient) replacePathParams(relativeURL string, pathParam PathParam) string {
	finalURL := relativeURL
	for k, v := range pathParam {
		finalURL = strings.ReplaceAll(relativeURL, fmt.Sprintf("{%s}", k), fmt.Sprintf("%v", v))
	}
	return w.BaseURL + "/" + finalURL
}

// NewGenericHTTPClient New a NewGenericHTTPClient instance
func NewGenericHTTPClient(baseURL string) *GenericHTTPClient {
	return NewGenericHTTPClientWithMiddlewares(baseURL, NewClient())
}

// NewGenericHTTPClientWithMiddlewares a GenericHTTPClient instance with a SimpleHTTP
func NewGenericHTTPClientWithMiddlewares(baseURL string, simpleHTTP *ClientWMiddlewares) *GenericHTTPClient {
	parsedURL, _ := url.Parse(baseURL) // TODO : handle parse error

	return &GenericHTTPClient{
		BaseURL:           parsedURL.String(),
		Serializer:        GeneralMultipartSerializer,
		SerializerForJSON: JSONBodySerializer,
		Deserializer:      JSONBodyDeserializer,
		genericClient:     simpleHTTP,
	}
}

func (w *GenericHTTPClient) GetSimpleHTTP() *ClientWMiddlewares {
	return w.genericClient
}

func GET[R any](client *GenericHTTPClient, url string) WithoutBody[R] {
	return DoNewRequest[R](client, http.MethodGet, url)
}

func DELETE[R any](client *GenericHTTPClient, url string) WithoutBody[R] {
	return DoNewRequest[R](client, http.MethodDelete, url)
}

func POST[T any, R any](client *GenericHTTPClient, url string) WithBody[T, R] {
	return WithJSONSerializer[T, R](client, http.MethodPost, url, "application/json", client.SerializerForJSON)
}

func PUT[T any, R any](client *GenericHTTPClient, url string) WithBody[T, R] {
	return WithJSONSerializer[T, R](client, http.MethodPost, url, "application/json", client.SerializerForJSON)
}

func PATCH[T any, R any](client *GenericHTTPClient, url string) WithBody[T, R] {
	return WithJSONSerializer[T, R](client, http.MethodPost, url, "application/json", client.SerializerForJSON)
}

func POSTMultipart[R any](client *GenericHTTPClient, url string) WithMultipart[R] {
	return WithMultipartSerializer[R](client, http.MethodPost, url, client.Serializer)
}

func PUTMultipart[R any](client *GenericHTTPClient, url string) WithMultipart[R] {
	return WithMultipartSerializer[R](client, http.MethodPost, url, client.Serializer)
}

func PATCHMultipart[R any](client *GenericHTTPClient, url string) WithMultipart[R] {
	return WithMultipartSerializer[R](client, http.MethodPost, url, client.Serializer)
}

func WithJSONSerializer[T any, R any](client *GenericHTTPClient, method string, relativeURL string, contentType string, bodySerializer BodySerializer) WithBody[T, R] {
	return WithBody[T, R](func(pathParam PathParam, body T, rawResp *R) *GenericCallbackWrapper[*APIResponse[R]] {
		return NewCallbackWrapper[*APIResponse[R]](func() *APIResponse[R] {
			var reqBody io.Reader
			if !IsNil(body) {
				var err error
				reqBody, err = bodySerializer(body)
				if err != nil {
					return &APIResponse[R]{ReqRespErr: ReqRespErr{Err: err}}
				}
			}

			ctx, cancel := client.GetSimpleHTTP().ContextTimeout()
			defer cancel()

			resp := client.genericClient.DoNewRequestWithBodyOptions(ctx, client.DefaultHeader.Clone(), method, client.replacePathParams(relativeURL, pathParam), reqBody, contentType)
			if resp.Err != nil {
				return &APIResponse[R]{ReqRespErr: *resp}
			}

			return decodeBody[R](client, &APIResponse[R]{ReqRespErr: *resp}, rawResp)
		})
	})
}

func WithMultipartSerializer[R any](client *GenericHTTPClient, method string, relativeURL string, multipartSerializer MultipartSerializer) WithMultipart[R] {
	return WithMultipart[R](func(pathParam PathParam, body *MultipartForm, rawResp *R) *GenericCallbackWrapper[*APIResponse[R]] {
		return NewCallbackWrapper[*APIResponse[R]](func() *APIResponse[R] {
			var (
				reqBody     io.Reader
				contentType string
			)

			if !IsNil(body) {
				var err error
				reqBody, contentType, err = multipartSerializer(body)
				if err != nil {
					return &APIResponse[R]{
						ReqRespErr: ReqRespErr{
							Err: err,
						},
					}
				}
			}

			ctx, cancel := client.GetSimpleHTTP().ContextTimeout()
			defer cancel()

			resp := client.genericClient.DoNewRequestWithBodyOptions(ctx, client.DefaultHeader.Clone(), method, client.replacePathParams(relativeURL, pathParam), reqBody, contentType)
			if resp.Err != nil {
				return &APIResponse[R]{ReqRespErr: *resp}
			}

			return decodeBody[R](client, &APIResponse[R]{ReqRespErr: *resp}, rawResp)
		})
	})
}

func DoNewRequest[R any](client *GenericHTTPClient, method string, relativeURL string) WithoutBody[R] {
	return WithoutBody[R](func(pathParam PathParam, rawResp *R) *GenericCallbackWrapper[*APIResponse[R]] {
		return NewCallbackWrapper[*APIResponse[R]](func() *APIResponse[R] {
			ctx, cancel := client.GetSimpleHTTP().ContextTimeout()
			defer cancel()

			resp := client.genericClient.DoNewRequest(ctx, client.DefaultHeader.Clone(), method, client.replacePathParams(relativeURL, pathParam))
			if resp.Err != nil {
				return &APIResponse[R]{ReqRespErr: *resp}
			}

			return decodeBody[R](client, &APIResponse[R]{ReqRespErr: *resp}, rawResp)
		})
	})
}

func decodeBody[R any](client *GenericHTTPClient, resp *APIResponse[R], rawResp *R) *APIResponse[R] {
	body, err := ioutil.ReadAll(resp.Resp.Body)
	if err != nil {
		resp.Err = err
		return resp
	}

	var decodedResp any
	decodedResp, resp.Err = client.Deserializer(body, rawResp)
	if resp.Err == nil {
		resp.Response = decodedResp.(*R) // decoding succeeded
	}

	return resp
}

func Add[T any](source []T, targets ...[]T) []T {
	srclen := len(source)
	max := srclen

	for _, slice := range targets {
		if slice == nil {
			continue
		}

		max += len(slice)
	}

	result := make([]T, max)
	for i, item := range source {
		result[i] = item
	}
	idx := srclen

	for _, slice := range targets {
		if slice == nil {
			continue
		}

		target := slice
		targetLen := len(target)
		for j, item := range target {
			result[idx+j] = item
		}

		idx += targetLen
	}

	return result
}

func Sub[T comparable](source, target []T) []T {
	idx := 0
	result := make([]T, len(source))
	targetMap := SliceToMap(true, target...)

	for _, item := range source {
		if _, has := targetMap[item]; !has {
			result[idx] = item
			idx++
		}
	}

	return result[:idx]
}

func Copy[T any](source []T) []T {
	if len(source) > 0 {
		return append(source[:0:0], source...)
	}

	return make([]T, 0)
}

func IsNil(source any) bool {
	if Kind(source) == reflect.Ptr {
		return reflect.ValueOf(source).IsNil()
	}

	return !reflect.ValueOf(source).IsValid()
}

func Kind(source any) reflect.Kind {
	return reflect.ValueOf(source).Kind()
}

func SliceToMap[T comparable, R any](defaultValue R, source ...T) map[T]R {
	result := make(map[T]R)
	for _, key := range source {
		if _, has := result[key]; !has {
			result[key] = defaultValue
		}
	}
	return result
}
