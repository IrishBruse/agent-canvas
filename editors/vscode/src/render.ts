import { execFile } from "node:child_process";
import { mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";

export async function renderCanvas(executable: string, source: string): Promise<string> {
  const dir = await mkdtemp(join(tmpdir(), "agent-canvas-"));
  const input = join(dir, "input.canvas");
  const output = join(dir, "output.html");
  try {
    await writeFile(input, source);
    await new Promise<void>((resolve, reject) => {
      execFile(executable, [input, output], { timeout: 15000 }, (err, _stdout, stderr) => {
        if (err) {
          reject(new Error((stderr || err.message).trim()));
          return;
        }
        resolve();
      });
    });
    return await readFile(output, "utf8");
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
}
