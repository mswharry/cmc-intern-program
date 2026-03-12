package model

type BatchCreateAssetsRequest struct {
	Assets []struct {
		Name   string `json:"name"`
		Type   string `json:"type"`
		Status string `json:"status,omitempty"`
	} `json:"assets"`
}

type BatchCreateAssetsResponse struct {
	Created int 	  `json:"created"`
	IDs     []string  `json:"ids"`
}

type BatchDeleteAssetsRequest struct {
	IDs []string `json:"ids"`
}

type BatchDeleteAssetsResponse struct {
	Deleted int `json:"deleted"`
	NotFound int `json:"not_found"`
}