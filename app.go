package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type App struct {
	Name                string `json:"name,omitempty"`
	ExtendedDescription string `json:"extended_description,omitempty"`
	ShortDescription    string `json:"short_description,omitempty"`
	DocumentationLinks  string `json:"documentation_links,omitempty"`
	IconUrl             string `json:"icon_url,omitempty"`
	LogoUrl             string `json:"logo_url,omitempty"`
	ScreenshotListUrls  string `json:"screenshot_list_urls,omitempty"`
	LogoThumbnailUrl    string `json:"logo_thumbnail_url,omitempty"`
	SsoStrategyType     string `json:"sso_strategy_type,omitempty"`
	RemoteServerType    string `json:"remote_server_type,omitempty"`
	FolderBehaviorType  string `json:"folder_behavior_type,omitempty"`
	ExternalHomepageUrl string `json:"external_homepage_url,omitempty"`
	MarketingYoutubeUrl string `json:"marketing_youtube_url,omitempty"`
	TutorialYoutubeUrl  string `json:"tutorial_youtube_url,omitempty"`
	AppType             string `json:"app_type,omitempty"`
	Featured            *bool  `json:"featured,omitempty"`
}

type AppCollection []App

type AppListParams struct {
	Cursor     string          `url:"cursor,omitempty" required:"false"`
	PerPage    int64           `url:"per_page,omitempty" required:"false"`
	SortBy     json.RawMessage `url:"sort_by,omitempty" required:"false"`
	Filter     json.RawMessage `url:"filter,omitempty" required:"false"`
	FilterGt   json.RawMessage `url:"filter_gt,omitempty" required:"false"`
	FilterGteq json.RawMessage `url:"filter_gteq,omitempty" required:"false"`
	FilterLike json.RawMessage `url:"filter_like,omitempty" required:"false"`
	FilterLt   json.RawMessage `url:"filter_lt,omitempty" required:"false"`
	FilterLteq json.RawMessage `url:"filter_lteq,omitempty" required:"false"`
	lib.ListParams
}

func (a *App) UnmarshalJSON(data []byte) error {
	type app App
	var v app
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = App(v)
	return nil
}

func (a *AppCollection) UnmarshalJSON(data []byte) error {
	type apps []App
	var v apps
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*a = AppCollection(v)
	return nil
}

func (a *AppCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
