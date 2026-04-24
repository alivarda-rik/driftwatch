# driftwatch

> A CLI tool that detects and reports configuration drift between deployed services and their declared state in YAML/TOML files.

---

## Installation

```bash
go install github.com/yourusername/driftwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftwatch.git
cd driftwatch
go build -o driftwatch .
```

---

## Usage

Point `driftwatch` at your declared config file and a running service endpoint to check for drift:

```bash
driftwatch check --config ./config/production.yaml --service http://localhost:8080/config
```

Compare two config files directly:

```bash
driftwatch diff --declared ./config/production.toml --live ./snapshots/live.toml
```

Watch for drift continuously and report changes:

```bash
driftwatch watch --config ./config/production.yaml --service http://localhost:8080/config --interval 30s
```

Example output:

```
[DRIFT DETECTED] 3 differences found in production.yaml
  ~ database.max_connections: declared=100, live=50
  ~ cache.ttl:                declared=300s, live=600s
  + feature_flags.dark_mode:  declared=true,  live=<missing>
```

---

## Configuration

`driftwatch` supports both YAML and TOML config formats. See [`docs/configuration.md`](docs/configuration.md) for a full reference.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)