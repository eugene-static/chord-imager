package dto

// NameInfo stores data from frontend
type NameInfo struct {
	Name    string `json:"name,omitempty"`
	Pattern string `json:"pattern"`
	Fret    int    `json:"fret"`
	Capo    bool   `json:"capo"`
}
type NameRequest struct {
	Event string    `json:"event"`
	Info  *NameInfo `json:"info"`
}

// Name stores data sending to frontend
type Name struct {
	Root     string `json:"root,omitempty"`
	Quality  string `json:"quality,omitempty"`
	Extended string `json:"extended,omitempty"`
	Altered  string `json:"altered,omitempty"`
	Omitted  string `json:"omitted,omitempty"`
}

// NameResponse stores main matching name and its inversions
type NameResponse struct {
	BaseName   *Name  `json:"base_name,omitempty"`
	Variations []Name `json:"variations,omitempty"`
}

type TabResponse struct {
	Tab string `json:"tab"`
}

type ImageResponse struct {
	URL string `json:"url"`
}
