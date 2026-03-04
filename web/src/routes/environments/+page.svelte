<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api.js";
  import type { Environment } from "$lib/types.js";

  let environments = $state<Environment[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let showCreateForm = $state(false);
  let editingEnvironment = $state<Environment | null>(null);

  // Form state
  let formData = $state({
    key: "",
    name: "",
    description: "",
    color: "#3b82f6",
    sort_order: 0,
    is_active: true,
  });

  // Delete confirmation
  let deletingEnvironment = $state<Environment | null>(null);

  async function loadEnvironments() {
    try {
      loading = true;
      const response = await api.getEnvironments();
      environments = response.data.sort((a, b) => a.sort_order - b.sort_order);
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load environments";
    } finally {
      loading = false;
    }
  }

  function resetForm() {
    formData = {
      key: "",
      name: "",
      description: "",
      color: "#3b82f6",
      sort_order: environments.length,
      is_active: true,
    };
    showCreateForm = false;
    editingEnvironment = null;
  }

  function startCreate() {
    resetForm();
    showCreateForm = true;
  }

  function startEdit(env: Environment) {
    formData = {
      key: env.key,
      name: env.name,
      description: env.description,
      color: env.color,
      sort_order: env.sort_order,
      is_active: env.is_active,
    };
    editingEnvironment = env;
    showCreateForm = true;
  }

  function generateKey() {
    if (formData.name) {
      formData.key = formData.name
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, "")
        .replace(/\s+/g, "-")
        .replace(/-+/g, "-")
        .replace(/^-|-$/g, "");
    }
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    try {
      if (editingEnvironment) {
        await api.updateEnvironment(editingEnvironment.key, formData);
      } else {
        await api.createEnvironment(formData);
      }
      await loadEnvironments();
      resetForm();
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to save environment";
    }
  }

  async function handleDelete() {
    if (!deletingEnvironment) return;

    try {
      await api.deleteEnvironment(deletingEnvironment.key);
      await loadEnvironments();
      deletingEnvironment = null;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to delete environment";
    }
  }

  onMount(() => {
    loadEnvironments();
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Environments</h1>
      <p class="mt-1 text-sm text-gray-600">
        Manage deployment environments for feature flags
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
      Add Environment
    </button>
  </div>

  <!-- Error message -->
  {#if error}
    <div class="rounded-md bg-red-50 p-4">
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

  <!-- Create/Edit form -->
  {#if showCreateForm}
    <div class="bg-white shadow rounded-lg">
      <div class="px-4 py-5 sm:p-6">
        <h3 class="text-lg font-medium leading-6 text-gray-900 mb-4">
          {editingEnvironment ? "Edit Environment" : "Create Environment"}
        </h3>
        <form onsubmit={handleSubmit} class="space-y-4">
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
              <input
                type="text"
                id="name"
                bind:value={formData.name}
                oninput={generateKey}
                required
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                placeholder="Production"
              />
            </div>
            <div>
              <label for="key" class="block text-sm font-medium text-gray-700">Key</label>
              <input
                type="text"
                id="key"
                bind:value={formData.key}
                required
                readonly={!!editingEnvironment}
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm {editingEnvironment ? 'bg-gray-50' : ''}"
                placeholder="production"
              />
            </div>
          </div>

          <div>
            <label for="description" class="block text-sm font-medium text-gray-700">Description</label>
            <textarea
              id="description"
              bind:value={formData.description}
              rows={2}
              class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              placeholder="Production environment for live users"
            ></textarea>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
            <div>
              <label for="color" class="block text-sm font-medium text-gray-700">Color</label>
              <input
                type="color"
                id="color"
                bind:value={formData.color}
                class="mt-1 block w-full h-10 border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
            <div>
              <label for="sort_order" class="block text-sm font-medium text-gray-700">Sort Order</label>
              <input
                type="number"
                id="sort_order"
                bind:value={formData.sort_order}
                min="0"
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              />
            </div>
            <div class="flex items-center pt-6">
              <input
                type="checkbox"
                id="is_active"
                bind:checked={formData.is_active}
                class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
              />
              <label for="is_active" class="ml-2 block text-sm text-gray-900">Active</label>
            </div>
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
              {editingEnvironment ? "Update" : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}

  <!-- Environments list -->
  <div class="bg-white shadow overflow-hidden sm:rounded-md">
    {#if loading}
      <div class="px-4 py-8 text-center">
        <div class="inline-flex items-center text-sm text-gray-500">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Loading environments...
        </div>
      </div>
    {:else if environments.length === 0}
      <div class="text-center py-8">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M19 11H5m14-7H3a2 2 0 00-2 2v12a2 2 0 002 2h16a2 2 0 002-2V6a2 2 0 00-2-2z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No environments</h3>
        <p class="mt-1 text-sm text-gray-500">Get started by creating your first environment.</p>
      </div>
    {:else}
      <ul class="divide-y divide-gray-200">
        {#each environments as environment (environment.id)}
          <li class="px-6 py-4">
            <div class="flex items-center justify-between">
              <div class="flex items-center space-x-4">
                <div
                  class="w-4 h-4 rounded-full border border-gray-200"
                  style="background-color: {environment.color}"
                ></div>
                <div>
                  <div class="flex items-center space-x-2">
                    <h3 class="text-sm font-medium text-gray-900">{environment.name}</h3>
                    <code class="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded">{environment.key}</code>
                    {#if !environment.is_active}
                      <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-red-100 text-red-800">
                        Inactive
                      </span>
                    {/if}
                  </div>
                  {#if environment.description}
                    <p class="mt-1 text-sm text-gray-500">{environment.description}</p>
                  {/if}
                  <p class="text-xs text-gray-500 mt-1">Sort order: {environment.sort_order}</p>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <button
                  type="button"
                  class="inline-flex items-center px-3 py-1 border border-gray-300 rounded-md text-xs font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  onclick={() => startEdit(environment)}
                >
                  Edit
                </button>
                <button
                  type="button"
                  class="inline-flex items-center px-3 py-1 border border-red-300 rounded-md text-xs font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                  onclick={() => (deletingEnvironment = environment)}
                >
                  Delete
                </button>
              </div>
            </div>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>

<!-- Delete confirmation modal -->
{#if deletingEnvironment}
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
        <h3 class="text-lg font-medium text-gray-900">Delete Environment</h3>
        <div class="mt-2">
          <p class="text-sm text-gray-500">
            Are you sure you want to delete the <strong>{deletingEnvironment.name}</strong> environment?
            This action cannot be undone and will affect all flags using this environment.
          </p>
        </div>
      </div>
      <div class="mt-5 flex justify-center space-x-3">
        <button
          type="button"
          class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm bg-white text-sm font-medium text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          onclick={() => (deletingEnvironment = null)}
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