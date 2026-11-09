import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { test } from "node:test";

import { renderCanvas } from "./render";

test("renders a canvas document to one html file", async () => {
  const bin = resolve(__dirname, "../bin/agent-canvas");
  const source = await readFile(resolve(__dirname, "../../../examples/release-health.canvas"), "utf8");
  const html = await renderCanvas(bin, source);
  assert.match(html, /^<!DOCTYPE html>/);
  assert.match(html, /Checkout reliability/);
  assert.match(html, /<style>/);
  assert.equal(html.split("<html").length - 1, 1);
});
