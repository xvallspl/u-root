package integration

import (
	"flag"
	"os"
	"testing"
)

var tests = []struct {
	name    string
	runTest func(t *testing.T)
	// Allow a user to manually pass in an initramfs as a flag for a given test.
	initramfs *string
}{
	// uinit_test.go
	{
		name:      "TestHelloWorld",
		initramfs: flag.String("TestHelloWorld_Initramfs", "", "specify a custom initramfs"),
		runTest:   RunTestHelloWorld,
	},
	{
		name:      "TestHelloWorldNegative",
		initramfs: flag.String("TestHelloWorldNegative_Initramfs", "", "specify a custom initramfs"),
		runTest:   RunTestHelloWorldNegative,
	},
	{
		name:      "TestScript",
		initramfs: flag.String("TestScript_Initramfs", "", "specify a custom initramfs"),
		runTest:   RunTestScript,
	},
	// dhclient_test.go
	{
		name:      "TestDhclient",
		initramfs: flag.String("TestDhclient_Initramfs", "", "specify a custom initramfs"),
		runTest:   RunTestDhclient,
	},
	{
		name:      "TestPxeboot",
		initramfs: flag.String("TestPxeboot_Initramfs", "", "specify a custom initramfs"),
		runTest:   RunTestPxeboot,
	},
	{
		name:      "TestQEMUDHCPTimesOut",
		initramfs: flag.String("TestQEMUDHCPTimesOut_Initramfs", "", "specify a custom initramfs"),
		runTest:   RunTestQEMUDHCPTimesOut,
	},
}

func TestIntegration(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.runTest == nil {
				t.Fatalf("No runner function found for %s", tt.name)
			}

			// If a custom initramfs was provided by user, prefer that.
			if len(*tt.initramfs) != 0 {
				os.Setenv("UROOT_INITRAMFS", *tt.initramfs)
			}

			tt.runTest(t)
		})
	}
}

// TODO try this with passing an initramfs in

// Add a README explaining 1) how to add a new test, 2) changed how to run tests
// With this we get the naming hierarchy for free - can do go test -test.run=Foo/TestScript
