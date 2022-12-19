# Duke
 
A tool for creating desktop applications using `Go` with an emphasis on performance

## Installation

```bash
$: go get https://github.com/aacebo/duke
```

## Development

### Chromium

[Install Google Depot Tools](https://chromium.googlesource.com/chromium/src/+/main/docs/mac_build_instructions.md#install)

#### Fetch Chromium Source Code

```bash
$: cd chromium && caffeinate fetch chromium
```

#### Build Chromium From Source

```bash
$: gn gen out/Default && autoninja -C out/Default chrome
```
