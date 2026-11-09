# Install agent-canvas into GOPATH/bin so the command is on PATH.
install:
    go install ./cmd/agent-canvas

# Copy SKILL.md to ~/.agents/skills/agent-canvas/SKILL.md.
install-skill:
    mkdir -p "$HOME/.agents/skills/agent-canvas"
    cp SKILL.md "$HOME/.agents/skills/agent-canvas/SKILL.md"

# Build the VS Code extension and install it into Cursor or VS Code.
install-extension:
    go build -o editors/vscode/bin/agent-canvas ./cmd/agent-canvas
    cd editors/vscode && npm install
    cd editors/vscode && npm test
    cd editors/vscode && npx --yes @vscode/vsce package --allow-missing-repository -o agent-canvas.vsix
    if command -v cursor >/dev/null; then cursor --install-extension editors/vscode/agent-canvas.vsix --force; \
    elif command -v code >/dev/null; then code --install-extension editors/vscode/agent-canvas.vsix --force; \
    else echo "cursor or code is not on PATH" >&2; exit 1; fi
