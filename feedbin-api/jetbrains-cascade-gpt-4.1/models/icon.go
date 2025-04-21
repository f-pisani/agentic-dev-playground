package models

type Icon struct {
	Host string `json:"host"`
	URL  string `json:"url"`
}

type IconsResponse []Icon
