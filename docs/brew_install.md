## Installing using Homebrew on Mac/Linux

You will need to install [Homebrew](https://brew.sh/) first.

### Install

First install the custom tap, then trust it. Homebrew 6.0+ refuses to load
formulae from third-party taps until they are explicitly trusted.

```
brew tap muquit/markdown-serve https://github.com/muquit/markdown-serve.git
brew trust muquit/markdown-serve
brew install markdown-serve
```

Or tap, trust and install in one go:
```
brew tap muquit/markdown-serve https://github.com/muquit/markdown-serve.git
brew trust muquit/markdown-serve
brew install muquit/markdown-serve/markdown-serve
```

### Upgrade
```
brew upgrade markdown-serve
```

### Uninstall
```
brew uninstall markdown-serve
```

### Remove the tap
```
brew untap muquit/markdown-serve
```
