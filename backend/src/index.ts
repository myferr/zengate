import { Elysia } from "elysia";
import { z } from "zod";
import { add, getAll, remove, loadVault, update } from "./vault";
import cors from "@elysiajs/cors";

const port = process.env.PORT || 3002;

await loadVault();

const VaultEntrySchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  username: z.string(),
  password: z.string(), // encrypted ciphertext, opaque to backend
});

const AddEntrySchema = VaultEntrySchema.omit({ id: true });

const allowedOriginsEnv = process.env.FRONTEND_ORIGINS || "";
const allowedOrigins = allowedOriginsEnv
  .split(",")
  .map((origin) => origin.trim())
  .filter(Boolean);

const app = new Elysia()
  .use(
    cors({
      origin: "*", // 🔓 Allow all origins — required for Cloudflare tunnels + Tauri
      methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
      credentials: false,
    })
  )

  .get("/vaults", () => {
    const entries = getAll();
    return VaultEntrySchema.array().parse(entries);
  })
  .post("/vaults", async (context: { body: unknown }) => {
    const parsed = AddEntrySchema.parse(context.body);
    const newEntry = await add(parsed);
    return VaultEntrySchema.parse(newEntry);
  })
  .delete("/vaults/:id", async ({ params, set }) => {
    try {
      await remove(params.id);
      return { ok: true };
    } catch (e) {
      set.status = 404;
      return { error: e instanceof Error ? e.message : "Entry not found" };
    }
  })
  .put("/vaults/:id", async ({ params, body, set }) => {
    try {
      const parsed = AddEntrySchema.parse(body); // validate incoming
      const updated = await update(params.id, parsed);
      return VaultEntrySchema.parse(updated);
    } catch (e) {
      set.status = 400;
      return {
        error: e instanceof Error ? e.message : "Failed to update entry",
      };
    }
  });

app.listen(port);
console.log(`🚀 Zengate backend running at http://localhost:${port}`);
