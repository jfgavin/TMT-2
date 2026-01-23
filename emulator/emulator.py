import argparse
import json
import socket
import pygame
import subprocess
import os
import sys
from pathlib import Path

# ---------------------------
# Assert CI mode
# ---------------------------
CI_MODE = os.environ.get("CI", "false").lower() == "true"

# ---------------------------
# Pygame setup
# ---------------------------
pygame.init()

# Initial grid dimension (this will update once JSON is loaded)
GRID_SIZE = 12

# Sidebar constants
SIDEBAR_RATIO = 0.25
SIDEBAR_PADDING = 16
ARROW_PADDING = 30

# Create a resizable window.
# We'll allocate main grid area dynamically based on whatever size the window is.
INITIAL_WINDOW_SIZE = 600
screen = pygame.display.set_mode(
    (INITIAL_WINDOW_SIZE * (1 + SIDEBAR_RATIO), INITIAL_WINDOW_SIZE),
    pygame.RESIZABLE
)

pygame.display.set_caption("TMT 2.0 Emulator")
clock = pygame.time.Clock()

# Fonts
font = pygame.font.SysFont("monospace", 24, bold=True)
button_font = pygame.font.SysFont("monospace", 18, bold=True)
sidebar_font = pygame.font.SysFont("monospace", 16)

# ---------------------------
# Socket listener setup
# ---------------------------
HOST = "127.0.0.1"
PORT = 5000

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.bind((HOST, PORT))
sock.listen(1)
print(f"Listening on {HOST}:{PORT} ...")

bin_path = Path.cwd().parent / Path("bin/tmt_bin")
print(bin_path)

if not bin_path.exists():
    print("Binary not found! Please build Go binary first...")
    exit(0)

go_proc = subprocess.Popen([str(bin_path)])

conn, addr = sock.accept()

# ---------------------------
# Read incoming JSON lines
# ---------------------------
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
                continue

print(f"Received {len(states)} game states.")
conn.close()
sock.close()


# --- Update GRID_SIZE based on first state ---
if states:
    firstState = states[0]
    if "Grid" in firstState:
        firstGrid = firstState.get("Grid", [])[0]
        GRID_SIZE = len(firstGrid)

# ---------------------------
# Main display loop
# ---------------------------
current_index = 0
running = True
while running:
    mouse = pygame.mouse.get_pos()
    click = False

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
        elif event.type == pygame.KEYDOWN and states:
            if event.key in (pygame.K_RIGHT, pygame.K_d):
                current_index = min(current_index + 1, len(states) - 1)
            elif event.key in (pygame.K_LEFT, pygame.K_a):
                current_index = max(current_index - 1, 0)
            elif event.key in (pygame.K_PAGEDOWN, pygame.K_PERIOD):
                current_index = min(current_index + 10, len(states) - 1)
            elif event.key in (pygame.K_PAGEUP, pygame.K_COMMA):
                current_index = max(current_index - 10, 0)
    if not states:
        screen.fill((0, 0, 0))
        pygame.display.flip()
        continue

    state = states[current_index]
    screen.fill((77, 77, 77))

    
    grid = state.get("Grid", [])

    # --- Dynamic scaling ---
    WINDOW_W, WINDOW_H = screen.get_size()
    SIDEBAR_W = int(WINDOW_W * SIDEBAR_RATIO)
    GRID_W = WINDOW_W - SIDEBAR_W
    GRID_H = WINDOW_H
    CELL_W = GRID_W / GRID_SIZE
    CELL_H = GRID_H / GRID_SIZE

    # --- Draw grid background ---
    pygame.draw.rect(screen, (255, 255, 255), (0, 0, WINDOW_W, WINDOW_H))

     # --- Draw grid lines ---
    line_color = (50, 50, 50)
    for y in range(GRID_SIZE + 1):
        pygame.draw.line(screen, line_color,
            (0, y * CELL_H),
            (WINDOW_W, y * CELL_H)
        )
    for x in range(GRID_SIZE + 1):
        pygame.draw.line(screen, line_color,
            (x * CELL_W, 0),
            (x * CELL_W, WINDOW_H)
        )

    # --- Draw agents ---
    agents = state.get("Agents", {})
    for uuid, agent in agents.items():
        pos = agent["Pos"]
        px = (pos["X"] + 0.5) * CELL_W
        py = (pos["Y"] + 0.5) * CELL_H
        radius = min(CELL_W, CELL_H) * 0.45
        pygame.draw.circle(screen, (255, 0, 0), (px, py), radius)

    # --- Sidebar background ---
    pygame.draw.rect(screen, (30, 30, 30), (GRID_W, 0, SIDEBAR_W, WINDOW_H))

    # --- Sidebar text ---
    sidebar_x = GRID_W + SIDEBAR_PADDING
    sidebar_y = SIDEBAR_PADDING
    sidebar_lines = [
        f"Iteration: {state.get('Iteration', '?')}",
        f"Turn: {state.get('Turn', '?')}",
    ]
    for idx, text in enumerate(sidebar_lines):
        rendered = sidebar_font.render(text, True, (255, 255, 255))
        screen.blit(rendered, (sidebar_x, sidebar_y + idx * 30))

    pygame.display.flip()
    clock.tick(30)

    if CI_MODE:
        break

pygame.quit()
