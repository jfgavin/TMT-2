import subprocess
import sys
from pathlib import Path

HERE = Path(__file__).resolve().parent
PARENT = HERE.parent
BIN_PATH = PARENT / "bin" / "tmt_bin"

def build_go_binary():
    print("Building Go binary...")
    result = subprocess.run(
        ["go", "build", "-o", str(BIN_PATH), str(PARENT)],
        capture_output=True,
        text=True,
    )

    if result.returncode != 0:
        print("Go build failed:", file=sys.stderr)
        print(result.stderr, file=sys.stderr)
        sys.exit(result.returncode)
