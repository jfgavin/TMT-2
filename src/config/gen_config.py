#!/usr/bin/env python3
import json
from pathlib import Path

CONF_DIR = Path(__file__).resolve().parent
INDENT = f"\t"

# Map JSON types to Go types
def go_type(value):
    if isinstance(value, int):
        return "int"
    if isinstance(value, float):
        return "float64"
    if isinstance(value, str):
        return "string"
    if isinstance(value, bool):
        return "bool"
    if isinstance(value, dict):
        return None  # nested struct
    if isinstance(value, list):
        return "[]" + go_type(value[0]) if value else "[]interface{}"
    return "interface{}"

def generate_struct(name, obj):
    lines = [f"type {name} struct {{"]
    nested = []

    f_len = len(max(obj, key=len))

    for field, v in obj.items():
        t = go_type(v)
        if t:
            lines.append(f"{INDENT}{field:<{f_len}} {t}")
        else:
            nested_name = field + "Config"
            lines.append(f"{INDENT}{field:<{f_len}} {nested_name}")
            nested.extend(generate_struct(nested_name, v))

    lines.append("}")
    lines.append("")
    return lines + nested

def format_value(v):
    if isinstance(v, str):
        return f"\"{v}\""
    if isinstance(v, bool):
        return "true" if v else "false"
    return str(v)

def generate_constructor_fields(obj, level):
    lines = []

    f_len = len(max(obj, key=len))

    for field, v in obj.items():
        if isinstance(v, dict):
            nested_name = field + "Config"
            lines.append(
                f"{INDENT * level}{field:<{f_len}}: {nested_name}{{"
            )
            lines.extend(generate_constructor_fields(v, level + 1))
            lines.append(f"{INDENT * level}}},")
        else:
            lines.append(
                f"{INDENT * level}{field:<{f_len}}: {format_value(v)},"
            )
    return lines

def generate_constructor(obj):
    lines = [
        "func NewConfig() Config {",
        f"{INDENT}return Config{{",
    ]

    for field, v in obj.items():
        struct_name = field + "Config"
        lines.append(f"{INDENT * 2}{field}: {struct_name}{{")
        lines.extend(generate_constructor_fields(v, 3))
        lines.append(f"{INDENT * 2}}},")

    lines.append(f"{INDENT}}}")
    lines.append("}")
    lines.append("")
    return lines

def generate_config():

    go_file = CONF_DIR / Path("config.go")
    json_file = CONF_DIR / Path("config.json")

    print("Generating Go Config from JSON...")

    with open(json_file) as f:
        config = json.load(f)

    lines = ["package config", ""]

    lines.extend(generate_struct("Config", config))
    lines.extend(generate_constructor(config))

    with open(go_file, "w") as f:
        f.write("\n".join(lines))


def main():
    generate_config()

if __name__ == "__main__":
    main()
