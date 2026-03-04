<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api.js";
  import type { ApiKey } from "$lib/types.js";

  let apiKeys = $state<ApiKey[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let showCreateForm = $state(false);
  let createdKey = $state<string | null>(null);

  // Form state
  let formData = $state({
    name: "",
    environment: "production",
    permissions: ["evaluate"],
  });

  // Delete confirmation
  let deletingApiKey = $state<ApiKey | null>(null);

  const environmentOptions = [
    { value: "production", label: "Production" },
    { value: "staging", label: "Staging" },
    { value: "development", label: "Development" },
  ];

  const permissionOptions = [
    { value: "evaluate", label: "Evaluate flags" },
    { value: "read", label: "Read configuration" },
  ];

  async function loadApiKeys() {
    try {
      loading = true;
      const response = await api.getApiKeys();
      apiKeys = response.data.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load API keys";
    } finally {
      loading = false;
    }
  }

  function resetForm() {
    formData = {
      name: "",
      environment: "production",
      permissions: ["evaluate"],
    };
    showCreateForm = false;
    createdKey = null;
  }

  function startCreate() {
    resetForm();
    showCreateForm = true;
  }

  function togglePermission(permission: string) {
    if (formData.permissions.includes(permission)) {
      formData.permissions = formData.permissions.filter(p => p !== permission);
    } else {
      formData.permissions = [...formData.permissions, permission];
    }
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    if (formData.permissions.length === 0) {
      error = "Please select at least one permission";
      return;
    }

    try {
      const response = await api.createApiKey(formData);
      // The response should include the raw_key for one-time display
      if (response && typeof response === 'object' && 'raw_key' in response) {
        createdKey = response.raw_key as string;
      }
      await loadApiKeys();
      formData.name = "";
      formData.permissions = ["evaluate"];
      error = null;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to create API key";
    }
  }

  async function handleDelete() {
    if (!deletingApiKey) return;

    try {
      await api.deleteApiKey(deletingApiKey.id);
      await loadApiKeys();
      deletingApiKey = null;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to delete API key";
    }
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  function formatLastUsed(dateString: string | null): string {
    if (!dateString) return "Never";
    return new Date(dateString).toLocaleDateString();
  }

  function copyToClipboard(text: string) {
    navigator.clipboard.writeText(text).then(() => {
      // Could add a toast notification here
    }).catch(() => {
      // Fallback for older browsers
      const textArea = document.createElement('textarea');
      textArea.value = text;
      document.body.appendChild(textArea);
      textArea.select();
      document.execCommand('copy');
      document.body.removeChild(textArea);
    });
  }

  onMount(() => {
    loadApiKeys();
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Settings</h1>
      <p class="mt-1 text-sm text-gray-600">
        Manage API keys and application configuration
      </p>
    </div>
  </div>

  <!-- API Keys Section -->
  <div class="bg-white shadow rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <div class="flex items-center justify-between mb-4">
        <div>
          <h3 class="text-lg font-medium leading-6 text-gray-900">API Keys</h3>
          <p class="mt-1 text-sm text-gray-500">
            Create and manage API keys for programmatic access to FlagDeck
          </p>
        </div>
        <button
          type="button"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          onclick={() => startCreate()}
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M12 4v16m8-8H4" />
          </svg>
          Create API Key
        </button>
      </div>

      <!-- Error message -->
      {#if error}
        <div class="rounded-md bg-red-50 p-4 mb-4">
          <div class="flex">
            <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
            <div class="ml-3">
              <p class="text-sm text-red-800">{error}</p>
            </div>
            <button
              type="button"
              class="ml-auto pl-3"
              onclick={() => (error = null)}
            >
              <span class="sr-only">Dismiss</span>
              <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      {/if}

      <!-- Created key display -->
      {#if createdKey}
        <div class="rounded-md bg-green-50 p-4 mb-4">
          <div class="flex">
            <svg class="h-5 w-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <div class="ml-3">
              <p class="text-sm font-medium text-green-800">API key created successfully!</p>
              <div class="mt-2 text-sm text-green-700">
                <p class="mb-2">
                  <strong>Important:</strong> This is the only time you'll see the full API key. Copy it now and store it securely.
                </p>
                <div class="flex items-center space-x-2">
                  <code class="bg-white px-3 py-2 rounded border font-mono text-sm text-gray-900 flex-1">
                    {createdKey}
                  </code>
                  <button
                    type="button"
                    onclick={() => copyToClipboard(createdKey)}
                    class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-green-700 bg-green-100 hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                  >
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                    Copy
                  </button>
                </div>
              </div>
            </div>
            <button
              type="button"
              class="ml-auto pl-3"
              onclick={() => (createdKey = null)}
            >
              <span class="sr-only">Dismiss</span>
              <svg class="h-5 w-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      {/if}

      <!-- Create form -->
      {#if showCreateForm}
        <div class="border border-gray-200 rounded-lg p-4 mb-6">
          <h4 class="text-sm font-medium text-gray-900 mb-4">Create New API Key</h4>
          <form onsubmit={handleSubmit} class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
              <input
                type="text"
                id="name"
                bind:value={formData.name}
                required
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                placeholder="Production Backend"
              />
              <p class="mt-1 text-xs text-gray-500">A descriptive name to help you identify this key</p>
            </div>

            <div>
              <label for="environment" class="block text-sm font-medium text-gray-700">Environment</label>
              <select
                id="environment"
                bind:value={formData.environment}
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                {#each environmentOptions as option}
                  <option value={option.value}>{option.label}</option>
                {/each}
              </select>
              <p class="mt-1 text-xs text-gray-500">The environment this key will have access to</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Permissions</label>
              <div class="space-y-2">
                {#each permissionOptions as permission}
                  <label class="flex items-center">
                    <input
                      type="checkbox"
                      checked={formData.permissions.includes(permission.value)}
                      onchange={() => togglePermission(permission.value)}
                      class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <span class="ml-2 text-sm text-gray-700">{permission.label}</span>
                  </label>
                {/each}
              </div>
              <p class="mt-1 text-xs text-gray-500">Select the permissions this API key should have</p>
            </div>

            <div class="flex justify-end space-x-3">
              <button
                type="button"
                class="px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                onclick={resetForm}
              >
                Cancel
              </button>
              <button
                type="submit"
                class="inline-flex justify-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Create API Key
              </button>
            </div>
          </form>
        </div>
      {/if}

      <!-- API Keys list -->
      {#if loading}
        <div class="py-8 text-center">
          <div class="inline-flex items-center text-sm text-gray-500">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Loading API keys...
          </div>
        </div>
      {:else if apiKeys.length === 0}
        <div class="text-center py-8">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">No API keys</h3>
          <p class="mt-1 text-sm text-gray-500">Get started by creating your first API key.</p>
        </div>
      {:else}
        <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
          <table class="min-w-full divide-y divide-gray-300">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
                  Name
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
                  Key Prefix
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
                  Environment
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
                  Permissions
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
                  Last Used
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wide">
                  Created
                </th>
                <th scope="col" class="relative px-6 py-3">
                  <span class="sr-only">Actions</span>
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              {#each apiKeys as apiKey (apiKey.id)}
                <tr>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{apiKey.name}</div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <code class="text-sm bg-gray-100 text-gray-800 px-2 py-1 rounded">{apiKey.key_prefix}...</code>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      {apiKey.environment}
                    </span>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex flex-wrap gap-1">
                      {#each apiKey.permissions as permission}
                        <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                          {permission}
                        </span>
                      {/each}
                    </div>
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {formatLastUsed(apiKey.last_used_at)}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {formatDate(apiKey.created_at)}
                  </td>
                  <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <button
                      type="button"
                      class="text-red-600 hover:text-red-900"
                      onclick={() => (deletingApiKey = apiKey)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      {/if}
    </div>
  </div>
</div>

<!-- Delete confirmation modal -->
{#if deletingApiKey}
  <div class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-lg max-w-md w-full p-6">
      <div class="flex items-center">
        <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
          <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
        </div>
      </div>
      <div class="mt-3 text-center">
        <h3 class="text-lg font-medium text-gray-900">Delete API Key</h3>
        <div class="mt-2">
          <p class="text-sm text-gray-500">
            Are you sure you want to delete the API key <strong>{deletingApiKey.name}</strong>?
            This action cannot be undone and will immediately revoke access for any applications using this key.
          </p>
        </div>
      </div>
      <div class="mt-5 flex justify-center space-x-3">
        <button
          type="button"
          class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm bg-white text-sm font-medium text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          onclick={() => (deletingApiKey = null)}
        >
          Cancel
        </button>
        <button
          type="button"
          class="inline-flex justify-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
          onclick={handleDelete}
        >
          Delete
        </button>
      </div>
    </div>
  </div>
{/if}