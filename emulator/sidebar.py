import dearpygui.dearpygui as dpg

class TMTSidebar:
    def __init__(self, save_sim=None):
        with dpg.collapsing_header(label="Simulation State", default_open=True):
            dpg.add_text("This is where information about the simulation will be printed", wrap=300, tag="state_text") 
            if save_sim is not None:
                dpg.add_button(label="Save Simulation", callback=save_sim)
        with dpg.collapsing_header(label="Metrics", default_open=True):
            dpg.add_text("Coord: (?, ?)", wrap=300, tag="coord_text")          
        with dpg.collapsing_header(label="Agent", tag="agent_header", default_open=True, show=False):
            dpg.add_text("", tag="agent_name")
            dpg.add_text("", tag="agent_energy")

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

    def update_agent(self, uuid=None, agent=None):
        if uuid is None or agent is None:
            return
        
        dpg.configure_item("agent_header", show=True)
        dpg.set_value("agent_name", "Name: " + agent["Name"])
        dpg.set_value("agent_energy", f"Energy: {agent["Energy"]}")
        

