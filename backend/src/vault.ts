import { randomUUID } from "crypto";
import { readFile, writeFile } from "fs/promises";
import { z } from "zod";

const VAULT_PATH = "./vault.json";

// Enforce that the password is base64-encoded and reasonably long
const encryptedPasswordSchema = z
  .string()
  .min(16, { message: "Password is too short to be encrypted" })
  .regex(/^[A-Za-z0-9+/=]+$/, { message: "Password must be base64 encoded" });

const entrySchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  username: z.string(),
  password: encryptedPasswordSchema,
});

export type VaultEntry = z.infer<typeof entrySchema>;

let vault: Record<string, VaultEntry> = {};

export async function loadVault() {
  try {
    const raw = await readFile(VAULT_PATH, "utf-8");
    const parsed = JSON.parse(raw);

    vault = {};
    for (const entry of parsed) {
      const valid = entrySchema.parse(entry);
      vault[valid.id] = valid;
    }
  } catch (e) {
    console.error("Failed loading vault or invalid entries:", e);
    vault = {};
  }
}

export async function saveVault() {
  const entries = Object.values(vault);
  for (const entry of entries) {
    entrySchema.parse(entry); // validate before saving
  }
  await writeFile(VAULT_PATH, JSON.stringify(entries, null, 2), "utf-8");
}

export function getAll(): VaultEntry[] {
  return Object.values(vault);
}

export async function add(entry: Omit<VaultEntry, "id">) {
  const validated = entrySchema.omit({ id: true }).parse(entry);
  const id = randomUUID();
  const newEntry = { ...validated, id };

  vault[id] = entrySchema.parse(newEntry);
  await saveVault();
  return newEntry;
}

export async function remove(id: string) {
  if (!vault[id]) {
    throw new Error("Entry not found");
  }
  delete vault[id];
  await saveVault();
}

export async function update(id: string, data: Omit<VaultEntry, "id">) {
  if (!vault[id]) {
    throw new Error("Entry not found");
  }

  const updated = entrySchema.parse({ ...data, id });
  vault[id] = updated;
  await saveVault();
  return updated;
}
