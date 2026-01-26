import dearpygui.dearpygui as dpg
import os
from grid import TMTGrid
from parser import TMTParser

class TMTEmulator():
    def __init__(self):
        self.WINDOW_WIDTH, self.WINDOW_HEIGHT = 1250, 950
        self.SIDEBAR_WIDTH = 350

        self.PADDING = 10
        self.MARGIN = 20

        self.CI_MODE = os.environ.get("CI", "false").lower() == "true"

        self.parser = TMTParser()

        self.__init_window()

    def __init_window(self):
        sim_size = min(self.WINDOW_HEIGHT - 2*self.PADDING, self.WINDOW_WIDTH - self.SIDEBAR_WIDTH - 3*self.PADDING)


        with dpg.window(label="TMT") as self.TMTWindow:
            with dpg.child_window(
                label="Simulation",
                tag="sim",
                height=sim_size,
                width=sim_size,
                no_scrollbar=True,
                border=False,
                pos=[self.PADDING, self.PADDING]
            ) as simulation_window:
                self.grid = TMTGrid(parent=simulation_window)
                state = self.parser.get_state(0)
                self.grid.draw_grid(state)

            with dpg.child_window(label="Sidebar", width=self.SIDEBAR_WIDTH, pos=[sim_size + 2*self.PADDING, self.PADDING], border=False, no_scrollbar=True):
                with dpg.collapsing_header(label="Controls", default_open=True) as self.controls_panel:
                    pass
                with dpg.collapsing_header(label="Info", default_open=True):
                    self.infotext = dpg.add_text("infofinfo", wrap=300)

            

        dpg.set_primary_window(self.TMTWindow, True)

    def viewport_config(self):
        # Necessary to set the window size
        dpg.create_viewport()
        dpg.setup_dearpygui()
        dpg.show_viewport()
        dpg.set_viewport_title("TMT Simulator 2.0")
        dpg.set_viewport_width(self.WINDOW_WIDTH)
        dpg.set_viewport_height(self.WINDOW_HEIGHT)
