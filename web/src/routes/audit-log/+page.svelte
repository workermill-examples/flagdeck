<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api.js";
  import type { AuditLogEntry } from "$lib/types.js";
  import AuditTimeline from "$lib/components/AuditTimeline.svelte";

  let entries = $state<AuditLogEntry[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  // Pagination
  let currentPage = $state(1);
  let totalEntries = $state(0);
  let entriesPerPage = $state(25);

  // Filters
  let resourceFilter = $state("");
  let actionFilter = $state("");

  const totalPages = $derived(Math.ceil(totalEntries / entriesPerPage));
  const hasNextPage = $derived(currentPage < totalPages);
  const hasPrevPage = $derived(currentPage > 1);

  // Filter options
  const resourceOptions = [
    { value: "", label: "All Resources" },
    { value: "flag", label: "Flags" },
    { value: "environment", label: "Environments" },
    { value: "segment", label: "Segments" },
    { value: "experiment", label: "Experiments" },
    { value: "api_key", label: "API Keys" },
  ];

  const actionOptions = [
    { value: "", label: "All Actions" },
    { value: "created", label: "Created" },
    { value: "updated", label: "Updated" },
    { value: "deleted", label: "Deleted" },
    { value: "toggled", label: "Toggled" },
  ];

  async function loadEntries() {
    try {
      loading = true;
      error = null;

      const params = new URLSearchParams();
      params.set("limit", entriesPerPage.toString());
      params.set("offset", ((currentPage - 1) * entriesPerPage).toString());

      if (resourceFilter) {
        params.set("resource", resourceFilter);
      }

      if (actionFilter) {
        // Convert action filter to full action format (e.g., "created" -> any action ending with ".created")
        // The backend expects specific action names like "flag.created", so we'll filter client-side if needed
        // For now, we'll send the action filter as-is and let the backend handle it
        params.set("action", actionFilter);
      }

      const response = await api.getAuditLog(params);
      entries = response.data || [];
      totalEntries = response.total || 0;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load audit log";
      entries = [];
      totalEntries = 0;
    } finally {
      loading = false;
    }
  }

  function handleFilterChange() {
    currentPage = 1; // Reset to first page when filters change
    loadEntries();
  }

  function goToPage(page: number) {
    if (page >= 1 && page <= totalPages) {
      currentPage = page;
      loadEntries();
    }
  }

  function nextPage() {
    if (hasNextPage) {
      goToPage(currentPage + 1);
    }
  }

  function prevPage() {
    if (hasPrevPage) {
      goToPage(currentPage - 1);
    }
  }

  function clearFilters() {
    resourceFilter = "";
    actionFilter = "";
    currentPage = 1;
    loadEntries();
  }

  // Filter entries client-side for more flexible action filtering since backend might not support partial matching
  const filteredEntries = $derived(() => {
    if (!actionFilter) return entries;

    return entries.filter(entry => {
      if (actionFilter === "created") return entry.action.endsWith(".created");
      if (actionFilter === "updated") return entry.action.endsWith(".updated");
      if (actionFilter === "deleted") return entry.action.endsWith(".deleted");
      if (actionFilter === "toggled") return entry.action.endsWith(".toggled");
      return entry.action.includes(actionFilter);
    });
  });

  onMount(() => {
    loadEntries();
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div>
    <h1 class="text-2xl font-bold text-gray-900">Audit Log</h1>
    <p class="mt-1 text-sm text-gray-600">
      Track changes and activity across your feature flags and configuration
    </p>
  </div>

  <!-- Filters -->
  <div class="bg-white shadow rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-3 lg:grid-cols-4">
        <div>
          <label for="resource-filter" class="block text-sm font-medium text-gray-700 mb-1">Resource Type</label>
          <select
            id="resource-filter"
            bind:value={resourceFilter}
            onchange={handleFilterChange}
            class="w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          >
            {#each resourceOptions as option}
              <option value={option.value}>{option.label}</option>
            {/each}
          </select>
        </div>

        <div>
          <label for="action-filter" class="block text-sm font-medium text-gray-700 mb-1">Action</label>
          <select
            id="action-filter"
            bind:value={actionFilter}
            onchange={handleFilterChange}
            class="w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          >
            {#each actionOptions as option}
              <option value={option.value}>{option.label}</option>
            {/each}
          </select>
        </div>

        <div>
          <label for="entries-per-page" class="block text-sm font-medium text-gray-700 mb-1">Per Page</label>
          <select
            id="entries-per-page"
            bind:value={entriesPerPage}
            onchange={() => { currentPage = 1; loadEntries(); }}
            class="w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          >
            <option value={10}>10</option>
            <option value={25}>25</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
          </select>
        </div>

        <div class="flex items-end">
          <button
            type="button"
            onclick={clearFilters}
            class="w-full inline-flex justify-center items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Clear Filters
          </button>
        </div>
      </div>

      {#if resourceFilter || actionFilter}
        <div class="mt-4 flex flex-wrap gap-2">
          {#if resourceFilter}
            <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              Resource: {resourceOptions.find(r => r.value === resourceFilter)?.label}
              <button
                type="button"
                onclick={() => { resourceFilter = ""; handleFilterChange(); }}
                class="ml-1 text-blue-600 hover:text-blue-800"
              >
                ×
              </button>
            </span>
          {/if}
          {#if actionFilter}
            <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-green-100 text-green-800">
              Action: {actionOptions.find(a => a.value === actionFilter)?.label}
              <button
                type="button"
                onclick={() => { actionFilter = ""; handleFilterChange(); }}
                class="ml-1 text-green-600 hover:text-green-800"
              >
                ×
              </button>
            </span>
          {/if}
        </div>
      {/if}
    </div>
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

  <!-- Content -->
  <div class="bg-white shadow rounded-lg">
    <div class="px-4 py-5 sm:p-6">
      <!-- Stats -->
      <div class="flex items-center justify-between mb-6">
        <div class="text-sm text-gray-600">
          {#if loading}
            Loading entries...
          {:else if filteredEntries.length === 0}
            No entries found
          {:else}
            Showing {((currentPage - 1) * entriesPerPage) + 1}-{Math.min(currentPage * entriesPerPage, totalEntries)} of {totalEntries} entries
          {/if}
        </div>

        {#if !loading && totalEntries > 0}
          <button
            type="button"
            onclick={loadEntries}
            class="inline-flex items-center px-3 py-1 border border-gray-300 rounded-md text-xs font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Refresh
          </button>
        {/if}
      </div>

      <!-- Timeline -->
      {#if loading}
        <div class="py-8 text-center">
          <div class="inline-flex items-center text-sm text-gray-500">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Loading audit log entries...
          </div>
        </div>
      {:else}
        <AuditTimeline entries={filteredEntries} />
      {/if}
    </div>

    <!-- Pagination -->
    {#if !loading && totalPages > 1}
      <div class="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <p class="text-sm text-gray-700">
              Page
              <span class="font-medium">{currentPage}</span>
              of
              <span class="font-medium">{totalPages}</span>
            </p>
          </div>

          <div class="flex space-x-2">
            <!-- Previous page -->
            <button
              type="button"
              onclick={prevPage}
              disabled={!hasPrevPage}
              class="relative inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-500 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M15 19l-7-7 7-7" />
              </svg>
              Previous
            </button>

            <!-- Page numbers -->
            {#each Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
              const startPage = Math.max(1, currentPage - 2);
              const endPage = Math.min(totalPages, startPage + 4);
              return Math.max(1, Math.min(endPage - 4, startPage)) + i;
            }).filter(page => page <= totalPages) as page}
              <button
                type="button"
                onclick={() => goToPage(page)}
                class="relative inline-flex items-center px-3 py-2 border text-sm font-medium rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 {
                  page === currentPage
                    ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                }"
              >
                {page}
              </button>
            {/each}

            <!-- Next page -->
            <button
              type="button"
              onclick={nextPage}
              disabled={!hasNextPage}
              class="relative inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-500 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Next
              <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>