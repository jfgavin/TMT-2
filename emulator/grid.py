import dearpygui.dearpygui as dpg
import copy

class TMTGrid:
    def __init__(self, parent, gridSize, initState):

        self.parent = parent
        self.grid_size = gridSize
        self.agents = []
        self.graves = []

        # Drawlist size will be set dynamically
        parent_width = dpg.get_item_width(parent)
        parent_height = dpg.get_item_height(parent)

        self.drawlist_tag = "grid"

        with dpg.drawlist(parent=parent, tag=self.drawlist_tag,
                            width=parent_width, height=parent_height):
            pass

        # Draw initial grid
        self.draw_blank_grid(initState)

    def draw_blank_grid(self, state=None):
        if state is not None:
            gs = self.grid_size

            # Clear old rectangles
            dpg.delete_item(self.drawlist_tag, children_only=True)

            # Draw Grid
            cw, ch = self._get_cell_size()

            for y in range(gs):
                for x in range(gs):
                    x0 = x * cw
                    y0 = y * ch
                    x1 = x0 + cw
                    y1 = y0 + ch

                    cell_id = dpg.draw_rectangle(
                        (x0, y0), (x1, y1),
                        color=(0, 0, 0),
                        parent=self.drawlist_tag,
                        tag=f"cell-{x}-{y}",
                    )

    def update_grid(self, state=None):
        if state is None:
            return

        # Resources
        resources = state.get("Resources", [])

        for entry in resources:
            pos, amt = entry["Pos"], entry["Value"]
            x, y = pos["X"], pos["Y"]
            fill=(255, 255, 0, amt*10)
            try:
                dpg.configure_item(f"cell-{x}-{y}", fill=fill)
            except:
                continue

        cw, ch = self._get_cell_size()

        # Graves
        for grave in self.graves:
            dpg.delete_item(grave)
        self.graves = []

        graves = state.get("Graves", [])        
        for grave in graves:
            pos = grave["Pos"]
            agx = (pos["X"] + 0.5) * cw
            agy = (pos["Y"] + 0.5) * ch
            radius = cw * 0.5

            grave_circle = dpg.draw_circle(
                (agx, agy),
                radius,
                fill=(0, 255, 255),
                parent=self.drawlist_tag,
            )
            self.graves.append(grave_circle)

        # Agents
        for agent in self.agents:
            dpg.delete_item(agent)
        self.agents = []
        
        agents = state.get("Agents", [])

        for uuid, agent in agents.items():
            pos = agent["Pos"]
            agx = (pos["X"] + 0.5) * cw
            agy = (pos["Y"] + 0.5) * ch
            radius = cw * 0.5

            agent_circle = dpg.draw_circle(
                (agx, agy),
                radius,
                fill=(255, 0, 0),
                parent=self.drawlist_tag,
                tag=f"agent-{uuid}",
            )

            self.agents.append(agent_circle)


    def colour_agent(self, uuid=None):
        if uuid is None:
            return

        for agent in self.agents:
            dpg.configure_item(agent, fill=(255, 0, 0))

        dpg.configure_item(f"agent-{uuid}", fill=(0, 255, 0))

    def _get_cell_size(self):
        canvas_width = dpg.get_item_width(self.drawlist_tag)
        canvas_height = dpg.get_item_height(self.drawlist_tag)
        cell_width = canvas_width / self.grid_size
        cell_height = canvas_height / self.grid_size
        return cell_width, cell_height