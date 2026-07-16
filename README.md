# 🎨 ASCII Art Generator

A simple command-line tool written in **Go** that converts text strings into stylized **ASCII art banners**. It supports multiple font styles and can render multi-line input.

---

## ✨ Features

- 🖋️ Converts any printable ASCII text into large ASCII-art banners
- 🎭 Supports multiple banner styles: `standard`, `shadow`, and `thinkertoy`
- 📄 Handles multi-line input (via `\n`)
- ⚠️ Clear error handling for invalid characters or unknown styles
- ⚡ Lightweight, no external dependencies

---

## 📦 Requirements

- [Go](https://go.dev/dl/) 1.18 or later
- Banner font files: `standard.txt`, `shadow.txt`, `thinkertoy.txt` (must be present in the project directory)

---

## 🚀 Installation

Clone the repository and navigate into the project folder:

```bash
git clone https://01.nextera.education/git/yelkahky/ascii-art
cd ascii-art
```

Make sure the banner files (`standard.txt`, `shadow.txt`, `thinkertoy.txt`) are located in the same directory as `main.go`.

---

## ▶️ Usage

Run the program using `go run .` followed by the text you want to convert, and an optional banner style:

```bash
go run . "<string>" [standard|shadow|thinkertoy]
```

### Arguments

| Argument     | Required | Description                                                    |
|--------------|----------|------------------------------------------------------------------|
| `<string>`   | ✅ Yes   | The text to convert into ASCII art                              |
| `style`      | ❌ No    | Banner style: `standard` (default), `shadow`, or `thinkertoy`   |

> 💡 If no style is provided, `standard` is used by default.

---

## 📝 Examples

**Basic usage (default style):**
```bash
go run . "Hello"
```

**Using the shadow style:**
```bash
go run . "Hello" shadow
```

**Using the thinkertoy style:**
```bash
go run . "Hello" thinkertoy
```

**Multi-line text:**
```bash
go run . "Hello\nWorld"
```

---

## ⚙️ How It Works

1. The program reads the input string and optional style from command-line arguments.
2. It loads the corresponding banner font file (e.g. `standard.txt`).
3. Each character is mapped to its ASCII-art representation using its printable ASCII index.
4. The text is rendered line-by-line, row-by-row, to build the final banner output.

---

## 🚫 Error Handling

The program returns clear errors in the following cases:

- ❌ Missing or too many arguments
- ❌ Unknown banner style (must be `standard`, `shadow`, or `thinkertoy`)
- ❌ Unreadable or missing banner file
- ❌ Non-printable ASCII characters in the input

---

## 📁 Project Structure

```
.
├── main.go          # Main program logic
├── standard.txt      # Standard banner font
├── shadow.txt         # Shadow banner font
├── thinkertoy.txt     # Thinkertoy banner font
└── README.md
```

---

## 👥 Authors

- **Youssef Elkahky**
- **Mohamed Elgohary**
- **Mohanad Ali**

