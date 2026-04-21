# Troubleshooting

## Quick diagnostics

1. Check `cli.exe --status`.
2. Confirm `logs/wrapper.log` is updating.

## Expected messages

### `[START] Write param successful`

Normal signal that launch parameters were applied.

### `[START] Process exit with code 0`

Normal process termination (including manual game close).

## Common issues

### Game starts with no preset effect

- verify `Install` completed successfully;
- ensure antivirus did not remove/block `service.exe`;
- run `Uninstall` -> `Install` again.

### Auto tune selected a wrong preset

- run `Apply Recommended Config` again;
- verify hardware profile is unchanged (RAM/CPU mode);
- switch preset manually via `Select Config` if needed.

## When opening an issue

Attach:

- `logs/wrapper.log`;
- environment details (launcher/Steam/EGS, selected preset, reproduction steps).
