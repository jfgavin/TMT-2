import json
import struct
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

    def _read_message(self, conn):
        # Read 4-byte length prefix
        len_bytes = conn.recv(4)
        if not len_bytes:
            return None
        length = struct.unpack(">I", len_bytes)[0]

        # Read the JSON payload exactly
        chunks = []
        received = 0
        while received < length:
            chunk = conn.recv(length - received)
            if not chunk:
                raise ConnectionError("Connection closed mid-message")
            chunks.append(chunk)
            received += len(chunk)
        return json.loads(b"".join(chunks))

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
                try:
                    msg = self._read_message(conn)
                    if msg is None:
                        break
                    if msg["Type"] == "metadata":
                        self.metadata = msg["Data"]
                    elif msg["Type"] == "state":
                        states.append(msg["Data"])
                except ConnectionError:
                    break

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

    def get_metadata(self):
        return self.metadata
