package consumer

import (
	"log"
	"testing"

	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/utils"
)

// V3HTTPMockProvider is the entrypoint for V3 http consumer tests
// This object is not thread safe
type V3HTTPMockProvider struct {
	*httpMockProvider
}

// NewV3Pact configures a new V3 HTTP Mock Provider for consumer tests
func NewV3Pact(config MockHTTPProviderConfig) (*V3HTTPMockProvider, error) {
	provider := &V3HTTPMockProvider{
		httpMockProvider: &httpMockProvider{
			config:               config,
			specificationVersion: models.V3,
		},
	}
	err := provider.configure()

	if err != nil {
		return nil, err
	}

	return provider, err
}

// AddInteraction to the pact
func (p *V3HTTPMockProvider) AddInteraction() *UnconfiguredV3Interaction {
	log.Println("[DEBUG] pact add V3 interaction")
	interaction := p.httpMockProvider.mockserver.NewInteraction("")

	i := &UnconfiguredV3Interaction{
		interaction: &Interaction{
			specificationVersion: models.V3,
			interaction:          interaction,
		},
		provider: p,
	}

	return i
}

type UnconfiguredV3Interaction struct {
	interaction *Interaction
	provider    *V3HTTPMockProvider
}

// Given specifies a provider state, may be called multiple times. Optional.
func (i *UnconfiguredV3Interaction) Given(state string) *UnconfiguredV3Interaction {
	i.interaction.interaction.Given(state)

	return i
}

// GivenWithParameter specifies a provider state with parameters, may be called multiple times. Optional.
func (i *UnconfiguredV3Interaction) GivenWithParameter(state models.ProviderState) *UnconfiguredV3Interaction {
	if len(state.Parameters) > 0 {
		i.interaction.interaction.GivenWithParameter(state.Name, state.Parameters)
	} else {
		i.interaction.interaction.Given(state.Name)
	}

	return i
}

type V3InteractionWithRequest struct {
	interaction *Interaction
	provider    *V3HTTPMockProvider
}

type V3RequestBuilder func(*V3InteractionWithRequestBuilder)

type V3InteractionWithRequestBuilder struct {
	interaction *Interaction
	provider    *V3HTTPMockProvider
}

// UponReceiving specifies the name of the test case. This becomes the name of
// the consumer/provider pair in the Pact file. Mandatory.
func (i *UnconfiguredV3Interaction) UponReceiving(description string) *UnconfiguredV3Interaction {
	i.interaction.interaction.UponReceiving(description)

	return i
}

// WithRequest provides a builder for the expected request
func (i *UnconfiguredV3Interaction) WithRequest(method Method, path string, builders ...V3RequestBuilder) *V3InteractionWithRequest {
	return i.WithRequestPathMatcher(method, matchers.String(path), builders...)
}

// WithRequestPathMatcher allows a matcher in the expected request path
func (i *UnconfiguredV3Interaction) WithRequestPathMatcher(method Method, path matchers.Matcher, builders ...V3RequestBuilder) *V3InteractionWithRequest {
	i.interaction.interaction.WithRequest(string(method), path)

	for _, builder := range builders {
		builder(&V3InteractionWithRequestBuilder{
			interaction: i.interaction,
			provider:    i.provider,
		})
	}

	return &V3InteractionWithRequest{
		interaction: i.interaction,
		provider:    i.provider,
	}
}

// Query specifies any query string on the expect request
func (i *V3InteractionWithRequestBuilder) Query(key string, values ...matchers.Matcher) *V3InteractionWithRequestBuilder {
	i.interaction.interaction.WithQuery(keyValuesToMapStringArrayInterface(key, values...))

	return i
}

// Header adds a header to the expected request
func (i *V3InteractionWithRequestBuilder) Header(key string, values ...matchers.Matcher) *V3InteractionWithRequestBuilder {
	i.interaction.interaction.WithRequestHeaders(keyValuesToMapStringArrayInterface(key, values...))

	return i
}

// Headers sets the headers on the expected request
func (i *V3InteractionWithRequestBuilder) Headers(headers matchers.HeadersMatcher) *V3InteractionWithRequestBuilder {
	i.interaction.interaction.WithRequestHeaders(headersMatcherToNativeHeaders(headers))

	return i
}

// JSONBody adds a JSON body to the expected request
func (i *V3InteractionWithRequestBuilder) JSONBody(body interface{}) *V3InteractionWithRequestBuilder {
	// TODO: Don't like panic, but not sure if there is a better builder experience?
	if err := validateMatchers(i.interaction.specificationVersion, body); err != nil {
		panic(err)
	}

	if s, ok := body.(string); ok {
		// Check if someone tried to add an object as a string representation
		// as per original allowed implementation, e.g.
		// { "foo": "bar", "baz": like("bat") }
		if utils.IsJSONFormattedObject(string(s)) {
			log.Println("[WARN] request body appears to be a JSON formatted object, " +
				"no matching will occur. Support for structured strings has been" +
				"deprecated as of 0.13.0. Please use the JSON() method instead")
		}
	}

	i.interaction.interaction.WithJSONRequestBody(body)

	return i
}

// BinaryBody adds a binary body to the expected request
func (i *V3InteractionWithRequestBuilder) BinaryBody(body []byte) *V3InteractionWithRequestBuilder {
	i.interaction.interaction.WithBinaryRequestBody(body)

	return i
}

// MultipartBody adds a multipart  body to the expected request
func (i *V3InteractionWithRequestBuilder) MultipartBody(contentType string, filename string, mimePartName string) *V3InteractionWithRequestBuilder {
	i.interaction.interaction.WithRequestMultipartFile(contentType, filename, mimePartName)

	return i
}

// Body adds general body to the expected request
func (i *V3InteractionWithRequestBuilder) Body(contentType string, body []byte) *V3InteractionWithRequestBuilder {
	// Check if someone tried to add an object as a string representation
	// as per original allowed implementation, e.g.
	// { "foo": "bar", "baz": like("bat") }
	if utils.IsJSONFormattedObject(string(body)) {
		log.Println("[WARN] request body appears to be a JSON formatted object, " +
			"no matching will occur. Support for structured strings has been" +
			"deprecated as of 0.13.0. Please use the JSON() method instead")
	}

	i.interaction.interaction.WithRequestBody(contentType, body)

	return i
}

// BodyMatch uses struct tags to automatically determine matchers from the given struct
func (i *V3InteractionWithRequestBuilder) BodyMatch(body interface{}) *V3InteractionWithRequestBuilder {
	i.interaction.interaction.WithJSONRequestBody(matchers.MatchV2(body))

	return i
}

// WillRespondWith sets the expected status and provides a response builder
func (i *V3InteractionWithRequest) WillRespondWith(status int, builders ...V3ResponseBuilder) *V3InteractionWithResponse {
	i.interaction.interaction.WithStatus(status)

	for _, builder := range builders {

		builder(&V3InteractionWithResponseBuilder{
			interaction: i.interaction,
			provider:    i.provider,
		})
	}

	return &V3InteractionWithResponse{
		interaction: i.interaction,
		provider:    i.provider,
	}
}

type V3ResponseBuilder func(*V3InteractionWithResponseBuilder)

type V3InteractionWithResponseBuilder struct {
	interaction *Interaction
	provider    *V3HTTPMockProvider
}

type V3InteractionWithResponse struct {
	interaction *Interaction
	provider    *V3HTTPMockProvider
}

// Header adds a header to the expected response
func (i *V3InteractionWithResponseBuilder) Header(key string, values ...matchers.Matcher) *V3InteractionWithResponseBuilder {
	i.interaction.interaction.WithResponseHeaders(keyValuesToMapStringArrayInterface(key, values...))

	return i
}

// Headers sets the headers on the expected response
func (i *V3InteractionWithResponseBuilder) Headers(headers matchers.HeadersMatcher) *V3InteractionWithResponseBuilder {
	i.interaction.interaction.WithResponseHeaders(headersMatcherToNativeHeaders(headers))

	return i
}

// JSONBody adds a JSON body to the expected response
func (i *V3InteractionWithResponseBuilder) JSONBody(body interface{}) *V3InteractionWithResponseBuilder {
	// TODO: Don't like panic, how to build a better builder here - nil return + log?
	if err := validateMatchers(i.interaction.specificationVersion, body); err != nil {
		panic(err)
	}

	if s, ok := body.(string); ok {
		// Check if someone tried to add an object as a string representation
		// as per original allowed implementation, e.g.
		// { "foo": "bar", "baz": like("bat") }
		if utils.IsJSONFormattedObject(string(s)) {
			log.Println("[WARN] response body appears to be a JSON formatted object, " +
				"no matching will occur. Support for structured strings has been" +
				"deprecated as of 0.13.0. Please use the JSON() method instead")
		}
	}
	i.interaction.interaction.WithJSONResponseBody(body)

	return i
}

// BinaryBody adds a binary body to the expected response
func (i *V3InteractionWithResponseBuilder) BinaryBody(body []byte) *V3InteractionWithResponseBuilder {
	i.interaction.interaction.WithBinaryResponseBody(body)

	return i
}

// MultipartBody adds a multipart  body to the expected response
func (i *V3InteractionWithResponseBuilder) MultipartBody(contentType string, filename string, mimePartName string) *V3InteractionWithResponseBuilder {
	i.interaction.interaction.WithResponseMultipartFile(contentType, filename, mimePartName)

	return i
}

// Body adds general body to the expected request
func (i *V3InteractionWithResponseBuilder) Body(contentType string, body []byte) *V3InteractionWithResponseBuilder {
	i.interaction.interaction.WithResponseBody(contentType, body)

	return i
}

// BodyMatch uses struct tags to automatically determine matchers from the given struct
func (i *V3InteractionWithResponseBuilder) BodyMatch(body interface{}) *V3InteractionWithResponseBuilder {
	i.interaction.interaction.WithJSONResponseBody(matchers.MatchV2(body))

	return i
}

// ExecuteTest runs the current test case against a Mock Service.
func (m *V3InteractionWithResponse) ExecuteTest(t *testing.T, integrationTest func(MockServerConfig) error) error {
	return m.provider.ExecuteTest(t, integrationTest)
}