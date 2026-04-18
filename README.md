# ynab-mcp

An [MCP](https://modelcontextprotocol.io) server for [YNAB](https://www.ynab.com/) (You Need A Budget). Connect Claude Desktop, Cursor, or any other MCP client to your YNAB account to query budgets, transactions, categories, and more.

## Install

### Homebrew (macOS)

```sh
brew install --cask JadoJodo/tap/ynab-mcp
```

The cask handles macOS Gatekeeper quarantine for you.

### Download a binary (Linux, Windows, or macOS without Homebrew)

Grab the archive for your OS/arch from the [Releases page](https://github.com/JadoJodo/ynab-mcp/releases), extract it, and move `ynab-mcp` somewhere on your `$PATH` (e.g. `/usr/local/bin`).

On macOS, a manually downloaded binary is quarantined by Gatekeeper. Remove the attribute after extracting:

```sh
xattr -d com.apple.quarantine ./ynab-mcp
```

The Homebrew cask above does this automatically.

### Build from source

```sh
go install github.com/JadoJodo/ynab-mcp@latest
```

## Configure

You'll need a YNAB personal access token — create one at [app.ynab.com/settings/developer](https://app.ynab.com/settings/developer).

Point your MCP client at the installed binary and pass the token via the `YNAB_API_TOKEN` environment variable. Example for Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "ynab": {
      "command": "/usr/local/bin/ynab-mcp",
      "env": {
        "YNAB_API_TOKEN": "your-token-here"
      }
    }
  }
}
```

Use the absolute path to the binary (`which ynab-mcp` will tell you). Restart Claude Desktop and the YNAB tools should appear.

## License

MIT — see [LICENSE](LICENSE).
