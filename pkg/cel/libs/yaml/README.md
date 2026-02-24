# YAML Library for CEL

This library provides YAML parsing capabilities in CEL expressions, allowing you to parse YAML from strings, local files, or URLs.

## Functions

### `yaml.Unmarshal(data string) -> dyn`

Parses YAML data from a string and returns the parsed structure.

**Example:**
```cel
yaml.Unmarshal("name: John\nage: 30")
// Returns: {"name": "John", "age": 30.0}
```

### `yaml.UnmarshalFromFile(path string) -> dyn`

Reads and parses a YAML file from the local filesystem.

**Example:**
```cel
yaml.UnmarshalFromFile("/etc/config.yaml")
// Returns the parsed content of the file
```

### `yaml.UnmarshalFromURL(url string) -> dyn`

Fetches and parses a YAML file from a URL.

**Example:**
```cel
yaml.UnmarshalFromURL("https://example.com/config.yaml")
// Returns the parsed content from the URL
```

## Usage in Policies

You can use the YAML library in authorization policies to parse configuration files or remote YAML resources:

```cel
// Parse a config file and check a value
yaml.UnmarshalFromFile("/etc/app/config.yaml").database.host == "localhost"

// Fetch remote configuration
yaml.UnmarshalFromURL("https://config.example.com/settings.yaml").feature_flags.enabled

// Parse inline YAML
yaml.Unmarshal("items:\n  - apple\n  - banana").items.size() > 0
```

## Notes

- Numbers in YAML are parsed as `float64` by default
- The library uses `sigs.k8s.io/yaml` for parsing, which handles both JSON and YAML formats
- Network requests for `UnmarshalFromURL` use the standard `http.Get` with no timeout configuration
- File paths for `UnmarshalFromFile` must be absolute or relative to the working directory
- All functions return errors as CEL errors if parsing or fetching fails
