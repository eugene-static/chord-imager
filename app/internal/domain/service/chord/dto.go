package chord

type NameInfo struct {
	Pattern string
	Fret    int
	Capo    bool
}

type TabInfo struct {
	Name string
	Info NameInfo
}

type ImageInfo struct {
	Name string
	Path string
	Info NameInfo
}

type Name struct {
	Root     string
	Quality  string
	Extended string
	Altered  string
	Omitted  string
}

type NameResponse struct {
	BaseName   *Name
	Variations []*Name
}

type TabResponse struct {
	Tab string
}

type ImageResponse struct {
	URL string
}
