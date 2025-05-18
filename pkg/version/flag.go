package version

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

// versionValue defines the type for version flag.
type versionValue int

// Define some constants.
const (
	// Version not set.
	VersionNotSet versionValue = 0
	// Version enabled.
	VersionEnabled versionValue = 1
	// Raw version.
	VersionRaw versionValue = 2
)

const (
	// String representing raw version.
	strRawVersion = "raw"
	// Name of the version flag.
	versionFlagName = "version"
)

// versionFlag defines the version flag.
var versionFlag = Version(versionFlagName, VersionNotSet, "Print version information and quit.")

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

// String implements the String method of the pflag.Value interface.
func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion // Return raw version string
	}
	return strconv.FormatBool(bool(*v == VersionEnabled))
}

// Set implements the Set method of the pflag.Value interface.
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionEnabled
	} else {
		*v = VersionNotSet
	}
	return err
}

// Type implements the Type method of the pflag.Value interface.
func (v *versionValue) Type() string {
	return "version"
}

// VersionVar defines a flag with the specified name and usage.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)

	// `--version` is equivalent to `--version=true`
	pflag.Lookup(name).NoOptDefVal = "true"
}

// Version wraps the VersionVar function.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

// AddFlags registers the flags of this package on any FlagSet, so they point to the same value as the global flags.
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}

// PrintAndExitIfRequested checks if the `--version` flag was passed, and if so, prints the version and exits.
func PrintAndExitIfRequested() {
	// Check the value of the version flag and print the corresponding information
	if *versionFlag == VersionRaw {
		fmt.Printf("%s\n", Get().Text())
		os.Exit(0)
	} else if *versionFlag == VersionEnabled {
		fmt.Printf("%s\n", Get().String())
		os.Exit(0)
	}
}
