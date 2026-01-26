import dearpygui.dearpygui as dpg

class TMTGrid:
    def __init__(self, parent, initState):

        self.parent = parent
        self.cell_ids = []

        # Drawlist size will be set dynamically
        parent_width = dpg.get_item_width(parent)
        parent_height = dpg.get_item_height(parent)

        self.drawlist_tag = "grid"

        with dpg.drawlist(parent=parent, tag=self.drawlist_tag,
                            width=parent_width, height=parent_height):
            pass

        # Draw initial grid
        self.__draw_blank_grid(initState)

    def __draw_blank_grid(self, state=None):
        canvas_width = dpg.get_item_width(self.drawlist_tag)
        canvas_height = dpg.get_item_height(self.drawlist_tag)
        if canvas_width <= 0 or canvas_height <= 0:
            return  # Skip drawing until valid

        if state is not None and "Grid" in state:

            grid = state.get("Grid", [])
            self.GRID_SIZE = len(grid[0])

            # Clear old rectangles
            self.cell_ids.clear()

            # Draw Grid
            cell_width = canvas_width / self.GRID_SIZE
            cell_height = canvas_height / self.GRID_SIZE

            # Colour Resources
            for y, row in enumerate(grid):
                row_ids = []
                for x, cell in enumerate(row):
                    x0 = x * cell_width
                    y0 = y * cell_height
                    x1 = x0 + cell_width
                    y1 = y0 + cell_height

                    resources = cell.get("Resources", 0)
                    rgba = (255, 255, 0, resources * 10)

                    cell_id = dpg.draw_rectangle(
                        (x0, y0), (x1, y1),
                        fill=(255, 255, 0, resources*10),
                        parent=self.drawlist_tag
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

        # Update grid
        grid = state["Grid"]
        for y, row in enumerate(grid):
            for x, cell in enumerate(row):
                resources = cell.get("Resources", 0)
                rgba = (255, 255, 0, resources*10)

                # Update rectangle color without recreating it
                dpg.configure_item(self.cell_ids[y][x], fill=rgba)

        