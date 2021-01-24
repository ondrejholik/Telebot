# Telebot
- Telegram bot for everyday usage.

## Install
- Build from source
- Prerequisite -> go version 1.14.0
- `git clone https://github.com/ondrejholik/telebot`
- `cd telebot`
- `go build`
- Change Token in config.toml to your bot token
- If you want special birthdays feature rename assets/example.birthdays.csv to assets/birthdays.csv and add your birthdays list

## Run
- `./telebot`

## Commands
- /start -- initialize your info into database
- /bd -- birthdays 
[] /w -- weather based on location
[] /s -- split

## Other features
- Send location -- set location to your coordinates
