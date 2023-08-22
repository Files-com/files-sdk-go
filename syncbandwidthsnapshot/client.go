package sync_bandwidth_snapshot

import (
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	lib "github.com/Files-com/files-sdk-go/v2/lib"
)

type Client struct {
	files_sdk.Config
}

func (c *Client) Create(params files_sdk.SyncBandwidthSnapshotCreateParams, opts ...files_sdk.RequestResponseOption) (syncBandwidthSnapshot files_sdk.SyncBandwidthSnapshot, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sync_bandwidth_snapshots", Params: params, Entity: &syncBandwidthSnapshot}, opts...)
	return
}

func Create(params files_sdk.SyncBandwidthSnapshotCreateParams, opts ...files_sdk.RequestResponseOption) (syncBandwidthSnapshot files_sdk.SyncBandwidthSnapshot, err error) {
	return (&Client{}).Create(params, opts...)
}

func (c *Client) CreateBatch(opts ...files_sdk.RequestResponseOption) (syncBandwidthSnapshotCollection files_sdk.SyncBandwidthSnapshotCollection, err error) {
	err = files_sdk.Resource(c.Config, lib.Resource{Method: "POST", Path: "/sync_bandwidth_snapshots/create_batch", Params: lib.Interface(), Entity: &syncBandwidthSnapshotCollection}, opts...)
	return
}

func CreateBatch(opts ...files_sdk.RequestResponseOption) (syncBandwidthSnapshotCollection files_sdk.SyncBandwidthSnapshotCollection, err error) {
	return (&Client{}).CreateBatch(opts...)
}
