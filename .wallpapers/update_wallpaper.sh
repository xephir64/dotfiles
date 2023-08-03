#!/bin/sh

DIR=$(dirname "$0")
HYPRLAND_CONFIG=$HOME/.config/hypr/hyprland.conf

wallpapers=()

for entry in "$DIR"/*
do
	if [ "$entry" != "$0" ]
	then
		wallpapers+=("$(basename "$entry")")
	fi
done

for i in "${!wallpapers[@]}"
do
	echo "$((i+1)). ${wallpapers[i]}"
done

read -p "Enter the number corresponding to the wallpaper that you want to apply: " selection

regex='^[0-9]+$'
if ! [[ $selection =~ $regex ]] || [ "$selection" -le 0 ] || [ "$selection" -gt "${#wallpapers[@]}" ]
then
	echo "Please enter a valid number."
	exit 1
fi

selected_wallpaper="${wallpapers[$((selection-1))]}"

if [ -f "$HYPRLAND_CONFIG" ]
then
	sed -i "s|^exec = wbg .*$|exec = wbg ~/.wallpapers/$selected_wallpaper|" "$HYPRLAND_CONFIG"
	echo "Wallpaper updated to $selected_wallpaper"
else
	echo "hyprland.conf not found in ~/.config/hypr/"
fi
