package chroot

import (
	"testing"

	"github.com/hashicorp/packer/packer"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"ami_name":   "foo",
		"source_ami": "foo",
		"region":     "us-east-1",
		// region validation logic is checked in ami_config_test
		"skip_region_validation": true,
	}
}

func TestBuilder_ImplementsBuilder(t *testing.T) {
	var raw interface{}
	raw = &Builder{}
	if _, ok := raw.(packer.Builder); !ok {
		t.Fatalf("Builder should be a builder")
	}
}

func TestBuilderPrepare_AMIName(t *testing.T) {
	var b Builder
	config := testConfig()

	// Test good
	config["ami_name"] = "foo"
	config["skip_region_validation"] = true
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}

	// Test bad
	config["ami_name"] = "foo {{"
	b = Builder{}
	warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}

	// Test bad
	delete(config, "ami_name")
	b = Builder{}
	warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestBuilderPrepare_ChrootMounts(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["chroot_mounts"] = nil
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestBuilderPrepare_ChrootMountsBadDefaults(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["chroot_mounts"] = [][]string{
		{"bad"},
	}
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}
}
func TestBuilderPrepare_SourceAmi(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["source_ami"] = ""
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatal("should have error")
	}

	config["source_ami"] = "foo"
	warnings, err = b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestBuilderPrepare_CommandWrapper(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["command_wrapper"] = "echo hi; {{.Command}}"
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}
}

func TestBuilderPrepare_CopyFiles(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}

	if len(b.config.CopyFiles) != 1 && b.config.CopyFiles[0] != "/etc/resolv.conf" {
		t.Errorf("Was expecting default value for copy_files.")
	}
}

func TestBuilderPrepare_CopyFilesNoDefault(t *testing.T) {
	b := &Builder{}
	config := testConfig()

	config["copy_files"] = []string{}
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err != nil {
		t.Errorf("err: %s", err)
	}

	if len(b.config.CopyFiles) > 0 {
		t.Errorf("Was expecting no default value for copy_files.")
	}
}

func TestBuilderPrepare_RootDeviceNameAndAMIMappings(t *testing.T) {
	var b Builder
	config := testConfig()

	config["root_device_name"] = "/dev/sda"
	config["ami_block_device_mappings"] = []interface{}{map[string]string{}}
	config["root_volume_size"] = 15
	warnings, err := b.Prepare(config)
	if len(warnings) == 0 {
		t.Fatal("Missing warning, stating block device mappings will be overwritten")
	} else if len(warnings) > 1 {
		t.Fatalf("excessive warnings: %#v", warnings)
	}
	if err != nil {
		t.Fatalf("should not have error: %s", err)
	}
}

func TestBuilderPrepare_AMIMappingsNoRootDeviceName(t *testing.T) {
	var b Builder
	config := testConfig()

	config["ami_block_device_mappings"] = []interface{}{map[string]string{}}
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatalf("should have error")
	}
}

func TestBuilderPrepare_RootDeviceNameNoAMIMappings(t *testing.T) {
	var b Builder
	config := testConfig()

	config["root_device_name"] = "/dev/sda"
	warnings, err := b.Prepare(config)
	if len(warnings) > 0 {
		t.Fatalf("bad: %#v", warnings)
	}
	if err == nil {
		t.Fatalf("should have error")
	}
}
