package restapi

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	resty "gopkg.in/resty.v1"
)

// ErrEntityNotFound error message which is returned when the entity cannot be found at the server
var ErrEntityNotFound = errors.New("failed to get resource from Instana API. 404 - Resource not found")

const contentTypeHeader = "Content-Type"
const encodingApplicationJSON = "application/json; charset=utf-8"

// RestClient interface to access REST resources of the Instana API
type RestClient interface {
	Get(resourcePath string) ([]byte, error)
	GetOne(id string, resourcePath string) ([]byte, error)
	Post(data InstanaDataObject, resourcePath string) ([]byte, error)
	PostWithID(data InstanaDataObject, resourcePath string) ([]byte, error)
	Put(data InstanaDataObject, resourcePath string) ([]byte, error)
	Delete(resourceID string, resourceBasePath string) error
	GetByQuery(resourcePath string, queryParams map[string]string) ([]byte, error)
	PostByQuery(resourcePath string, queryParams map[string]string) ([]byte, error)
	PutByQuery(resourcePath string, is string, queryParams map[string]string) ([]byte, error)
}

type apiRequest struct {
	method          string
	url             string
	request         resty.Request
	responseChannel chan *apiResponse
	ctx             context.Context
}

type apiResponse struct {
	data []byte
	err  error
}

// NewClient creates a new instance of the Instana REST API client
func NewClient(apiToken string, host string, skipTlsVerification bool) RestClient {
	restyClient := resty.New()
	if skipTlsVerification {
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //nolint:gosec
	}

	throttleRate := time.Second / 5 //5 write requests per second
	throttledRequests := make(chan *apiRequest, 1000)
	client := &restClientImpl{
		apiToken:          apiToken,
		host:              host,
		restyClient:       restyClient,
		throttledRequests: throttledRequests,
		throttleRate:      throttleRate,
	}

	go client.processThrottledRequests()
	return client
}

type restClientImpl struct {
	apiToken          string
	host              string
	restyClient       *resty.Client
	throttledRequests chan *apiRequest
	throttleRate      time.Duration
}

var emptyResponse = make([]byte, 0)

// Get request data via HTTP GET for the given resourcePath
func (client *restClientImpl) Get(resourcePath string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	req := client.createRequest()
	return client.executeRequest(resty.MethodGet, url, req)
}

// Get request data via HTTP GET for the given resourcePath and query parameters
func (client *restClientImpl) GetByQuery(resourcePath string, queryParams map[string]string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	req := client.createRequest()
	client.appendQueryParameters(req, queryParams)
	return client.executeRequest(resty.MethodGet, url, req)
}

// GetOne request the resource with the given ID
func (client *restClientImpl) GetOne(id string, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, id)
	req := client.createRequest()
	return client.executeRequest(resty.MethodGet, url, req)
}

// Post executes a HTTP PUT request to create or update the given resource
func (client *restClientImpl) Post(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	req := client.createRequest().SetHeader(contentTypeHeader, encodingApplicationJSON).SetBody(data)
	return client.executeRequestWithThrottling(resty.MethodPost, url, req)
}

// PostWithID executes a HTTP PUT request to create or update the given resource using the ID from the InstanaDataObject in the resource path
func (client *restClientImpl) PostWithID(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, data.GetIDForResourcePath())
	req := client.createRequest().SetHeader(contentTypeHeader, encodingApplicationJSON).SetBody(data)
	return client.executeRequestWithThrottling(resty.MethodPost, url, req)
}

// Put executes a HTTP PUT request to create or update the given resource
func (client *restClientImpl) Put(data InstanaDataObject, resourcePath string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, data.GetIDForResourcePath())
	req := client.createRequest().SetHeader(contentTypeHeader, encodingApplicationJSON).SetBody(data)
	return client.executeRequestWithThrottling(resty.MethodPut, url, req)
}

// Delete executes a HTTP DELETE request to delete the resource with the given ID
func (client *restClientImpl) Delete(resourceID string, resourceBasePath string) error {
	url := client.buildResourceURL(resourceBasePath, resourceID)
	req := client.createRequest()
	_, err := client.executeRequestWithThrottling(resty.MethodDelete, url, req)
	return err
}

// PostByQuery executes a HTTP POST request to create the resource by providing the data a query parameters
func (client *restClientImpl) PostByQuery(resourcePath string, queryParams map[string]string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	req := client.createRequest()
	client.appendQueryParameters(req, queryParams)
	return client.executeRequest(resty.MethodPost, url, req)
}

// PutByQuery executes a HTTP PUT request to update the resource with the given ID by providing the data a query parameters
func (client *restClientImpl) PutByQuery(resourcePath string, id string, queryParams map[string]string) ([]byte, error) {
	url := client.buildResourceURL(resourcePath, id)
	req := client.createRequest()
	client.appendQueryParameters(req, queryParams)
	return client.executeRequest(resty.MethodPut, url, req)
}

func (client *restClientImpl) createRequest() *resty.Request {
	//get path to root directory from runtime executor
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Join(filepath.Dir(b), "../..")

	//open CHANGELOG.md from root directory (only file storing updated version number)
	terraformProviderVersion := ""
	file, err := os.Open(basepath + "/CHANGELOG.md")
	if err != nil {
		log.Println("Error: couldn't open file", basepath+"/CHANGELOG.md", err)
		return client.restyClient.R().SetHeader("Accept", "application/json").SetHeader("Authorization", fmt.Sprintf("apiToken %s", client.apiToken)).SetHeader("user-agent", terraformProviderVersion)
	}
	defer file.Close()

	//read lines from CHANGELOG.md until first line starting with ##
	scanner := bufio.NewScanner(file)
	for !strings.Contains(scanner.Text(), "##") {
		scanner.Scan()
	}
	if err := scanner.Err(); err != nil {
		log.Println("Error: Couldn't scan file", err)
	} else {
		//Read version number from first line  of CHANGELOG.md starting with ##
		terraformProviderVersion = scanner.Text()
		terraformProviderVersion = strings.Split(terraformProviderVersion, "]")[0]
		terraformProviderVersion = strings.Split(terraformProviderVersion, "[")[1]
	}

	//return client with headers needed for every call
	return client.restyClient.R().SetHeader("Accept", "application/json").SetHeader("Authorization", fmt.Sprintf("apiToken %s", client.apiToken)).SetHeader("user-agent", terraformProviderVersion)
}

func (client *restClientImpl) executeRequestWithThrottling(method string, url string, req *resty.Request) ([]byte, error) {
	responseChannel := make(chan *apiResponse)
	ctx, cancel := context.WithCancel(context.Background())
	defer close(responseChannel)
	defer cancel()

	client.throttledRequests <- &apiRequest{
		method:          method,
		url:             url,
		request:         *req,
		ctx:             ctx,
		responseChannel: responseChannel,
	}

	select {
	case r := <-responseChannel:
		return r.data, r.err
	case <-time.After(30 * time.Second):
		return nil, errors.New("API request timed out")
	}
}

func (client *restClientImpl) processThrottledRequests() {
	throttle := time.NewTicker(client.throttleRate).C
	for req := range client.throttledRequests {
		<-throttle
		go client.handleThrottledAPIRequest(req)
	}
}

func (client *restClientImpl) handleThrottledAPIRequest(req *apiRequest) {
	data, err := client.executeRequest(req.method, req.url, &req.request)
	responseMessage := &apiResponse{
		data: data,
		err:  err,
	}
	select {
	case <-req.ctx.Done():
		return
	default:
		req.responseChannel <- responseMessage
	}
}

func (client *restClientImpl) executeRequest(method string, url string, req *resty.Request) ([]byte, error) {
	log.Printf("[DEBUG] Call %s %s\n", method, url)
	resp, err := req.Execute(method, url)
	// Log request body at INFO level
	if req.Body != nil {
		// Convert request body to JSON string for proper logging
		jsonBytes, err := json.MarshalIndent(req.Body, "", "  ")
		if err != nil {
			log.Printf("[INFO] Request Body (could not marshal to JSON): %v\n", req.Body)
		} else {
			log.Printf("[INFO] Request Body: %s\n", string(jsonBytes))
		}
	}
	if err != nil {
		if resp == nil {
			return emptyResponse, fmt.Errorf("failed to send HTTP %s request to Instana API; %s", method, err)
		}
		return emptyResponse, fmt.Errorf("failed to send HTTP %s request to Instana API; status code = %d; status message = %s; Headers %s, %s", method, resp.StatusCode(), resp.Status(), resp.Header(), err)
	}
	// Log response body at INFO level
	log.Printf("[INFO] Response Status: %d %s\n", resp.StatusCode(), resp.Status())

	// Try to parse response body as JSON for better formatting
	var responseObj interface{}
	responseBody := resp.Body()
	if err := json.Unmarshal(responseBody, &responseObj); err != nil {
		// If not valid JSON, log as is
		log.Printf("[INFO] Response Body: %s\n", string(responseBody))
	} else {
		// If valid JSON, pretty print it
		prettyJSON, err := json.MarshalIndent(responseObj, "", "  ")
		if err != nil {
			log.Printf("[INFO] Response Body: %s\n", string(responseBody))
		} else {
			log.Printf("[INFO] Response Body: %s\n", string(prettyJSON))
		}
	}
	statusCode := resp.StatusCode()
	if statusCode == 404 {
		return emptyResponse, ErrEntityNotFound
	}
	if statusCode < 200 || statusCode >= 300 {
		return emptyResponse, fmt.Errorf("failed to send HTTP %s request to Instana API; status code = %d; status message = %s; Headers %s\nBody: %s", method, statusCode, resp.Status(), resp.Header(), resp.Body())
	}
	return resp.Body(), nil
}

func (client *restClientImpl) appendQueryParameters(req *resty.Request, queryParams map[string]string) {
	for k, v := range queryParams {
		req.QueryParam.Add(k, v)
	}
}

func (client *restClientImpl) buildResourceURL(resourceBasePath string, id string) string {
	pattern := "%s/%s"
	if strings.HasSuffix(resourceBasePath, "/") {
		pattern = "%s%s"
	}
	resourcePath := fmt.Sprintf(pattern, resourceBasePath, id)
	return client.buildURL(resourcePath)
}

func (client *restClientImpl) buildURL(resourcePath string) string {
	return fmt.Sprintf("https://%s%s", client.host, resourcePath)
}
