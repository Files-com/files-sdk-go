package files_sdk

import (
	"encoding/json"

	lib "github.com/Files-com/files-sdk-go/v3/lib"
)

// App is a Files.com API resource.
type App struct {
	// The type of the App.
	AppType string `json:"app_type,omitempty" path:"app_type,omitempty" url:"app_type,omitempty"`
	// Collection of named links to documentation.
	DocumentationLinks interface{} `json:"documentation_links,omitempty" path:"documentation_links,omitempty" url:"documentation_links,omitempty"`
	// Long description for the in-App landing page.
	ExtendedDescription string `json:"extended_description,omitempty" path:"extended_description,omitempty" url:"extended_description,omitempty"`
	// Long form description of the App.
	ExtendedDescriptionForMarketingSite string `json:"extended_description_for_marketing_site,omitempty" path:"extended_description_for_marketing_site,omitempty" url:"extended_description_for_marketing_site,omitempty"`
	// Link to external homepage.
	ExternalHomepageUrl string `json:"external_homepage_url,omitempty" path:"external_homepage_url,omitempty" url:"external_homepage_url,omitempty"`
	// Is featured on the App listing?
	Featured *bool `json:"featured,omitempty" path:"featured,omitempty" url:"featured,omitempty"`
	// Associated Folder Behavior type, if any.
	FolderBehaviorType string `json:"folder_behavior_type,omitempty" path:"folder_behavior_type,omitempty" url:"folder_behavior_type,omitempty"`
	// App icon.
	IconUrl string `json:"icon_url,omitempty" path:"icon_url,omitempty" url:"icon_url,omitempty"`
	// Logo thumbnail for the App.
	LogoThumbnailUrl string `json:"logo_thumbnail_url,omitempty" path:"logo_thumbnail_url,omitempty" url:"logo_thumbnail_url,omitempty"`
	// Full size logo for the App.
	LogoUrl string `json:"logo_url,omitempty" path:"logo_url,omitempty" url:"logo_url,omitempty"`
	// Marketing introdution of the App.
	MarketingIntro string `json:"marketing_intro,omitempty" path:"marketing_intro,omitempty" url:"marketing_intro,omitempty"`
	// Marketing video page.
	MarketingYoutubeUrl string `json:"marketing_youtube_url,omitempty" path:"marketing_youtube_url,omitempty" url:"marketing_youtube_url,omitempty"`
	// Name of the App.
	Name string `json:"name,omitempty" path:"name,omitempty" url:"name,omitempty"`
	// Package manager install command.
	PackageManagerInstallCommand string `json:"package_manager_install_command,omitempty" path:"package_manager_install_command,omitempty" url:"package_manager_install_command,omitempty"`
	// Associated Remote Server type, if any.
	RemoteServerType string `json:"remote_server_type,omitempty" path:"remote_server_type,omitempty" url:"remote_server_type,omitempty"`
	// Screenshots of the App.
	ScreenshotListUrls []string `json:"screenshot_list_urls,omitempty" path:"screenshot_list_urls,omitempty" url:"screenshot_list_urls,omitempty"`
	// Link to SDK installation instructions.
	SdkInstallationInstructionsLink string `json:"sdk_installation_instructions_link,omitempty" path:"sdk_installation_instructions_link,omitempty" url:"sdk_installation_instructions_link,omitempty"`
	// Short description of the App.
	ShortDescription string `json:"short_description,omitempty" path:"short_description,omitempty" url:"short_description,omitempty"`
	// Associated SSO Strategy type, if any.
	SsoStrategyType string `json:"sso_strategy_type,omitempty" path:"sso_strategy_type,omitempty" url:"sso_strategy_type,omitempty"`
	// Associated SIEM type, if any.
	SiemType string `json:"siem_type,omitempty" path:"siem_type,omitempty" url:"siem_type,omitempty"`
	// Tutorial video page.
	TutorialYoutubeUrl string `json:"tutorial_youtube_url,omitempty" path:"tutorial_youtube_url,omitempty" url:"tutorial_youtube_url,omitempty"`
}

// Identifier no path or id

// AppCollection is a list of App resources.
type AppCollection []App

// AppListParams contains the request parameters for GET /apps.
type AppListParams struct {
	// If set, sort records by the specified field in either `asc` or `desc` direction. Valid fields are `name` and `app_type`.
	SortBy interface{} `url:"sort_by,omitempty" json:"sort_by,omitempty" path:"sort_by"`
	// If set, return records where the specified field is equal to the supplied value. Valid fields are `name` and `app_type`. Valid field combinations are `[ name, app_type ]`.
	Filter interface{} `url:"filter,omitempty" json:"filter,omitempty" path:"filter"`
	// If set, return records where the specified field is prefixed by the supplied value. Valid fields are `name`.
	FilterPrefix interface{} `url:"filter_prefix,omitempty" json:"filter_prefix,omitempty" path:"filter_prefix"`
	// ListParams controls cursor-based pagination: PerPage and Cursor are sent to the API, and MaxPages, which is not sent, limits the number of pages requested when positive; zero allows all pages.
	ListParams
}

// UnmarshalJSON decodes an API resource. A decoding error leaves the receiver unchanged.
func (a *App) UnmarshalJSON(data []byte) error {
	type app App
	var v app
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, map[string]interface{}{})
	}

	*a = App(v)
	return nil
}

// UnmarshalJSON decodes a list of API resources. A decoding error leaves the receiver unchanged.
func (a *AppCollection) UnmarshalJSON(data []byte) error {
	type apps AppCollection
	var v apps
	if err := json.Unmarshal(data, &v); err != nil {
		return lib.ErrorWithOriginalResponse{}.ProcessError(data, err, []map[string]interface{}{})
	}

	*a = AppCollection(v)
	return nil
}

// ToSlice returns a new slice containing the resources as interface values.
func (a *AppCollection) ToSlice() *[]interface{} {
	ret := make([]interface{}, len(*a))
	for i, v := range *a {
		ret[i] = v
	}

	return &ret
}
