package omgo

import "time"

// DailyEnsembleData contains daily ensemble forecast data from multiple model members.
// Each field is a 2D slice: first dimension is the ensemble member, second is time.
type DailyEnsembleData struct {
	// Time contains timestamps for each day (at 00:00).
	Times []time.Time `json:"-"`

	// Temperature
	Temperature2mMax  [][]float64 `json:"temperature_2m_max,omitempty"`
	Temperature2mMin  [][]float64 `json:"temperature_2m_min,omitempty"`
	Temperature2mMean [][]float64 `json:"temperature_2m_mean,omitempty"`

	// Apparent temperature
	ApparentTemperatureMax  [][]float64 `json:"apparent_temperature_max,omitempty"`
	ApparentTemperatureMin  [][]float64 `json:"apparent_temperature_min,omitempty"`
	ApparentTemperatureMean [][]float64 `json:"apparent_temperature_mean,omitempty"`

	// Precipitation
	PrecipitationSum   []float64 `json:"precipitation_sum,omitempty"`
	RainSum            []float64 `json:"rain_sum,omitempty"`
	ShowersSum         []float64 `json:"showers_sum,omitempty"`
	SnowfallSum        []float64 `json:"snowfall_sum,omitempty"`
	PrecipitationHours []float64 `json:"precipitation_hours,omitempty"`

	// Precipitation probability
	PrecipitationProbabilityMax  []float64 `json:"precipitation_probability_max,omitempty"`
	PrecipitationProbabilityMin  []float64 `json:"precipitation_probability_min,omitempty"`
	PrecipitationProbabilityMean []float64 `json:"precipitation_probability_mean,omitempty"`

	// Wind
	WindSpeed10mMax          []float64 `json:"wind_speed_10m_max,omitempty"`
	WindGusts10mMax          []float64 `json:"wind_gusts_10m_max,omitempty"`
	WindDirection10mDominant []float64 `json:"wind_direction_10m_dominant,omitempty"`

	// Radiation (MJ/m²)
	ShortwaveRadiationSum []float64 `json:"shortwave_radiation_sum,omitempty"`

	// Evapotranspiration
	ET0FAOEvapotranspiration []float64 `json:"et0_fao_evapotranspiration,omitempty"`

	// UV Index
	UVIndexMax         []float64 `json:"uv_index_max,omitempty"`
	UVIndexClearSkyMax []float64 `json:"uv_index_clear_sky_max,omitempty"`
}

// MemberCount returns the number of ensemble members (first dimension of 2D arrays).
// Returns 0 if no ensemble data is available.
func (d *DailyEnsembleData) MemberCount() int {
	if len(d.Temperature2mMax) > 0 {
		return len(d.Temperature2mMax)
	}
	if len(d.Temperature2mMin) > 0 {
		return len(d.Temperature2mMin)
	}
	if len(d.Temperature2mMean) > 0 {
		return len(d.Temperature2mMean)
	}
	return 0
}

// DayCount returns the number of forecast days (second dimension).
// Returns 0 if no data or no members.
func (d *DailyEnsembleData) DayCount() int {
	if len(d.Temperature2mMax) > 0 && len(d.Temperature2mMax[0]) > 0 {
		return len(d.Temperature2mMax[0])
	}
	return 0
}

// CountMembersAboveThreshold counts how many members exceed the given threshold
// for the specified metric on the specified day index.
//
// This is the core method for weather trading: if 20 out of 31 members
// predict temperature > 80°F, your probability is 20/31 ≈ 64.5%.
func (d *DailyEnsembleData) CountMembersAboveThreshold(day int, threshold float64, metric DailyMetric) int {
	data := d.getMetric(metric)
	if data == nil || len(data) == 0 {
		return 0
	}

	count := 0
	for _, member := range data {
		if day < len(member) && member[day] > threshold {
			count++
		}
	}
	return count
}

// CountMembersBelowThreshold counts how many members are below the threshold.
func (d *DailyEnsembleData) CountMembersBelowThreshold(day int, threshold float64, metric DailyMetric) int {
	data := d.getMetric(metric)
	if data == nil || len(data) == 0 {
		return 0
	}

	count := 0
	for _, member := range data {
		if day < len(member) && member[day] < threshold {
			count++
		}
	}
	return count
}

// ProbabilityAbove returns the fraction of members exceeding the threshold
// for the specified metric on the specified day (0.0 to 1.0).
func (d *DailyEnsembleData) ProbabilityAbove(day int, threshold float64, metric DailyMetric) float64 {
	total := d.MemberCount()
	if total == 0 {
		return 0.0
	}
	return float64(d.CountMembersAboveThreshold(day, threshold, metric)) / float64(total)
}

// ProbabilityBelow returns the fraction of members below the threshold.
func (d *DailyEnsembleData) ProbabilityBelow(day int, threshold float64, metric DailyMetric) float64 {
	total := d.MemberCount()
	if total == 0 {
		return 0.0
	}
	return float64(d.CountMembersBelowThreshold(day, threshold, metric)) / float64(total)
}

func (d *DailyEnsembleData) getMetric(metric DailyMetric) [][]float64 {
	switch metric {
	case DailyTemperature2mMax:
		return d.Temperature2mMax
	case DailyTemperature2mMin:
		return d.Temperature2mMin
	case DailyTemperature2mMean:
		return d.Temperature2mMean
	case DailyApparentTemperatureMax:
		return d.ApparentTemperatureMax
	case DailyApparentTemperatureMin:
		return d.ApparentTemperatureMin
	case DailyApparentTemperatureMean:
		return d.ApparentTemperatureMean
	default:
		return nil
	}
}
