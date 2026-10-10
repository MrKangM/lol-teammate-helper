// Starts the Vite dev server for `wails dev`.
//
// Wails launches the frontend watcher without an interactive stdin. Vite exits
// silently (code 0) when stdin is closed, so Wails would wait forever for the
// dev server. Keeping a pipe open on the child's stdin avoids that.
import { spawn } from "node:child_process"
import { fileURLToPath } from "node:url"

const viteBin = fileURLToPath(new URL("../node_modules/vite/bin/vite.js", import.meta.url))
const args = ["--host", "localhost", "--port", "5173", "--strictPort", ...process.argv.slice(2)]

const child = spawn(process.execPath, [viteBin, ...args], { stdio: ["pipe", "inherit", "inherit"] })

child.on("exit", (code, signal) => process.exit(code ?? (signal ? 1 : 0)))
for (const sig of ["SIGINT", "SIGTERM", "SIGBREAK"]) {
  process.on(sig, () => child.kill())
}
