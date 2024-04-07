package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CurrentCondition struct {
	FeelsLikeC       string `json:"FeelsLikeC"`
	FeelsLikeF       string `json:"FeelsLikeF"`
	CloudCover       string `json:"cloudcover"`
	Humidity         string `json:"humidity"`
	LocalObsDateTime string `json:"localObsDateTime"`
	ObservationTime  string `json:"observation_time"`
	PrecipInches     string `json:"precipInches"`
	PrecipMM         string `json:"precipMM"`
	Pressure         string `json:"pressure"`
	PressureInches   string `json:"pressureInches"`
	TempC            string `json:"temp_C"`
	TempF            string `json:"temp_F"`
	UVIndex          string `json:"uvIndex"`
	Visibility       string `json:"visibility"`
	VisibilityMiles  string `json:"visibilityMiles"`
	WeatherCode      string `json:"weatherCode"`
	Winddir16        string `json:"winddir16Point"`
	WinddirDegree    string `json:"winddirDegree"`
	WindspeedKmph    string `json:"windspeedKmph"`
}

type NearestArea struct {
	AreaName []struct {
		Value string `json:"value"`
	} `json:"areaName"`
	Country []struct {
		Value string `json:"value"`
	} `json:"country"`
	WeatherUrl []struct {
		Value string `json:"value"`
	} `json:"weatherUrl"`
}

type Astronomy struct {
	MoonIllumination string `json:"moon_illumination"`
	MoonPhase        string `json:"moon_phase"`
	Moonrise         string `json:"moonrise"`
	Moonset          string `json:"moonset"`
	Sunrise          string `json:"sunrise"`
	Sunset           string `json:"sunset"`
}

type Weather struct {
	Astronomy   []Astronomy `json:"astronomy"`
	AvgTempC    string      `json:"avgtempC"`
	AvgTempF    string      `json:"avgtempF"`
	Date        string      `json:"date"`
	MaxTempC    string      `json:"maxtempC"`
	MaxTempF    string      `json:"maxtempF"`
	MinTempC    string      `json:"mintempC"`
	MinTempF    string      `json:"mintempF"`
	SunHour     string      `json:"sunHour"`
	TotalSnowCm string      `json:"totalSnow_cm"`
	UVIndex     string      `json:"uvIndex"`
}

type WeatherData struct {
	CurrentCondition []CurrentCondition `json:"current_condition"`
	NearestArea      []NearestArea      `json:"nearest_area"`
	Weather          []Weather          `json:"weather"`
}

// Waybar json structure documentation: https://github.com/Alexays/Waybar/wiki/Module:-Custom#return-type
type WaybarWeather struct {
	Text    string `json:"text"`    // Displayed value
	Tooltip string `json:"tooltip"` // Displayed info when the cursor is over the displayed value
}

func get_weather_emoji(code string) string {
	// Weather codes available here: https://github.com/chubin/wttr.in/blob/master/lib/constants.py
	WWO_CODE := map[string]string{
		"113": "☀️",
		"116": "🌤️",
		"119": "☁️",
		"122": "☁️☁️",
		"143": "🌫️",
		"176": "🌦️",
		"179": "🌧️",
		"182": "🌨️",
		"185": "🌨️",
		"200": "⛈️",
		"227": "🌨️",
		"230": "🌨️",
		"248": "🌫️",
		"260": "🌫️",
		"263": "🌦️",
		"266": "🌧️",
		"281": "🌨️",
		"284": "🌨️",
		"293": "🌧️",
		"296": "🌧️",
		"299": "🌧️",
		"302": "🌧️",
		"305": "🌧️",
		"308": "🌧️",
		"311": "🌨️",
		"314": "🌨️",
		"317": "🌨️",
		"320": "🌨️",
		"323": "🌨️",
		"326": "🌨️",
		"329": "🌨️",
		"332": "🌨️",
		"335": "🌨️",
		"338": "🌨️",
		"350": "🌨️",
		"353": "🌦️",
		"356": "🌧️",
		"359": "🌧️",
		"362": "🌧️",
		"365": "🌧️",
		"368": "🌨️",
		"371": "🌨️",
		"374": "🌧️",
		"377": "🌨️",
		"386": "⛈️",
		"389": "⛈️",
		"392": "⛈️",
		"395": "🌨️",
	}
	// Check if the condition exists in the map
	emoji, found := WWO_CODE[code]
	if found {
		return emoji
	}
	return "❓" // Default emoji
}

func get_moon_phase_emoji(moon_phase string) string {
	MOON_PHASES := map[string]string{
		"New":             "🌑",
		"Waxing Crescent": "🌒",
		"First Quarter":   "🌓",
		"Waxing Gibbous":  "🌔",
		"Full":            "🌕",
		"Waning Gibbous":  "🌖",
		"Third Quarter":   "🌗",
		"Waning Crescent": "🌘",
	}
	emoji, found := MOON_PHASES[moon_phase]
	if found {
		return emoji
	}
	return "❓"
}

func check_args_length(arg string, max_len int) int {
	if len(arg) > max_len {
		return -1
	}
	return 0
}

func main() {
	lang_ptr := flag.String("lang", "en", "Language")
	city_ptr := flag.String("city", "", "City")
	unit_ptr := flag.String("unit", "C", "Temperature Unit")

	flag.Parse()

	lang := *lang_ptr
	city := *city_ptr
	unit := *unit_ptr

	if check_args_length(lang, 2) == -1 {
		fmt.Fprintf(os.Stderr, "The --lang argument can have a maximum of 2 characters (e.g., en).")
		os.Exit(1)
	}

	if check_args_length(unit, 1) == -1 {
		fmt.Fprintf(os.Stderr, "The --unit argument can have a maximum of 1 character (e.g., C).")
		os.Exit(1)
	}

	weather_url := "https://wttr.in/" + city + "?format=j2&lang=" + lang
	resp, err := http.Get(weather_url)
	if err != nil {
		panic(err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	var weather_data WeatherData
	json.Unmarshal(body, &weather_data)

	var temperature string
	switch unit {
	case "C":
		temperature = weather_data.CurrentCondition[0].FeelsLikeC
		break
	case "F":
		temperature = weather_data.CurrentCondition[0].FeelsLikeF
		break
	}
	weather_emoji := get_weather_emoji(weather_data.CurrentCondition[0].WeatherCode)

	var waybar_weather WaybarWeather
	waybar_weather.Text = weather_emoji + temperature + "°" + unit
	waybar_weather.Tooltip =
		"City: " + weather_data.NearestArea[0].AreaName[0].Value + "\n" +
			"Wind: " + weather_data.CurrentCondition[0].WindspeedKmph + " km/h \n" +
			"Moon: " + get_moon_phase_emoji(weather_data.Weather[0].Astronomy[0].MoonPhase)

	json_waybar, err := json.Marshal(waybar_weather)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println(string(json_waybar))
}
