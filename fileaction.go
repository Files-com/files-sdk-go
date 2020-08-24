package files_sdk

import (
	"encoding/json"
)

type FileAction struct {
}

type FileActionCollection []FileAction

// Copy file/folder
type FileActionCopyParams struct {
	Path        string `url:"-,omitempty"`
	Destination string `url:"destination,omitempty"`
	Structure   *bool  `url:"structure,omitempty"`
}

// Move file/folder
type FileActionMoveParams struct {
	Path        string `url:"-,omitempty"`
	Destination string `url:"destination,omitempty"`
}

// Begin file upload
type FileActionBeginUploadParams struct {
	Path         string `url:"-,omitempty"`
	MkdirParents *bool  `url:"mkdir_parents,omitempty"`
	Part         int    `url:"part,omitempty"`
	Parts        int    `url:"parts,omitempty"`
	Ref          string `url:"ref,omitempty"`
	Restart      int    `url:"restart,omitempty"`
	WithRename   *bool  `url:"with_rename,omitempty"`
}

func (f *FileAction) UnmarshalJSON(data []byte) error {
	type fileAction FileAction
	var v fileAction
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileAction(v)
	return nil
}

func (f *FileActionCollection) UnmarshalJSON(data []byte) error {
	type fileActions []FileAction
	var v fileActions
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*f = FileActionCollection(v)
	return nil
}
