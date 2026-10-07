import { spawn, spawnSync } from "node:child_process";
import { root, environment, npm } from "./env.mjs";
const env = environment();
const children = [];
let stopping = false;
function stop(code = 0) {
  if (stopping) return;
  stopping = true;
  for (const child of children) child.kill("SIGTERM");
  setTimeout(() => process.exit(code), 500).unref();
}
function run(cmd, args, cwd) {
  const p = spawn(cmd, args, {
    cwd,
    env,
    stdio: "inherit",
    shell: process.platform === "win32" && cmd === npm,
  });
  children.push(p);
  p.on("error", (e) => {
    console.error(e.message);
    stop(1);
  });
  p.on("exit", (code) => {
    if (!stopping) {
      console.error("Service stopped. Check PostgreSQL and .env.");
      stop(code || 1);
    }
  });
}
console.log(
  "Rantaya local: " + (env.APP_URL || "http://localhost:5173") + " — Go API: " +
    (env.API_ADDR || "127.0.0.1:8080"),
);
const binary =
  root +
  "/backend/bin/ruang-api" +
  (process.platform === "win32" ? ".exe" : "");
const build = spawnSync(
  "go",
  ["build", "-buildvcs=false", "-o", binary, "./cmd/api"],
  { cwd: root + "/backend", env, stdio: "inherit" },
);
if (build.error || build.status !== 0) {
  console.error(build.error?.message || "Go build failed");
  process.exit(1);
}
run(binary, [], root);
run(npm, ["run", "dev"], root + "/web");
process.on("SIGINT", () => stop());
process.on("SIGTERM", () => stop());
