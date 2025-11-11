import json
import socket
import pygame
import subprocess
from pathlib import Path

# ---------------------------
# Pygame setup
# ---------------------------
pygame.init()
GRID_SIZE = 8
CELL_SIZE = 100
SIDEBAR_SIZE = 300
SIDEBAR_PADDING = 20
SCREEN_SIZE = GRID_SIZE * CELL_SIZE
screen = pygame.display.set_mode((SCREEN_SIZE + SIDEBAR_SIZE, SCREEN_SIZE))
pygame.display.set_caption("Manic Maell-Strom")
clock = pygame.time.Clock()
font = pygame.font.SysFont("monospace", 30, bold=True)
button_font = pygame.font.SysFont("monospace", 24, bold=True)
sidebar_font = pygame.font.SysFont("monospace", 20)

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
for name in ["floor", "portal", "shop", "shop_looted"]:
    TILE_IMAGES[name] = load_scaled_image(ASSETS_DIR / f"{name}.png", (CELL_SIZE, CELL_SIZE))
TILE_MAP = {"empty": "floor", "exit": "portal", "shop": "shop", "looted_shop": "shop_looted"}

# Hero textures
for name in ["dwarf", "elf", "troll", "orc"]:
    HERO_IMAGES[name[0].upper()] = load_scaled_image(ASSETS_DIR / f"{name}.png", (CELL_SIZE, CELL_SIZE))

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

go_proc = subprocess.Popen(["./TMT-2"])
print("Launched Go game...")

conn, addr = sock.accept()
print(f"Connected by {addr}")

# ---------------------------
# Read incoming JSON lines
# ---------------------------
grids = []
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
                grid = json.loads(line.strip())
                grids.append(grid)
            except json.JSONDecodeError:
                continue

print(f"Received {len(grids)} game states.")
conn.close()
sock.close()

# ---------------------------
# Helper functions
# ---------------------------
def draw_arrow(rect, direction="left", color=(255, 153, 0)):
    x, y, w, h = rect
    if direction == "left":
        points = [(x + w, y), (x + w, y + h), (x, y + h / 2)]
    else:
        points = [(x, y), (x, y + h), (x + w, y + h / 2)]
    pygame.draw.polygon(screen, color, points)

def arrow_clicked(rect, mouse_pos):
    x, y, w, h = rect
    return x <= mouse_pos[0] <= x + w and y <= mouse_pos[1] <= y + h

# ---------------------------
# Main display loop
# ---------------------------
current_index = 0
selected_hero = None
running = True
while running:
    mouse = pygame.mouse.get_pos()
    click = False

    for event in pygame.event.get():
        if event.type == pygame.QUIT:
            running = False
        elif event.type == pygame.MOUSEBUTTONDOWN and event.button == 1:
            click = True

    if not grids:
        screen.fill((0, 0, 0))
        pygame.display.flip()
        continue

    grid = grids[current_index]
    screen.fill((77, 77, 77))
    selected_hero = None

    for y, row in enumerate(grid["board"]):
        for x, tile in enumerate(row):
            rect = pygame.Rect(x * CELL_SIZE, y * CELL_SIZE, CELL_SIZE, CELL_SIZE)
            screen.blit(floor_image, (x * CELL_SIZE, y * CELL_SIZE))
            tile_type = tile["type"]
            if tile_type == "hero":
                hero_char = tile.get("name", "?")
                img = HERO_IMAGES.get(hero_char.upper())
                if rect.collidepoint(mouse):
                    for hero in grid["heroes"]:
                        if hero["name"][0].upper() == hero_char.upper():
                            selected_hero = hero
            else:
                img = TILE_IMAGES[TILE_MAP.get(tile_type, "floor")]
            if img:
                screen.blit(img, (x * CELL_SIZE, y * CELL_SIZE))
            pygame.draw.rect(screen, (50, 50, 50), rect, 1)

    # Arrows
    arrow_w, arrow_h = 40, 40
    padding = 30
    left_rect = (SCREEN_SIZE + SIDEBAR_SIZE / 2 - arrow_w - padding, SCREEN_SIZE - arrow_h - padding, arrow_w, arrow_h)
    right_rect = (SCREEN_SIZE + SIDEBAR_SIZE / 2 + padding, SCREEN_SIZE - arrow_h - padding, arrow_w, arrow_h)
    draw_arrow(left_rect, "left")
    draw_arrow(right_rect, "right")

    if click:
        if arrow_clicked(left_rect, mouse) and current_index > 0:
            current_index -= 1
        elif arrow_clicked(right_rect, mouse) and current_index < len(grids) - 1:
            current_index += 1

    # Sidebar
    sidebar_x = SCREEN_SIZE + SIDEBAR_PADDING
    sidebar_y = SIDEBAR_PADDING
    iteration_text = sidebar_font.render(f"Iteration: {grid.get('iteration', '?')}", True, (255, 255, 255))
    turn_text = sidebar_font.render(f"Turn: {grid.get('turn', '?')}", True, (255, 255, 255))
    screen.blit(iteration_text, (sidebar_x, sidebar_y))
    screen.blit(turn_text, (sidebar_x, sidebar_y + SIDEBAR_PADDING))

    if selected_hero:
        stats_y = sidebar_y + 3 * SIDEBAR_PADDING
        for key in ["name", "energy", "looted"]:
            if key in selected_hero:
                text = sidebar_font.render(f"{key.title()}: {selected_hero[key]}", True, (255, 255, 255))
                screen.blit(text, (sidebar_x, stats_y))
                stats_y += SIDEBAR_PADDING

    pygame.display.flip()
    clock.tick(30)

pygame.quit()
