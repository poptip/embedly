package embedly

const (
	TypeHTML  = "html"
	TypeText  = "text"
	TypeImage = "image"
	TypeVideo = "video"
	TypeAudio = "audio"
	TypeRSS   = "rss"
	TypeXML   = "xml"
	TypeAtom  = "atom"
	TypeJSON  = "json"
	TypePPT   = "ptt"
	TypeLink  = "link"
	TypeError = "error"
)

type Options struct {
	MaxWidth     int
	MaxHeight    int
	Width        int
	Words        int
	Chars        int
	WMode        bool
	AllowScripts bool
	NoStyle      bool
	Autoplay     bool
	VideoSrc     bool
	Frame        bool
	Secure       bool
}

type Response struct {
	OriginalURL     string    `json:"original_url"`
	URL             string    `json:"url"`
	Type            string    `json:"type"`
	ErrorMessage    string    `json:"error_message,omitempty"`
	ErrorCode       int       `json:"error_code"`
	Safe            bool      `json:"safe"`
	SafeType        string    `json:"safe_type,omitempty"`
	SafeMessage     string    `json:"safe_message,omitempty"`
	CacheAge        int       `json:"cache_age,omitempty"`
	ProviderName    string    `json:"provider_name"`
	ProviderURL     string    `json:"provider_url"`
	ProviderDisplay string    `json:"provider_display"`
	FaviconURL      string    `json:"favicon_url"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Authors         []Author  `json:"authors"`
	Media           Media     `json:"media"`
	Published       int64     `json:"published"`
	Offset          int64     `json:"offset"`
	Lead            string    `json:"lead"`
	Content         string    `json:"content"`
	Keywords        []Keyword `json:"keywords"`
	Entities        []Entity  `json:"entities"`
	Images          []Image   `json:"images"`
}

type Author struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Media struct {
	Type   string `json:"type"`
	URL    string `json:"url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
	HTML   string `json:"html,omitempty"`
}

type Keyword struct {
	Score int    `json:"score"`
	Name  string `json:"name"`
}

type Entity struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
}

type Image struct {
	Caption string  `json:"caption"`
	URL     string  `json:"url"`
	Width   int     `json:"width"`
	Height  int     `json:"height"`
	Colors  []Color `json:"colors"`
	Entropy float64 `json:"entropy"`
	Size    int     `json:"size"`
}

type Color struct {
	Color  []int   `json:"color"`
	Weight float64 `json:"weight"`
}
