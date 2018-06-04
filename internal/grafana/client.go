/*
Copyright 2018 The Chronologist Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grafana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const basePath = "/api"

// Client is a grafana HTTP API client.
//
// Client implements Annotator interface using grafana HTTP API.
type Client struct {
	host   string
	apiKey string

	client *http.Client
}

// SaveAnnotation saves annotation to grafana, either creating or updating it.
//
// See:
// - http://docs.grafana.org/v4.6/http_api/annotations/#create-annotation
// - http://docs.grafana.org/v4.6/http_api/annotations/#update-annotation
func (c *Client) SaveAnnotation(ctx context.Context, annotation Annotation) error {
	if annotation.ID == 0 {
		return c.createAnnotation(ctx, annotation)
	}
	return c.updateAnnotation(ctx, annotation)
}

func (c *Client) createAnnotation(ctx context.Context, annotation Annotation) error {
	b, err := json.Marshal(annotation)
	if err != nil {
		return errors.Wrap(err, "encode request to json")
	}

	u := fmt.Sprintf("%s%s/annotations", c.host, basePath)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "create request")
	}
	req = c.enrichRequest(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("got response %s", resp.Status)
	}

	return nil
}

func (c *Client) updateAnnotation(ctx context.Context, annotation Annotation) error {
	b, err := json.Marshal(annotation)
	if err != nil {
		return errors.Wrap(err, "encode request to json")
	}

	u := fmt.Sprintf("%s%s/annotations/%d", c.host, basePath, annotation.ID)
	req, err := http.NewRequest(http.MethodPut, u, bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "create request")
	}
	req = c.enrichRequest(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("got response %s", resp.Status)
	}

	return nil
}

// GetAnnotations fetches annotations from grafana using optional query params.
//
// See: http://docs.grafana.org/v4.6/http_api/annotations/#find-annotations
func (c *Client) GetAnnotations(ctx context.Context, in GetAnnotationsParams) (Annotations, error) {
	query := url.Values{}
	if len(in.Tags) > 0 {
		query["tags"] = in.Tags
	}

	u := fmt.Sprintf("%s%s%s?%s", c.host, basePath, "/annotations", query.Encode())
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, errors.Wrap(err, "create request")
	}
	req = c.enrichRequest(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("got response %s", resp.Status)
	}

	var aa Annotations
	if err = json.NewDecoder(resp.Body).Decode(&aa); err != nil {
		return nil, errors.Wrap(err, "decode response body from json")
	}

	return aa, nil
}

// DeleteAnnotation deletes annotation from grafana by its id.
//
// See: http://docs.grafana.org/v4.6/http_api/annotations/#delete-annotation-by-id
func (c *Client) DeleteAnnotation(ctx context.Context, id int) error {
	u := fmt.Sprintf("%s%s/annotations/%d", c.host, basePath, id)
	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return errors.Wrap(err, "create request")
	}
	req = c.enrichRequest(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("got response %s", resp.Status)
	}

	return nil
}

func (c *Client) enrichRequest(ctx context.Context, req *http.Request) *http.Request {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// NewClient returns a new grafana client.
func NewClient(host, apiKey string) *Client {
	return &Client{
		host:   host,
		apiKey: apiKey,
		client: http.DefaultClient,
	}
}
