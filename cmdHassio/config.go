package cmdHassio

import (
	"encoding/json"
	"github.com/MickMake/GoUnify/Only"
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

// missingParents records ParentNames that NewDevice could not resolve. Every
// caller treats ok=false as "skip this entity" and returns no error, so an
// unregistered parent used to fail completely silently: no discovery payload
// published, no error set, no trace in the log. Tracked here so the condition
// is reported once at INFO instead of vanishing.
var missingParents = map[string]bool{}

func (m *Mqtt) NewDevice(config EntityConfig) (bool, Device) {
	var ok bool
	var ret Device

	for range Only.Once {
		var parent Device
		if parent, ok = m.MqttDevices[config.ParentName]; !ok {
			if !missingParents[config.ParentName] {
				missingParents[config.ParentName] = true
				m.logger.Info("Unknown parentDevice: %s - will ignore. Discovery payload for '%s' not published.\n",
					config.ParentName, config.FullId)
			}
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

		// deviceKey == config.ParentName means this device is the parent itself
		// (empty DeviceGroup, or the "inverter" default). Inheriting the parent's
		// via_device there would point the device at itself, which Home Assistant
		// rejects with "A device can not be its own via device". A root device
		// carries no via_device at all.
		viaDevice := parent.ViaDevice
		connections := [][]string{
			{ m.EntityPrefix, JoinStringsForId(m.EntityPrefix, config.ParentName) },
		}
		if deviceKey == config.ParentName {
			viaDevice = ""
		} else {
			connections = append(connections, []string{
				JoinStringsForId(m.EntityPrefix, config.ParentName),
				JoinStringsForId(m.EntityPrefix, deviceKey),
			})
		}

		ret = Device {
			ConfigurationUrl: parent.ConfigurationUrl,
			Connections:      connections,
			Identifiers:      []string{ JoinStringsForId(m.EntityPrefix, deviceKey) },
			Manufacturer:     manu,
			Model:            modl,
			Name:             deviceName,
			SuggestedArea:    parent.SuggestedArea,
			SwVersion:        parent.SwVersion,
			ViaDevice:        viaDevice,
		}

		m.MqttDevices[deviceKey] = ret
		ok = true
	}

	return ok, ret
}
