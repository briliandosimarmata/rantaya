import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
export const root = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  "..",
);
export function environment() {
  const env = { ...process.env };
  const file = path.join(root, ".env");
  if (fs.existsSync(file))
    for (const raw of fs.readFileSync(file, "utf8").split(/\r?\n/)) {
      const line = raw.trim();
      if (!line || line.startsWith("#")) continue;
      const match = /^(?:export\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*=\s*(.*)$/.exec(
        line,
      );
      if (!match) throw new Error("Invalid .env line");
      let value = match[2];
      if (
        (value.startsWith('"') && value.endsWith('"')) ||
        (value.startsWith("'") && value.endsWith("'"))
      )
        value = value.slice(1, -1);
      if (env[match[1]] === undefined) env[match[1]] = value;
    }
  if (!env.UPLOAD_DIR || !path.isAbsolute(env.UPLOAD_DIR))
    env.UPLOAD_DIR = path.resolve(root, env.UPLOAD_DIR || "var/uploads");
  return env;
}
export const npm = process.platform === "win32" ? "npm.cmd" : "npm";
