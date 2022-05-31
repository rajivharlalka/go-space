package utils

type Imgur_data struct {
	Data struct {
		ID          string        `json:"id"`
		Title       interface{}   `json:"title"`
		Description interface{}   `json:"description"`
		Datetime    int           `json:"datetime"`
		Type        string        `json:"type"`
		Animated    bool          `json:"animated"`
		Width       int           `json:"width"`
		Height      int           `json:"height"`
		Size        int           `json:"size"`
		Views       int           `json:"views"`
		Bandwidth   int           `json:"bandwidth"`
		Vote        interface{}   `json:"vote"`
		Favorite    bool          `json:"favorite"`
		Nsfw        interface{}   `json:"nsfw"`
		Section     interface{}   `json:"section"`
		AccountURL  interface{}   `json:"account_url"`
		AccountID   int           `json:"account_id"`
		IsAd        bool          `json:"is_ad"`
		InMostViral bool          `json:"in_most_viral"`
		HasSound    bool          `json:"has_sound"`
		Tags        []interface{} `json:"tags"`
		AdType      int           `json:"ad_type"`
		AdURL       string        `json:"ad_url"`
		Edited      string        `json:"edited"`
		InGallery   bool          `json:"in_gallery"`
		Deletehash  string        `json:"deletehash"`
		Name        string        `json:"name"`
		Link        string        `json:"link"`
		Mp4         string        `json:"mp4"`
		Gifv        string        `json:"gifv"`
		Hls         string        `json:"hls"`
		Mp4Size     int           `json:"mp4_size"`
		Looping     bool          `json:"looping"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}
