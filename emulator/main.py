import argparse
from gobuild import build_go_binary, BIN_PATH
import dearpygui.dearpygui as dpg
from emulator import TMTEmulator


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