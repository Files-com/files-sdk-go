# PartnerConnection

`files_sdk.PartnerConnection` is a Files.com API resource. Its client package is
`github.com/Files-com/files-sdk-go/v3/partnerconnection`, named `partner_connection`.

This package has no API operations of its own.

## Model

`files_sdk.PartnerConnection` has these fields from the API schema:

| Field | Go type | JSON key | Description |
| --- | --- | --- | --- |
| `Id` | `int64` | `id` | Relationship ID used with DELETE /partner_sites/:id to disconnect. |
| `Role` | `string` | `role` | This Partner's role in this connection. A host shares local files with the connected site; a guest accesses files shared by the connected site. |
| `SiteId` | `int64` | `site_id` | ID of the connected site. |
| `SiteName` | `string` | `site_name` | Name of the connected site. |
| `MountPath` | `string` | `mount_path` | File API path to the connected Host's mount on this site when role is guest. Null when role is host. File access remains subject to the caller's permissions and the Host Partner's grants. |
