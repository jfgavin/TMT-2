import dearpygui.dearpygui as dpg

class TMTControls:
    def __init__(self, start_value=0, index_change=None, mouse_move=None):
        self.value = start_value
        self.index_change = index_change
        self.mouse_move = mouse_move
        self._setup_key_handlers()

    def _setup_key_handlers(self):
        with dpg.handler_registry():
            dpg.add_key_press_handler(key=dpg.mvKey_Left,  callback=self._on_left)
            dpg.add_key_press_handler(key=dpg.mvKey_Right, callback=self._on_right)
            dpg.add_mouse_move_handler(callback=self.mouse_move)

    def _on_left(self, sender, app_data):
        if self.value > 0:
            self.value -= 1
        if self.index_change:
            self.index_change(self.value)

    def _on_right(self, sender, app_data):
        self.value += 1
        if self.index_change:
            self.index_change(self.value)

    def get_value(self):
        return self.value
