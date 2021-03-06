package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/keptn/go-utils/pkg/api/models"
)

// APIHandler handles projects
type APIHandler struct {
	BaseURL    string
	AuthToken  string
	AuthHeader string
	HTTPClient *http.Client
	Scheme     string
}

// NewAuthenticatedAPIHandler returns a new APIHandler that authenticates at the api-service endpoint via the provided token
func NewAuthenticatedAPIHandler(baseURL string, authToken string, authHeader string, httpClient *http.Client, scheme string) *APIHandler {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	httpClient.Transport = getClientTransport()

	baseURL = strings.TrimPrefix(baseURL, "http://")
	baseURL = strings.TrimPrefix(baseURL, "https://")
	return &APIHandler{
		BaseURL:    baseURL,
		AuthHeader: authHeader,
		AuthToken:  authToken,
		HTTPClient: httpClient,
		Scheme:     scheme,
	}
}

func (a *APIHandler) getBaseURL() string {
	return a.BaseURL
}

func (a *APIHandler) getAuthToken() string {
	return a.AuthToken
}

func (a *APIHandler) getAuthHeader() string {
	return a.AuthHeader
}

func (a *APIHandler) getHTTPClient() *http.Client {
	return a.HTTPClient
}

// SendEvent sends an event to Keptn
func (a *APIHandler) SendEvent(event models.KeptnContextExtendedCE) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(event)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(a.Scheme+"://"+a.getBaseURL()+"/v1/event", bodyStr, a)
}

// TriggerEvaluation triggers a new evaluation
func (a *APIHandler) TriggerEvaluation(project, stage, service string, evaluation models.Evaluation) (*models.EventContext, *models.Error) {
	bodyStr, err := json.Marshal(evaluation)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	return postWithEventContext(a.Scheme+"://"+a.getBaseURL()+"/v1/project/"+project+"/stage/"+stage+"/service/"+service+"/evaluation", bodyStr, a)
}

// GetEvent returns an event specified by keptnContext and eventType
//
// Deprecated: this function is deprecated and should be replaced with the GetEvents function
func (a *APIHandler) GetEvent(keptnContext string, eventType string) (*models.KeptnContextExtendedCE, *models.Error) {
	return getEvent(a.Scheme+"://"+a.getBaseURL()+"/v1/event?keptnContext="+keptnContext+"&type="+eventType+"&pageSize=10", a)
}

func getEvent(uri string, api APIService) (*models.KeptnContextExtendedCE, *models.Error) {

	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, api)

	resp, err := api.getHTTPClient().Do(req)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		if len(body) > 0 {
			var cloudEvent models.KeptnContextExtendedCE
			err = json.Unmarshal(body, &cloudEvent)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}

			return &cloudEvent, nil
		}

		return nil, nil
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return nil, &respErr
}

// CreateProject creates a new project
func (a *APIHandler) CreateProject(project models.CreateProject) (string, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(a.Scheme+"://"+a.getBaseURL()+"/shipyard-controller/v1/project", bodyStr, a)
}

// UpdateProject updates project
func (a *APIHandler) UpdateProject(project models.CreateProject) (string, *models.Error) {
	bodyStr, err := json.Marshal(project)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return put(a.Scheme+"://"+a.getBaseURL()+"/shipyard-controller/v1/project", bodyStr, a)
}

// DeleteProject deletes a project
func (a *APIHandler) DeleteProject(project models.Project) (*models.DeleteProjectResponse, *models.Error) {
	resp, err := delete(a.Scheme+"://"+a.getBaseURL()+"/shipyard-controller/v1/project/"+project.ProjectName, a)
	if err != nil {
		return nil, err
	}
	deletePrjResponse := &models.DeleteProjectResponse{}
	if err2 := json.Unmarshal([]byte(resp), deletePrjResponse); err2 != nil {
		msg := "Could not decode DeleteProjectResponse: " + err2.Error()
		return nil, &models.Error{
			Message: &msg,
		}
	}
	return deletePrjResponse, nil
}

// CreateService creates a new service
func (a *APIHandler) CreateService(project string, service models.CreateService) (string, *models.Error) {
	bodyStr, err := json.Marshal(service)
	if err != nil {
		return "", buildErrorResponse(err.Error())
	}
	return post(a.Scheme+"://"+a.getBaseURL()+"/shipyard-controller/v1/project/"+project+"/service", bodyStr, a)
}

// DeleteProject deletes a project
func (a *APIHandler) DeleteService(project, service string) (*models.DeleteServiceResponse, *models.Error) {
	resp, err := delete(a.Scheme+"://"+a.getBaseURL()+"/shipyard-controller/v1/project/"+project+"/service/"+service, a)
	if err != nil {
		return nil, err
	}
	deleteSvcResponse := &models.DeleteServiceResponse{}
	if err2 := json.Unmarshal([]byte(resp), deleteSvcResponse); err2 != nil {
		msg := "Could not decode DeleteServiceResponse: " + err2.Error()
		return nil, &models.Error{
			Message: &msg,
		}
	}
	return deleteSvcResponse, nil
}

// GetMetadata retrieve keptn MetaData information
func (a *APIHandler) GetMetadata() (*models.Metadata, *models.Error) {
	//return get(s.Scheme+"://"+s.getBaseURL()+"/v1/metadata", nil, s)

	req, err := http.NewRequest("GET", a.Scheme+"://"+a.getBaseURL()+"/v1/metadata", nil)
	req.Header.Set("Content-Type", "application/json")
	addAuthHeader(req, a)

	resp, err := a.getHTTPClient().Do(req)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		if len(body) > 0 {
			var respMetadata models.Metadata
			err = json.Unmarshal(body, &respMetadata)
			if err != nil {
				return nil, buildErrorResponse(err.Error())
			}

			return &respMetadata, nil
		}

		return nil, nil
	}

	var respErr models.Error
	err = json.Unmarshal(body, &respErr)
	if err != nil {
		return nil, buildErrorResponse(err.Error())
	}

	return nil, &respErr

}
