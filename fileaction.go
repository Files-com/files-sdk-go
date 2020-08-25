package files_sdk

import (
	"encoding/json"
)

type FileAction struct {
}

type FileActionCollection []FileAction

// Copy file/folder
type FileActionCopyParams struct {
	Path        string `url:"-,omitempty" required:"true"`
	Destination string `url:"destination,omitempty" required:"true"`
	Structure   *bool  `url:"structure,omitempty" required:"false"`
}

// Move file/folder
type FileActionMoveParams struct {
	Path        string `url:"-,omitempty" required:"true"`
	Destination string `url:"destination,omitempty" required:"true"`
}

// Begin file upload
type FileActionBeginUploadParams struct {
	Path         string `url:"-,omitempty" required:"true"`
	MkdirParents *bool  `url:"mkdir_parents,omitempty" required:"false"`
	Part         int    `url:"part,omitempty" required:"false"`
	Parts        int    `url:"parts,omitempty" required:"false"`
	Ref          string `url:"ref,omitempty" required:"false"`
	Restart      int    `url:"restart,omitempty" required:"false"`
	WithRename   *bool  `url:"with_rename,omitempty" required:"false"`
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
