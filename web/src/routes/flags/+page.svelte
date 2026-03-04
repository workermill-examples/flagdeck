<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "../../lib/api.js";
  import FlagToggle from "../../lib/components/FlagToggle.svelte";
  import type { Flag, Environment } from "../../lib/types.js";

  let flags = $state<Flag[]>([]);
  let environments = $state<Environment[]>([]);
  let isLoading = $state(true);
  let error = $state<string | null>(null);
  let searchQuery = $state("");
  let selectedType = $state("");
  let selectedTag = $state("");

  // Derived states for filtering
  const availableTags = $derived.by(() => {
    const allTags = flags.flatMap((flag) => flag.tags);
    return [...new Set(allTags)].sort();
  });

  const filteredFlags = $derived.by(() => {
    return flags.filter((flag) => {
      const matchesSearch =
        searchQuery === "" ||
        flag.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        flag.key.toLowerCase().includes(searchQuery.toLowerCase()) ||
        flag.description.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesType = selectedType === "" || flag.type === selectedType;

      const matchesTag = selectedTag === "" || flag.tags.includes(selectedTag);

      return matchesSearch && matchesType && matchesTag;
    });
  });

  onMount(async () => {
    await loadData();
  });

  async function loadData() {
    try {
      isLoading = true;
      error = null;

      const [flagsResponse, environmentsResponse] = await Promise.all([
        api.getFlags(),
        api.getEnvironments(),
      ]);

      flags = flagsResponse.data;
      environments = environmentsResponse.data.sort(
        (a, b) => a.sort_order - b.sort_order,
      );
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load data";
      console.error("Failed to load flags:", err);
    } finally {
      isLoading = false;
    }
  }

  async function handleToggle() {
    // Refresh flags after toggle
    await loadData();
  }

  function getActiveEnvironmentsCount(flag: Flag): number {
    return Object.values(flag.environments).filter((env) => env.enabled).length;
  }

  function navigateToFlag(flagKey: string) {
    window.location.href = `/flags/${flagKey}`;
  }
</script>

<svelte:head>
  <title>Feature Flags - FlagDeck</title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Feature Flags</h1>
      <p class="text-gray-600 mt-1">
        Manage your application's feature flags and their deployment
        configurations
      </p>
    </div>
    <a
      href="/flags/create"
      class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
    >
      <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 20 20">
        <path
          fill-rule="evenodd"
          d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
          clip-rule="evenodd"
        ></path>
      </svg>
      Create Flag
    </a>
  </div>

  <!-- Search and Filters -->
  <div class="bg-white shadow rounded-lg p-4">
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div>
        <label for="search" class="block text-sm font-medium text-gray-700 mb-1"
          >Search</label
        >
        <input
          id="search"
          type="text"
          placeholder="Search flags..."
          bind:value={searchQuery}
          class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
        />
      </div>

      <div>
        <label
          for="type-filter"
          class="block text-sm font-medium text-gray-700 mb-1">Type</label
        >
        <select
          id="type-filter"
          bind:value={selectedType}
          class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="">All Types</option>
          <option value="boolean">Boolean</option>
          <option value="string">String</option>
          <option value="number">Number</option>
          <option value="json">JSON</option>
        </select>
      </div>

      <div>
        <label
          for="tag-filter"
          class="block text-sm font-medium text-gray-700 mb-1">Tag</label
        >
        <select
          id="tag-filter"
          bind:value={selectedTag}
          class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
        >
          <option value="">All Tags</option>
          {#each availableTags as tag}
            <option value={tag}>{tag}</option>
          {/each}
        </select>
      </div>

      <div class="flex items-end">
        <button
          onclick={() => {
            searchQuery = "";
            selectedType = "";
            selectedTag = "";
          }}
          class="w-full px-3 py-2 text-sm text-gray-600 bg-gray-100 border border-gray-300 rounded-md hover:bg-gray-200 focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
        >
          Clear Filters
        </button>
      </div>
    </div>
  </div>

  <!-- Loading State -->
  {#if isLoading}
    <div class="flex items-center justify-center py-12">
      <div
        class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"
      ></div>
    </div>
  {:else if error}
    <!-- Error State -->
    <div class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <svg
          class="flex-shrink-0 h-5 w-5 text-red-400"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path
            fill-rule="evenodd"
            d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
            clip-rule="evenodd"
          ></path>
        </svg>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error</h3>
          <p class="text-sm text-red-700 mt-1">{error}</p>
        </div>
      </div>
      <button
        onclick={loadData}
        class="mt-3 text-sm text-red-600 hover:text-red-800 underline"
      >
        Try again
      </button>
    </div>
  {:else if filteredFlags.length === 0}
    <!-- Empty State -->
    <div class="text-center py-12">
      <svg
        class="mx-auto h-12 w-12 text-gray-400"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M3 21v-4m0 0V5a2 2 0 012-2h6.586a1 1 0 01.707.293l2.414 2.414a1 1 0 01.293.707V13a2 2 0 01-2 2H5a2 2 0 01-2 2z"
        />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No flags found</h3>
      <p class="mt-1 text-sm text-gray-500">
        {#if searchQuery || selectedType || selectedTag}
          Try adjusting your filters to see more results.
        {:else}
          Get started by creating your first feature flag.
        {/if}
      </p>
    </div>
  {:else}
    <!-- Flags Table -->
    <div class="bg-white shadow rounded-lg overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Flag
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Type
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Status
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Environment Toggles
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Tags
              </th>
              <th
                class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            {#each filteredFlags as flag (flag.key)}
              <tr class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div>
                    <div class="text-sm font-medium text-gray-900">
                      {flag.name}
                    </div>
                    <div class="text-sm text-gray-500">
                      {flag.key}
                    </div>
                    {#if flag.description}
                      <div class="text-xs text-gray-400 mt-1 max-w-xs truncate">
                        {flag.description}
                      </div>
                    {/if}
                  </div>
                </td>

                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
										{flag.type === 'boolean'
                      ? 'bg-blue-100 text-blue-800'
                      : flag.type === 'string'
                        ? 'bg-green-100 text-green-800'
                        : flag.type === 'number'
                          ? 'bg-purple-100 text-purple-800'
                          : 'bg-yellow-100 text-yellow-800'}"
                  >
                    {flag.type}
                  </span>
                </td>

                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center space-x-2">
                    <span
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
											{flag.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}"
                    >
                      {flag.is_active ? "Active" : "Inactive"}
                    </span>
                    <span class="text-xs text-gray-500">
                      {getActiveEnvironmentsCount(flag)}/{environments.length} envs
                    </span>
                  </div>
                </td>

                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center space-x-3">
                    {#each environments as env}
                      <div class="flex items-center space-x-1">
                        <div
                          class="w-3 h-3 rounded-full"
                          style="background-color: {env.color}"
                        ></div>
                        <span class="text-xs text-gray-600">{env.key}</span>
                        <FlagToggle
                          {flag}
                          environment={env.key}
                          onToggle={handleToggle}
                        />
                      </div>
                    {/each}
                  </div>
                </td>

                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex flex-wrap gap-1">
                    {#each flag.tags as tag}
                      <span
                        class="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-gray-100 text-gray-800"
                      >
                        {tag}
                      </span>
                    {/each}
                  </div>
                </td>

                <td
                  class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium"
                >
                  <button
                    onclick={() => navigateToFlag(flag.key)}
                    class="text-blue-600 hover:text-blue-900"
                  >
                    Edit
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    <!-- Summary -->
    <div class="text-sm text-gray-500 text-center">
      Showing {filteredFlags.length} of {flags.length} flags
    </div>
  {/if}
</div>
