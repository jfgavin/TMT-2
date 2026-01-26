import dearpygui.dearpygui as dpg

class TMTControls:
    def __init__(self, start_value=0, on_change=None):
        self.value = start_value
        self.on_change = on_change
        self._setup_key_handlers()

    def _setup_key_handlers(self):
        with dpg.handler_registry():
            dpg.add_key_press_handler(key=dpg.mvKey_Left,  callback=self._on_left)
            dpg.add_key_press_handler(key=dpg.mvKey_Right, callback=self._on_right)

    def _on_left(self, sender, app_data):
        self.value -= 1
        if self.on_change:
            self.on_change(self.value)

    def _on_right(self, sender, app_data):
        self.value += 1
        if self.on_change:
            self.on_change(self.value)

    def get_value(self):
        return self.value
