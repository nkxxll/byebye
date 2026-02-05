package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Version        string              `json:"version"`
	WindowManagers map[string]WMConfig `json:"windowManagers"`
}

type WMConfig struct {
	DisplayServer string    `json:"displayServer"`
	Lock          ActionSet `json:"lock"`
	Logout        []string  `json:"logout"`
	Sleep         []string  `json:"sleep"`
	Suspend       []string  `json:"suspend"`
	Hibernate     []string  `json:"hibernate"`
	Shutdown      []string  `json:"shutdown"`
	Restart       []string  `json:"restart"`
}

type ActionSet struct {
	Commands []string
	ByServer map[string][]string
}

func (a *ActionSet) UnmarshalJSON(data []byte) error {
	var list []string
	if err := json.Unmarshal(data, &list); err == nil {
		a.Commands = list
		a.ByServer = nil
		return nil
	}

	var byServer map[string][]string
	if err := json.Unmarshal(data, &byServer); err == nil {
		a.Commands = nil
		a.ByServer = byServer
		return nil
	}

	return errors.New("lock action must be a list of commands or a display-server map")
}

func (a ActionSet) MarshalJSON() ([]byte, error) {
	if len(a.Commands) > 0 {
		return json.Marshal(a.Commands)
	}
	if len(a.ByServer) > 0 {
		return json.Marshal(a.ByServer)
	}
	return []byte("[]"), nil
}

type Environment struct {
	Desktop    string
	Session    string
	Display    DisplayServer
	Wayland    bool
	X11        bool
	IsHeadless bool
}

type DisplayServer int

const (
	displayUnknown DisplayServer = iota
	displayX11
	displayWayland
)

func (d DisplayServer) String() string {
	switch d {
	case displayX11:
		return "X11"
	case displayWayland:
		return "Wayland"
	default:
		return "unknown"
	}
}

func detectEnvironment() Environment {
	desktop := os.Getenv("XDG_CURRENT_DESKTOP")
	session := os.Getenv("DESKTOP_SESSION")
	sessionType := os.Getenv("XDG_SESSION_TYPE")
	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	x11Display := os.Getenv("DISPLAY")

	env := Environment{
		Desktop: desktop,
		Session: session,
	}

	if sessionType == "wayland" || waylandDisplay != "" {
		env.Display = displayWayland
		env.Wayland = true
	}

	if sessionType == "x11" || x11Display != "" {
		env.Display = displayX11
		env.X11 = true
	}

	if !env.Wayland && !env.X11 {
		env.Display = displayUnknown
		env.IsHeadless = true
	}

	return env
}

func detectWindowManager(env Environment) string {
	if containsAny(env.Desktop, "Hyprland") || os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") != "" {
		return "Hyprland"
	}
	if containsAny(env.Desktop, "sway") || os.Getenv("SWAYSOCK") != "" {
		return "Sway"
	}
	if containsAny(env.Desktop, "KDE", "Plasma") || containsAny(env.Session, "kde", "plasma") {
		return "KDE"
	}
	if containsAny(env.Desktop, "GNOME") || containsAny(env.Session, "gnome") {
		return "GNOME"
	}
	if containsAny(env.Desktop, "XFCE") || containsAny(env.Session, "xfce") {
		return "XFCE"
	}
	if containsAny(env.Desktop, "i3") || containsAny(env.Session, "i3") {
		return "i3"
	}
	if containsAny(env.Desktop, "LXDE", "LXQt") || containsAny(env.Session, "lxde", "lxqt") {
		return "LXDE"
	}
	if env.Display == displayX11 {
		return "generic"
	}
	if env.Display == displayWayland {
		return "generic"
	}
	return "generic"
}

func containsAny(value string, needles ...string) bool {
	for _, needle := range needles {
		if needle == "" {
			continue
		}
		if containsFold(value, needle) {
			return true
		}
	}
	return false
}

func containsFold(value, needle string) bool {
	if value == "" || needle == "" {
		return false
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}

func defaultConfigData() (Config, error) {
	var cfg Config
	if err := json.Unmarshal([]byte(defaultConfig), &cfg); err != nil {
		return Config{}, fmt.Errorf("parse default config: %w", err)
	}
	return cfg, nil
}

func loadConfigFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	return cfg, nil
}

func mergeConfigs(base Config, override Config) Config {
	if override.Version != "" {
		base.Version = override.Version
	}

	if base.WindowManagers == nil {
		base.WindowManagers = map[string]WMConfig{}
	}

	for name, overrideWM := range override.WindowManagers {
		baseWM, ok := base.WindowManagers[name]
		if !ok {
			base.WindowManagers[name] = overrideWM
			continue
		}

		if overrideWM.DisplayServer != "" {
			baseWM.DisplayServer = overrideWM.DisplayServer
		}

		if len(overrideWM.Lock.Commands) > 0 || len(overrideWM.Lock.ByServer) > 0 {
			baseWM.Lock = overrideWM.Lock
		}
		if len(overrideWM.Logout) > 0 {
			baseWM.Logout = overrideWM.Logout
		}
		if len(overrideWM.Sleep) > 0 {
			baseWM.Sleep = overrideWM.Sleep
		}
		if len(overrideWM.Suspend) > 0 {
			baseWM.Suspend = overrideWM.Suspend
		}
		if len(overrideWM.Hibernate) > 0 {
			baseWM.Hibernate = overrideWM.Hibernate
		}
		if len(overrideWM.Shutdown) > 0 {
			baseWM.Shutdown = overrideWM.Shutdown
		}
		if len(overrideWM.Restart) > 0 {
			baseWM.Restart = overrideWM.Restart
		}

		base.WindowManagers[name] = baseWM
	}

	return base
}

func configPaths() []string {
	paths := []string{
		"/etc/byebye/config.json",
	}

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, ".config", "byebye", "config.json"))
	}

	if env := os.Getenv("BYEBYE_CONFIG"); env != "" {
		paths = append(paths, env)
	}

	return paths
}

func loadConfig() (Config, error) {
	cfg, err := defaultConfigData()
	if err != nil {
		return Config{}, err
	}

	for _, path := range configPaths() {
		loaded, err := loadConfigFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return Config{}, err
		}
		cfg = mergeConfigs(cfg, loaded)
	}

	if err := validateConfig(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func validateConfig(cfg Config) error {
	if cfg.Version == "" {
		return errors.New("config version is required")
	}
	if len(cfg.WindowManagers) == 0 {
		return errors.New("config must define at least one window manager")
	}
	return nil
}

func (c Config) getCommands(wm string, action string, display DisplayServer) []string {
	if c.WindowManagers == nil {
		return nil
	}

	if wmConfig, ok := c.WindowManagers[wm]; ok {
		if cmds := commandsForAction(wmConfig, action, display); len(cmds) > 0 {
			return cmds
		}
	}

	if generic, ok := c.WindowManagers["generic"]; ok {
		return commandsForAction(generic, action, display)
	}

	return nil
}

func commandsForAction(cfg WMConfig, action string, display DisplayServer) []string {
	switch action {
	case "lock":
		return resolveLockCommands(cfg.Lock, display)
	case "logout":
		return cfg.Logout
	case "sleep":
		return cfg.Sleep
	case "suspend":
		return cfg.Suspend
	case "hibernate":
		return cfg.Hibernate
	case "shutdown":
		return cfg.Shutdown
	case "restart":
		return cfg.Restart
	default:
		return nil
	}
}

func resolveLockCommands(lock ActionSet, display DisplayServer) []string {
	if len(lock.Commands) > 0 {
		return lock.Commands
	}

	if len(lock.ByServer) == 0 {
		return nil
	}

	switch display {
	case displayWayland:
		if cmds := lock.ByServer["Wayland"]; len(cmds) > 0 {
			return cmds
		}
	case displayX11:
		if cmds := lock.ByServer["X11"]; len(cmds) > 0 {
			return cmds
		}
	}

	return nil
}
