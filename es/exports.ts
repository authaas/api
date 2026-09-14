import { walk } from "@std/fs/walk"

const configPath = "deno.json"
const config = JSON.parse(await Deno.readTextFile(configPath))

const exports: Record<string, string> = {}

for await (const entry of walk(".", {
  exts: [".ts"],
  skip: [/^exports\.ts$/, /_test\.ts$/],
})) {
  exports[`./${entry.path.slice(0, -".ts".length)}`] = `./${entry.path}`
}

const sorted = Object.fromEntries(
  Object.entries(exports).sort(([a], [b]) => a.localeCompare(b)),
)

await Deno.writeTextFile(
  configPath,
  `${JSON.stringify({ ...config, exports: sorted }, null, 2)}\n`,
)
