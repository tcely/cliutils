#!/bin/bash
# Restart Bluetooth (and optionally wifi too) then connect to paired devices.
#
# Usage:
#     restart-bluetooth.sh [wifi]
#
# Requires:
#     awk
#     networksetup
#     sleep
#     xargs
#     https://github.com/breiter/blueutil
#     https://github.com/lapfelix/BluetoothConnector

get_bt_state() {
    local _awk_program='{print $NF; exit;}'

    blueutil status | awk "$_awk_program"
}

[ "on" = "$(get_bt_state)" ] || exit 0

get_wifi_dev() {
    local _awk_program='$NF ~ /^Wi-Fi$/ {getline; print $2; exit;}'

    networksetup -listallhardwareports | awk "$_awk_program"
}

if [ "wifi" = "$1" ]; then
    blueutil off

    _wifi_dev="$(get_wifi_dev)"
    if [ -n "$_wifi_dev" ]; then
	networksetup -setairportpower "$_wifi_dev" off
	sleep 2
	networksetup -setairportpower "$_wifi_dev" on
    fi

    blueutil on
else
    blueutil restart
fi

BluetoothConnector | awk 'NR>=5 {print $1}' | xargs -P 2 -n 1 BluetoothConnector
