package cmdHassio

import (
	"testing"
)

// A device that is its own parent must not carry via_device, otherwise Home
// Assistant rejects the discovery payload with
// "A device can not be its own via device".
func TestNewDevice_RootDeviceHasNoViaDevice(t *testing.T) {
	m := &Mqtt{
		EntityPrefix: "GoSungrow",
		MqttDevices:  map[string]Device{},
	}
	// Parent registered the way SetDeviceConfig does for swname == parentId.
	m.MqttDevices["options"] = Device{
		Connections: [][]string{{"GoSungrow", "GoSungrow-options"}},
		Identifiers: []string{"GoSungrow-options"},
		Name:        "Options",
	}

	// Empty DeviceGroup => deviceKey == ParentName => this IS the parent.
	ok, dev := m.NewDevice(EntityConfig{ParentName: "options"})
	if !ok {
		t.Fatal("NewDevice returned ok=false")
	}
	if dev.ViaDevice != "" {
		t.Errorf("root device via_device = %q, want empty", dev.ViaDevice)
	}
	for _, c := range dev.Connections {
		if c[0] == c[1] {
			t.Errorf("degenerate self-connection %v", c)
		}
	}
}

// A child device must still inherit its parent's via_device so the device tree
// keeps its nesting.
func TestNewDevice_ChildInheritsViaDevice(t *testing.T) {
	m := &Mqtt{
		EntityPrefix: "GoSungrow",
		MqttDevices:  map[string]Device{},
	}
	m.MqttDevices["parent"] = Device{
		Identifiers: []string{"GoSungrow-parent"},
		Name:        "Parent",
		ViaDevice:   "GoSungrow-root",
	}

	ok, dev := m.NewDevice(EntityConfig{ParentName: "parent", DeviceGroup: "plant"})
	if !ok {
		t.Fatal("NewDevice returned ok=false")
	}
	if dev.ViaDevice != "GoSungrow-root" {
		t.Errorf("child via_device = %q, want %q", dev.ViaDevice, "GoSungrow-root")
	}
	if len(dev.Connections) != 2 {
		t.Errorf("child connections = %v, want 2 entries", dev.Connections)
	}
}

// SetOption for an id that was never created must return an error rather than
// dereferencing the nil EntityConfig that EntityConfig() returns for a miss.
func TestSetOptionUnknownIDReturnsError(t *testing.T) {
	m := &Mqtt{UserOptions: Options{}, MqttDevices: map[string]Device{}}
	m.UserOptions.New()

	if err := m.SetOption("nope", "info"); err == nil {
		t.Error("SetOption with unknown id returned nil error, want error")
	}
}

// NewDevice reports failure by returning ok=false, and SelectPublishConfig
// then breaks out of its single-iteration loop without setting an error. This
// documents that silent path: an option whose ParentName is not registered is
// never published, and the caller sees success.
func TestOptionWithUnregisteredParentIsNotPublished(t *testing.T) {
	m := &Mqtt{UserOptions: Options{}, MqttDevices: map[string]Device{}}
	m.UserOptions.New()

	if err := m.CreateOption("loglevel", "Log Level", nil, "info", "debug"); err != nil {
		t.Fatalf("CreateOption returned %v, want nil", err)
	}
	if _, ok := m.MqttDevices["options"]; ok {
		t.Error("options device unexpectedly registered")
	}
}
