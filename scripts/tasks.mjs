import { spawnSync } from "node:child_process";
import { writeFileSync } from "node:fs";
import { root, environment, npm } from "./env.mjs";
const env = environment();
function run(command, args, folder) {
  const p = spawnSync(command, args, {
    cwd: root + "/" + folder,
    env,
    stdio: "inherit",
    shell: process.platform === "win32" && command === npm,
  });
  if (p.error) {
    console.error(p.error.message);
    process.exit(1);
  }
  if (p.status !== 0) process.exit(p.status || 1);
}
const task = process.argv[2];
if (task === "setup") {
  run(npm, ["ci"], "web");
  run("go", ["mod", "download"], "backend");
} else if (task === "check") {
  run(npm, ["run", "check"], "web");
  run("go", ["vet", "./..."], "backend");
  run("go", ["test", "./..."], "backend");
} else if (task === "build") {
  run(npm, ["run", "build"], "web");
  writeFileSync(
    root + "/web/build/origin.txt",
    env.APP_URL || "http://localhost:5173",
  );
  writeFileSync(
    root + "/web/build/allowed-origins.txt",
    env.APP_ALLOWED_ORIGINS || "",
  );
  run(
    "go",
    [
      "build",
      "-buildvcs=false",
      "-o",
      "bin/ruang-api" + (process.platform === "win32" ? ".exe" : ""),
      "./cmd/api",
    ],
    "backend",
  );
} else if (task === "test") {
  run("go", ["test", "./...", "-count=1", "-v"], "backend");
} else if (task === "test:e2e") {
  run(npm, ["run", "test:e2e", "--", ...process.argv.slice(3)], "web");
} else throw new Error("Unknown task");
