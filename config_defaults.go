package main

const defaultConfig = `{
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
    },
    "Sway": {
      "displayServer": "Wayland",
      "lock": ["swaylock -c 000000 -e"],
      "logout": ["swaymsg exit"],
      "sleep": ["swaylock -c 000000 &", "systemctl suspend"],
      "suspend": ["swaylock -c 000000 &", "systemctl suspend"],
      "hibernate": ["swaylock -c 000000 &", "systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    },
    "KDE": {
      "displayServer": "auto",
      "lock": ["qdbus org.freedesktop.ScreenSaver /ScreenSaver Lock"],
      "logout": ["qdbus org.kde.ksmserver /KSMServer logout 0 1 1"],
      "sleep": ["qdbus org.freedesktop.PowerManagement /org/freedesktop/PowerManagement Suspend"],
      "suspend": ["qdbus org.freedesktop.PowerManagement /org/freedesktop/PowerManagement Suspend"],
      "hibernate": ["qdbus org.freedesktop.PowerManagement /org/freedesktop/PowerManagement Hibernate"],
      "shutdown": ["qdbus org.kde.ksmserver /KSMServer logout 1 1 1"],
      "restart": ["qdbus org.kde.ksmserver /KSMServer logout 1 2 1"]
    },
    "GNOME": {
      "displayServer": "auto",
      "lock": ["loginctl lock-session", "gnome-screensaver-command -l"],
      "logout": ["gio launch gnome-session-logout", "loginctl terminate-session"],
      "sleep": ["systemctl suspend"],
      "suspend": ["systemctl suspend"],
      "hibernate": ["systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    },
    "XFCE": {
      "displayServer": "auto",
      "lock": ["xfce4-screensaver-command -l"],
      "logout": ["xfce4-session-logout"],
      "sleep": ["xfce4-power-manager-settings", "systemctl suspend"],
      "suspend": ["xfce4-power-manager-settings", "systemctl suspend"],
      "hibernate": ["xfce4-power-manager-settings", "systemctl hibernate"],
      "shutdown": ["xfce4-session-logout --halt"],
      "restart": ["xfce4-session-logout --reboot"]
    },
    "i3": {
      "displayServer": "X11",
      "lock": ["i3lock -c 000000", "i3lock-color -c 000000", "i3lock-fancy"],
      "logout": ["i3-msg exit"],
      "sleep": ["i3lock -c 000000 &", "systemctl suspend"],
      "suspend": ["i3lock -c 000000 &", "systemctl suspend"],
      "hibernate": ["i3lock -c 000000 &", "systemctl hibernate"],
      "shutdown": ["systemctl poweroff"],
      "restart": ["systemctl reboot"]
    },
    "LXDE": {
      "displayServer": "X11",
      "lock": ["lxlock", "slock"],
      "logout": ["lxsession-logout"],
      "sleep": ["systemctl suspend"],
      "suspend": ["systemctl suspend"],
      "hibernate": ["systemctl hibernate"],
      "shutdown": ["lxsession-logout --shutdown"],
      "restart": ["lxsession-logout --reboot"]
    },
    "generic": {
      "displayServer": "auto",
      "lock": {
        "X11": ["i3lock -c 000000", "i3lock-color -c 000000", "slock", "xdg-screensaver lock", "xlock", "xtrlock"],
        "Wayland": ["swaylock -c 000000", "gtklock", "waylock", "wlopm"]
      },
      "logout": ["loginctl terminate-session", "systemctl --user exit", "kill -15 $PPID"],
      "sleep": ["systemctl suspend", "loginctl suspend", "pm-suspend"],
      "suspend": ["systemctl suspend", "loginctl suspend", "pm-suspend"],
      "hibernate": ["systemctl hibernate", "loginctl hibernate", "pm-hibernate"],
      "shutdown": ["shutdown -h now", "systemctl poweroff", "loginctl poweroff"],
      "restart": ["shutdown -r now", "systemctl reboot", "loginctl reboot"]
    }
  }
}`
