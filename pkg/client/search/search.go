package search

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v4/opensearchtransport"
)

func NewClient(conf config.SearchConfig, runtime config.Runtime) (*opensearchapi.Client, error) {
	clientAddresses := []string{}
	for _, addr := range conf.ClientAddress {
		clientAddresses = append(clientAddresses, addr.String())
	}

	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Addresses:     clientAddresses,
			Username:      conf.Username,
			Password:      conf.Password,
			EnableMetrics: true,
			Logger: func() opensearchtransport.Logger {
				if runtime.IsLocalEnvironment() {
					return &opensearchtransport.ColorLogger{
						Output: os.Stdout,
					}
				}
				return &opensearchtransport.JSONLogger{
					Output: os.Stdout,
				}
			}(),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating search client: %w", err)
	}

	return client, nil
}
