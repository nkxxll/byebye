package main

import (
    "testing"
)

func TestGetDisplay(t *testing.T)  {
    a := getDisplay()
    t.Log(a)
    if a != x11 && a != wayland {
        t.Errorf("The Display is wrong")
    }
}
