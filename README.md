# YouTimeTrack

[![build](https://github.com/bullshitsoftware/youtimetrack/actions/workflows/ci.yml/badge.svg)](https://github.com/bullshitsoftware/youtimetrack/actions/workflows/ci.yml)
[![License: WTFPL](https://img.shields.io/badge/License-WTFPL-brightgreen.svg)](http://www.wtfpl.net/about/)
[![Coverage Status](https://coveralls.io/repos/github/bullshitsoftware/youtimetrack/badge.svg?branch=master)](https://coveralls.io/github/bullshitsoftware/youtimetrack?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bullshitsoftware/youtimetrack)](https://goreportcard.com/report/github.com/bullshitsoftware/youtimetrack)
[![LoC](https://tokei.rs/b1/github/bullshitsoftware/youtimetrack)](https://github.com/bullshitsoftware/youtimetrack)
[![Release](https://img.shields.io/github/release/bullshitsoftware/youtimetrack.svg?style=flat-square)](https://github.com/bullshitsoftware/youtimetrack/releases/latest)

CLI tools to help you manage your timetracks at YouTrack.

Main features:

- Check timetracks against your calendar settings, which include custom
  - Working day duration
  - Half-holiday duration
  - Weekends
  - Half-holiday
  - Public holidays
  - Vacation
  - And even extra working days!
- Show details about your timetracks
- Add/delete timetracks
- Written in Rust

## Install

Download suitable release from [here](https://github.com/bullshitsoftware/youtimetrack/releases) and place needed
binaries somewhere in your system.

## Usage

### Configuration

Initialize application config file with `yttconf` command or place
[sample config](https://github.com/bullshitsoftware/youtimetrack/blob/master/internal/app/config_stub.json) to 
`$HOME/.config/ytt/config.json`. 

```
❯ bin/yttconf
Created /home/user/.config/ytt/config.json edit it with your favorite text editor
```

Adjust it to your needs.

### Check your summary stats

Usage

```
❯ bin/yttsum --help
Usage of bin/yttsum:
  -end value
        end date (2006-01-02), default to current month end
  -start value
        start date (2006-01-02), default to current month start
```

Example

```
❯ bin/yttsum 
50h / 48h / 168h (worked / today / month)
```

### Check details

Usage

```
❯ bin/yttdet --help
Usage of bin/yttdet:
  -end value
        end date (2006-01-02), default to current month end
  -start value
        start date (2006-01-02), default to current month start
```

Example

```
❯ bin/yttdet --start 2022-07-08
2022-07-08      1h      XY-777  Do somethig good
116-74718               Did something good
```

### Log time

Usage

```
❯ bin/yttadd --help
Usage of bin/yttadd type issue duration comment:
  - type, work type prefix (e.g., develop)
  - issue, issue number (e.g., XY-123)
  - duration, spent time in YouTrack format (e.g., 1h 30m)
  - comment, work item description (e.g., did something cool)
```

Example

```
❯ bin/yttadd deve XY-777 "1h 30m" "did something"
Time tracked
```

### Delete time

Usage
```
❯ bin/yttdel --help
Usage of bin/yttdel issue item_id:
  - issue, issue number (e.g., XY-123)
  - item_id, work item id (e.g., 110-12312)
```

Example
```
❯ bin/yttdel XY-777 116-74718
Item deleted
```
