import { Plugin } from "@opencode/plugin";
import path from "node:path";

export default Plugin.define({
  id: "run-slop",
  async setup(ctx) {
    const projectRoot = ctx.location.project.directory;
    await ctx.tool.transform((editor) => {
      editor.add({
        name: "run_slop",
        description:
          "Run a SlopLang (.slop) file with the Sloperator interpreter. Returns stdout, stderr, and exitCode.",
        input: {
          type: "object",
          properties: {
            path: {
              type: "string",
              description:
                "Path to the .slop file to run, absolute or relative to the project root.",
            },
          },
          required: ["path"],
          additionalProperties: false,
        },
        output: {
          type: "object",
          properties: {
            stdout: { type: "string" },
            stderr: { type: "string" },
            exitCode: { type: "number" },
          },
          required: ["stdout", "stderr", "exitCode"],
          additionalProperties: false,
        },
        options: { codemode: false },
        execute: async (input) => {
          const { path: filePath } = input as { path: string };
          const resolved = path.isAbsolute(filePath)
            ? filePath
            : path.join(projectRoot, filePath);
          const proc = Bun.spawn(
            ["uv", "run", "--no-active", "sloperator.py", resolved],
            {
              cwd: projectRoot,
              stdout: "pipe",
              stderr: "pipe",
            },
          );
          const [stdout, stderr, exitCode] = await Promise.all([
            new Response(proc.stdout).text(),
            new Response(proc.stderr).text(),
            proc.exited,
          ]);
          const output = { stdout, stderr, exitCode };
          return {
            output,
            content: JSON.stringify(output),
          };
        },
      });
    });
  },
});
