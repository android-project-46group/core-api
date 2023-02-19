package model

//nolint:tagliatelle
type Member struct {
	ID         int     `json:"id"`
	Group      string  `json:"group,omitempty"`
	Name       string  `json:"name"`
	Birthday   string  `json:"birthday"`
	Height     float64 `json:"height"`
	BloodType  string  `json:"blood_type"`
	Generation string  `json:"generation"`
	BlogURL    string  `json:"blog_url"`
	ImgURL     string  `json:"img_url"`
	LeftAt     string  `json:"left_at,omitempty"`
}
