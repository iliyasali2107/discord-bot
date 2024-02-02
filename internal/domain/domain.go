package domain

type Location struct {
	Name           string  `json:"name"`
	Region         string  `json:"region"`
	Country        string  `json:"country"`
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	TzID           string  `json:"tz_id"`
	LocaltimeEpoch int64   `json:"localtime_epoch"`
	Localtime      string  `json:"localtime"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}

type Current struct {
	LastUpdatedEpoch int64     `json:"last_updated_epoch"`
	LastUpdated      string    `json:"last_updated"`
	TempC            float64   `json:"temp_c"`
	TempF            float64   `json:"temp_f"`
	IsDay            int       `json:"is_day"`
	Condition        Condition `json:"condition"`
	WindMPH          float64   `json:"wind_mph"`
	WindKPH          float64   `json:"wind_kph"`
	WindDegree       int       `json:"wind_degree"`
	WindDir          string    `json:"wind_dir"`
	PressureMB       float64   `json:"pressure_mb"`
	PressureIN       float64   `json:"pressure_in"`
	PrecipMM         float64   `json:"precip_mm"`
	PrecipIN         float64   `json:"precip_in"`
	Humidity         int       `json:"humidity"`
	Cloud            int       `json:"cloud"`
	FeelslikeC       float64   `json:"feelslike_c"`
	FeelslikeF       float64   `json:"feelslike_f"`
	VisKM            float64   `json:"vis_km"`
	VisMiles         float64   `json:"vis_miles"`
	UV               float64   `json:"uv"`
	GustMPH          float64   `json:"gust_mph"`
	GustKPH          float64   `json:"gust_kph"`
}

type WeatherData struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
}
