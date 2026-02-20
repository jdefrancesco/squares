
# Squares

**Squares** is a simple fun/addictive game where the player start as a tiny square whose mission is to eat other squares with a smaller size. The catch is that your square grows after consuming each square.  becomes more more difficult.

## How to play

Squares and circles spawn from the edges and drift toward you as the difficulty ramps up.

- You are the rotating black square in the center.
- Move your square with the mouse.
- Eat smaller squares to grow and increase your score.
- Hitting a larger square ends the game.

### Circles

- **Green circle**: grants invincibility for a short time.
- **Black circle**: instant death on contact.

## Controls

- **Mouse**: move
- **Left click**: dash (short burst; has a cooldown)
- **P** or **Esc**: pause / resume
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
