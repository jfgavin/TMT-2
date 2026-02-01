import dearpygui.dearpygui as dpg
import copy

class TMTGrid:
    def __init__(self, parent, initState):

        self.parent = parent
        self.cell_ids = []
        self.agents = []

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
        canvas_width = dpg.get_item_width(self.drawlist_tag)
        canvas_height = dpg.get_item_height(self.drawlist_tag)
        if canvas_width <= 0 or canvas_height <= 0:
            return  # Skip drawing until valid

        if state is not None and "Grid" in state:

            grid = state.get("Grid", [])
            self.GRID_SIZE = len(grid[0])
            self._cell_content = [[{} for _ in range(self.GRID_SIZE)] for _ in range(self.GRID_SIZE)]


            # Clear old rectangles
            self.cell_ids.clear()
            dpg.delete_item(self.drawlist_tag, children_only=True)

            # Draw Grid
            cell_width = canvas_width / self.GRID_SIZE
            cell_height = canvas_height / self.GRID_SIZE

            # Colour Resources
            for y, row in enumerate(grid):
                row_ids = []
                for x, cell in enumerate(row):
                    self._cell_content[y][x] = copy.deepcopy(cell)

                    x0 = x * cell_width
                    y0 = y * cell_height
                    x1 = x0 + cell_width
                    y1 = y0 + cell_height

                    resources = cell.get("Resources", 0)
                    rgba = (255, 255, 0, resources * 10)

                    cell_id = dpg.draw_rectangle(
                        (x0, y0), (x1, y1),
                        color=(0, 0, 0),
                        fill=(255, 255, 0, resources*10),
                        parent=self.drawlist_tag,
                        tag=f"cell-{x}-{y}",
                    )
                    row_ids.append(cell_id)
                self.cell_ids.append(row_ids)

    def update_grid(self, state=None):
        if state is None or "Grid" not in state:
            return

        canvas_width = dpg.get_item_width(self.drawlist_tag)
        canvas_height = dpg.get_item_height(self.drawlist_tag)
        cell_width = canvas_width / self.GRID_SIZE
        cell_height = canvas_height / self.GRID_SIZE

        # Redraw Agents
        for agent in self.agents:
            dpg.delete_item(agent)
        self.agents.clear()

        freshAgents = state.get("Agents", {})
        for uuid, agent in freshAgents.items():
            pos = agent["Pos"]
            agx = (pos["X"] + 0.5) * cell_width
            agy = (pos["Y"] + 0.5) * cell_height
            radius = cell_width * 0.5

            agent_circle = dpg.draw_circle(
                (agx, agy),
                radius,
                fill=(255, 0, 0),
                parent=self.drawlist_tag,
                tag=f"agent-{uuid}",
            )

            self.agents.append(agent_circle)

        # Update grid
        grid = state["Grid"]
        for y, row in enumerate(grid):
            for x, cell in enumerate(row):
                if cell != self._cell_content[y][x]:
                    resources = cell.get("Resources", 0)
                
                    rgba = (255, 255, 0, resources*10)
                    dpg.configure_item(self.cell_ids[y][x], fill=rgba)
                    
                    self._cell_content[y][x] = copy.deepcopy(cell)

    def colour_agent(self, uuid=None):
        if uuid is None:
            return

        for agent in self.agents:
            dpg.configure_item(agent, fill=(255, 0, 0))

        dpg.configure_item(f"agent-{uuid}", fill=(0, 255, 0))