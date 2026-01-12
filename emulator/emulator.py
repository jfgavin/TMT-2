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
GRID_SIZE = 12
CELL_SIZE = 60
SIDEBAR_SIZE = 220
SIDEBAR_PADDING = 16
SCREEN_SIZE = GRID_SIZE * CELL_SIZE
ARROW_PADDING = 30
screen = pygame.display.set_mode((SCREEN_SIZE + SIDEBAR_SIZE, SCREEN_SIZE))
pygame.display.set_caption("TMT 2.0 Emulator")
clock = pygame.time.Clock()
font = pygame.font.SysFont("monospace", 24, bold=True)
button_font = pygame.font.SysFont("monospace", 18, bold=True)
sidebar_font = pygame.font.SysFont("monospace", 16)

# ---------------------------
# Load assets
# ---------------------------
ASSETS_DIR = Path("pygame_assets") / "textures"
TILE_IMAGES = {}
HERO_IMAGES = {}

def load_scaled_image(path, size):
    img = pygame.image.load(path).convert_alpha()
    return pygame.transform.scale(img, size)

# Tile textures
for name in ["floor"]:
    TILE_IMAGES[name] = load_scaled_image(ASSETS_DIR / f"{name}.png", (CELL_SIZE, CELL_SIZE))


# Easy-access floor image
floor_image = TILE_IMAGES["floor"]

# ---------------------------
# Socket listener setup
# ---------------------------
HOST = "127.0.0.1"
PORT = 5000

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.bind((HOST, PORT))
sock.listen(1)
print(f"Listening on {HOST}:{PORT} ...")

go_proc = subprocess.Popen(["../bin/tmt_bin"])

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

if states:
    firstState = states[0]
    if "Grid" in firstState:
        firstGrid = firstState.get("Grid", [])[0]
        GRID_SIZE = len(firstGrid)
        SCREEN_SIZE = GRID_SIZE * CELL_SIZE
        screen = pygame.display.set_mode((SCREEN_SIZE + SIDEBAR_SIZE, SCREEN_SIZE))

# ---------------------------
# Helper functions
# ---------------------------
def draw_arrow(rect, direction="left", color=(255, 153, 0)):
    x, y, w, h = rect
    if direction == "left":
        points = [(x + w, y), (x + w, y + h), (x, y + h / 2)]
    elif direction == "right":
        points = [(x, y), (x, y + h), (x + w, y + h / 2)]
    elif direction == "up":
        points = [(x + w / 2, y), (x, y + h), (x + w, y + h)]
    elif direction == "down":
        points = [(x, y), (x + w, y), (x + w / 2, y + h)]
    else:
        # default to right if unknown
        points = [(x, y), (x, y + h), (x + w, y + h / 2)]
    pygame.draw.polygon(screen, color, points)

def arrow_clicked(rect, mouse_pos):
    x, y, w, h = rect
    return x <= mouse_pos[0] <= x + w and y <= mouse_pos[1] <= y + h

# ---------------------------
# Main display loop
# ---------------------------
current_index = 0
running = True
while running:
    mouse = pygame.mouse.get_pos()
    click = False
    layer_delta = 0

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
        elif event.type == pygame.MOUSEBUTTONDOWN and event.button == 1:
            click = True
        elif event.type == pygame.KEYDOWN and states:
            # Single-step navigation
            if event.key in (pygame.K_RIGHT, pygame.K_d):
                current_index = min(current_index + 1, len(states) - 1)
            elif event.key in (pygame.K_LEFT, pygame.K_a):
                current_index = max(current_index - 1, 0)
            # Jump 10 states at a time
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
    selected_hero = None

    # Choose which board to render: new API provides `layers` (list of boards)
    grid = None
    if "Grid" in state:
        grid = state.get("Grid", [])

    for y, row in enumerate(grid):
        for x, tile in enumerate(row):
            rect = pygame.Rect(x * CELL_SIZE, y * CELL_SIZE, CELL_SIZE, CELL_SIZE)
            pos = (x * CELL_SIZE, y * CELL_SIZE)

            # --- Base floor ---
            screen.blit(floor_image, pos)

            # --- Tile border ---
            pygame.draw.rect(screen, (50, 50, 50), rect, 1)


    # Game State Arrows
    arrow_w, arrow_h = 40, 40
    left_rect = (SCREEN_SIZE + SIDEBAR_SIZE / 2 - arrow_w - ARROW_PADDING, SCREEN_SIZE - arrow_h - ARROW_PADDING, arrow_w, arrow_h)
    right_rect = (SCREEN_SIZE + SIDEBAR_SIZE / 2 + ARROW_PADDING, SCREEN_SIZE - arrow_h - ARROW_PADDING, arrow_w, arrow_h)
    draw_arrow(left_rect, "left")
    draw_arrow(right_rect, "right")

    # Sidebar
    sidebar_x = SCREEN_SIZE + SIDEBAR_PADDING
    sidebar_y = SIDEBAR_PADDING
    sidebar_lines = [
        f"Iteration: {state.get('Iteration', '?')}",
        f"Turn: {state.get('Turn', '?')}",
    ]

    for idx, text in enumerate(sidebar_lines):
        rendered = sidebar_font.render(text, True, (255, 255, 255))
        screen.blit(rendered, (sidebar_x, sidebar_y + idx * SIDEBAR_PADDING))

    pygame.display.flip()
    clock.tick(30)

    if CI_MODE:
        break

pygame.quit()
