
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

## Screenshots

TODO
