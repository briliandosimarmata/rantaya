import { spawn } from "node:child_process";
import { readFileSync } from "node:fs";
import { root, environment } from "./env.mjs";

const env = environment();
const origin = new URL(env.APP_URL || "http://localhost:5173");
const builtOrigin = readFileSync(root + "/web/build/origin.txt", "utf8").trim();
const builtAllowedOrigins = readFileSync(
  root + "/web/build/allowed-origins.txt",
  "utf8",
).trim();
if (
  builtOrigin !== origin.origin ||
  builtAllowedOrigins !== (env.APP_ALLOWED_ORIGINS || "").trim()
) {
  throw new Error(
    "APP_URL or APP_ALLOWED_ORIGINS differs from the frontend build. Run npm run build again.",
  );
}
env.HOST ||= "127.0.0.1";
env.PORT ||= origin.port || (origin.protocol === "https:" ? "443" : "80");
env.BODY_SIZE_LIMIT ||= "10M";
const children = [];
let stopping = false;
function stop(code = 0) {
  if (stopping) return;
  stopping = true;
  for (const child of children) child.kill("SIGTERM");
  setTimeout(() => process.exit(code), 500).unref();
}
for (const [cmd, args] of [
  [
    root +
      "/backend/bin/ruang-api" +
      (process.platform === "win32" ? ".exe" : ""),
    [],
  ],
  [process.execPath, ["web/build"]],
]) {
  const p = spawn(cmd, args, { cwd: root, env, stdio: "inherit" });
  children.push(p);
  p.on("error", (e) => {
    console.error(e.message);
    stop(1);
  });
  p.on("exit", (code) => {
    if (!stopping) stop(code || 1);
  });
}
console.log("Rantaya preview: " + origin.origin);
process.on("SIGINT", () => stop());
process.on("SIGTERM", () => stop());
