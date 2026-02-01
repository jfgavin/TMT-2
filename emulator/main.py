import argparse
import subprocess
import sys
from pathlib import Path
from paths import PARENT, BIN_PATH
import dearpygui.dearpygui as dpg
from emulator import TMTEmulator

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

def main():
    # Check build flag
    parser = argparse.ArgumentParser()
    parser.add_argument("-b", "--build", action="store_true", help="Build Go binary before running")
    args = parser.parse_args()

    if args.build or not BIN_PATH.exists():
        build_go_binary()

    # Run dpg emulator
    dpg.create_context()

    emulator = TMTEmulator()
    emulator.viewport_config()

    dpg.start_dearpygui()

if __name__ == "__main__":
    main()