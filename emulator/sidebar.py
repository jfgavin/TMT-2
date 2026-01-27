import dearpygui.dearpygui as dpg

class TMTSidebar:
    def __init__(self, save_sim=None):
        with dpg.collapsing_header(label="Simulation State", default_open=True):
            dpg.add_text("This is where information about the simulation will be printed", wrap=300, tag="state_text") 
            if save_sim is not None:
                dpg.add_button(label="Save Simulation", callback=save_sim)
        with dpg.collapsing_header(label="Tile Info", tag="tile", default_open=True):
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
        tile = grid[y][x]

        for key, value in tile.items():
            dpg.add_text(f"{key}: {value}", parent="tile")

    def update_agent(self, uuid=None, agent=None):
        if uuid is None or agent is None:
            return
        
        dpg.configure_item("agent_header", show=True)
        dpg.set_value("agent_name", "Name: " + agent["Name"])
        dpg.set_value("agent_energy", f"Energy: {agent["Energy"]}")
        

