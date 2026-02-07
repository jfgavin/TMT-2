import os
import argparse
import dearpygui.dearpygui as dpg
from emulator.emulator import TMTEmulator
from emulator.build_go import build_go_binary, BIN_PATH
from src.config.gen_config import generate_config, CONF_DIR



def main():
    # Check build flag
    parser = argparse.ArgumentParser()
    parser.add_argument("-b", "--build", action="store_true", help="Build Go binary before running")
    parser.add_argument("-gc", "--gen-config", action="store_true", help="Regenerate config.go from config.JSON")
    args = parser.parse_args()

    if args.gen_config:
        generate_config()

    if args.build or not BIN_PATH.exists():
        build_go_binary()

    if os.environ.get("CI") == "true" or os.environ.get("DISPLAY") is None:
        print("Emulator running in CI - exiting early...")
        return

    # Run dpg emulator
    dpg.create_context()

    emulator = TMTEmulator()
    emulator.viewport_config()

    dpg.start_dearpygui()

if __name__ == "__main__":
    main()