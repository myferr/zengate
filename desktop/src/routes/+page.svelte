<script lang="ts">
  import "../app.css";
  import { onMount, onDestroy } from "svelte";
  import { writable, get } from "svelte/store";
  import { fade } from "svelte/transition";

  import { Pencil, Clipboard, Trash, RotateCw } from "lucide-svelte";

  function base64ToArrayBuffer(base64: string): Uint8Array {
    const binary_string = atob(base64);
    const len = binary_string.length;
    const bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
      bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes;
  }
  function arrayBufferToBase64(buffer: ArrayBuffer) {
    let binary = "";
    const bytes = new Uint8Array(buffer);
    for (let b of bytes) binary += String.fromCharCode(b);
    return btoa(binary);
  }

  let cryptoKey: CryptoKey;

  const apiUrl = writable(localStorage.getItem("apiUrl") || "");
  const encryptionKeyBase64 = writable(
    localStorage.getItem("encryptionKeyBase64") || ""
  );

  type VaultEntry = {
    id: string;
    name: string;
    username: string;
    password: string;
    decryptedPassword?: string | undefined;
  };

  const vaults = writable<VaultEntry[]>([]);

  let editingUsernameId: string | null = null;
  let editingPasswordId: string | null = null;
  let editingUsernameValue = "";
  let editingPasswordValue = "";

  let newName = "";
  let newUsername = "";
  let newPassword = "";

  let hasSettings = false;
  let loadingVault = false;
  let errorMsg = "";

  async function importKey() {
    try {
      const rawKey = base64ToArrayBuffer(get(encryptionKeyBase64));
      cryptoKey = await crypto.subtle.importKey(
        "raw",
        rawKey,
        "AES-GCM",
        false,
        ["encrypt", "decrypt"]
      );
    } catch {
      errorMsg = "Invalid Encryption Key - must be base64 32 bytes";
      throw new Error(errorMsg);
    }
  }

  async function encryptPassword(plaintext: string): Promise<string> {
    const encoder = new TextEncoder();
    const data = encoder.encode(plaintext);
    const iv = crypto.getRandomValues(new Uint8Array(12));
    const encrypted = await crypto.subtle.encrypt(
      { name: "AES-GCM", iv },
      cryptoKey,
      data
    );
    const combined = new Uint8Array(iv.length + encrypted.byteLength);
    combined.set(iv, 0);
    combined.set(new Uint8Array(encrypted), iv.length);
    return arrayBufferToBase64(combined.buffer);
  }

  async function decryptPassword(ciphertextBase64: string): Promise<string> {
    try {
      const combined = base64ToArrayBuffer(ciphertextBase64);
      const iv = combined.subarray(0, 12);
      const data = combined.subarray(12);
      const decrypted = await crypto.subtle.decrypt(
        { name: "AES-GCM", iv },
        cryptoKey,
        data
      );
      const decoder = new TextDecoder();
      return decoder.decode(decrypted);
    } catch {
      return "[Decryption failed]";
    }
  }

  async function loadVaults() {
    errorMsg = "";
    loadingVault = true;
    try {
      const res = await fetch(`${get(apiUrl)}/vaults`);
      if (!res.ok) {
        throw new Error(`Failed to fetch vault: ${res.status}`);
      }
      const data: VaultEntry[] = await res.json();

      const decrypted = await Promise.all(
        data.map(async (entry) => {
          const decryptedPassword = await decryptPassword(entry.password);
          return { ...entry, decryptedPassword };
        })
      );
      vaults.set(decrypted);
      hasSettings = true;
    } catch (e) {
      if (e instanceof Error) {
        errorMsg = e.message || "Failed to load vault";
      } else {
        errorMsg = "Failed to load vault";
      }
    }
    loadingVault = false;
  }

  async function saveSettings() {
    errorMsg = "";
    try {
      await importKey();
      localStorage.setItem("apiUrl", get(apiUrl));
      localStorage.setItem("encryptionKeyBase64", get(encryptionKeyBase64));
      await loadVaults();
    } catch (e) {
      errorMsg = "Invalid encryption key or API URL";
    }
  }

  async function addEntry() {
    if (!newName || !newUsername || !newPassword) {
      errorMsg = "Please fill all fields";
      return;
    }
    errorMsg = "";
    const encrypted = await encryptPassword(newPassword);

    const res = await fetch(`${get(apiUrl)}/vaults`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: newName,
        username: newUsername,
        password: encrypted,
      }),
    });
    if (!res.ok) {
      errorMsg = "Failed to add entry";
      return;
    }
    const added = await res.json();
    const decryptedPassword = await decryptPassword(added.password);
    vaults.update((vs) => [...vs, { ...added, decryptedPassword }]);
    newName = "";
    newUsername = "";
    newPassword = "";
  }

  async function deleteEntry(id: string) {
    console.log("deleteEntry called for id:", id);

    try {
      const res = await fetch(`${get(apiUrl)}/vaults/${id}`, {
        method: "DELETE",
      });
      console.log("DELETE response status:", res.status);
      const text = await res.text();
      console.log("DELETE response body:", text);

      if (!res.ok) {
        errorMsg = "Failed to delete entry";
        console.error("Delete failed:", text);
        return;
      }
      vaults.update((vs) => vs.filter((v) => v.id !== id));
    } catch (err) {
      console.error("Delete request error:", err);
      errorMsg = "Failed to delete entry due to network error";
    }
  }

  function startEditUsername(id: string, currentUsername: string) {
    editingUsernameId = id;
    editingUsernameValue = currentUsername;
  }

  async function saveUsername(id: string, newUsername: string) {
    const all = get(vaults);
    const entry = all.find((v) => v.id === id);
    if (!entry) return;

    const res = await fetch(`${get(apiUrl)}/vaults/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: entry.name,
        username: newUsername,
        password: entry.password,
      }),
    });

    if (!res.ok) {
      errorMsg = "Failed to update username";
      return;
    }
    const updated = await res.json();
    // Update local vaults immediately
    vaults.update((vs) =>
      vs.map((v) =>
        v.id === id
          ? { ...updated, decryptedPassword: entry.decryptedPassword }
          : v
      )
    );
    editingUsernameId = null;
    editingUsernameValue = "";
  }

  function startEditPassword(id: string, currentPassword: string) {
    editingPasswordId = id;
    editingPasswordValue = currentPassword;
  }

  async function savePassword(id: string, newPasswordPlaintext: string) {
    const all = get(vaults);
    const entry = all.find((v) => v.id === id);
    if (!entry) return;

    const encrypted = await encryptPassword(newPasswordPlaintext);

    const res = await fetch(`${get(apiUrl)}/vaults/${id}`, {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        name: entry.name,
        username: entry.username,
        password: encrypted,
      }),
    });

    if (!res.ok) {
      errorMsg = "Failed to update password";
      return;
    }
    const updated = await res.json();
    const decryptedPassword = await decryptPassword(updated.password);
    // Update local vaults immediately
    vaults.update((vs) =>
      vs.map((v) => (v.id === id ? { ...updated, decryptedPassword } : v))
    );
    editingPasswordId = null;
    editingPasswordValue = "";
  }

  async function copyText(text: string) {
    try {
      await navigator.clipboard.writeText(text);
      alert("Copied!");
    } catch {
      alert("Copy failed");
    }
  }

  // Poll vault every 30 seconds for real-time updates
  let intervalId: ReturnType<typeof setInterval>;
  onMount(() => {
    if (get(apiUrl) && get(encryptionKeyBase64)) {
      saveSettings();
    }
    intervalId = setInterval(() => {
      loadVaults();
    }, 30000);
  });

  onDestroy(() => {
    clearInterval(intervalId);
  });
</script>

<div class="min-h-screen bg-gray-50 p-6 max-w-4xl mx-auto">
  <h1 class="text-4xl font-bold mb-6 text-center text-gray-900">
    Zengate Vault
  </h1>

  {#if !hasSettings}
    <section
      class="max-w-lg mx-auto bg-white p-6 rounded shadow-md space-y-4"
      in:fade
      out:fade
    >
      <h2 class="text-xl font-semibold mb-2">Setup</h2>
      {#if errorMsg}
        <p class="text-red-600 font-semibold" in:fade out:fade>{errorMsg}</p>
      {/if}
      <label class="block">
        <span class="text-gray-700 font-medium">API Tunnel URL</span>
        <input
          type="url"
          placeholder="https://yourtunnel.example.com"
          class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-orange-500 focus:ring-orange-500"
          bind:value={$apiUrl}
        />
      </label>

      <label class="block">
        <span class="text-gray-700 font-medium"
          >Encryption Key (base64 32 bytes)</span
        >
        <input
          type="text"
          placeholder="Base64 key here"
          class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 focus:border-orange-500 focus:ring-orange-500"
          bind:value={$encryptionKeyBase64}
        />
      </label>

      <button
        on:click={saveSettings}
        class="w-full mt-4 bg-orange-600 text-white py-2 rounded hover:bg-orange-700 transition p-2 hover:bg-muted/50 hover:cursor-pointer"
      >
        Save & Load Vault
      </button>
    </section>
  {:else}
    {#if $apiUrl && $encryptionKeyBase64}
      <section class="mb-6" in:fade out:fade>
        <button
          on:click={loadVaults}
          class="flex items-center gap-2 bg-orange-500 text-white px-4 py-2 rounded hover:bg-orange-600 transition p-2 hover:bg-muted/50 hover:cursor-pointer"
        >
          <RotateCw class="w-4 h-4" />
          Reload Vault
        </button>
      </section>
    {/if}

    <section
      class="mb-6 max-w-3xl mx-auto bg-white p-6 rounded shadow-md"
      in:fade
      out:fade
    >
      {#if errorMsg}
        <p class="text-red-600 font-semibold mb-2" in:fade out:fade>
          {errorMsg}
        </p>
      {/if}

      <h2 class="text-xl font-semibold mb-4">Add New Entry</h2>
      <form on:submit|preventDefault={addEntry} class="space-y-4">
        <input
          type="text"
          placeholder="Name (e.g. github.com)"
          class="w-full rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring focus:ring-orange-400"
          bind:value={newName}
          required
        />
        <input
          type="text"
          placeholder="Username"
          class="w-full rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring focus:ring-orange-400"
          bind:value={newUsername}
          required
        />
        <input
          type="password"
          placeholder="Password"
          class="w-full rounded border border-gray-300 px-3 py-2 focus:outline-none focus:ring focus:ring-orange-400"
          bind:value={newPassword}
          required
        />
        <button
          type="submit"
          class="w-full bg-orange-600 text-white py-2 rounded hover:bg-orange-700 transition p-2 hover:bg-muted/50 hover:cursor-pointer"
        >
          Add Entry
        </button>
      </form>
    </section>

    <section
      class="max-w-4xl mx-auto bg-white p-6 rounded shadow-md"
      in:fade
      out:fade
    >
      <h2 class="text-xl font-semibold mb-4">Vault Entries</h2>

      {#if loadingVault}
        <p>Loading vault...</p>
      {:else if $vaults.length === 0}
        <p class="italic text-gray-500">Vault is empty.</p>
      {:else}
        <ul class="space-y-4">
          {#each $vaults as entry (entry.id)}
            <li
              class="border border-gray-300 rounded p-4 flex flex-col md:flex-row md:items-center md:justify-between gap-4"
              in:fade
              out:fade
            >
              <div
                class="flex flex-col md:flex-row md:items-center gap-4 flex-grow"
              >
                <div class="font-semibold">{entry.name}</div>

                {#if editingUsernameId === entry.id}
                  <input
                    type="text"
                    class="border border-gray-300 rounded px-2 py-1 w-40"
                    bind:value={editingUsernameValue}
                    on:keydown={(e) => {
                      if (e.key === "Enter")
                        saveUsername(entry.id, editingUsernameValue);
                      if (e.key === "Escape") {
                        editingUsernameId = null;
                        editingUsernameValue = "";
                      }
                    }}
                  />
                  <button
                    on:click={() =>
                      saveUsername(entry.id, editingUsernameValue)}
                    class="ml-2 px-3 py-1 bg-orange-600 text-white rounded hover:bg-orange-700 transition p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                  >
                    Save
                  </button>
                  <button
                    on:click={() => {
                      editingUsernameId = null;
                      editingUsernameValue = "";
                    }}
                    class="ml-2 px-3 py-1 bg-gray-400 text-white rounded hover:bg-gray-500 transition p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                  >
                    Cancel
                  </button>
                {:else}
                  <div class="flex items-center gap-2">
                    <span>Username: {entry.username}</span>
                    <button
                      on:click={() =>
                        startEditUsername(entry.id, entry.username)}
                      class="ml-2 p-1 text-orange-600 hover:text-orange-800 p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                      title="Edit Username"
                    >
                      <Pencil class="w-4 h-4" />
                    </button>
                    <button
                      on:click={() => copyText(entry.username)}
                      class="ml-1 p-1 text-green-600 hover:text-green-800 p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                      title="Copy Username"
                    >
                      <Clipboard class="w-4 h-4" />
                    </button>
                  </div>
                {/if}
              </div>

              <div
                class="flex flex-col md:flex-row md:items-center gap-4 flex-grow"
              >
                {#if editingPasswordId === entry.id}
                  <input
                    type="text"
                    class="border border-gray-300 rounded px-2 py-1 w-40"
                    bind:value={editingPasswordValue}
                    on:keydown={(e) => {
                      if (e.key === "Enter")
                        savePassword(entry.id, editingPasswordValue);
                      if (e.key === "Escape") {
                        editingPasswordId = null;
                        editingPasswordValue = "";
                      }
                    }}
                  />
                  <button
                    on:click={() =>
                      savePassword(entry.id, editingPasswordValue)}
                    class="ml-2 px-3 py-1 bg-orange-600 text-white rounded hover:bg-orange-700 transition p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                  >
                    Save
                  </button>
                  <button
                    on:click={() => {
                      editingPasswordId = null;
                      editingPasswordValue = "";
                    }}
                    class="ml-2 px-3 py-1 bg-gray-400 text-white rounded hover:bg-gray-500 transition p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                  >
                    Cancel
                  </button>
                {:else}
                  <div class="flex items-center gap-2">
                    <span>Password: ••••••••</span>
                    <button
                      on:click={() =>
                        startEditPassword(
                          entry.id,
                          entry.decryptedPassword || ""
                        )}
                      class="ml-2 p-1 text-orange-600 hover:text-orange-800 p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                      title="Edit Password"
                    >
                      <Pencil class="w-4 h-4" />
                    </button>
                    <button
                      on:click={() => copyText(entry.decryptedPassword || "")}
                      class="ml-1 p-1 text-green-600 hover:text-green-800 p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center"
                      title="Copy Password"
                    >
                      <Clipboard class="w-4 h-4" />
                    </button>
                  </div>
                {/if}

                <button
                  on:click={() => deleteEntry(entry.id)}
                  class="ml-4 bg-red-600 text-white rounded hover:bg-red-700 transition p-2 hover:bg-muted/50 hover:cursor-pointer flex items-center justify-center gap-1"
                  title="Delete Entry"
                >
                  <Trash class="w-4 h-4" />
                  <span>Remove</span>
                </button>
              </div>
            </li>
          {/each}
        </ul>
      {/if}
    </section>
  {/if}
</div>
