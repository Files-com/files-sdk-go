package files_sdk

import (
	"encoding/json"
	"time"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// ChatSession is a Files.com API resource.
type ChatSession struct {
	// Chat Session ID.
	Id string `json:"id,omitempty" path:"id,omitempty" url:"id,omitempty"`
	// Short AI-generated chat title.
	Title string `json:"title,omitempty" path:"title,omitempty" url:"title,omitempty"`
	// User ID.
	UserId int64 `json:"user_id,omitempty" path:"user_id,omitempty" url:"user_id,omitempty"`
	// AI Task ID. Present when the conversation was started by an AI Task.
	AiTaskId int64 `json:"ai_task_id,omitempty" path:"ai_task_id,omitempty" url:"ai_task_id,omitempty"`
	// Workspace ID. `0` means the default workspace.
	WorkspaceId int64 `json:"workspace_id,omitempty" path:"workspace_id,omitempty" url:"workspace_id,omitempty"`
	// Most recent chat activity date/time.
	LastActiveAt *time.Time `json:"last_active_at,omitempty" path:"last_active_at,omitempty" url:"last_active_at,omitempty"`
	// Chat session creation date/time.
	CreatedAt *time.Time `json:"created_at,omitempty" path:"created_at,omitempty" url:"created_at,omitempty"`
	// Visible conversation messages in this chat session. For performance reasons, this is not provided when listing chat sessions.
	Messages []ChatMessage `json:"messages,omitempty" path:"messages,omitempty" url:"messages,omitempty"`
}

// Identifier returns the resource ID.
func (c ChatSession) Identifier() interface{} {
	return c.Id
}

// ChatSessionCollection is a list of ChatSession resources.
type ChatSessionCollection []ChatSession

// ChatSessionListParams contains the request parameters for GET /chat_sessions.
type ChatSessionListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `id`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `ai_task_id`, `user_id` or `workspace_id`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// ChatSessionFindParams contains the request parameters for GET /chat_sessions/{id}.
type ChatSessionFindParams struct {
	// Chat Session ID. Required; sent in the request path.
	Id string `url:"-,omitempty" json:"-,omitempty" path:"id"`
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (c *ChatSession) UnmarshalJSON(data []byte) error {
	type chatSession ChatSession
	var v chatSession
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*c = ChatSession(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (c *ChatSessionCollection) UnmarshalJSON(data []byte) error {
	type chatSessions ChatSessionCollection
	var v chatSessions
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*c = ChatSessionCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (c *ChatSessionCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*c))
	for i, v := range *c {
		ret[i] = v
	}

	return &ret
}
