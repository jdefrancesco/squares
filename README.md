
# Squares

**Squares** is a mouse-only survival game: eat, grow, dodge, and beat your personal best. There are no power-ups or dash controls to learn.

## How to play

You start as a rotating green square, 22 pixels wide, in a 640×480 playfield. Your square follows the mouse and stays inside the playfield. Survive and eat as many squares as you can; there is no win condition.

- **Amber squares: eat.** Each smaller square gives **1 point**, grows the player, and triggers a brief pulse, "+1", and chime.
- **Dark squares: avoid.** Equal-sized or larger squares end the round on contact. A danger square turns amber when you grow larger than it.
- **Growth tapers.** The initial reward is **5% of the eaten square's side length + 0.8 pixels**, scaled down as your size approaches **64 pixels**. You can keep scoring without growing beyond that limit.
- **Rotation is visual.** Your collision box is an axis-aligned square at 90% of the player's side length. Edge contact counts, and moving the mouse quickly cannot skip a danger square.

### Pressure and fair spawning

- The opening two seconds spawn only food. Difficulty then rises with active time and score, reaching full pressure after `75 - 1.5 × score` seconds of active play.
- Spawn intervals decrease from roughly **0.8 seconds** to **0.32 seconds**, while squares get faster and dangerous squares become more common.
- Squares enter from the edges and head toward your position at spawn time. They travel in straight lines afterward. Small edge markers show where incoming squares will enter.
- New squares provide at least **0.85 seconds** before they can hit your current position. You can still collide sooner by moving toward one.
- Spawning limits the number of active dangers and checks nearby escape routes over the next **1.5 seconds**. If another threat would block those routes, it spawns food instead. These safeguards do not protect you from moving into existing danger.
- After **two consecutive danger spawns**, the next spawn is food. Eating a square also resets that counter.
- The HUD shows **Score** and **Best**. Death shows the result immediately; **one click** starts another round.

### Personal best

Your best score survives restarts and is saved when a round ends or the window closes. The record is stored in `Squares/best-score.json` under your OS user configuration directory (on macOS, `~/Library/Application Support/Squares/best-score.json`). If storage is unavailable, you can still play and retain a best score for the current session.

## Controls

- **Mouse**: move
- **Left click**: retry after game over (release a held click first)
- **P** or **Esc**: pause / resume; movement, spawns, difficulty, and countdown timers freeze while paused
- **R**: retry after game over
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
