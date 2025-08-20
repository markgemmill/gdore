package environ

import (
	"github.com/BurntSushi/toml"
	"github.com/markgemmill/appdirs"
	"github.com/markgemmill/pathlib"
)

type Config struct {
	DefaultUser string `toml:"default_user,omitempty"`
	OutputDir   string `toml:"out_dir"`
	Headless    bool   `toml:"headless_browser"`
}

type Environ struct {
	HomeDir     pathlib.Path
	LogDir      pathlib.Path
	TmpDir      pathlib.Path
	ConfigFile  pathlib.Path
	OutputDir   pathlib.Path
	DefaultUser string
	Headless    bool
}

func (env Environ) LogFile(logname string) pathlib.Path {
	return env.LogDir.Join(logname)
}

func (env Environ) DumpConfig() Config {
	return Config{
		DefaultUser: env.DefaultUser,
		Headless:    env.Headless,
		OutputDir:   env.OutputDir.String(),
	}
}

func DefaultEnviron(homeDir pathlib.Path) (Environ, error) {
	tmp, err := pathlib.NewTempDir("gdore-")
	if err != nil {
		return Environ{}, err
	}

	outDir := pathlib.NewPath(".", 0770).Resolve()

	env := Environ{
		HomeDir:     homeDir,
		LogDir:      homeDir.Join("logs"),
		ConfigFile:  homeDir.Join("config.toml"),
		TmpDir:      tmp,
		Headless:    true,
		OutputDir:   outDir,
		DefaultUser: "",
	}

	return env, nil
}

func DumpConfig(env Environ) error {
	cfg := Config{
		DefaultUser: env.DefaultUser,
		OutputDir:   env.OutputDir.String(),
		Headless:    env.Headless,
	}
	content, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = env.ConfigFile.Write(content)
	if err != nil {
		return err
	}

	return nil
}

func loadEnviron(homeDir pathlib.Path) (Environ, error) {
	env, err := DefaultEnviron(homeDir)
	if err != nil {
		return env, err
	}
	var cfg Config

	if env.ConfigFile.Exists() {
		content, err := env.ConfigFile.Read()
		if err != nil {
			return env, err
		}
		toml.Unmarshal(content, &cfg)
		// set config values to environment
		env.OutputDir = pathlib.NewPath(cfg.OutputDir, 0770).Resolve()
		env.DefaultUser = cfg.DefaultUser
		env.Headless = cfg.Headless

	} else {
		// no existing config, then write the defaults
		err := DumpConfig(env)
		if err != nil {
			return env, err
		}
	}

	return env, nil
}

func CreateEnvironment() (Environ, error) {
	// find our home directory
	appDirPath := appdirs.UserDataDir("gdore", "")

	home := pathlib.NewPath(appDirPath, 0770)
	err := home.MkDirs()
	if err != nil {
		return Environ{}, err
	}

	// load our environment
	env, err := loadEnviron(home)
	if err != nil {
		return env, err
	}

	// make sure our folders exist
	err = env.LogDir.MkDirs()
	if err != nil {
		return env, err
	}

	return env, nil
}
