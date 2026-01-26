import dearpygui.dearpygui as dpg
import os
from grid import TMTGrid
from parser import TMTParser
from controls import TMTControls
from metrics import TMTMetrics

class TMTEmulator():
    def __init__(self):
        self.WINDOW_WIDTH, self.WINDOW_HEIGHT = 1440, 1080
        self.SIDEBAR_WIDTH = 350

        self.PADDING = 10
        self.MARGIN = 20

        self.CI_MODE = os.environ.get("CI", "false").lower() == "true"

        self.parser = TMTParser()
        self.controls = TMTControls(index_change=self._on_index_change, mouse_move=self._on_mouse_move)

        self.INDEX = 0
        self.GRID_SIZE = self.parser.get_grid_size()

        self._init_window()

        self._on_index_change(self.INDEX)

    def _init_window(self):
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
                self.grid = TMTGrid(parent=simulation_window, initState=self.parser.get_state(self.INDEX))

            with dpg.child_window(label="Sidebar", tag="side", width=self.SIDEBAR_WIDTH, pos=[sim_size + 2*self.PADDING, self.PADDING], border=False, no_scrollbar=True) as sidebar:
                
                self.metrics = TMTMetrics(parent=sidebar)


        dpg.set_primary_window(self.TMTWindow, True)

    # =====
    # Viewport Config
    # =====

    def viewport_config(self):
        # Necessary to set the window size
        dpg.create_viewport()
        dpg.setup_dearpygui()
        dpg.show_viewport()
        dpg.set_viewport_title("TMT 2.0 Simulator")
        dpg.set_viewport_width(self.WINDOW_WIDTH)
        dpg.set_viewport_height(self.WINDOW_HEIGHT)
        dpg.set_viewport_resize_callback(self._on_viewport_resize)


    def _on_viewport_resize(self, sender, app_data):
        new_width, new_height, _, _ = app_data

        # Compute new simulation size based on viewport size
        sim_size = min(new_height - 2*self.PADDING, new_width - self.SIDEBAR_WIDTH - 3*self.PADDING)

        # Resize simulation child window
        dpg.configure_item("sim", width=sim_size, height=sim_size)

        # Move sidebar to the new position
        dpg.configure_item("side", pos=[sim_size + 2*self.PADDING, self.PADDING])

        dpg.configure_item(self.grid.drawlist_tag, width=sim_size, height=sim_size)
        state = self.parser.get_state(self.INDEX)
        self.grid.draw_blank_grid(state)
        self.grid.update_grid(state)

    # =====
    # Top-level Controls
    # =====

    def _on_index_change(self, new_index):
        """
            When a new intex is given, the corresponding state is visualised
        """
        # prevent out-of-range access
        if new_index < 0 or new_index >= len(self.parser.states):
            return

        self.INDEX = new_index

        state = self.parser.get_state(self.INDEX)
        self.grid.update_grid(state)
        self.metrics.update_state_metrics(state)

    def _on_mouse_move(self):
        """
            Passes the clicked grid-coordinate, else nothing
        """
        mx, my = dpg.get_mouse_pos(local=False)
        sim_x, sim_y = dpg.get_item_pos("sim")
        sim_w = dpg.get_item_width("sim")
        sim_h = dpg.get_item_height("sim")

        corr_x = mx - sim_x
        corr_y = my - sim_y

        if  corr_x < 0 or corr_y < 0 or corr_x > sim_w or corr_y > sim_h:
            coord = "(?, ?)"
        else:
            x = int((corr_x / sim_w) * self.GRID_SIZE)
            y = int((corr_y / sim_h) * self.GRID_SIZE)
            coord = (x, y)

        self.metrics.update_coord(coord)
