package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type App struct {
	Name                string          `json:"name,omitempty" path:"name"`
	ExtendedDescription string          `json:"extended_description,omitempty" path:"extended_description"`
	ShortDescription    string          `json:"short_description,omitempty" path:"short_description"`
	DocumentationLinks  json.RawMessage `json:"documentation_links,omitempty" path:"documentation_links"`
	IconUrl             string          `json:"icon_url,omitempty" path:"icon_url"`
	LogoUrl             string          `json:"logo_url,omitempty" path:"logo_url"`
	ScreenshotListUrls  []string        `json:"screenshot_list_urls,omitempty" path:"screenshot_list_urls"`
	LogoThumbnailUrl    string          `json:"logo_thumbnail_url,omitempty" path:"logo_thumbnail_url"`
	SsoStrategyType     string          `json:"sso_strategy_type,omitempty" path:"sso_strategy_type"`
	RemoteServerType    string          `json:"remote_server_type,omitempty" path:"remote_server_type"`
	FolderBehaviorType  string          `json:"folder_behavior_type,omitempty" path:"folder_behavior_type"`
	ExternalHomepageUrl string          `json:"external_homepage_url,omitempty" path:"external_homepage_url"`
	MarketingYoutubeUrl string          `json:"marketing_youtube_url,omitempty" path:"marketing_youtube_url"`
	TutorialYoutubeUrl  string          `json:"tutorial_youtube_url,omitempty" path:"tutorial_youtube_url"`
	AppType             string          `json:"app_type,omitempty" path:"app_type"`
	Featured            *bool           `json:"featured,omitempty" path:"featured"`
}

// Identifier no path or id

type AppCollection []App

type AppListParams struct {
	SortBy       json.RawMessage `url:"sort_by,omitempty" required:"false" json:"sort_by,omitempty" path:"sort_by"`
	Filter       json.RawMessage `url:"filter,omitempty" required:"false" json:"filter,omitempty" path:"filter"`
	FilterPrefix json.RawMessage `url:"filter_prefix,omitempty" required:"false" json:"filter_prefix,omitempty" path:"filter_prefix"`
	ListParams
}

func (a *App) UnmarshalJSON(data []byte) error {
	type app App
	var v app
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = App(v)
	return nil
}

func (a *AppCollection) UnmarshalJSON(data []byte) error {
	type apps AppCollection
	var v apps
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
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
