import json
from pathlib import Path
from src.config.gen_config import generate_config, JSON_FILE
from emulator.build_go import build_go_binary, BIN_PATH
from emulator.parser import TMTParser

RESULT_DIR = Path(__file__).resolve().parent / Path("_results")

def main():
    base = json.loads(JSON_FILE.read_text())

    agents = [4, 8, 16]

    RESULT_DIR.mkdir(exist_ok=True)

    for i, na in enumerate(agents):
        # Edit JSON field
        cfg = base.copy()
        cfg["Serv"]["NumAgents"] = na

        # Write back
        with open(JSON_FILE, "w") as f:
            f.write(json.dumps(cfg, indent=2))

        # Regenerate config
        generate_config()

        # Rebuild binary
        build_go_binary()

        # Parse output
        parser = TMTParser()

        # Print num states
        print(parser.get_num_states())
    
if __name__ == "__main__":
    main()
