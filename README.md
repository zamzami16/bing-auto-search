# Bing Auto Search - Configuration Guide

## Folder Structure Example

```text
bing-auto-search/
├── bin/
│   ├── desktop-bing-auto.exe
│   └── config/
│       └── config.json
├── internal/
│   └── ...
├── scripts/
│   └── build-windows.ps1
├── README.md
└── ...
```

> **Place your `config.json` in `bin/config/` so it is always found relative to the executable.**

## Example `config.json`

```json
{
  "global_setting": {
    "delay": { "min": 3, "max": 7 },
    "total_search": 2,
    "scroll": { "min": 500, "max": 1000 },
    "total_scroll": { "min": 2, "max": 5 }
  },
  "data": [
    {
      "name": "desktop-1",
      "configs": [{ "pos_x": 400, "pos_y": 50, "d_x": 150, "d_y": 5 }]
    },
    {
      "name": "desktop-2",
      "configs": [
        { "pos_x": 1300, "pos_y": 50, "d_x": 150, "d_y": 5 },
        { "pos_x": 500, "pos_y": 200, "d_x": 150, "d_y": 5 }
      ]
    },
    {
      "name": "desktop-3",
      "configs": [{ "pos_x": 800, "pos_y": 100, "d_x": 150, "d_y": 5 }]
    }
  ]
}
```

## Configuration Fields

- **global_setting.delay.min/max**: Minimum and maximum delay (seconds) between actions.
- **global_setting.total_search**: Number of search cycles per desktop/view.
- **global_setting.scroll.min/max**: Scroll amount per action (mouse wheel notches).
- **global_setting.total_scroll.min/max**: Number of scroll actions per search.
- **data**: List of desktops. Each desktop has a name and a list of view configs.
- **configs.pos_x/pos_y**: Center position for mouse actions (pixels).
- **configs.d_x/d_y**: Random offset range for mouse position (pixels).

## Tips

- Adjust `delay` for more human-like timing.
- Add or remove desktops/views as needed for your workflow.
- Make sure the config path is correct and relative to the executable (see folder structure above).

---

For more details, see the code comments or ask for advanced configuration help.
