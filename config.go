package readeef

import (
	"os"
	"time"

	"code.google.com/p/gcfg"
)

var apiversion = 1

type Config struct {
	Readeef struct {
		Debug bool
	}
	API struct {
		Version int
	}
	Timeout struct {
		Connect   string
		ReadWrite string `gcfg:"read-write"`

		Converted struct {
			Connect   time.Duration
			ReadWrite time.Duration
		}
	}
	DB struct {
		Driver  string
		Connect string
	}
	Auth struct {
		Secret          string
		IgnoreURLPrefix []string `gcfg:"ignore-url-prefix"`
	}
	Hubbub struct {
		CallbackURL  string `gcfg:"callback-url"` // http://www.example.com
		RelativePath string `gcfg:"relative-path"`
		From         string
	}
	Updater struct {
		Interval string

		Converted struct {
			Interval time.Duration
		}
	}
}

func ReadConfig(path ...string) (Config, error) {
	def, err := defaultConfig()

	if err != nil {
		return Config{}, err
	}

	c := def

	if len(path) != 0 {
		err = gcfg.ReadFileInto(&c, path[0])

		if err != nil && !os.IsNotExist(err) {
			return Config{}, err
		}
	}

	c.API.Version = apiversion

	if d, err := time.ParseDuration(c.Timeout.Connect); err == nil {
		c.Timeout.Converted.Connect = d
	} else {
		c.Timeout.Converted.Connect = time.Second
	}

	if d, err := time.ParseDuration(c.Timeout.ReadWrite); err == nil {
		c.Timeout.Converted.ReadWrite = d
	} else {
		c.Timeout.Converted.ReadWrite = time.Second
	}

	if d, err := time.ParseDuration(c.Updater.Interval); err == nil {
		c.Updater.Converted.Interval = d
	} else {
		c.Updater.Converted.Interval = 30 * time.Minute
	}
	return c, nil
}

func defaultConfig() (Config, error) {
	var def Config

	err := gcfg.ReadStringInto(&def, cfg)

	if err != nil {
		return Config{}, err
	}

	def.API.Version = apiversion
	return def, nil
}

var cfg string = `
[readeef]
	debug = true
[db]
	driver = sqlite3
	connect = file:./readeef.sqlite3?cache=shared&mode=rwc
[timeout]
	connect = 1s
	read-write = 2s
[hubbub]
	relative-path = /hubbub
	from = readeef
`
