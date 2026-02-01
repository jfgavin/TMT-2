import json
import socket
import subprocess
from gobuild import build_go_binary, BIN_PATH
from pathlib import Path
from datetime import datetime

class TMTParser():
    def __init__(self):
        self.HOST = "127.0.0.1"
        self.PORT = 5000

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

        path = Path.cwd() / "saves" / filename
        with open(path, "w") as f:
            json.dump(self.states, f, indent=2)

    def get_grid_size(self):
        init_state = self.states[0]

        grid = init_state.get("Grid", [])
        return len(grid[0])
