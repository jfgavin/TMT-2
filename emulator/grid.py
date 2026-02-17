import dearpygui.dearpygui as dpg
from pathlib import Path

from emulator.gobuild import EMU_DIR

ASSET_DIR = EMU_DIR / Path("assets")


class TMTGrid:
    GRAVE_TEXTURE_TAG = "grave_img"

    def __init__(self, parent, grid_size, init_state):
        self.parent = parent
        self.grid_size = grid_size

        self.agents = []
        self.graves = []
        self.drawlist_tag = "grid"

        self._load_textures()
        self._create_drawlist()
        self.draw_blank_grid(init_state)

    # ------------------------------------------------------------------ #
    # Initialization Helpers
    # ------------------------------------------------------------------ #

    def _load_textures(self):
        """Load required textures once."""
        grave_path = ASSET_DIR / "grave.png"
        width, height, _, data = dpg.load_image(str(grave_path))

        with dpg.texture_registry():
            dpg.add_static_texture(
                width, height, data, tag=self.GRAVE_TEXTURE_TAG
            )

    def _create_drawlist(self):
        parent_width = dpg.get_item_width(self.parent)
        parent_height = dpg.get_item_height(self.parent)

        with dpg.drawlist(
            parent=self.parent,
            tag=self.drawlist_tag,
            width=parent_width,
            height=parent_height,
        ):
            pass

    # ------------------------------------------------------------------ #
    # Grid Rendering
    # ------------------------------------------------------------------ #

    def draw_blank_grid(self, state=None):
        if state is None:
            return

        dpg.delete_item(self.drawlist_tag, children_only=True)

        cell_w, cell_h = self._get_cell_size()

        for y in range(self.grid_size):
            for x in range(self.grid_size):
                self._draw_cell(x, y, cell_w, cell_h)

    def _draw_cell(self, x, y, cell_w, cell_h):
        x0, y0 = x * cell_w, y * cell_h
        x1, y1 = x0 + cell_w, y0 + cell_h

        dpg.draw_rectangle(
            (x0, y0),
            (x1, y1),
            color=(0, 0, 0),
            parent=self.drawlist_tag,
            tag=f"cell-{x}-{y}",
        )

    # ------------------------------------------------------------------ #
    # State Updates
    # ------------------------------------------------------------------ #

    def update_grid(self, state=None):
        if state is None:
            return

        cell_w, cell_h = self._get_cell_size()

        self._update_resources(state.get("Resources", []))
        self._update_graves(state.get("Graves", []), cell_w, cell_h)
        self._update_agents(state.get("Agents", {}), cell_w, cell_h)

    def _update_resources(self, resources):
        for x, y, val in resources:
            fill = (255, 255, 0, val * 10)

            cell_tag = f"cell-{x}-{y}"
            if dpg.does_item_exist(cell_tag):
                dpg.configure_item(cell_tag, fill=fill)

    def _update_graves(self, graves, cell_w, cell_h):
        self._clear_items(self.graves)

        for x, y, _ in graves:
            x0, y0 = x * cell_w, y * cell_h
            x1, y1 = x0 + cell_w, y0 + cell_h

            grave_sprite = dpg.draw_image(
                self.GRAVE_TEXTURE_TAG,
                (x0, y0),
                (x1, y1),
                parent=self.drawlist_tag,
            )

            self.graves.append(grave_sprite)

    def _update_agents(self, agents, cell_w, cell_h):
        self._clear_items(self.agents)

        for uuid, agent in agents.items():
            pos = agent.get("Pos", {})
            x = pos.get("X")
            y = pos.get("Y")

            if x is None or y is None:
                continue

            center_x = (x + 0.5) * cell_w
            center_y = (y + 0.5) * cell_h
            radius = cell_w * 0.5

            tag = f"agent-{uuid}"

            agent_circle = dpg.draw_circle(
                (center_x, center_y),
                radius,
                fill=(255, 0, 0),
                parent=self.drawlist_tag,
                tag=tag,
            )

            self.agents.append(agent_circle)

    # ------------------------------------------------------------------ #
    # Utilities
    # ------------------------------------------------------------------ #

    def colour_agent(self, uuid=None):
        if uuid is None:
            return

        for agent in self.agents:
            dpg.configure_item(agent, fill=(255, 0, 0))

        selected_tag = f"agent-{uuid}"
        if dpg.does_item_exist(selected_tag):
            dpg.configure_item(selected_tag, fill=(0, 255, 0))

    def _clear_items(self, items):
        for item in items:
            if dpg.does_item_exist(item):
                dpg.delete_item(item)
        items.clear()

    def _get_cell_size(self):
        canvas_width = dpg.get_item_width(self.drawlist_tag)
        canvas_height = dpg.get_item_height(self.drawlist_tag)
        return (
            canvas_width / self.grid_size,
            canvas_height / self.grid_size,
        )
