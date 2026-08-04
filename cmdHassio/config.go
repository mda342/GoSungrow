package cmdHassio

import (
	"encoding/json"
	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdLog"
)


type Config struct {
	Entry        string       `json:"~,omitempty" required:"false"`
	Name         string       `json:"name,omitempty" required:"false"`
	UniqueId     string       `json:"unique_id,omitempty" required:"false"`
	StateTopic   string       `json:"state_topic,omitempty" required:"true"`
	DeviceConfig DeviceConfig `json:"device,omitempty" required:"false"`
}
type DeviceConfig struct {
	Identifiers  []string `json:"identifiers,omitempty" required:"false"`
	SwVersion    string   `json:"sw_version,omitempty" required:"false"`
	Name         string   `json:"name,omitempty" required:"false"`
	Manufacturer string   `json:"manufacturer,omitempty" required:"false"`
	Model        string   `json:"model,omitempty" required:"false"`
}

func (c *Config) Json() string {
	j, _ := json.Marshal(*c)
	return string(j)
}

var DeviceGroupLabels = map[string]string{
	"inverter": "Inverter",
	"grid":     "Grid",
	"load":     "Load",
	"battery":  "Battery",
	"plant":    "Plant",
}

func (m *Mqtt) NewDevice(config EntityConfig) (bool, Device) {
	var ok bool
	var ret Device

	for range Only.Once {
		var parent Device
		if parent, ok = m.MqttDevices[config.ParentName]; !ok {
			cmdLog.LogPrintDate("Unknown parentDevice: %s - will ignore.\n", config.ParentName)
			break
		}

		manu := parent.Manufacturer
		if manu == "" {
			manu = m.DeviceName
		}
		modl := parent.Model
		if modl == "" {
			modl = m.DeviceName
		}

		// Build device key: parent name for "inverter" (default), parent_group for others
		deviceKey := config.ParentName
		if config.DeviceGroup != "" && config.DeviceGroup != "inverter" {
			deviceKey = JoinStringsForId(config.ParentName, config.DeviceGroup)
		}

		// Return cached device if already created
		if cached, cachedOk := m.MqttDevices[deviceKey]; cachedOk {
			return cachedOk, cached
		}

		// Build device name
		deviceName := JoinStrings(m.EntityPrefix, config.ParentName, "-", parent.Name)
		if config.DeviceGroup != "" && config.DeviceGroup != "inverter" {
			groupLabel := DeviceGroupLabels[config.DeviceGroup]
			if groupLabel == "" {
				groupLabel = config.DeviceGroup
			}
			deviceName = JoinStrings(m.EntityPrefix, config.ParentName, "-", parent.Name, "-", groupLabel)
		}

		ret = Device {
			ConfigurationUrl: parent.ConfigurationUrl,
			Connections:      [][]string {
				{ m.EntityPrefix, JoinStringsForId(m.EntityPrefix, config.ParentName) },
				{ JoinStringsForId(m.EntityPrefix, config.ParentName), JoinStringsForId(m.EntityPrefix, deviceKey) },
			},
			Identifiers:      []string{ JoinStringsForId(m.EntityPrefix, deviceKey) },
			Manufacturer:     manu,
			Model:            modl,
			Name:             deviceName,
			SuggestedArea:    parent.SuggestedArea,
			SwVersion:        parent.SwVersion,
			ViaDevice:        parent.ViaDevice,
		}

		m.MqttDevices[deviceKey] = ret
		ok = true
	}

	return ok, ret
}
