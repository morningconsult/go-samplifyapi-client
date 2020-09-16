package samplify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// ErrIncorrectEnvironemt ...
var ErrIncorrectEnvironemt = errors.New("one of dev/uat/prod only are allowed")

// ClientOptions to use while creating a new Client
var (
	DevClientOptions = &ClientOptions{
		APIBaseURL: "https://api.dev.pe.dynata.com/sample/v1",
		AuthURL:    "https://api.dev.pe.dynata.com/auth/v1",
		StatusURL:  "https://api.dev.pe.dynata.com/status",
		GatewayURL: "https://api.dev.pe.dynata.com/status/gateway",
	}
	UATClientOptions = &ClientOptions{
		APIBaseURL: "https://api.uat.pe.dynata.com/sample/v1",
		AuthURL:    "https://api.uat.pe.dynata.com/auth/v1",
		StatusURL:  "https://api.uat.pe.dynata.com/status",
		GatewayURL: "https://api.uat.pe.dynata.com/status/gateway",
	}
	ProdClientOptions = &ClientOptions{
		APIBaseURL: "https://api.researchnow.com/sample/v1",
		AuthURL:    "https://api.researchnow.com/auth/v1",
		StatusURL:  "https://api.researchnow.com/status",
		GatewayURL: "https://api.researchnow.com/status/gateway",
	}
)

// ErrSessionExpired ... Returns if both Access and Refresh tokens are expired
var ErrSessionExpired = errors.New("session expired")

const defaulttimeout = 20

// ClientOptions ...
type ClientOptions struct {
	APIBaseURL string `conform:"trim"`
	AuthURL    string `conform:"trim"`
	StatusURL  string `conform:"trim"`
	GatewayURL string `conform:"trim"`
	Timeout    *int
	HTTPClient httpClient
}

// Client is used to make API requests to the Samplify API.
type Client struct {
	Credentials TokenRequest
	Auth        TokenResponse
	Options     *ClientOptions
}

// GetInvoicesSummary ...
func (c *Client) GetInvoicesSummary(options *QueryOptions) (*APIResponse, error) {
	path := fmt.Sprintf("/projects/invoices/summary%s", query2String(options))
	return c.request(http.MethodGet, c.Options.APIBaseURL, path, nil)
}

// GetInvoicesSummaryWithContext ...
func (c *Client) GetInvoicesSummaryWithContext(ctx context.Context, options *QueryOptions) (*APIResponse, error) {
	path := fmt.Sprintf("/projects/invoices/summary%s", query2String(options))
	return c.requestWithContext(ctx, http.MethodGet, c.Options.APIBaseURL, path, nil)
}

// CreateProject ...
func (c *Client) CreateProject(project *CreateProjectCriteria) (*ProjectResponse, error) {
	f := func(res *ProjectResponse) error {
		return c.requestAndParseResponse(http.MethodPost, "/projects", project, res)
	}

	return c.createProject(project, f)
}

// CreateProjectWithContext ...
func (c *Client) CreateProjectWithContext(ctx context.Context, project *CreateProjectCriteria) (*ProjectResponse, error) {
	f := func(res *ProjectResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, "/projects", project, res)
	}

	return c.createProject(project, f)
}

func (c *Client) createProject(
	project *CreateProjectCriteria,
	requestAndParseResponse func(res *ProjectResponse) error,
) (*ProjectResponse, error) {
	if err := Validate(project); err != nil {
		return nil, err
	}

	var res ProjectResponse

	if err := requestAndParseResponse(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateProject ...
func (c *Client) UpdateProject(project *UpdateProjectCriteria) (*ProjectResponse, error) {
	f := func(path string, res *ProjectResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, project, res)
	}

	return c.updateProject(project, f)
}

// UpdateProjectWithContext ...
func (c *Client) UpdateProjectWithContext(ctx context.Context, project *UpdateProjectCriteria) (*ProjectResponse, error) {
	f := func(path string, res *ProjectResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, project, res)
	}

	return c.updateProject(project, f)
}

func (c *Client) updateProject(
	project *UpdateProjectCriteria,
	requestAndParseResponse func(path string, res *ProjectResponse) error,
) (*ProjectResponse, error) {
	if err := Validate(project); err != nil {
		return nil, err
	}

	var res ProjectResponse
	path := fmt.Sprintf("/projects/%s", project.ExtProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// BuyProject ...
func (c *Client) BuyProject(extProjectID string, buy []*BuyProjectCriteria) (*BuyProjectResponse, error) {
	f := func(path string, res *BuyProjectResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, buy, res)
	}

	return c.buyProject(extProjectID, buy, f)
}

// BuyProjectWithContext ...
func (c *Client) BuyProjectWithContext(
	ctx context.Context,
	extProjectID string,
	buy []*BuyProjectCriteria,
) (*BuyProjectResponse, error) {
	f := func(path string, res *BuyProjectResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, buy, res)
	}

	return c.buyProject(extProjectID, buy, f)
}

func (c *Client) buyProject(
	extProjectID string,
	buy []*BuyProjectCriteria,
	requestAndParseResponse func(path string, res *BuyProjectResponse) error,
) (*BuyProjectResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	if err := Validate(buy); err != nil {
		return nil, err
	}

	var res BuyProjectResponse
	path := fmt.Sprintf("/projects/%s/buy", extProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CloseProject ...
func (c *Client) CloseProject(extProjectID string) (*CloseProjectResponse, error) {
	f := func(path string, res *CloseProjectResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, nil, res)
	}

	return c.closeProject(extProjectID, f)
}

// CloseProjectWithContext ...
func (c *Client) CloseProjectWithContext(ctx context.Context, extProjectID string) (*CloseProjectResponse, error) {
	f := func(path string, res *CloseProjectResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, nil, res)
	}

	return c.closeProject(extProjectID, f)
}

func (c *Client) closeProject(
	extProjectID string,
	requestAndParseResponse func(path string, res *CloseProjectResponse) error,
) (*CloseProjectResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res CloseProjectResponse
	path := fmt.Sprintf("/projects/%s/close", extProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAllProjects ...
func (c *Client) GetAllProjects(options *QueryOptions) (*GetAllProjectsResponse, error) {
	f := func(path string, res *GetAllProjectsResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getAllProjects(options, f)
}

// GetAllProjectsWithContext ...
func (c *Client) GetAllProjectsWithContext(ctx context.Context, options *QueryOptions) (*GetAllProjectsResponse, error) {
	f := func(path string, res *GetAllProjectsResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getAllProjects(options, f)
}

func (c *Client) getAllProjects(
	options *QueryOptions,
	requestAndParseResponse func(path string, res *GetAllProjectsResponse) error,
) (*GetAllProjectsResponse, error) {
	var res GetAllProjectsResponse
	path := fmt.Sprintf("/projects%s", query2String(options))

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetProjectBy returns project by id
func (c *Client) GetProjectBy(extProjectID string) (*ProjectResponse, error) {
	f := func(path string, res *ProjectResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getProjectBy(extProjectID, f)
}

// GetProjectByWithContext returns project by id
func (c *Client) GetProjectByWithContext(ctx context.Context, extProjectID string) (*ProjectResponse, error) {
	f := func(path string, res *ProjectResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getProjectBy(extProjectID, f)
}

func (c *Client) getProjectBy(
	extProjectID string,
	requestAndParseResponse func(path string, res *ProjectResponse) error,
) (*ProjectResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res ProjectResponse
	path := fmt.Sprintf("/projects/%s", extProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetProjectReport returns a project's report based on observed data from actual panelists.
func (c *Client) GetProjectReport(extProjectID string) (*ProjectReportResponse, error) {
	f := func(path string, res *ProjectReportResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getProjectReport(extProjectID, f)
}

// GetProjectReportWithContext returns a project's report based on observed data from actual panelists.
func (c *Client) GetProjectReportWithContext(ctx context.Context, extProjectID string) (*ProjectReportResponse, error) {
	f := func(path string, res *ProjectReportResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getProjectReport(extProjectID, f)
}

func (c *Client) getProjectReport(
	extProjectID string,
	requestAndParseResponse func(path string, res *ProjectReportResponse) error,
) (*ProjectReportResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res ProjectReportResponse
	path := fmt.Sprintf("/projects/%s/report", extProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// AddLineItem ...
func (c *Client) AddLineItem(extProjectID string, lineItem *CreateLineItemCriteria) (*LineItemResponse, error) {
	requestAndParseResponse := func(path string, res *LineItemResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, lineItem, res)
	}

	return c.addLineItem(extProjectID, lineItem, requestAndParseResponse)
}

// AddLineItemWithContext ...
func (c *Client) AddLineItemWithContext(
	ctx context.Context,
	extProjectID string,
	lineItem *CreateLineItemCriteria,
) (*LineItemResponse, error) {
	requestAndParseResponse := func(path string, res *LineItemResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, lineItem, res)
	}

	return c.addLineItem(extProjectID, lineItem, requestAndParseResponse)
}

func (c *Client) addLineItem(
	extProjectID string,
	lineItem *CreateLineItemCriteria,
	requestAndParseResponse func(path string, res *LineItemResponse) error,
) (*LineItemResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	if err := Validate(lineItem); err != nil {
		return nil, err
	}

	if err := ValidateSchedule(&lineItem.DaysInField, lineItem.FieldSchedule); err != nil {
		return nil, err
	}

	var res LineItemResponse
	path := fmt.Sprintf("/projects/%s/lineItems", extProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateLineItem ...
func (c *Client) UpdateLineItem(
	extProjectID string,
	extLineItemID string,
	lineItem *UpdateLineItemCriteria,
) (*LineItemResponse, error) {
	f := func(path string, res *LineItemResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, lineItem, res)
	}

	return c.updateLineItem(extProjectID, extLineItemID, lineItem, f)
}

// UpdateLineItemContext ...
func (c *Client) UpdateLineItemContext(
	ctx context.Context,
	extProjectID string,
	extLineItemID string,
	lineItem *UpdateLineItemCriteria,
) (*LineItemResponse, error) {
	f := func(path string, res *LineItemResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, lineItem, res)
	}

	return c.updateLineItem(extProjectID, extLineItemID, lineItem, f)
}

func (c *Client) updateLineItem(
	extProjectID string,
	extLineItemID string,
	lineItem *UpdateLineItemCriteria,
	requestAndParseResponse func(path string, res *LineItemResponse) error,
) (*LineItemResponse, error) {
	if err := ValidateNotEmpty(extProjectID, extLineItemID); err != nil {
		return nil, err
	}

	if err := Validate(lineItem); err != nil {
		return nil, err
	}

	if err := ValidateSchedule(lineItem.DaysInField, lineItem.FieldSchedule); err != nil {
		return nil, err
	}

	var res LineItemResponse
	path := fmt.Sprintf("/projects/%s/lineItems/%s", extProjectID, extLineItemID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateLineItemState ... Changes the state of the line item based on provided action.
func (c *Client) UpdateLineItemState(
	extProjectID string,
	extLineItemID string,
	action Action,
) (*UpdateLineItemStateResponse, error) {
	f := func(path string, res *UpdateLineItemStateResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, nil, res)
	}

	return c.updateLineItemState(extProjectID, extLineItemID, action, f)
}

// UpdateLineItemStateWithContext ... Changes the state of the line item based on provided action.
func (c *Client) UpdateLineItemStateWithContext(
	ctx context.Context,
	extProjectID string,
	extLineItemID string,
	action Action,
) (*UpdateLineItemStateResponse, error) {
	f := func(path string, res *UpdateLineItemStateResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, nil, res)
	}

	return c.updateLineItemState(extProjectID, extLineItemID, action, f)
}

func (c *Client) updateLineItemState(
	extProjectID string,
	extLineItemID string,
	action Action,
	requestAndParseResponse func(path string, res *UpdateLineItemStateResponse) error,
) (*UpdateLineItemStateResponse, error) {
	if err := ValidateNotEmpty(extProjectID, extLineItemID); err != nil {
		return nil, err
	}

	if err := ValidateAction(action); err != nil {
		return nil, err
	}

	var res UpdateLineItemStateResponse
	path := fmt.Sprintf("/projects/%s/lineItems/%s/%s", extProjectID, extLineItemID, action)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// LaunchLineItem utility function to launch a line item
func (c *Client) LaunchLineItem(pid, lid string) (*UpdateLineItemStateResponse, error) {
	return c.UpdateLineItemState(pid, lid, ActionLaunched)
}

// LaunchLineItemWithContext utility function to launch a line item
func (c *Client) LaunchLineItemWithContext(ctx context.Context, pid, lid string) (*UpdateLineItemStateResponse, error) {
	return c.UpdateLineItemStateWithContext(ctx, pid, lid, ActionLaunched)
}

// PauseLineItem utility function to pause a lineitem
func (c *Client) PauseLineItem(pid, lid string) (*UpdateLineItemStateResponse, error) {
	return c.UpdateLineItemState(pid, lid, ActionPaused)
}

// PauseLineItemWithContext utility function to pause a lineitem
func (c *Client) PauseLineItemWithContext(ctx context.Context, pid, lid string) (*UpdateLineItemStateResponse, error) {
	return c.UpdateLineItemStateWithContext(ctx, pid, lid, ActionPaused)
}

// CloseLineItem utility function to close a lineitem
func (c *Client) CloseLineItem(pid, lid string) (*UpdateLineItemStateResponse, error) {
	return c.UpdateLineItemState(pid, lid, ActionClosed)
}

// CloseLineItemWithContext utility function to close a lineitem
func (c *Client) CloseLineItemWithContext(ctx context.Context, pid, lid string) (*UpdateLineItemStateResponse, error) {
	return c.UpdateLineItemStateWithContext(ctx, pid, lid, ActionClosed)
}

// SetQuotaCellStatus ... Changes the state of the line item based on provided action.
func (c *Client) SetQuotaCellStatus(
	extProjectID string,
	extLineItemID string,
	quotaCellID string,
	action Action,
) (*QuotaCellResponse, error) {
	f := func(path string, res *QuotaCellResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, nil, res)
	}

	return c.setQuotaCellStatus(extProjectID, extLineItemID, quotaCellID, action, f)
}

// SetQuotaCellStatusWithContext ... Changes the state of the line item based on provided action.
func (c *Client) SetQuotaCellStatusWithContext(
	ctx context.Context,
	extProjectID string,
	extLineItemID string,
	quotaCellID string,
	action Action,
) (*QuotaCellResponse, error) {
	f := func(path string, res *QuotaCellResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, nil, res)
	}

	return c.setQuotaCellStatus(extProjectID, extLineItemID, quotaCellID, action, f)
}

func (c *Client) setQuotaCellStatus(
	extProjectID string,
	extLineItemID string,
	quotaCellID string,
	action Action,
	requestAndParseResponse func(path string, res *QuotaCellResponse) error,
) (*QuotaCellResponse, error) {
	if err := ValidateNotEmpty(extProjectID, extLineItemID, quotaCellID); err != nil {
		return nil, err
	}

	if err := ValidateAction(action); err != nil {
		return nil, err
	}

	var res QuotaCellResponse
	path := fmt.Sprintf("/projects/%s/lineItems/%s/quotaCells/%s/%s", extProjectID, extLineItemID, quotaCellID, action)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAllLineItems ...
func (c *Client) GetAllLineItems(extProjectID string, options *QueryOptions) (*GetAllLineItemsResponse, error) {
	f := func(path string, res *GetAllLineItemsResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getAllLineItems(extProjectID, options, f)
}

// GetAllLineItemsWithContext ...
func (c *Client) GetAllLineItemsWithContext(
	ctx context.Context,
	extProjectID string,
	options *QueryOptions,
) (*GetAllLineItemsResponse, error) {
	f := func(path string, res *GetAllLineItemsResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getAllLineItems(extProjectID, options, f)
}

func (c *Client) getAllLineItems(
	extProjectID string,
	options *QueryOptions,
	requestAndParseResponse func(path string, res *GetAllLineItemsResponse) error,
) (*GetAllLineItemsResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res GetAllLineItemsResponse
	path := fmt.Sprintf("/projects/%s/lineItems%s", extProjectID, query2String(options))

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetLineItemBy ...
func (c *Client) GetLineItemBy(extProjectID, extLineItemID string) (*LineItemResponse, error) {
	f := func(path string, res *LineItemResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getLineItemBy(extProjectID, extLineItemID, f)
}

// GetLineItemByWithContext ...
func (c *Client) GetLineItemByWithContext(
	ctx context.Context,
	extProjectID string,
	extLineItemID string,
) (*LineItemResponse, error) {
	f := func(path string, res *LineItemResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getLineItemBy(extProjectID, extLineItemID, f)
}

func (c *Client) getLineItemBy(
	extProjectID string,
	extLineItemID string,
	requestAndParseResponse func(path string, res *LineItemResponse) error,
) (*LineItemResponse, error) {
	if err := ValidateNotEmpty(extProjectID, extLineItemID); err != nil {
		return nil, err
	}

	var res LineItemResponse
	path := fmt.Sprintf("/projects/%s/lineItems/%s", extProjectID, extLineItemID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetFeasibility ... Returns the feasibility for all the line items of the requested project. Takes 20 - 120
// seconds to execute. Check the `GetFeasibilityResponse.Feasibility.Status` field value to see if it is
// FeasibilityStatusReady ("READY") or FeasibilityStatusProcessing ("PROCESSING")
// If GetFeasibilityResponse.Feasibility.Status == FeasibilityStatusProcessing, call this function again in 2 mins.
func (c *Client) GetFeasibility(extProjectID string, options *QueryOptions) (*GetFeasibilityResponse, error) {
	f := func(path string, res *GetFeasibilityResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getFeasibility(extProjectID, options, f)
}

// GetFeasibilityWithContext ... Returns the feasibility for all the line items of the requested project.
// Takes 20 - 120 seconds to execute. Check the `GetFeasibilityResponse.Feasibility.Status` field value
// to see if it is FeasibilityStatusReady ("READY") or FeasibilityStatusProcessing ("PROCESSING")
// If GetFeasibilityResponse.Feasibility.Status == FeasibilityStatusProcessing, call this function again in 2 mins.
func (c *Client) GetFeasibilityWithContext(
	ctx context.Context,
	extProjectID string,
	options *QueryOptions,
) (*GetFeasibilityResponse, error) {
	f := func(path string, res *GetFeasibilityResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getFeasibility(extProjectID, options, f)
}

func (c *Client) getFeasibility(
	extProjectID string,
	options *QueryOptions,
	requestAndParseResponse func(path string, res *GetFeasibilityResponse) error,
) (*GetFeasibilityResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res GetFeasibilityResponse
	path := fmt.Sprintf("/projects/%s/feasibility%s", extProjectID, query2String(options))

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetInvoice ... Get the invoice of the requested project
func (c *Client) GetInvoice(extProjectID string, options *QueryOptions) (*APIResponse, error) {
	path := fmt.Sprintf("/projects/%s/invoices", extProjectID)
	return c.request(http.MethodGet, c.Options.APIBaseURL, path, nil)
}

// GetInvoiceWithContext ... Get the invoice of the requested project
func (c *Client) GetInvoiceWithContext(
	ctx context.Context,
	extProjectID string,
	options *QueryOptions,
) (*APIResponse, error) {
	path := fmt.Sprintf("/projects/%s/invoices", extProjectID)
	return c.requestWithContext(ctx, http.MethodGet, c.Options.APIBaseURL, path, nil)
}

// UploadReconcile ...  Upload the Request correction file
func (c *Client) UploadReconcile(
	extProjectID string,
	file multipart.File,
	fileName string,
	message string,
	options *QueryOptions,
) (*APIResponse, error) {
	return c.sendFormData(
		c.Options.APIBaseURL,
		http.MethodPost,
		fmt.Sprintf("/projects/%s/reconcile", extProjectID),
		c.Auth.AccessToken,
		file,
		fileName,
		message,
	)
}

// UploadReconcileWithContext ...  Upload the Request correction file
func (c *Client) UploadReconcileWithContext(
	ctx context.Context,
	extProjectID string,
	file multipart.File,
	fileName string,
	message string,
	options *QueryOptions,
) (*APIResponse, error) {
	return c.sendFormDataWithContext(
		ctx,
		c.Options.APIBaseURL,
		http.MethodPost,
		fmt.Sprintf("/projects/%s/reconcile", extProjectID),
		c.Auth.AccessToken,
		file,
		fileName,
		message,
	)
}

// GetCountries ... Get the list of supported countries and languages in each country.
func (c *Client) GetCountries(options *QueryOptions) (*GetCountriesResponse, error) {
	var res GetCountriesResponse
	path := fmt.Sprintf("/countries%s", query2String(options))

	if err := c.requestAndParseResponse(http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetCountriesWithContext ... Get the list of supported countries and languages in each country.
func (c *Client) GetCountriesWithContext(ctx context.Context, options *QueryOptions) (*GetCountriesResponse, error) {
	var res GetCountriesResponse
	path := fmt.Sprintf("/countries%s", query2String(options))

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAttributes ... Get the list of supported attributes for a country and
// language. This data is required to build up the Quota Plan.
func (c *Client) GetAttributes(countryCode, languageCode string, options *QueryOptions) (*GetAttributesResponse, error) {
	f := func(path string, res *GetAttributesResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getAttributes(countryCode, languageCode, options, f)
}

// GetAttributesWithContext ... Get the list of supported attributes for a
// country and language. This data is required to build up the Quota Plan.
func (c *Client) GetAttributesWithContext(
	ctx context.Context,
	countryCode string,
	languageCode string,
	options *QueryOptions,
) (*GetAttributesResponse, error) {
	f := func(path string, res *GetAttributesResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getAttributes(countryCode, languageCode, options, f)
}

func (c *Client) getAttributes(
	countryCode string,
	languageCode string,
	options *QueryOptions,
	requestAndParseResponse func(path string, res *GetAttributesResponse) error,
) (*GetAttributesResponse, error) {
	if err := ValidateNotEmpty(countryCode, languageCode); err != nil {
		return nil, err
	}

	if err := IsCountryCodeOrEmpty(countryCode); err != nil {
		return nil, err
	}

	if err := IsLanguageCodeOrEmpty(languageCode); err != nil {
		return nil, err
	}

	var res GetAttributesResponse
	path := fmt.Sprintf("/attributes/%s/%s%s", countryCode, languageCode, query2String(options))

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetSurveyTopics ... Get the list of supported Survey Topics for a project.
// This data is required to setup a project.
func (c *Client) GetSurveyTopics(options *QueryOptions) (*GetSurveyTopicsResponse, error) {
	var res GetSurveyTopicsResponse
	path := fmt.Sprintf("/categories/surveyTopics%s", query2String(options))

	if err := c.requestAndParseResponse(http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetSurveyTopicsWithContext ... Get the list of supported Survey Topics for a
// project. This data is required to setup a project.
func (c *Client) GetSurveyTopicsWithContext(ctx context.Context, options *QueryOptions) (*GetSurveyTopicsResponse, error) {
	var res GetSurveyTopicsResponse
	path := fmt.Sprintf("/categories/surveyTopics%s", query2String(options))

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetSources ... Get the list of all the Sample sources
func (c *Client) GetSources(options *QueryOptions) (*GetSampleSourceResponse, error) {
	var res GetSampleSourceResponse
	path := fmt.Sprintf("/sources%s", query2String(options))

	if err := c.requestAndParseResponse(http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetSourcesWithContext ... Get the list of all the Sample sources
func (c *Client) GetSourcesWithContext(ctx context.Context, options *QueryOptions) (*GetSampleSourceResponse, error) {
	var res GetSampleSourceResponse
	path := fmt.Sprintf("/sources%s", query2String(options))

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetEvents ... Returns the list of all events that have occurred for your
// company account. Most recent events occur at the top of the list.
func (c *Client) GetEvents(options *QueryOptions) (*GetEventListResponse, error) {
	var res GetEventListResponse
	path := fmt.Sprintf("/events%s", query2String(options))

	if err := c.requestAndParseResponse(http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetEventsWithContext ... Returns the list of all events that have occurred
// for your company account. Most recent events occur at the top of the list.
func (c *Client) GetEventsWithContext(ctx context.Context, options *QueryOptions) (*GetEventListResponse, error) {
	var res GetEventListResponse
	path := fmt.Sprintf("/events%s", query2String(options))

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetEventBy ... Returns the requested event based on the eventID
func (c *Client) GetEventBy(eventID string) (*GetEventResponse, error) {
	var res GetEventResponse
	path := fmt.Sprintf("/events/%s", eventID)

	if err := c.requestAndParseResponse(http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetEventByWithContext ... Returns the requested event based on the eventID
func (c *Client) GetEventByWithContext(ctx context.Context, eventID string) (*GetEventResponse, error) {
	var res GetEventResponse
	path := fmt.Sprintf("/events/%s", eventID)

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// AcceptEvent ...
func (c *Client) AcceptEvent(event *Event) error {
	if event.Actions == nil || len(event.Actions.AcceptURL) == 0 {
		return ErrEventActionNotApplicable
	}

	_, err := c.request(http.MethodPost, event.Actions.AcceptURL, "", nil)
	return err
}

// AcceptEventWithContext ...
func (c *Client) AcceptEventWithContext(ctx context.Context, event *Event) error {
	if event.Actions == nil || len(event.Actions.AcceptURL) == 0 {
		return ErrEventActionNotApplicable
	}

	_, err := c.requestWithContext(ctx, http.MethodPost, event.Actions.AcceptURL, "", nil)
	return err
}

// RejectEvent ...
func (c *Client) RejectEvent(event *Event) error {
	if event.Actions == nil || len(event.Actions.RejectURL) == 0 {
		return ErrEventActionNotApplicable
	}

	_, err := c.request(http.MethodPost, event.Actions.RejectURL, "", nil)
	return err
}

// RejectEventWithContext ...
func (c *Client) RejectEventWithContext(ctx context.Context, event *Event) error {
	if event.Actions == nil || len(event.Actions.RejectURL) == 0 {
		return ErrEventActionNotApplicable
	}

	_, err := c.requestWithContext(ctx, http.MethodPost, event.Actions.RejectURL, "", nil)
	return err
}

// GetDetailedProjectReport returns a project's detailed report based on
// observed data from actual panelists.
func (c *Client) GetDetailedProjectReport(extProjectID string) (*DetailedProjectReportResponse, error) {
	f := func(path string, res *DetailedProjectReportResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getDetailedProjectReport(extProjectID, f)
}

// GetDetailedProjectReportWithContext returns a project's detailed report
// based on observed data from actual panelists.
func (c *Client) GetDetailedProjectReportWithContext(
	ctx context.Context,
	extProjectID string,
) (*DetailedProjectReportResponse, error) {
	f := func(path string, res *DetailedProjectReportResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getDetailedProjectReport(extProjectID, f)
}

func (c *Client) getDetailedProjectReport(
	extProjectID string,
	requestAndParseResponse func(path string, res *DetailedProjectReportResponse) error,
) (*DetailedProjectReportResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res DetailedProjectReportResponse
	path := fmt.Sprintf("/projects/%s/detailedReport", extProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetDetailedLineItemReport returns a lineitems's report with quota cell level
// stats based on observed data from actual panelists.
func (c *Client) GetDetailedLineItemReport(extProjectID, extLineItemID string) (*DetailedLineItemReportResponse, error) {
	f := func(path string, res *DetailedLineItemReportResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.getDetailedLineItemReport(extProjectID, extLineItemID, f)
}

// GetDetailedLineItemReportWithContext returns a lineitems's report with quota
// cell level stats based on observed data from actual panelists.
func (c *Client) GetDetailedLineItemReportWithContext(
	ctx context.Context,
	extProjectID string,
	extLineItemID string,
) (*DetailedLineItemReportResponse, error) {
	f := func(path string, res *DetailedLineItemReportResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.getDetailedLineItemReport(extProjectID, extLineItemID, f)
}

func (c *Client) getDetailedLineItemReport(
	extProjectID string,
	extLineItemID string,
	requestAndParseResponse func(path string, res *DetailedLineItemReportResponse) error,
) (*DetailedLineItemReportResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res DetailedLineItemReportResponse
	path := fmt.Sprintf("/projects/%s/lineItems/%s/detailedReport", extProjectID, extLineItemID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetUserInfo gives information about the user that is currently logged in.
func (c *Client) GetUserInfo() (*UserResponse, error) {
	var res UserResponse

	if err := c.requestAndParseResponse(http.MethodGet, "/users/info", nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetUserInfoWithContext gives information about the user that is currently
// logged in.
func (c *Client) GetUserInfoWithContext(ctx context.Context) (*UserResponse, error) {
	var res UserResponse

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, "/users/info", nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CompanyUsers gives information about the user that is currently logged in.
func (c *Client) CompanyUsers() (*CompanyUsersResponse, error) {
	var res CompanyUsersResponse

	if err := c.requestAndParseResponse(http.MethodGet, "/users", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CompanyUsersWithContext gives information about the user that is currently logged in.
func (c *Client) CompanyUsersWithContext(ctx context.Context) (*CompanyUsersResponse, error) {
	var res CompanyUsersResponse

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, "/users", nil, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// TeamsInfo gives information about the user that is currently logged in.
func (c *Client) TeamsInfo() (*TeamsResponse, error) {
	var res TeamsResponse

	if err := c.requestAndParseResponse(http.MethodGet, "/teams", nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// TeamsInfoWithContext gives information about the user that is currently
// logged in.
func (c *Client) TeamsInfoWithContext(ctx context.Context) (*TeamsResponse, error) {
	var res TeamsResponse

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, "/teams", nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Roles returns the roles specified in the filter.
func (c *Client) Roles(options *QueryOptions) (*RolesResponse, error) {
	var res RolesResponse
	path := fmt.Sprintf("/roles%s", query2String(options))

	if err := c.requestAndParseResponse(http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// RolesWithContext returns the roles specified in the filter.
func (c *Client) RolesWithContext(ctx context.Context, options *QueryOptions) (*RolesResponse, error) {
	var res RolesResponse
	path := fmt.Sprintf("/roles%s", query2String(options))

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ProjectPermissions gives information about the user that is currently
// logged in.
func (c *Client) ProjectPermissions(extProjectID string) (*ProjectPermissionsResponse, error) {
	f := func(path string, res *ProjectPermissionsResponse) error {
		return c.requestAndParseResponse(http.MethodGet, path, nil, res)
	}

	return c.projectPermissions(extProjectID, f)
}

// ProjectPermissionsWithContext gives information about the user that is
// currently logged in.
func (c *Client) ProjectPermissionsWithContext(
	ctx context.Context,
	extProjectID string,
) (*ProjectPermissionsResponse, error) {
	f := func(path string, res *ProjectPermissionsResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res)
	}

	return c.projectPermissions(extProjectID, f)
}

func (c *Client) projectPermissions(
	extProjectID string,
	requestAndParseResponse func(path string, res *ProjectPermissionsResponse) error,
) (*ProjectPermissionsResponse, error) {
	if err := ValidateNotEmpty(extProjectID); err != nil {
		return nil, err
	}

	var res ProjectPermissionsResponse
	path := fmt.Sprintf("/projects/%s/permissions", extProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpsertProjectPermissions gives information about the user that is currently
// logged in.
func (c *Client) UpsertProjectPermissions(permissions *UpsertPermissionsCriteria) (*ProjectPermissionsResponse, error) {
	f := func(path string, res *ProjectPermissionsResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, permissions, res)
	}

	return c.upsertProjectPermissions(permissions, f)
}

// UpsertProjectPermissionsWithContext gives information about the user that is
// currently logged in.
func (c *Client) UpsertProjectPermissionsWithContext(
	ctx context.Context,
	permissions *UpsertPermissionsCriteria,
) (*ProjectPermissionsResponse, error) {
	f := func(path string, res *ProjectPermissionsResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, permissions, res)
	}

	return c.upsertProjectPermissions(permissions, f)
}

func (c *Client) upsertProjectPermissions(
	permissions *UpsertPermissionsCriteria,
	requestAndParseResponse func(path string, res *ProjectPermissionsResponse) error,
) (*ProjectPermissionsResponse, error) {
	if err := Validate(permissions); err != nil {
		return nil, err
	}

	var res ProjectPermissionsResponse
	path := fmt.Sprintf("/projects/%s/permissions", permissions.ExtProjectID)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetStudyMetadata returns study metadata property info
func (c *Client) GetStudyMetadata() (*StudyMetadataResponse, error) {
	var res StudyMetadataResponse

	if err := c.requestAndParseResponse(http.MethodGet, "/studyMetadata", nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetStudyMetadataWithContext returns study metadata property info
func (c *Client) GetStudyMetadataWithContext(ctx context.Context) (*StudyMetadataResponse, error) {
	var res StudyMetadataResponse

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, "/studyMetadata", nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// CreateTemplate ...
func (c *Client) CreateTemplate(template *TemplateCriteria) (*TemplateResponse, error) {
	f := func(res *TemplateResponse) error {
		return c.requestAndParseResponse(http.MethodPost, "/templates/quotaPlan", template, res)
	}

	return c.createTemplate(template, f)
}

// CreateTemplateWithContext ...
func (c *Client) CreateTemplateWithContext(ctx context.Context, template *TemplateCriteria) (*TemplateResponse, error) {
	f := func(res *TemplateResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, "/templates/quotaPlan", template, res)
	}

	return c.createTemplate(template, f)
}

func (c *Client) createTemplate(
	template *TemplateCriteria,
	requestAndParseResponse func(res *TemplateResponse) error,
) (*TemplateResponse, error) {
	if err := Validate(template); err != nil {
		return nil, err
	}

	var res TemplateResponse

	if err := requestAndParseResponse(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdateTemplate ...
func (c *Client) UpdateTemplate(id int, template *TemplateCriteria) (*TemplateResponse, error) {
	f := func(path string, res *TemplateResponse) error {
		return c.requestAndParseResponse(http.MethodPost, path, template, res)
	}

	return c.updateTemplate(id, template, f)
}

// UpdateTemplateWithContext ...
func (c *Client) UpdateTemplateWithContext(
	ctx context.Context,
	id int,
	template *TemplateCriteria,
) (*TemplateResponse, error) {
	f := func(path string, res *TemplateResponse) error {
		return c.requestAndParseResponseWithContext(ctx, http.MethodPost, path, template, res)
	}

	return c.updateTemplate(id, template, f)
}

func (c *Client) updateTemplate(
	id int,
	template *TemplateCriteria,
	requestAndParseResponse func(path string, res *TemplateResponse) error,
) (*TemplateResponse, error) {
	if err := Validate(template); err != nil {
		return nil, err
	}

	var res TemplateResponse
	path := fmt.Sprintf("/templates/quotaPlan/%d", id)

	if err := requestAndParseResponse(path, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetTemplateList ...
func (c *Client) GetTemplateList(country string, lang string, options *QueryOptions) (*TemplatesResponse, error) {
	var res TemplatesResponse
	path := fmt.Sprintf("/templates/quotaPlan/%s/%s%s", country, lang, query2String(options))

	if err := c.requestAndParseResponse(http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetTemplateListWithContext ...
func (c *Client) GetTemplateListWithContext(
	ctx context.Context,
	country string,
	lang string,
	options *QueryOptions,
) (*TemplatesResponse, error) {
	var res TemplatesResponse
	path := fmt.Sprintf("/templates/quotaPlan/%s/%s%s", country, lang, query2String(options))

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodGet, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteTemplate ...
func (c *Client) DeleteTemplate(id int) (*AppError, error) {
	var res AppError
	path := fmt.Sprintf("/templates/quotaPlan/%d", id)

	if err := c.requestAndParseResponse(http.MethodDelete, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteTemplateWithContext ...
func (c *Client) DeleteTemplateWithContext(ctx context.Context, id int) (*AppError, error) {
	var res AppError
	path := fmt.Sprintf("/templates/quotaPlan/%d", id)

	if err := c.requestAndParseResponseWithContext(ctx, http.MethodDelete, path, nil, res); err != nil {
		return nil, err
	}

	return &res, nil
}

// RefreshToken ...
func (c *Client) RefreshToken() error {
	f := func(url string, req interface{}) (*APIResponse, error) {
		return c.sendRequest(url, http.MethodPost, "/token/refresh", "", req)
	}

	return c.refreshToken(f)
}

// RefreshTokenWithContext ...
func (c *Client) RefreshTokenWithContext(ctx context.Context) error {
	f := func(url string, req interface{}) (*APIResponse, error) {
		return c.sendRequestWithContext(ctx, url, http.MethodPost, "/token/refresh", "", req)
	}

	return c.refreshToken(f)
}

func (c *Client) refreshToken(sendRequest func(url string, req interface{}) (*APIResponse, error)) error {
	if c.Auth.RefreshTokenExpired() {
		return ErrSessionExpired
	}

	req := struct {
		ClientID     string `json:"clientId"`
		RefreshToken string `json:"refreshToken"`
	}{
		ClientID:     c.Credentials.ClientID,
		RefreshToken: c.Auth.RefreshToken,
	}

	ar, err := sendRequest(c.Options.AuthURL, &req)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(ar.Body, &c.Auth); err != nil {
		return err
	}

	t := time.Now()
	c.Auth.Acquired = &t

	return nil
}

// Logout ...
func (c *Client) Logout() error {
	if c.Auth.AccessTokenExpired() {
		return nil
	}

	req := struct {
		ClientID     string `json:"clientId"`
		RefreshToken string `json:"refreshToken"`
		AccessToken  string `json:"accessToken"`
	}{
		ClientID:     c.Credentials.ClientID,
		RefreshToken: c.Auth.RefreshToken,
		AccessToken:  c.Auth.AccessToken,
	}

	_, err := c.sendRequest(c.Options.AuthURL, http.MethodPost, "/logout", "", req)
	return err
}

// LogoutWithContext ...
func (c *Client) LogoutWithContext(ctx context.Context) error {
	if c.Auth.AccessTokenExpired() {
		return nil
	}

	req := struct {
		ClientID     string `json:"clientId"`
		RefreshToken string `json:"refreshToken"`
		AccessToken  string `json:"accessToken"`
	}{
		ClientID:     c.Credentials.ClientID,
		RefreshToken: c.Auth.RefreshToken,
		AccessToken:  c.Auth.AccessToken,
	}

	_, err := c.sendRequestWithContext(ctx, c.Options.AuthURL, http.MethodPost, "/logout", "", req)
	return err
}

// GetAuth ...
func (c *Client) GetAuth() (TokenResponse, error) {
	if err := c.requestAndParseToken(); err != nil {
		return TokenResponse{}, err
	}

	return c.Auth, nil
}

// GetAuthWithContext ...
func (c *Client) GetAuthWithContext(ctx context.Context) (TokenResponse, error) {
	if err := c.requestAndParseTokenWithContext(ctx); err != nil {
		return TokenResponse{}, err
	}

	return c.Auth, nil
}

func (c *Client) requestAndParseResponse(method, url string, body interface{}, resObj interface{}) error {
	ar, err := c.request(method, c.Options.APIBaseURL, url, body)
	if err != nil {
		if ar != nil {
			json.Unmarshal(ar.Body, &resObj)
		}
		return err
	}
	err = json.Unmarshal(ar.Body, &resObj)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) requestAndParseResponseWithContext(
	ctx context.Context,
	method string,
	url string,
	body interface{},
	resObj interface{},
) error {
	ar, err := c.requestWithContext(ctx, method, c.Options.APIBaseURL, url, body)
	if err != nil {
		if ar != nil {
			json.Unmarshal(ar.Body, &resObj)
		}
		return err
	}

	if err = json.Unmarshal(ar.Body, &resObj); err != nil {
		return err
	}

	return nil
}

func (c *Client) request(method, host, url string, body interface{}) (*APIResponse, error) {
	err := c.validateTokens()
	if err != nil {
		return nil, err
	}
	ar, err := c.sendRequest(host, method, url, c.Auth.AccessToken, body)
	errResp, ok := err.(*ErrorResponse)
	if ok && errResp.HTTPCode == http.StatusUnauthorized {
		err := c.requestAndParseToken()
		if err != nil {
			return nil, err
		}
		return c.sendRequest(host, method, url, c.Auth.AccessToken, body)
	}
	return ar, err
}

func (c *Client) requestWithContext(ctx context.Context, method, host, url string, body interface{}) (*APIResponse, error) {
	err := c.validateTokensWithContext(ctx)
	if err != nil {
		return nil, err
	}

	ar, err := c.sendRequestWithContext(ctx, host, method, url, c.Auth.AccessToken, body)
	if err != nil {
		var errResp *ErrorResponse

		if errors.As(err, &errResp) && errResp.HTTPCode == http.StatusUnauthorized {
			if err = c.requestAndParseToken(); err != nil {
				return nil, err
			}

			return c.sendRequestWithContext(ctx, host, method, url, c.Auth.AccessToken, body)
		}

		return nil, err
	}

	return ar, err
}

func (c *Client) requestAndParseToken() error {
	t := time.Now()
	ar, err := c.sendRequest(c.Options.AuthURL, http.MethodPost, "/token/password", "", c.Credentials)
	if err != nil {
		return err
	}
	err = json.Unmarshal(ar.Body, &c.Auth)
	if err != nil {
		return err
	}
	c.Auth.Acquired = &t
	return nil
}

func (c *Client) requestAndParseTokenWithContext(ctx context.Context) error {
	ar, err := c.sendRequestWithContext(ctx, c.Options.AuthURL, http.MethodPost, "/token/password", "", c.Credentials)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(ar.Body, &c.Auth); err != nil {
		return err
	}

	t := time.Now()
	c.Auth.Acquired = &t

	return nil
}

// validateTokens ...
func (c *Client) validateTokens() error {
	if c.Auth.AccessTokenExpired() {
		err := c.RefreshToken()
		if err != nil {
			err := c.requestAndParseToken()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// validateTokensWithContext ...
func (c *Client) validateTokensWithContext(ctx context.Context) error {
	if c.Auth.AccessTokenExpired() {
		err := c.RefreshTokenWithContext(ctx)
		if err != nil {
			err := c.requestAndParseTokenWithContext(ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewClient returns an API client.
// If options is nil, UATClientOptions will be used.
func NewClient(clientID, username, password string, options *ClientOptions) *Client {
	if options == nil {
		options = UATClientOptions
	}

	if options.Timeout == nil {
		timeout := defaulttimeout
		options.Timeout = &timeout
	}

	if options.HTTPClient == nil {
		options.HTTPClient = &http.Client{
			Timeout: time.Second * time.Duration(*options.Timeout),
		}
	}

	return &Client{
		Credentials: TokenRequest{
			ClientID: clientID,
			Username: username,
			Password: password,
		},
		Options: options,
	}
}

// SetOptions ...
func (c *Client) SetOptions(env string, timeout int) error {
	switch env {
	case "dev":
		c.Options = DevClientOptions
	case "uat":
		c.Options = UATClientOptions
	case "prod":
		c.Options = ProdClientOptions
	default:
		return ErrIncorrectEnvironemt
	}

	if timeout == 0 {
		timeout = defaulttimeout
		c.Options.Timeout = &timeout
	}

	c.Options.HTTPClient = &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}

	return nil
}

// NewClientFromEnv returns an API client.
func NewClientFromEnv(clientID, username, passsword string, env string, timeout int) (*Client, error) {
	client := &Client{
		Credentials: TokenRequest{
			ClientID: clientID,
			Username: username,
			Password: passsword,
		},
	}

	if err := client.SetOptions(env, timeout); err != nil {
		return nil, err
	}

	return client, nil
}

// GetHealthyStatus ... Get the healthy status on API
func (c *Client) GetHealthyStatus() (*APIResponse, error) {
	return c.request(http.MethodGet, c.Options.GatewayURL, "", nil)
}
