package msg

type DownloadResultMsg struct {
	URL string `json:"url"`
	Redirects []string `json:"redirects"`
}