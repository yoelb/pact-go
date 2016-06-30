package dsl

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/pact-foundation/pact-go/types"
	"github.com/pact-foundation/pact-go/utils"
)

func TestPact_findConsumers(t *testing.T) {
	port := setupMockBroker()
	request := &types.VerifyRequest{
		Tags:      []string{"dev", "prod"},
		BrokerURL: fmt.Sprintf("http://localhost:%d", port),
	}
	err := findConsumers("bobby", request)
	fmt.Println(err)
}

// Pretend to be a Broker for fetching Pacts
func setupMockBroker() int {
	port, _ := utils.GetFreePort()
	mux := http.NewServeMux()

	// Find latest 'bobby' consumers (no tag)
	// curl --user pactuser:pact -H "accept: application/hal+json" "http://pact.onegeek.com.au/pacts/provider/bobby/latest"
	mux.HandleFunc("/pacts/provider/bobby/latest", func(w http.ResponseWriter, req *http.Request) {
		log.Println("[DEBUG] get pacts for provider 'bobby'")
		fmt.Fprintf(w, `{"_links":{"self":{"href":"http://localhost:%d/pacts/provider/bobby/latest","title":"Latest pact versions for the provider bobby"},"provider":{"href":"http://localhost:%d/pacticipants/bobby","title":"bobby"},"pacts":[{"href":"http://localhost:%d/pacts/provider/bobby/consumer/billy/version/1.0.0","title":"Pact between billy (v1.0.0) and bobby","name":"billy"}]}}`, port, port, port)
		w.Header().Add("Content-Type", "application/hal+json")
	})

	// Find 'bobby' consumers for tag 'prod'
	// curl --user pactuser:pact -H "accept: application/hal+json" "http://pact.onegeek.com.au/pacts/provider/bobby/latest/sit4"
	mux.HandleFunc("/pacts/provider/bobby/latest/prod", func(w http.ResponseWriter, req *http.Request) {
		log.Println("[DEBUG] get all pacts for provider 'bobby' where the tag 'prod' exists")
		fmt.Fprintf(w, `{"_links":{"self":{"href":"http://localhost:%d/pacts/provider/bobby/latest/dev","title":"Latest pact versions for the provider bobby with tag 'dev'"},"provider":{"href":"http://localhost:%d/pacticipants/bobby","title":"bobby"},"pacts":[{"href":"http://localhost:%d/pacts/provider/bobby/consumer/billy/version/1.0.0","title":"Pact between billy (v1.0.0) and bobby","name":"billy"}]}}`, port, port, port)
		w.Header().Add("Content-Type", "application/hal+json")
	})

	// Find 'bobby' consumers for tag 'dev'
	// curl --user pactuser:pact -H "accept: application/hal+json" "http://pact.onegeek.com.au/pacts/provider/bobby/latest/sit4"
	mux.HandleFunc("/pacts/provider/bobby/latest/dev", func(w http.ResponseWriter, req *http.Request) {
		log.Println("[DEBUG] get all pacts for provider 'bobby' where the tag 'dev' exists")
		fmt.Fprintf(w, `{"_links":{"self":{"href":"http://localhost:%d/pacts/provider/bobby/latest/dev","title":"Latest pact versions for the provider bobby with tag 'dev'"},"provider":{"href":"http://localhost:%d/pacticipants/bobby","title":"bobby"},"pacts":[{"href":"http://localhost:%d/pacts/provider/bobby/consumer/billy/version/1.0.0","title":"Pact between billy (v1.0.0) and bobby","name":"billy"}]}}`, port, port, port)
		w.Header().Add("Content-Type", "application/hal+json")
	})

	// Consumer Pact
	// curl -v --user pactuser:pact -H "accept: application/json" http://pact.onegeek.com.au/pacts/provider/bobby/consumer/billy/version/1.0.0
	mux.HandleFunc("/pacts/provider/bobby/consumer/billy/version/1.0.0", func(w http.ResponseWriter, req *http.Request) {
		log.Println("[DEBUG] get all pacts for provider 'bobby' where the tag 'dev' exists")
		fmt.Fprintf(w, `{"consumer":{"name":"billy"},"provider":{"name":"bobby"},"interactions":[{"description":"Some name for the test","provider_state":"Some state","request":{"method":"GET","path":"/foobar"},"response":{"status":200,"headers":{"Content-Type":"application/json"}}},{"description":"Some name for the test","provider_state":"Some state2","request":{"method":"GET","path":"/bazbat"},"response":{"status":200,"headers":{},"body":[[{"colour":"red","size":10,"tag":[["jumper","shirt"],["jumper","shirt"]]}]],"matchingRules":{"$.body":{"min":1},"$.body[*].*":{"match":"type"},"$.body[*]":{"min":1},"$.body[*][*].*":{"match":"type"},"$.body[*][*].colour":{"match":"regex","regex":"red|green|blue"},"$.body[*][*].size":{"match":"type"},"$.body[*][*].tag":{"min":2},"$.body[*][*].tag[*].*":{"match":"type"},"$.body[*][*].tag[*][0]":{"match":"type"},"$.body[*][*].tag[*][1]":{"match":"type"}}}}],"metadata":{"pactSpecificationVersion":"2.0.0"},"updatedAt":"2016-06-11T13:11:33+00:00","createdAt":"2016-06-09T12:46:42+00:00","_links":{"self":{"title":"Pact","name":"Pact between billy (v1.0.0) and bobby","href":"http://localhost:%d/pacts/provider/bobby/consumer/billy/version/1.0.0"},"pb:consumer":{"title":"Consumer","name":"billy","href":"http://localhost:%d/pacticipants/billy"},"pb:provider":{"title":"Provider","name":"bobby","href":"http://localhost:%d/pacticipants/bobby"},"pb:latest-pact-version":{"title":"Pact","name":"Latest version of this pact","href":"http://localhost:%d/pacts/provider/bobby/consumer/billy/latest"},"pb:previous-distinct":{"title":"Pact","name":"Previous distinct version of this pact","href":"http://localhost:%d/pacts/provider/bobby/consumer/billy/version/1.0.0/previous-distinct"},"pb:diff-previous-distinct":{"title":"Diff","name":"Diff with previous distinct version of this pact","href":"http://localhost:%d/pacts/provider/bobby/consumer/billy/version/1.0.0/diff/previous-distinct"},"pb:pact-webhooks":{"title":"Webhooks for the pact between billy and bobby","href":"http://localhost:%d/webhooks/provider/bobby/consumer/billy"},"pb:tag-prod-version":{"title":"Tag this version as 'production'","href":"http://localhost:%d/pacticipants/billy/versions/1.0.0/tags/prod"},"pb:tag-version":{"title":"Tag version","href":"http://localhost:%d/pacticipants/billy/versions/1.0.0/tags/{tag}"},"curies":[{"name":"pb","href":"http://localhost:%d/doc/{rel}","templated":true}]}}`, port, port, port, port, port, port, port, port, port, port)
		w.Header().Add("Content-Type", "application/hal+json")
	})

	go http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	return port
}
