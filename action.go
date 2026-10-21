package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// Action is a Files.com API resource.
type Action struct {
	// Action ID.
	Id int64 `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Path. This must be slash-delimited, but it must neither start nor end with a slash. Maximum of 5000 characters.
	Path string `json:"path,omitempty" path:"path,omitempty" url:"path,omitempty"`
	// Action occurrence date/time.
	When *time.Time `json:"when,omitempty" path:"when,omitempty" url:"when,omitempty"`
	// The destination path for this action, if applicable.
	Destination string `json:"destination,omitempty" path:"destination,omitempty" url:"destination,omitempty"`
	// Friendly displayed output.
	Display string `json:"display,omitempty" path:"display,omitempty" url:"display,omitempty"`
	// IP Address that performed this action.
	Ip string `json:"ip,omitempty" path:"ip,omitempty" url:"ip,omitempty"`
	// The source path for this action, if applicable.
	Source string `json:"source,omitempty" path:"source,omitempty" url:"source,omitempty"`
	// Targets.
	Targets interface{} `json:"targets,omitempty" path:"targets,omitempty" url:"targets,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// Username.
	Username string `json:"username,omitempty" path:"username,omitempty" url:"username,omitempty"`
	// true if this change was performed by a user on a parent site.
	UserIsFromParentSite *bool `json:"user_is_from_parent_site,omitempty" path:"user_is_from_parent_site,omitempty" url:"user_is_from_parent_site,omitempty"`
	// Type of action.
	Action string `json:"action,omitempty" path:"action,omitempty" url:"action,omitempty"`
	// Failure type. If action was a user login or session failure, why did it fail?
	FailureType string `json:"failure_type,omitempty" path:"failure_type,omitempty" url:"failure_type,omitempty"`
	// Interface on which this action occurred.
	Interface string `json:"interface,omitempty" path:"interface,omitempty" url:"interface,omitempty"`
}

// Identifier returns the resource ID.
func (a Action) Identifier() interface{} {
	return a.Id
}

// ActionCollection is a list of Action resources.
type ActionCollection []Action

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *Action) UnmarshalJSON(data []byte) error {
	type action Action
	var v action
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = Action(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *ActionCollection) UnmarshalJSON(data []byte) error {
	type actions ActionCollection
	var v actions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = ActionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *ActionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
