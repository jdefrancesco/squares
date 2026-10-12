
# Squares

**Squares** is a survival game: move a small square with the mouse, eat smaller squares to grow, and avoid fatal collisions as the difficulty increases.

## How to play

You start as a rotating green square, 22 pixels wide, in a 640×480 playfield. Your square follows the mouse, with its center kept inside the playfield. Survive and eat as many squares as you can; there is no win condition.

- **Smaller black squares**: eat them by touching them. Each square gives **1 point** and increases your side length by **5% of the eaten square's side length + 0.8 pixels**.
- **Equal-sized or larger black squares**: touching them ends the game unless you have square protection from a dash or boost. While protected, you can eat squares of any size for the same score and growth rewards.
- **Collision boundaries**: the player's hitbox is an axis-aligned square with 90% of its visible side length. Rotation is visual; it does not rotate the hitbox. Touching hitbox edges counts as a collision.

### Circles

- **Teal circle**: grants **2.5 seconds** of protection against squares. Collecting another refreshes that timer; durations do not stack. Circles give no points or growth.
- **Black circle**: instant death on contact, even during dash or boost protection.

### Dash and difficulty

- **Dash** grants **0.25 seconds** of protection against squares, with a **0.9-second cooldown** between activations. It does not increase movement speed; your square still follows the mouse.
- Entities appear just outside a random edge, aimed roughly at your position when they spawn. They travel in straight lines afterward.
- Difficulty rises with both active play time and score. Spawns become more frequent (from roughly **0.85 seconds** apart down to **0.25 seconds**), hazards become more likely, boosts become less likely, and large squares become more threatening.
- After **six consecutive non-edible spawns**, the next spawn is forced to be a smaller square. Eating any square also resets that counter.
- The HUD shows your score and remaining square-protection time. A ring around the player indicates active protection.

## Controls

- **Mouse**: move
- **Left click**: activate dash protection (release and click again to reactivate once ready)
- **P** or **Esc**: pause / resume; movement, spawns, difficulty, and countdown timers freeze while paused
- **R**: restart (when game over)
- **Q**: quit

## Run / Build

Requires Go and a working graphics environment supported by [Ebiten](https://ebitengine.org/).

Run directly:

```sh
go run ./cmd/squares
```

Build a binary into `bin/`:

```sh
make build
./bin/squares
```

## macOS App Bundle

You can build a proper macOS `.app` bundle (so it shows up like a normal app in Finder / Launchpad):

```sh
make macos-app
open dist/Squares.app
```

You can override the bundle id and version:

```sh
make macos-app BUNDLE_ID=com.yourname.squares VERSION=1.0.0
```

### App icon

To include a Dock/Finder icon, add a 1024×1024 PNG at:

`assets/icon.png`

If you just want a reasonable placeholder icon quickly, generate one:

```sh
make icon
```

Then re-run `make macos-app`. The build script will generate an `.icns` and place it inside the app bundle.

If `assets/icon.png` is missing, the app will still build but will use the default generic icon.

### Notes on distribution

For distributing to other machines without Gatekeeper warnings, you typically need to code-sign and notarize the `.app` using an Apple Developer ID certificate. This repo currently only creates the unsigned `.app` bundle.

## Screenshots

TODO
