# byebye Configuration

`byebye` loads defaults that cover common Linux window managers and desktop
environments. You can override them with JSON config files.

## Config Locations

The config loader merges files in this order (later overrides earlier):

1. Built-in defaults
2. `/etc/byebye/config.json`
3. `~/.config/byebye/config.json`
4. `$BYEBYE_CONFIG`

## Format

```json
{
  "version": "1.0",
  "windowManagers": {
    "Hyprland": {
      "displayServer": "Wayland",
      "lock": ["hyprctl dispatch exec 'hyprlock'"],
      "logout": ["hyprctl dispatch exit"],
      "sleep": ["systemctl suspend"],
      "suspend": ["systemctl suspend"],
      "hibernate": ["systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    }
  }
}
```

`lock` can be either a list of commands or a map keyed by display server:

```json
{
  "windowManagers": {
    "generic": {
      "lock": {
        "X11": ["i3lock -c 000000", "slock"],
        "Wayland": ["swaylock -c 000000", "gtklock"]
      }
    }
  }
}
```

## Override Rules

- If you set an action list, it replaces the defaults for that action.
- If you override a single WM, other WMs keep their defaults.
- WM-specific overrides take precedence over the `generic` section.

## Examples

Custom lock screen for Hyprland:

```json
{
  "windowManagers": {
    "Hyprland": {
      "lock": ["swaylock -c 1e1e2e"]
    }
  }
}
```

Add a notification before suspend in Sway:

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

Add a sleep action override:

```json
{
  "windowManagers": {
    "generic": {
      "sleep": ["systemctl suspend"]
    }
  }
}
```

Override generic logout chain:

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
