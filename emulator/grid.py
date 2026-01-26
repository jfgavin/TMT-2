import dearpygui.dearpygui as dpg
import numpy as np

class TMTGrid:
    def __init__(self, parent, initState):

        self.parent = parent

        # Canvas size
        canvas_width = dpg.get_item_width(parent)
        canvas_height = dpg.get_item_height(parent)

        # Get grid
        grid = initState.get("Grid", [])
        self.GRID_SIZE = len(grid[0])

        # Texture data
        self.texture_data = np.zeros((self.GRID_SIZE, self.GRID_SIZE, 4), dtype=np.float32)

        # Create texture registry and dynamic texture
        with dpg.texture_registry(show=False):
            self.texture_id = dpg.add_dynamic_texture(
                width=self.GRID_SIZE,
                height=self.GRID_SIZE,
                default_value=self.texture_data
            )

        # Drawlist to show texture (also draw agents here)
        self.drawlist_tag = "grid_drawlist"
        with dpg.drawlist(parent=parent, width=canvas_width, height=canvas_height, tag=self.drawlist_tag):
            dpg.draw_image(
                self.texture_id,
                pmin=(0, 0),
                pmax=(canvas_width, canvas_height)
            )

        # Draw grid lines once
        self._draw_grid_lines()


        # Initialize texture
        self.update_texture(grid)
        

    def _draw_grid_lines(self):
        canvas_width = dpg.get_item_width(self.drawlist_tag)
        canvas_height = dpg.get_item_height(self.drawlist_tag)
        cell_width = canvas_width / self.GRID_SIZE
        cell_height = canvas_height / self.GRID_SIZE

        # Vertical lines
        for x in range(self.GRID_SIZE + 1):
            dpg.draw_line(
                (x * cell_width, 0),
                (x * cell_width, canvas_height),
                color=(0, 0, 0, 255),
                parent=self.drawlist_tag
            )

        # Horizontal lines
        for y in range(self.GRID_SIZE + 1):
            dpg.draw_line(
                (0, y * cell_height),
                (canvas_width, y * cell_height),
                color=(0, 0, 0, 255),
                parent=self.drawlist_tag
            )

    def update_texture(self, grid):
        for y, row in enumerate(grid):
            for x, cell in enumerate(row):
                resources = cell.get("Resources", 0)
                alpha = min(1.0, resources / 25)
                self.texture_data[y, x] = [1.0, 1.0, 0.0, alpha]

        # Upload to DearPyGui
        dpg.set_value(self.texture_id, self.texture_data)

    def update_grid(self, state=None):
        if state is None or "Grid" not in state:
            return

        grid = state["Grid"]
        self.update_texture(grid)

        canvas_width = dpg.get_item_width(self.drawlist_tag)
        canvas_height = dpg.get_item_height(self.drawlist_tag)
        cell_width = canvas_width / self.GRID_SIZE
        cell_height = canvas_height / self.GRID_SIZE

        # Draw Agents
        agents = state.get("Agents", {})
        for uuid, agent in agents.items():
            dpg.delete_item(f"agent-{uuid}")
            pos = agent["Pos"]
            agx = (pos["X"] + 0.5) * cell_width
            agy = (pos["Y"] + 0.5) * cell_height
            radius = cell_width * 0.5

            dpg.draw_circle(
                (agx, agy),
                radius,
                color=(255, 0, 0),
                fill=(255, 0, 0),
                parent=self.drawlist_tag,
                tag=f"agent-{uuid}",
            )