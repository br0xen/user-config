package userConfig

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

// GeneralConfig is the basic ConfigFile
// All configs made with package userConfig will have this file
// All additional config files will have an entry in 'ConfigFiles' here
type GeneralConfig struct {
	Name        string            `toml:"-"`
	Path        string            `toml:"-"`
	ConfigFiles []string          `toml:"additional_config"`
	Values      map[string]string `toml:"general"`
}

// NewGeneralConfig generates a General Config struct
func NewGeneralConfig(name, path string) (*GeneralConfig, error) {
	gf := &GeneralConfig{Name: name, Path: path}
	gf.ConfigFiles = []string{}
	gf.Values = make(map[string]string)

	// Check if file exists
	//var f os.FileInfo
	var err error
	if _, err = os.Stat(gf.GetFullPath()); os.IsNotExist(err) {
		if err = gf.Save(); err != nil {
			return gf, err
		}
	}

	if err := gf.Load(); err != nil {
		return gf, err
	}
	return gf, nil
}

/** START of ConfigFile Interface Implementation **/

// GetName returns the name of this config file
func (gf *GeneralConfig) GetName() string {
	return gf.Name
}

// GetPath returns the path of this config file
func (gf *GeneralConfig) GetPath() string {
	return gf.Path
}

// Load loads config files into the config
func (gf *GeneralConfig) Load() error {
	if strings.TrimSpace(gf.Name) == "" || strings.TrimSpace(gf.Path) == "" {
		return errors.New("Invalid ConfigFile Name: " + gf.GetFullPath())
	}

	// Config files end with .conf
	tomlData, err := ioutil.ReadFile(gf.GetFullPath())
	if err != nil {
		return err
	}
	if _, err := toml.Decode(string(tomlData), &gf); err != nil {
		return err
	}
	return nil
}

// Save writes the config to file(s)
func (gf *GeneralConfig) Save() error {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(gf); err != nil {
		return err
	}
	return ioutil.WriteFile(gf.GetFullPath(), buf.Bytes(), 0644)
}

// Set sets a key/value pair in gf, if unable to save, revert to old value
// (and return the error)
func (gf *GeneralConfig) Set(k, v string) error {
	oldVal := gf.Values[k]
	gf.Values[k] = v
	if err := gf.Save(); err != nil {
		gf.Values[k] = oldVal
		return err
	}
	return nil
}

// Get gets a key/value pair from gf
func (gf *GeneralConfig) Get(k string) string {
	return gf.Values[k]
}

// GetFullPath returns the full path & filename to the config file
func (gf *GeneralConfig) GetFullPath() string {
	return gf.Path + "/" + gf.Name + ".conf"
}

/** END of ConfigFile Interface Implementation **/

// Additional General Config Functions

// HasConfigFile returns true if GeneralConfig knows about <name>.conf
func (gf *GeneralConfig) HasConfigFile(name string) bool {
	for _, v := range gf.ConfigFiles {
		if v == name {
			return true
		}
	}
	return false
}

// AddConfigFile adds the config file <name>.conf
/*
func (gf *GeneralConfig) AddConfigFile(name string) (ConfigFile, error) {
	// Check if file exists
	var f os.FileInfo
	var err error
	cf := ConfigFile{}
	if f, err = os.Stat(cf.GetFullPath()); os.IsNotExist(err) {
		if err = cf.Save(); err != nil {
			return cf, err
		}
	}
	if gf.HasConfigFile() {
		// We already know about this file... So just return it
		return gf.GetConfigFile(name)
	}
	gf.ConfigFiles = append(gf.ConfigFiles, name)
	return cf, nil
}

// GetConfigFile returns an additional config file from the config directory
func (gf *GeneralConfig) GetConfigFile(name string) (ConfigFile, error) {

}
*/
