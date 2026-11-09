import { existsSync } from "node:fs";
import { dirname, join } from "node:path";

import * as vscode from "vscode";

import { renderCanvas } from "./render";

const viewType = "agentCanvas.preview";

export function activate(context: vscode.ExtensionContext): void {
  context.subscriptions.push(
    vscode.window.registerCustomEditorProvider(viewType, new CanvasPreview(context), {
      webviewOptions: { retainContextWhenHidden: true },
      supportsMultipleEditorsPerDocument: false,
    }),
    vscode.commands.registerCommand("agentCanvas.openPreview", async (uri?: vscode.Uri) => {
      const target = uri ?? vscode.window.activeTextEditor?.document.uri;
      if (!target) {
        return;
      }
      await vscode.commands.executeCommand("vscode.openWith", target, viewType);
    }),
    vscode.commands.registerCommand("agentCanvas.openSource", async () => {
      const input = vscode.window.tabGroups.activeTabGroup.activeTab?.input;
      if (!(input instanceof vscode.TabInputCustom) || input.viewType !== viewType) {
        return;
      }
      await vscode.commands.executeCommand("vscode.openWith", input.uri, "default", {
        viewColumn: vscode.ViewColumn.Beside,
        preserveFocus: false,
      });
    }),
  );
}

export function deactivate(): void {}

class CanvasPreview implements vscode.CustomTextEditorProvider {
  constructor(private readonly context: vscode.ExtensionContext) {}

  resolveCustomTextEditor(
    document: vscode.TextDocument,
    webviewPanel: vscode.WebviewPanel,
  ): void {
    const webview = webviewPanel.webview;
    webview.options = { enableScripts: true };
    webview.html = shellHtml(webview);

    let generation = 0;
    let pending = document.getText();
    let timer: ReturnType<typeof setTimeout> | undefined;

    const schedule = (text: string) => {
      pending = text;
      if (timer) {
        clearTimeout(timer);
      }
      timer = setTimeout(() => {
        const ticket = ++generation;
        const source = pending;
        void renderCanvas(resolveExecutable(this.context), source).then(
          (html) => {
            if (ticket === generation) {
              void webview.postMessage({ type: "html", html });
            }
          },
          (err: unknown) => {
            if (ticket === generation) {
              void webview.postMessage({ type: "error", message: renderError(err) });
            }
          },
        );
      }, 80);
    };

    const subscriptions: vscode.Disposable[] = [
      vscode.workspace.onDidChangeTextDocument((event) => {
        if (event.document.uri.toString() === document.uri.toString()) {
          schedule(event.document.getText());
        }
      }),
    ];
    if (document.uri.scheme === "file") {
      const watcher = vscode.workspace.createFileSystemWatcher(
        new vscode.RelativePattern(vscode.Uri.file(dirname(document.uri.fsPath)), basename(document.uri)),
      );
      subscriptions.push(
        watcher,
        watcher.onDidChange(async () => {
          if (document.isDirty) {
            return;
          }
          const bytes = await vscode.workspace.fs.readFile(document.uri);
          schedule(Buffer.from(bytes).toString("utf8"));
        }),
      );
    }

    webviewPanel.onDidDispose(() => {
      if (timer) {
        clearTimeout(timer);
      }
      for (const item of subscriptions) {
        item.dispose();
      }
    });

    schedule(document.getText());
  }
}

function resolveExecutable(context: vscode.ExtensionContext): string {
  const configured = vscode.workspace.getConfiguration("agentCanvas").get("executable", "").trim();
  if (configured) {
    return configured;
  }
  const name = process.platform === "win32" ? "agent-canvas.exe" : "agent-canvas";
  const bundled = join(context.extensionPath, "bin", name);
  if (existsSync(bundled)) {
    return bundled;
  }
  return name;
}

function renderError(err: unknown): string {
  const message = err instanceof Error ? err.message : String(err);
  if (message.includes("ENOENT")) {
    return "agent-canvas was not found. Run just install-extension from the repo, or set agentCanvas.executable.";
  }
  return message;
}

function basename(uri: vscode.Uri): string {
  const slash = uri.path.lastIndexOf("/");
  return slash === -1 ? uri.path : uri.path.slice(slash + 1);
}

function shellHtml(webview: vscode.Webview): string {
  const nonce = String(Date.now());
  const csp = [
    "default-src 'none'",
    `style-src 'unsafe-inline'`,
    `script-src 'nonce-${nonce}'`,
    `frame-src ${webview.cspSource}`,
  ].join("; ");
  return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="Content-Security-Policy" content="${csp}">
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
  html, body { margin: 0; height: 100%; background: #181818; color: #f0f0f0; }
  iframe { border: 0; width: 100%; height: 100%; background: #181818; }
  #error { display: none; margin: 0; padding: 16px; white-space: pre-wrap; font: 13px/1.45 ui-monospace, monospace; color: #f85149; }
  body.error iframe { display: none; }
  body.error #error { display: block; }
</style>
</head>
<body>
<iframe id="preview" title="Canvas preview"></iframe>
<pre id="error"></pre>
<script nonce="${nonce}">
  const frame = document.getElementById("preview");
  const error = document.getElementById("error");
  window.addEventListener("message", (event) => {
    const data = event.data;
    if (!data || (data.type !== "html" && data.type !== "error")) {
      return;
    }
    if (data.type === "html") {
      const top = frame.contentWindow ? frame.contentWindow.scrollY : 0;
      frame.onload = () => {
        if (frame.contentWindow) {
          frame.contentWindow.scrollTo(0, top);
        }
      };
      document.body.classList.remove("error");
      frame.srcdoc = data.html;
      return;
    }
    if (data.type === "error") {
      document.body.classList.add("error");
      error.textContent = data.message;
    }
  });
</script>
</body>
</html>`;
}
