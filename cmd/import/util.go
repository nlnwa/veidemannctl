package importcmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	configV1 "github.com/nlnwa/veidemann-api/go/config/v1"
	"golang.org/x/net/html"
)

// UriChecker checks if a uri is reachable
type UriChecker struct {
	*http.Client
}

// Check checks if a uri is reachable and returns the uri if it is reachable
// If the uri is not reachable, it returns an error
// If the uri is redirected with 301, it returns the redirected uri
func (uc *UriChecker) Check(uri string) (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodHead, uri, nil)
	if err != nil {
		return "", err
	}
	resp, err := uc.Client.Do(req)
	if err != nil {
		var uerr *url.Error
		if errors.As(err, &uerr) && uerr.Timeout() {
			return "", fmt.Errorf("timeout")
		}
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) {
			return "", fmt.Errorf("no such host: %s", dnsErr.Name)
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusMovedPermanently {
		u, err := resp.Location()
		if err != nil {
			return "", err
		}
		return u.String(), nil
	}

	if resp.StatusCode < 400 {
		return uri, nil
	}

	return uri, fmt.Errorf("%s", resp.Status)
}

// GetTitle returns the title of the uri
func (uc *UriChecker) GetTitle(uri string) string {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, uri, nil)
	if err != nil {
		return ""
	}
	resp, err := uc.Client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return ""
	}
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = strings.TrimSpace(n.FirstChild.Data)
			return
		}
		if n.Type == html.ElementNode && n.Data == "body" {
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title
}

type seedDesc struct {
	EntityId          string            `json:"entityId,omitempty" yaml:"entityId,omitempty"`
	EntityName        string            `json:"entityName,omitempty" yaml:"entityName,omitempty"`
	EntityDescription string            `json:"entityDescription,omitempty" yaml:"entityDescription,omitempty"`
	EntityLabel       []*configV1.Label `json:"entityLabel,omitempty" yaml:"entityLabel,omitempty"`
	Uri               string            `json:"uri,omitempty" yaml:"uri,omitempty"`
	SeedDescription   string            `json:"seedDescription,omitempty" yaml:"seedDescription,omitempty"`
	SeedLabel         []*configV1.Label `json:"seedLabel,omitempty" yaml:"seedLabel,omitempty"`

	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	crawlJobRef []*configV1.ConfigRef
}

func (sd *seedDesc) String() string {
	b, _ := json.Marshal(sd)
	return string(b)
}

func (sd *seedDesc) toEntity() *configV1.ConfigObject {
	return &configV1.ConfigObject{
		ApiVersion: "v1",
		Kind:       configV1.Kind_crawlEntity,
		Id:         sd.EntityId,
		Meta: &configV1.Meta{
			Name:        sd.EntityName,
			Description: sd.EntityDescription,
			Label:       sd.EntityLabel,
		},
	}
}

func (sd *seedDesc) toSeed() *configV1.ConfigObject {
	return &configV1.ConfigObject{
		ApiVersion: "v1",
		Kind:       configV1.Kind_seed,
		Meta: &configV1.Meta{
			Name:        sd.Uri,
			Description: sd.SeedDescription,
			Label:       sd.SeedLabel,
		},
		Spec: &configV1.ConfigObject_Seed{
			Seed: &configV1.Seed{
				EntityRef: &configV1.ConfigRef{
					Kind: configV1.Kind_crawlEntity,
					Id:   sd.EntityId,
				},
				JobRef: sd.crawlJobRef,
			},
		},
	}
}
