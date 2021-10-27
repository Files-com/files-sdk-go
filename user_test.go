package files_sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_UnmarshalJSON(t *testing.T) {
	json := `{"admin_group_ids": [1, 2, 3]}`

	user := User{}

	err := user.UnmarshalJSON([]byte(json))
	assert.NoError(t, err)
	assert.Equal(t, []int64{1, 2, 3}, user.AdminGroupIds)

	json = `{"admin_group_ids": []}`

	user = User{}

	err = user.UnmarshalJSON([]byte(json))
	assert.NoError(t, err)
	assert.Equal(t, []int64{}, user.AdminGroupIds)
}
