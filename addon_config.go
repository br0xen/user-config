package userConfig

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// AddonConfig is an additional ConfigFile
type AddonConfig struct {
	Name   string                       `toml:"-"`
	Path   string                       `toml:"-"`
	Values map[string]map[string]string `toml:"-"`
}

// NewAddonConfig generates a Additional Config struct
func NewAddonConfig(name, path string) (*AddonConfig, error) {
	af := &AddonConfig{Name: name, Path: path}
	af.Values = make(map[string]map[string]string)

	// Check if file exists
	//var f os.FileInfo
	var err error
	if _, err = os.Stat(af.GetFullPath()); os.IsNotExist(err) {
		if err = af.Save(); err != nil {
			return af, err
		}
	}

	if err := af.Load(); err != nil {
		return af, err
	}
	return af, nil
}

/** START of ConfigFile Interface Implementation **/

// GetName returns the name of this config file
func (af *AddonConfig) GetName() string {
	return af.Name
}

// GetPath returns the path of this config file
func (af *AddonConfig) GetPath() string {
	return af.Path
}

// Load loads config files into the config
func (af *AddonConfig) Load() error {
	if strings.TrimSpace(af.Name) == "" || strings.TrimSpace(af.Path) == "" {
		return errors.New("Invalid ConfigFile Name: " + af.GetFullPath())
	}

	// Config files end with .conf
	tomlData, err := ioutil.ReadFile(af.GetFullPath())
	if err != nil {
		return err
	}
	fmt.Println(tomlData)
	// TODO: Figure out loading this into the struct
	//if _, err := toml.Decode(string(tomlData), &af); err != nil {
	//	return err
	//}
	return nil
}

// Save writes the config to file(s)
func (af *AddonConfig) Save() error {
	buf := new(bytes.Buffer)
	// TODO: Figure out writing struct to buf
	//if err := toml.NewEncoder(buf).Encode(af); err != nil {
	//	return err
	//}
	return ioutil.WriteFile(af.GetFullPath(), buf.Bytes(), 0644)
}

// Set sets a key/value pair in af, if unable to save, revert to old value
// (and return the error)
func (af *AddonConfig) Set(category, k, v string) error {
	if _, ok := af.Values[category]; !ok {
		af.Values[category] = make(map[string]string)
	}
	oldVal := af.Values[category][k]
	af.Values[category][k] = v
	if err := af.Save(); err != nil {
		af.Values[category][k] = oldVal
		return err
	}
	return nil
}

// Get gets a key/value pair from af
func (af *AddonConfig) Get(category, k string) string {
	if _, ok := af.Values[category]; !ok {
		return ""
	}
	return af.Values[category][k]
}

// GetFullPath returns the full path & filename to the config file
func (af *AddonConfig) GetFullPath() string {
	return af.Path + "/" + af.Name + ".conf"
}

/** END of ConfigFile Interface Implementation **/
