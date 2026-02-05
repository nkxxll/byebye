# Linux Compatibility Plan for byebye

## Strategy Overview

**Primary Detection Path**: Window Manager/Desktop Environment → WM-specific tools + generic fallbacks  
**Secondary Path**: Init system fallbacks (systemd → openrc → runit, etc.)

Detect the WM/DE first, then use its native tools. Only fall back to generic tools if WM-specific ones aren't available.

---

## Part 1: Window Manager/Desktop Environment Detection

### Detection Method
1. Check `XDG_CURRENT_DESKTOP` env var (primary)
2. Check `DESKTOP_SESSION` env var (fallback)
3. Check active window manager via X11/Wayland protocols
4. Check running processes (`ps aux | grep`)

### Detection Priority

```
1. Hyprland (wayland-only)
2. Sway (wayland-specific, i3-compatible)
3. KDE Plasma (X11 + Wayland)
4. GNOME (X11 + Wayland)
5. XFCE (X11 + Wayland)
6. i3 (X11-only)
7. OpenBox (X11)
8. Other X11 window managers
9. Wayland compositors (generic)
10. Fallback (unknown/headless)
```

---

## Part 2: WM/DE-Specific Command Mapping

### Hyprland (Wayland)

**Detection**: `XDG_CURRENT_DESKTOP=Hyprland` or `HYPRLAND_INSTANCE_SIGNATURE` set

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | `hyprctl dispatch exec 'hyprlock'`             |
| **Logout** | `hyprctl dispatch exit`                        |
| **Suspend**| `systemctl suspend`                 |
| **Hibernate** | `systemctl hibernate`          |
| **Shutdown** | `systemctl poweroff`                          |
| **Restart** | `systemctl reboot`                             |

**Notes**: Hyprlock for lock screen, hyprctl for WM control. Hyprctl exit handles logout directly.

---

### Sway (Wayland)

**Detection**: `XDG_CURRENT_DESKTOP=sway` or `SWAYSOCK` env var set

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | `swaylock -c 000000 -e`                        |
| **Logout** | `swaymsg exit`                                  |
| **Suspend**| `swaylock -c 000000 & systemctl suspend`       |
| **Hibernate** | `swaylock -c 000000 & systemctl hibernate` |
| **Shutdown** | `systemctl poweroff`                          |
| **Restart** | `systemctl reboot`                             |

**Notes**: swaymsg for WM control, swaylock for lock. Consider user's preferred colors.

---

### KDE Plasma (X11/Wayland)

**Detection**: `XDG_CURRENT_DESKTOP` contains `KDE` or `Plasma`

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | `qdbus org.freedesktop.ScreenSaver /ScreenSaver Lock` |
| **Logout** | `qdbus org.kde.ksmserver /KSMServer logout 0 1 1` |
| **Suspend**| `qdbus org.freedesktop.PowerManagement /org/freedesktop/PowerManagement Suspend` |
| **Hibernate** | `qdbus org.freedesktop.PowerManagement /org/freedesktop/PowerManagement Hibernate` |
| **Shutdown** | `qdbus org.kde.ksmserver /KSMServer logout 1 1 1` |
| **Restart** | `qdbus org.kde.ksmserver /KSMServer logout 1 2 1` |

**Notes**: KDE uses D-Bus. Fallback to systemctl if D-Bus unavailable.

---

### GNOME (X11/Wayland)

**Detection**: `XDG_CURRENT_DESKTOP` contains `GNOME`

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | `loginctl lock-session` or `gnome-screensaver-command -l` |
| **Logout** | `gio launch gnome-session-logout` or `loginctl terminate-session` |
| **Suspend**| `systemctl suspend`                             |
| **Hibernate** | `systemctl hibernate`                        |
| **Shutdown** | `systemctl poweroff`                          |
| **Restart** | `systemctl reboot`                             |

**Notes**: GNOME prefers systemd for power, loginctl for session management.

---

### XFCE (X11/Wayland)

**Detection**: `XDG_CURRENT_DESKTOP` contains `XFCE`

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | `xfce4-screensaver-command -l`                 |
| **Logout** | `xfce4-session-logout` or `kill -15 $PPID`     |
| **Suspend**| `xfce4-power-manager-settings` or `systemctl suspend` |
| **Hibernate** | `xfce4-power-manager-settings` or `systemctl hibernate` |
| **Shutdown** | `xfce4-session-logout --halt`                 |
| **Restart** | `xfce4-session-logout --reboot`                |

**Notes**: XFCE has dedicated session tools. Fall back to systemctl for power ops.

---

### i3 (X11-only)

**Detection**: `XDG_CURRENT_DESKTOP=i3` or `DISPLAY` set + `i3` running

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | `i3lock -c 000000` or `i3lock-fancy` or `i3lock-color` |
| **Logout** | `i3-msg exit` or `kill -15 $PPID`              |
| **Suspend**| `i3lock -c 000000 & systemctl suspend`         |
| **Hibernate** | `i3lock -c 000000 & systemctl hibernate`    |
| **Shutdown** | `systemctl poweroff`                          |
| **Restart** | `systemctl reboot`                             |

**Notes**: i3lock for locking (multiple variants available). i3msg for WM commands.

---

### LXDE/LXQt (X11)

**Detection**: `XDG_CURRENT_DESKTOP` contains `LXDE` or `LXQt`

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | `lxlock` or `slock`                             |
| **Logout** | `lxsession-logout` or `loginctl terminate-session` |
| **Suspend**| `systemctl suspend`                             |
| **Hibernate** | `systemctl hibernate`                        |
| **Shutdown** | `lxsession-logout --shutdown`                 |
| **Restart** | `lxsession-logout --reboot`                    |

---

### Generic X11 (No specific DE detected)

**Detection**: `DISPLAY` set but no known DE detected

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | See X11 Fallback Chain below                    |
| **Logout** | `kill -15 $PPID` or `loginctl terminate-session` |
| **Suspend**| `systemctl suspend`                             |
| **Hibernate** | `systemctl hibernate`                        |
| **Shutdown** | `systemctl poweroff`                          |
| **Restart** | `systemctl reboot`                             |

---

### Generic Wayland (No specific DE detected)

**Detection**: `WAYLAND_DISPLAY` set but no known WM detected

| Action     | Command(s)                                      |
|------------|------------------------------------------------|
| **Lock**   | See Wayland Fallback Chain below                |
| **Logout** | `loginctl terminate-session`                    |
| **Suspend**| `systemctl suspend`                             |
| **Hibernate** | `systemctl hibernate`                        |
| **Shutdown** | `systemctl poweroff`                          |
| **Restart** | `systemctl reboot`                             |

---

## Part 3: Generic Fallback Chains

### Lock Screen Fallback Chain

#### X11 Fallback Priority:
1. User's configured locker (if available)
2. `i3lock` (portable, highly compatible)
3. `i3lock-color` (enhanced version)
4. `slock` (simple, minimal)
5. `xdg-screensaver lock`
6. `xlock`
7. `xtrlock`

#### Wayland Fallback Priority:
1. User's configured locker (if available)
2. `swaylock` (most compatible)
3. `gtklock` (GTK-based, DE-agnostic)
4. `waylock` (Rust-based, minimal)
5. `wlopm` (simple screen blanking)
6. Fallback: Blank screen + disable input

---

### Logout Fallback Chain

**Priority Order**:
1. WM-specific exit command (handled above)
2. `loginctl terminate-session $XDG_SESSION_ID` (systemd-logind)
3. `systemctl --user exit` (user session manager)
4. Kill parent shell gracefully: `kill -15 $PPID`
5. Kill all user processes gracefully: `killall -15 -u $USER`
6. Last resort (avoid): `kill -9` with selected processes

---

### Suspend/Hibernate Fallback Chain

**For systems with systemd**:
1. `systemctl suspend` or `systemctl hibernate`
2. `loginctl suspend` or `loginctl hibernate`
3. `pm-suspend` or `pm-hibernate` (legacy)

**For systems without systemd**:
1. `rc-service openrc-run suspend` (OpenRC)
2. `/etc/init.d/suspend` (if exists)
3. Direct kernel interface: write to `/sys/power/state`

---

### Shutdown/Restart Fallback Chain

**For all systems**:
1. `shutdown -h now` (Shutdown)
2. `shutdown -r now` (Restart)
3. `systemctl poweroff` or `systemctl reboot` (systemd)
4. `loginctl poweroff` or `loginctl reboot` (systemd-logind)
5. `halt -p` or `reboot` (POSIX)

---

## Part 4: Implementation Architecture

### Detection Module
```go
type Environment struct {
    WindowManager   string      // "Hyprland", "Sway", "KDE", etc.
    DisplayServer   string      // "X11", "Wayland", "Unknown"
    InitSystem      string      // "systemd", "openrc", "runit", etc.
    AvailableTools  map[string]bool // cached command availability
}

func detectEnvironment() Environment
func (e Environment) hasCommand(cmd string) bool
```

### Command Builder Module
```go
func (e Environment) getCommandsForAction(action string) []string
// Returns: []string of commands to try in order for this action
// First respects WM/DE-specific commands, then falls back to generic chains
```

### Executor Module
```go
func executeCommandChain(commands []string) error
// Tries each command in order, returns on first success
// Logs which command was attempted and its result
```

---

## Part 5: Configuration System

### Overview

The tool will use a **default configuration** (compiled into binary) as fallback, with optional user overrides via JSON config file.

**Config file location** (in order of precedence):
1. `$BYEBYE_CONFIG` (environment variable)
2. `~/.config/byebye/config.json`
3. `/etc/byebye/config.json` (system-wide)
4. Compiled-in defaults

---

### Default Configuration Structure

The default config will be embedded as a JSON constant in the binary, containing all WM/DE-specific mappings and fallback chains from Parts 2-3.

```json
{
  "version": "1.0",
  "detectionPriority": [
    "Hyprland",
    "Sway",
    "KDE",
    "GNOME",
    "XFCE",
    "i3",
    "LXDE",
    "generic"
  ],
  "windowManagers": {
    "Hyprland": {
      "displayServer": "Wayland",
      "lock": ["hyprctl dispatch exec 'hyprlock'"],
      "logout": ["hyprctl dispatch exit"],
      "suspend": ["hyprlock &", "systemctl suspend"],
      "hibernate": ["hyprlock &", "systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    },
    "Sway": {
      "displayServer": "Wayland",
      "lock": ["swaylock -c 000000 -e"],
      "logout": ["swaymsg exit"],
      "suspend": ["swaylock -c 000000 &", "systemctl suspend"],
      "hibernate": ["swaylock -c 000000 &", "systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    },
    "KDE": {
      "displayServer": "auto",
      "lock": ["qdbus org.freedesktop.ScreenSaver /ScreenSaver Lock"],
      "logout": ["qdbus org.kde.ksmserver /KSMServer logout 0 1 1"],
      "suspend": ["qdbus org.freedesktop.PowerManagement /org/freedesktop/PowerManagement Suspend"],
      "hibernate": ["qdbus org.freedesktop.PowerManagement /org/freedesktop/PowerManagement Hibernate"],
      "shutdown": ["qdbus org.kde.ksmserver /KSMServer logout 1 1 1"],
      "restart": ["qdbus org.kde.ksmserver /KSMServer logout 1 2 1"]
    },
    "GNOME": {
      "displayServer": "auto",
      "lock": ["loginctl lock-session", "gnome-screensaver-command -l"],
      "logout": ["gio launch gnome-session-logout", "loginctl terminate-session"],
      "suspend": ["systemctl suspend"],
      "hibernate": ["systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    },
    "XFCE": {
      "displayServer": "auto",
      "lock": ["xfce4-screensaver-command -l"],
      "logout": ["xfce4-session-logout"],
      "suspend": ["xfce4-power-manager-settings", "systemctl suspend"],
      "hibernate": ["xfce4-power-manager-settings", "systemctl hibernate"],
      "shutdown": ["xfce4-session-logout --halt"],
      "restart": ["xfce4-session-logout --reboot"]
    },
    "i3": {
      "displayServer": "X11",
      "lock": ["i3lock -c 000000", "i3lock-color -c 000000", "i3lock-fancy"],
      "logout": ["i3-msg exit"],
      "suspend": ["i3lock -c 000000 &", "systemctl suspend"],
      "hibernate": ["i3lock -c 000000 &", "systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    },
    "LXDE": {
      "displayServer": "X11",
      "lock": ["lxlock", "slock"],
      "logout": ["lxsession-logout"],
      "suspend": ["systemctl suspend"],
      "hibernate": ["systemctl hibernate"],
      "shutdown": ["lxsession-logout --shutdown"],
      "restart": ["lxsession-logout --reboot"]
    },
    "generic": {
      "displayServer": "auto",
      "lock": {
        "X11": ["i3lock -c 000000", "i3lock-color", "slock", "xdg-screensaver lock"],
        "Wayland": ["swaylock -c 000000", "gtklock", "waylock", "wlopm"]
      },
      "logout": ["loginctl terminate-session", "systemctl --user exit", "kill -15 $PPID"],
      "suspend": ["systemctl suspend", "loginctl suspend", "pm-suspend"],
      "hibernate": ["systemctl hibernate", "loginctl hibernate", "pm-hibernate"],
      "shutdown": ["shutdown -h now", "systemctl poweroff", "loginctl poweroff"],
      "restart": ["shutdown -r now", "systemctl reboot", "loginctl reboot"]
    }
  }
}
```

---

### User Configuration Override

Users can create `~/.config/byebye/config.json` to override any defaults:

#### Example 1: Custom lock screen command
```json
{
  "windowManagers": {
    "Hyprland": {
      "lock": ["swaylock -c 1e1e2e"]
    }
  }
}
```

#### Example 2: Custom suspend with notification
```json
{
  "windowManagers": {
    "Sway": {
      "suspend": [
        "notify-send 'Suspending...'",
        "swaylock -c 000000 &",
        "systemctl suspend"
      ]
    }
  }
}
```

#### Example 3: Add custom logout with cleanup
```json
{
  "windowManagers": {
    "generic": {
      "logout": [
        "pkill -f myapp",
        "loginctl terminate-session",
        "kill -15 $PPID"
      ]
    }
  }
}
```

#### Example 4: Override entire action for multiple WMs
```json
{
  "windowManagers": {
    "Hyprland": {
      "logout": ["hyprctl dispatch exit"]
    },
    "Sway": {
      "logout": ["swaymsg exit"]
    }
  }
}
```

#### Example 5: Custom fallback chain
```json
{
  "windowManagers": {
    "generic": {
      "lock": {
        "X11": ["my-custom-locker", "i3lock -c 000000", "slock"],
        "Wayland": ["my-wayland-locker", "swaylock", "gtklock"]
      }
    }
  }
}
```

---

### Configuration Merging Logic

**Merging rules** (lowest to highest priority):
1. Compiled-in defaults
2. System-wide config `/etc/byebye/config.json`
3. User config `~/.config/byebye/config.json`
4. Environment override `$BYEBYE_CONFIG`

**Action-level merging**:
- If user specifies commands for an action, they **replace** defaults completely
- If user specifies commands for a specific WM, other WMs use defaults
- If user specifies WM-specific override, it takes precedence over generic

Example:
```go
// Default has: ["cmd1", "cmd2", "cmd3"]
// User config has: ["userCmd1"]
// Result: ["userCmd1"]  // User replaces completely

// To extend (not replace), user should include defaults:
// User config has: ["userCmd1", "cmd1", "cmd2", "cmd3"]
// Result: ["userCmd1", "cmd1", "cmd2", "cmd3"]
```

---

### Implementation Details

### Configuration Module
```go
type Config struct {
    Version string
    WindowManagers map[string]WMConfig
}

type WMConfig struct {
    DisplayServer string
    Lock      []string
    Logout    []string
    Suspend   []string
    Hibernate []string
    Shutdown  []string
    Restart   []string
}

func loadConfig() Config
// 1. Start with compiled defaults
// 2. Merge system-wide config if exists
// 3. Merge user config if exists
// 4. Merge env override if set
// 5. Return merged config

func (c Config) getCommands(wm, action string) []string
// Returns commands for action on given WM
// Falls back to generic if WM not found
```

### File Handling
```go
func mergeConfigs(base, override Config) Config
// Deep merge, respecting override rules

func loadConfigFile(path string) (Config, error)
// Unmarshal JSON with validation

func validateConfig(c Config) error
// Ensure required fields exist, commands are valid
```

---

### Default Config Embedding

Store the default config as a const string in a separate file:

```go
// config_defaults.go
package main

const defaultConfig = `{
  "version": "1.0",
  ...
}`
```

Then load it:
```go
func getDefaultConfig() Config {
    json.Unmarshal([]byte(defaultConfig), &config)
    return config
}
```

---

### Documentation

Create `CONFIG.md` with:
- Configuration file format specification
- Examples for common use cases
- How to override specific actions
- How to extend default chains
- Validation rules and error messages

---

## Part 6: Implementation Priority

1. **Phase 1**: Detection module + default config struct (Hyprland, Sway, KDE)
2. **Phase 2**: Config file loading + merging logic
3. **Phase 3**: Add GNOME, XFCE, i3, LXDE to defaults
4. **Phase 4**: Command execution with user config
5. **Phase 5**: Error handling, logging, validation
6. **Phase 6**: Documentation + testing

---

## Expected Result

After implementation:
- ✅ Works out-of-box for most Linux users (sensible defaults)
- ✅ Hyprland users get optimized experience by default
- ✅ Sway/KDE/GNOME/XFCE/i3 users get native DE experience by default
- ✅ Power users can customize everything via JSON config
- ✅ All init systems supported (systemd, openrc, runit)
- ✅ Both X11 and Wayland work seamlessly
- ✅ Config file allows extending/replacing any action
- ✅ Clear error messages when actions unavailable
- ✅ Fallback chains ensure maximum success rate
