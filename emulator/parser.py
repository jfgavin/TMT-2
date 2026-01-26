import json
import socket
import subprocess
from pathlib import Path

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

        bin_path = Path.cwd().parent / Path("bin/tmt_bin")

        if not bin_path.exists():
            print("Binary not found! Please build Go binary first...")
            exit(0)

        go_proc = subprocess.Popen([str(bin_path)])
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

    def get_grid_size(self):
        init_state = self.states[0]

        grid = init_state.get("Grid", [])
        return len(grid[0])
