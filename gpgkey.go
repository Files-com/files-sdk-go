package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

type GpgKey struct {
	Id                 int64      `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty" path:"expires_at,omitempty" url:"expires_at,omitempty"`
	Name               string     `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	UserId             int64      `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	PublicKey          string     `json:"public_key,omitempty" path:"public_key,omitempty" url:"public_key,omitempty"`
	PrivateKey         string     `json:"private_key,omitempty" path:"private_key,omitempty" url:"private_key,omitempty"`
	PrivateKeyPassword string     `json:"private_key_password,omitempty" path:"private_key_password,omitempty" url:"private_key_password,omitempty"`
}

func (g GpgKey) Identifier() interface{} {
	return g.Id
}

type GpgKeyCollection []GpgKey

type GpgKeyListParams struct {
	UserId int64                  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	SortBy map[string]interface{} `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	ListParams
}

type GpgKeyFindParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

type GpgKeyCreateParams struct {
	UserId             int64  `url:"user_id,omitempty" required:"false" json:"user_id,omitempty" path:"user_id"`
	PublicKey          string `url:"public_key,omitempty" required:"false" json:"public_key,omitempty" path:"public_key"`
	PrivateKey         string `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
	PrivateKeyPassword string `url:"private_key_password,omitempty" required:"false" json:"private_key_password,omitempty" path:"private_key_password"`
	Name               string `url:"name,omitempty" required:"true" json:"name,omitempty" path:"name"`
}

type GpgKeyUpdateParams struct {
	Id                 int64  `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
	PublicKey          string `url:"public_key,omitempty" required:"false" json:"public_key,omitempty" path:"public_key"`
	PrivateKey         string `url:"private_key,omitempty" required:"false" json:"private_key,omitempty" path:"private_key"`
	PrivateKeyPassword string `url:"private_key_password,omitempty" required:"false" json:"private_key_password,omitempty" path:"private_key_password"`
	Name               string `url:"name,omitempty" required:"false" json:"name,omitempty" path:"name"`
}

type GpgKeyDeleteParams struct {
	Id int64 `url:"-,omitempty" required:"false" json:"-,omitempty" path:"id"`
}

func (g *GpgKey) UnmarshalJSON(data []byte) error {
	type gpgKey GpgKey
	var v gpgKey
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*g = GpgKey(v)
	return nil
}

func (g *GpgKeyCollection) UnmarshalJSON(data []byte) error {
	type gpgKeys GpgKeyCollection
	var v gpgKeys
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*g = GpgKeyCollection(v)
	return nil
}

func (g *GpgKeyCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*g))
	for i, v := range *g {
		ret[i] = v
	}

	return &ret
}
