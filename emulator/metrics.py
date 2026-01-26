import dearpygui.dearpygui as dpg

class TMTMetrics:
    def __init__(self, parent):
        self.parent = parent

        with dpg.collapsing_header(label="Simulation State", default_open=True):
            dpg.add_text("This is where information about the simulation will be printed", wrap=300, tag="state_text") 
        with dpg.collapsing_header(label="Metrics", default_open=True):
            dpg.add_text("Coord: (?, ?)", wrap=300, tag="coord_text")          


    def update_state_metrics(self, state=None):
        if state is None:
            return
        
        dpg.set_value(
            "state_text",
            f"Iteration: {state.get("Iteration", "?")}\nTurn: {state.get("Turn", "?")}\n"
            )
    
    def update_coord(self, coord=None):
        if coord is None:
            return

        dpg.set_value("coord_text", f"Coord: {coord}")

