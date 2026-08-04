package cmdHassio

import (
	"github.com/MickMake/GoSungrow/iSolarCloud/api"
	"github.com/MickMake/GoSungrow/iSolarCloud/api/GoStruct/valueTypes"
	"fmt"
	"github.com/MickMake/GoUnify/Only"
)


type EntityConfig struct {
	// Type          string
	Name          string
	SubName       string

	ParentId      string
	ParentName    string

	UniqueId      string
	FullId        string
	Units         string
	ValueName     string
	DeviceClass   string
	StateClass    string
	Icon          string

	Value         *valueTypes.UnitValue
	Point         *api.Point
	ValueTemplate string

	UpdateFreq    string
	LastReset     string
	LastResetValueTemplate string

	IgnoreUpdate  bool

	EntityCategory    string
	EnabledByDefault  *bool
	DeviceGroup       string

	haType        string
	Options       []string
}

func (config *EntityConfig) FixConfig() {

	for range Only.Once {
		// mdi:power-socket-au
		// mdi:solar-power
		// mdi:home-lightning-bolt-outline
		// mdi:transmission-tower
		// mdi:transmission-tower-export
		// mdi:transmission-tower-import
		// mdi:transmission-tower-off
		// mdi:home-battery-outline
		// mdi:lightning-bolt
		// mdi:check-circle-outline | mdi:arrow-right-bold
		// mdi:transmission-tower

		// Set default ValueTemplate
		switch {
			case config.Point.IsBool():
				fallthrough
			case config.Value.IsBool():
				fallthrough
			case config.IsBinarySensor():
				config.ValueTemplate = SetDefault(config.ValueTemplate, "{{ value_json.value }}")

			case config.Value.IsFloat():
				if !config.Value.Valid {
					config.IgnoreUpdate = true
				}
				cnv := "| float"
				if config.Value.String() == "" {
					cnv = ""
				}
				vj := "value"
				if config.ValueName != "" {
					vj = config.ValueName
				}
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s %s }}", vj, cnv))

			case config.Value.IsTypeDateTime():
				vj := "value"
				if config.ValueName != "" {
					vj = config.ValueName
				}
				value, _, err := valueTypes.ParseDateTime(config.Value.String())
				if err == nil {
					config.Value.SetString(value.Local().Format(valueTypes.DateTimeFullLayout))
					config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | as_datetime }}", vj))
				} else {
					config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", vj))
				}

			case config.Value.IsInt():
				vj := "value"
				if config.ValueName != "" {
					vj = config.ValueName
				}
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s | int }}", vj))

			default:
				vj := "value"
				if config.ValueName != "" {
					vj = config.ValueName
				}
				config.ValueTemplate = SetDefault(config.ValueTemplate, fmt.Sprintf("{{ value_json.%s }}", vj))
		}

		// Set DeviceClass & Icon
		switch {
			case config.Point.IsBool():
				fallthrough
			case config.Value.IsBool():
				fallthrough
			case config.IsBinarySensor():
				config.DeviceClass = SetDefault(config.DeviceClass, "power")
				config.Icon = SetDefault(config.Icon, "mdi:check-circle-outline")

			case config.Units == "kWp":
				// Installed/design capacity (kilowatt-peak). HA has no device class
				// for kWp and rejects kWp together with device_class "power", so
				// leave the device class unset.
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:solar-panel-large")

			case config.Value.TypeValue == "Power":
				fallthrough
			case config.Units == "MW":
				fallthrough
			case config.Units == "kW":
				fallthrough
			case config.Units == "W":
				config.DeviceClass = SetDefault(config.DeviceClass, "power")
				config.Icon = SetDefault(config.Icon, "mdi:lightning-bolt")

			case config.Value.TypeValue == "Energy":
				fallthrough
			case config.Units == "MWh":
				fallthrough
			case config.Units == "kWh":
				fallthrough
			case config.Units == "Wh":
				config.DeviceClass = SetDefault(config.DeviceClass, "energy")
				config.Icon = SetDefault(config.Icon, "mdi:transmission-tower")

			case config.Units == "var":
				fallthrough
			case config.Units == "kvar":
				config.DeviceClass = SetDefault(config.DeviceClass, "reactive_power")
				config.Icon = SetDefault(config.Icon, "mdi:lightning-bolt")

			case config.Units == "VA":
				config.DeviceClass = SetDefault(config.DeviceClass, "apparent_power")
				config.Icon = SetDefault(config.Icon, "mdi:lightning-bolt")

			case config.Units == "Hz":
				config.DeviceClass = SetDefault(config.DeviceClass, "frequency")
				config.Icon = SetDefault(config.Icon, "mdi:sine-wave")

			case config.Units == "V":
				config.DeviceClass = SetDefault(config.DeviceClass, "voltage")
				config.Icon = SetDefault(config.Icon, "mdi:current-dc")

			case config.Units == "A":
				config.DeviceClass = SetDefault(config.DeviceClass, "current")
				config.Icon = SetDefault(config.Icon, "mdi:current-ac")

			case config.Units == "°F":
				fallthrough
			case config.Units == "F":
				fallthrough
			case config.Units == "℉":
				config.DeviceClass = SetDefault(config.DeviceClass, "temperature")
				config.Units = "°F"
				config.Icon = SetDefault(config.Icon, "mdi:thermometer")

			case config.Units == "°C":
				fallthrough
			case config.Units == "C":
				fallthrough
			case config.Units == "℃":
				config.DeviceClass = SetDefault(config.DeviceClass, "temperature")
				config.Units = "°C"
				config.Icon = SetDefault(config.Icon, "mdi:thermometer")

			case config.Icon == "mdi:home-battery-outline":
				fallthrough
			case config.Icon == "mdi:battery":
				config.DeviceClass = SetDefault(config.DeviceClass, "battery")
				// config.Icon = SetDefault(config.Icon, "mdi:percent") // mdi:home-battery-outline

			case config.Value.TypeValue == "Percent":
				fallthrough
			case config.Units == "%":
				// @TODO - Not supported in older versions of HA.
				// config.DeviceClass = SetDefault(config.DeviceClass, "battery")
				config.Icon = SetDefault(config.Icon, "mdi:percent") // mdi:home-battery-outline

			case config.Value.TypeValue == "DateTime":
				config.DeviceClass = SetDefault(config.DeviceClass, "timestamp") // date
				config.Icon = SetDefault(config.Icon, "mdi:clock-outline")

			case config.Units == "h":
				// config.DeviceClass = SetDefault(config.DeviceClass, "timestamp") // date
				config.Icon = SetDefault(config.Icon, "mdi:clock-outline")

			case config.Units == "kg":
				config.DeviceClass = SetDefault(config.DeviceClass, "weight")
				config.Icon = SetDefault(config.Icon, "mdi:weight")

			case config.Units == "km":
				config.DeviceClass = SetDefault(config.DeviceClass, "distance")
				config.Icon = SetDefault(config.Icon, "mdi:map-marker-distance")

			case config.Units == "W/㎡":
				fallthrough
			case config.Units == "W/m2":
				config.DeviceClass = SetDefault(config.DeviceClass, "irradiance")
				config.Icon = SetDefault(config.Icon, "mdi:weather-sunny")
				config.Units = "W/m²"

			case config.Units == "Wh/㎡":
				fallthrough
			case config.Units == "Wh/m2":
				// HA's irradiance device class only accepts W/m². Wh/m² is an
				// energy-per-area daily total, so leave the device class unset.
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:weather-sunny")
				config.Units = "Wh/m²"

			case config.Value.TypeValue == "Currency":
				fallthrough
			case config.Units == "AUD":
				fallthrough
			case config.Units == "$":
				config.DeviceClass = SetDefault(config.DeviceClass, "monetary")
				config.Icon = SetDefault(config.Icon, "mdi:currency-usd")

				// p13013 - power_factor

			case config.Units == "GPS":
				// config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "mdi:crosshairs-gps")

			default:
				config.DeviceClass = SetDefault(config.DeviceClass, "")
				config.Icon = SetDefault(config.Icon, "")
		}

		switch {
			case config.DeviceClass == "timestamp":
				// Timestamp sensors must not have a state_class.
				config.StateClass = ""
				config.LastReset = ""
				config.LastResetValueTemplate = ""

			case config.Point.IsBoot():
				config.StateClass = "measurement"
				config.LastReset = ""
				config.LastResetValueTemplate = ""

			case config.Point.IsDaily():
				fallthrough
			case config.Point.IsMonthly():
				fallthrough
			case config.Point.IsYearly():
				fallthrough
			case config.Point.IsTotal():
				if config.DeviceClass == "energy" {
					// HA requires total_increasing for energy sensors (resets are auto-detected).
					config.StateClass = "total_increasing"
					config.LastReset = ""
					config.LastResetValueTemplate = ""
				} else {
					config.StateClass = "total"
					config.LastResetValueTemplate = SetDefault(config.LastResetValueTemplate, "{{ value_json.last_reset | as_datetime }}")
				}

			case config.Point.Is5Minute():
				fallthrough
			case config.Point.Is15Minute():
				fallthrough
			case config.Point.Is30Minute():
				fallthrough
			case config.Point.IsInstant():
				fallthrough
			default:
				if config.DeviceClass == "energy" {
					// Energy sensors with unknown UpdateFreq default to total_increasing (cumulative).
					config.StateClass = "total_increasing"
				} else if config.DeviceClass == "timestamp" {
					// Timestamp sensors must not have a state_class.
					config.StateClass = ""
				} else {
					config.StateClass = "measurement"
				}
				config.LastReset = ""
				config.LastResetValueTemplate = ""
		}

		// StateClass implies a numeric value. If the value is a string
		// (e.g. "Matthew Wilson", "GMT+10", "Residential PV"), clear StateClass
		// to avoid HA errors about non-numeric values on state_class sensors.
		if config.StateClass != "" && config.Value != nil && !config.Value.IsFloat() && !config.Value.IsInt() && !config.Value.IsBool() {
			config.StateClass = ""
		}

		// Classify entity category and enabled_by_default.
		// Metadata sensors (plant name, timezone, GPS, capacity, etc.) go into
		// "diagnostic" so they don't clutter the default entity list.
		falseVal := false
		switch {
			// Purely informational metadata — diagnostic and hidden by default
			case config.FullId == "p80001":    // ps_name
				fallthrough
			case config.FullId == "p80003":    // plant_name
				fallthrough
			case config.FullId == "p80004":    // location / address
				fallthrough
			case config.FullId == "p80009":    // ps_timezone
				fallthrough
			case config.FullId == "p13039":    // latitude
				fallthrough
			case config.FullId == "p13040":    // longitude
				fallthrough
			case config.FullId == "p13013":    // plant_type
				fallthrough
			case config.FullId == "p13005":    // connect_type
				fallthrough
			case config.FullId == "p13007":    // capacity (MWp — static)
				config.EntityCategory = SetDefault(config.EntityCategory, "diagnostic")
				config.EnabledByDefault = &falseVal

			// Operational metadata — diagnostic but visible by default
			case config.FullId == "p83041":    // ps_status / running state
				fallthrough
			case config.FullId == "p83001":    // ps_switch
				fallthrough
			case config.FullId == "p83004":    // ps_connection
				config.EntityCategory = SetDefault(config.EntityCategory, "diagnostic")
		}

		// Classify device group for sub-device routing.
		// Entities are grouped into logical sub-devices under the parent.
		if config.DeviceGroup == "" && config.Point != nil {
			switch config.Point.GroupName {
				case "MPPT Information":
					config.DeviceGroup = "inverter"
				case "Other Information":
					config.DeviceGroup = "inverter"
				case "Grid Information":
					config.DeviceGroup = "grid"
				case "Load Information":
					config.DeviceGroup = "load"
				case "Battery Information":
					config.DeviceGroup = "battery"
				case "Overview":
					config.DeviceGroup = "plant"
				default:
					config.DeviceGroup = "inverter"
			}
		}
	}
}

func SetDefault(value string, def string) string {
	if value == "" {
		value = def
	}
	return value
}
