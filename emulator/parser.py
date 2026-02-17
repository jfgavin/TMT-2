import json
import socket
import subprocess
from emulator.gobuild import build_go_binary, EMU_DIR, BIN_PATH
from pathlib import Path
from datetime import datetime

class TMTParser():
    def __init__(self):
        self.HOST = "127.0.0.1"
        self.PORT = 5000

        self.metadata = ""
        self.states = self._parse_bin_json()

    def _parse_bin_json(self):
                
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.bind((self.HOST, self.PORT))
        sock.listen(1)
        print(f"Listening on {self.HOST}:{self.PORT} ...")

        if not BIN_PATH.exists():
            build_go_binary()

        assert BIN_PATH.exists()

        go_proc = subprocess.Popen([str(BIN_PATH)])
        conn, addr = sock.accept()

        got_meta = False

        states = []
        with conn:
            buffer = ""
            while True:
                data = conn.recv(4096)
                if not data:
                    break
                buffer += data.decode("utf-8")
                while "\n" in buffer:
                    line, buffer = buffer.split("\n", 1)
                    try:
                        if not got_meta:
                            self.metadata = json.loads(line.strip())
                            got_meta = True
                        else:
                            state = json.loads(line.strip())
                            states.append(state)
                    except json.JSONDecodeError:
                        pass

        print(f"Received {len(states)} game states.")
        conn.close()
        sock.close()

        return states

    def get_state(self, index=0):
        return self.states[index]

    def save_simulation(self):
        date = datetime.today().strftime("%Y-%m-%d_%H-%M-%S")
        filename = f"sim_{date}.json"

        path = EMU_DIR / "saves" / filename
        with open(path, "w") as f:
            json.dump(self.states, f, indent=2)

        print(f"Saved simulation to {str(path)}")

    def get_grid_size(self):
        return self.metadata.get("GridSize", 0)
