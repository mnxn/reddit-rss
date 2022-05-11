package reddit

import (
	"fmt"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetMyTrophies(t *testing.T) {
	url := fmt.Sprintf("%s/api/v1/me/trophies", baseAuthURL)
	mockResponseFromFile(url, "test_data/award/my_trophies.json")
	defer httpmock.DeactivateAndReset()

	client := NoAuthClient
	trophies, err := client.GetMyTrophies()
	assert.NoError(t, err)
	assert.Equal(t, len(trophies), 1)
}
