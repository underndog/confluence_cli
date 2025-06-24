package req

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPage_JSONMarshal(t *testing.T) {
	page := Page{
		SpaceID:  "TEST",
		ParentId: "123",
		Status:   "current",
		Title:    "Test Page",
		Body: Body{
			Representation: "storage",
			Value:          "<p>Test content</p>",
		},
	}

	jsonData, err := json.Marshal(page)
	assert.NoError(t, err)

	var unmarshaledPage Page
	err = json.Unmarshal(jsonData, &unmarshaledPage)
	assert.NoError(t, err)

	assert.Equal(t, page.SpaceID, unmarshaledPage.SpaceID)
	assert.Equal(t, page.ParentId, unmarshaledPage.ParentId)
	assert.Equal(t, page.Status, unmarshaledPage.Status)
	assert.Equal(t, page.Title, unmarshaledPage.Title)
	assert.Equal(t, page.Body.Representation, unmarshaledPage.Body.Representation)
	assert.Equal(t, page.Body.Value, unmarshaledPage.Body.Value)
}

func TestUpdatePagePayload_JSONMarshal(t *testing.T) {
	payload := UpdatePagePayload{
		ID:     "123",
		Status: "current",
		Title:  "Updated Page",
		Body: UpdatePageBody{
			Representation: "storage",
			Value:          "<p>Updated content</p>",
		},
		Version: UpdatePageVersion{
			Number:  2,
			Message: "Updated via CLI",
		},
	}

	jsonData, err := json.Marshal(payload)
	assert.NoError(t, err)

	var unmarshaledPayload UpdatePagePayload
	err = json.Unmarshal(jsonData, &unmarshaledPayload)
	assert.NoError(t, err)

	assert.Equal(t, payload.ID, unmarshaledPayload.ID)
	assert.Equal(t, payload.Status, unmarshaledPayload.Status)
	assert.Equal(t, payload.Title, unmarshaledPayload.Title)
	assert.Equal(t, payload.Body.Representation, unmarshaledPayload.Body.Representation)
	assert.Equal(t, payload.Body.Value, unmarshaledPayload.Body.Value)
	assert.Equal(t, payload.Version.Number, unmarshaledPayload.Version.Number)
	assert.Equal(t, payload.Version.Message, unmarshaledPayload.Version.Message)
}

func TestBody_JSONMarshal(t *testing.T) {
	body := Body{
		Representation: "storage",
		Value:          "<p>Test content</p>",
	}

	jsonData, err := json.Marshal(body)
	assert.NoError(t, err)

	var unmarshaledBody Body
	err = json.Unmarshal(jsonData, &unmarshaledBody)
	assert.NoError(t, err)

	assert.Equal(t, body.Representation, unmarshaledBody.Representation)
	assert.Equal(t, body.Value, unmarshaledBody.Value)
}

func TestUpdatePageVersion_JSONMarshal(t *testing.T) {
	version := UpdatePageVersion{
		Number:  3,
		Message: "Test update",
	}

	jsonData, err := json.Marshal(version)
	assert.NoError(t, err)

	var unmarshaledVersion UpdatePageVersion
	err = json.Unmarshal(jsonData, &unmarshaledVersion)
	assert.NoError(t, err)

	assert.Equal(t, version.Number, unmarshaledVersion.Number)
	assert.Equal(t, version.Message, unmarshaledVersion.Message)
}
