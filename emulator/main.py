import dearpygui.dearpygui as dpg
from emulator import TMTEmulator

dpg.create_context()

emulator = TMTEmulator()
emulator.viewport_config()

dpg.start_dearpygui()