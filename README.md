# M😝ji - An Emoji Programming Language

Moji is a fun programming language that uses emojis instead of traditional keywords. It's designed to make programming more visual and intuitive.

## Features

- 🎁 Variable declarations
- 📢 Print statements
- 🔀 If statements
- ↩️ Else clauses
- 🔄 While loops
- ⚖️ Equality comparisons
- ▶️ Greater than
- ◀️ Less than
- ✅ True
- ⛔️ False

## Example

```lox
🎁 age 👉 25;
🔀 (age ⚖️ 18) {
    📢 "You are an adult!";
} ↩️ {
    📢 "You are a minor!";
}
```

## Running the Interpreter

To run a Moji script:

```bash
go run . path/to/script.mji
```

Or start the interactive prompt:

```bash
go run .
```

## Development

This is a Go implementation of the Moji programming language. The interpreter is built using:

- Scanner (Lexer)
- Parser
- Evaluator

## License

MIT License
