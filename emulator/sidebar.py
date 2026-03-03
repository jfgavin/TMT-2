import dearpygui.dearpygui as dpg
import numpy as np

class TMTSidebar:
    def __init__(self, save_sim=None, dt=0.1):
        self.dt = dt
        with dpg.collapsing_header(label="Simulation State", tag="sim_state", default_open=True):
            dpg.add_text(f"Iteration: 0", tag="iter_text")
            dpg.add_text(f"Turn: 0", tag="turn_text")
            if save_sim is not None:
                dpg.add_button(label="Save Simulation", callback=save_sim)
        with dpg.collapsing_header(label="Tile Info", tag="tile", default_open=True):
            dpg.add_text("Coord: (?, ?)", wrap=300, tag="coord_text")          
        with dpg.collapsing_header(label="Agent", tag="agent_header", default_open=True, show=False):
            dpg.add_text("", tag="agent_name")
            dpg.add_text("", tag="agent_energy")
        with dpg.plot(tag="model_plot", label="TMT Model Output", width=320, height=280, show=False):
            # optionally create legend
            dpg.add_plot_legend()

            # REQUIRED: create x and y axes
            dpg.add_plot_axis(dpg.mvXAxis, label="x", tag="x_axis")
            dpg.add_plot_axis(dpg.mvYAxis, label="y", tag="y_axis")

            # series belong to a y axis
            dpg.add_line_series([0], [0], tag="tmt_line_series", label="Model Output", parent="y_axis")


    def update_state_metrics(self, state=None):
        if state is None:
            return

        iter_num = state.get("Iteration", "?")
        turn_num = state.get("Turn", "?")

        dpg.set_value("iter_text", f"Iteration: {iter_num}")
        dpg.set_value("turn_text", f"Turn: {turn_num}")
    
    def update_coord(self, coord=None, valid=False, state=None):
        if coord is None:
            return

        # Check if coord has changed
        if dpg.does_item_exist("coord"):
            old_coord = dpg.get_item_user_data("coord")
            if coord == old_coord:
                return

        # Delete old children
        dpg.delete_item("tile", children_only=True)

        # Add coordinate
        dpg.add_text(f"Coord: {coord}", tag="coord", parent="tile", user_data=coord)

        # Add tile metrics
        if not valid or state is None or "Grid" not in state:
            return

        grid = state["Grid"]
        x, y = coord
        try:
            tile = grid[y][x]
        except:
            # Coordinate is not valid within the grid
            return

        for key, value in tile.items():
            dpg.add_text(f"{key}: {value}", parent="tile")

    def update_agent(self, uuid=None, agent=None):
        dpg.delete_item("agent_header", children_only=True)

        if uuid is None:
            dpg.configure_item("agent_header", show=False)
            dpg.configure_item("model_plot", show=False)
            return
        elif agent is None:
            # Agent is dead if UUID selected but not in state
            dpg.add_text("DEAD", parent="agent_header")
        else:
            for key, value in agent.items():
                if key != "ModelOutput":
                    dpg.add_text(f"{key}: {value}", parent="agent_header")

            model_output = agent["ModelOutput"]
            x_axis = np.arange(0, len(model_output) * self.dt, self.dt)
            dpg.set_value("tmt_line_series", [x_axis, model_output])
            dpg.fit_axis_data("x_axis")
            dpg.fit_axis_data("y_axis")

        dpg.configure_item("agent_header", show=True)
        dpg.configure_item("model_plot", show=True)
