<script lang="ts">
import { onMount } from "svelte";
import { api } from "$lib/api.js";
import type { Flag, Environment, Experiment, AuditLogEntry } from "$lib/types.js";
import StatCard from "$lib/components/StatCard.svelte";

let flags = $state<Flag[]>([]);
let environments = $state<Environment[]>([]);
let experiments = $state<Experiment[]>([]);
let auditEntries = $state<AuditLogEntry[]>([]);
let isLoading = $state(true);
let error = $state<string>("");

// Computed stats
let stats = $derived({
  totalFlags: flags.length,
  activeFlags: flags.filter(f =>
    f.is_active && Object.values(f.environments).some(e => e.enabled)
  ).length,
  totalEnvironments: environments.length,
  totalExperiments: experiments.length,
  runningExperiments: experiments.filter(e => e.status === "running").length,
});

async function loadDashboardData() {
  isLoading = true;
  error = "";

  try {
    const [flagsResponse, environmentsResponse, experimentsResponse, auditResponse] = await Promise.all([
      api.getFlags(),
      api.getEnvironments(),
      api.getExperiments(),
      api.getAuditLog(new URLSearchParams({ limit: "10" })),
    ]);

    flags = flagsResponse.data;
    environments = environmentsResponse.data;
    experiments = experimentsResponse.data;
    auditEntries = auditResponse.data;
  } catch (err) {
    error = err instanceof Error ? err.message : "Failed to load dashboard data";
    console.error("Dashboard load error:", err);
  } finally {
    isLoading = false;
  }
}

onMount(() => {
  loadDashboardData();
});

function formatTimeAgo(timestamp: string): string {
  const date = new Date(timestamp);
  const now = new Date();
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) return `${diffInSeconds}s ago`;
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`;
  return `${Math.floor(diffInSeconds / 86400)}d ago`;
}

function getActionIcon(action: string): string {
  if (action.includes("created")) return "➕";
  if (action.includes("updated")) return "📝";
  if (action.includes("deleted")) return "🗑️";
  if (action.includes("toggled")) return "🔄";
  return "📋";
}
</script>

<svelte:head>
  <title>Dashboard - FlagDeck</title>
</svelte:head>

<div class="min-h-full">
  <!-- Header -->
  <div class="bg-white border-b border-gray-200">
    <div class="px-6 py-4">
      <div class="flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
        <button
          onclick={loadDashboardData}
          disabled={isLoading}
          class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
        >
          <svg class="w-4 h-4 mr-2 {isLoading ? 'animate-spin' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Refresh
        </button>
      </div>
    </div>
  </div>

  <div class="p-6 space-y-6">
    <!-- Error State -->
    {#if error}
      <div class="bg-red-50 border border-red-200 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm text-red-800">{error}</p>
          </div>
        </div>
      </div>
    {/if}

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
      <StatCard
        label="Total Flags"
        value={stats.totalFlags}
        icon="🏁"
        loading={isLoading}
      />
      <StatCard
        label="Active Flags"
        value={stats.activeFlags}
        icon="✅"
        loading={isLoading}
      />
      <StatCard
        label="Environments"
        value={stats.totalEnvironments}
        icon="🌍"
        loading={isLoading}
      />
      <StatCard
        label="Total Experiments"
        value={stats.totalExperiments}
        icon="🧪"
        loading={isLoading}
      />
      <StatCard
        label="Running Experiments"
        value={stats.runningExperiments}
        icon="🔬"
        loading={isLoading}
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Recent Audit Log -->
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-semibold text-gray-900">Recent Activity</h2>
        </div>
        <div class="p-6">
          {#if isLoading}
            <div class="space-y-3">
              {#each Array(5) as _}
                <div class="animate-pulse flex space-x-4">
                  <div class="rounded-full bg-gray-200 h-10 w-10"></div>
                  <div class="flex-1 space-y-2 py-1">
                    <div class="h-4 bg-gray-200 rounded w-3/4"></div>
                    <div class="h-3 bg-gray-200 rounded w-1/2"></div>
                  </div>
                </div>
              {/each}
            </div>
          {:else if auditEntries.length === 0}
            <p class="text-gray-500 text-center py-8">No recent activity</p>
          {:else}
            <div class="space-y-4">
              {#each auditEntries as entry (entry.id)}
                <div class="flex items-start space-x-3">
                  <div class="flex-shrink-0">
                    <div class="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center text-sm">
                      {getActionIcon(entry.action)}
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900">
                      {entry.user_email}
                    </p>
                    <p class="text-sm text-gray-500">
                      {entry.action.replace(/\./g, " ").replace(/_/g, " ")} {entry.resource}
                      <span class="font-mono">{entry.resource_id}</span>
                    </p>
                    <p class="text-xs text-gray-400 mt-1">
                      {formatTimeAgo(entry.created_at)}
                    </p>
                  </div>
                </div>
              {/each}
            </div>
            <div class="mt-4 pt-4 border-t border-gray-200">
              <a
                href="/audit-log"
                class="text-sm text-blue-600 hover:text-blue-800 font-medium"
              >
                View all activity →
              </a>
            </div>
          {/if}
        </div>
      </div>

      <!-- Flag Status Overview -->
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-semibold text-gray-900">Flag Status Overview</h2>
        </div>
        <div class="p-6">
          {#if isLoading}
            <div class="space-y-3">
              {#each Array(5) as _}
                <div class="animate-pulse flex justify-between items-center">
                  <div class="h-4 bg-gray-200 rounded w-1/2"></div>
                  <div class="h-6 bg-gray-200 rounded w-16"></div>
                </div>
              {/each}
            </div>
          {:else if flags.length === 0}
            <p class="text-gray-500 text-center py-8">No flags found</p>
          {:else}
            <div class="space-y-3">
              {#each flags.slice(0, 8) as flag (flag.key)}
                <div class="flex items-center justify-between">
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900 truncate">
                      {flag.name}
                    </p>
                    <p class="text-xs text-gray-500 font-mono">
                      {flag.key}
                    </p>
                  </div>
                  <div class="flex items-center space-x-2">
                    <div class="flex items-center space-x-1">
                      {#each Object.entries(flag.environments) as [envKey, envConfig]}
                        <div
                          class="w-3 h-3 rounded-full {envConfig.enabled && flag.is_active ? 'bg-green-400' : 'bg-gray-300'}"
                          title="{envKey}: {envConfig.enabled && flag.is_active ? 'enabled' : 'disabled'}"
                        ></div>
                      {/each}
                    </div>
                  </div>
                </div>
              {/each}
            </div>
            <div class="mt-4 pt-4 border-t border-gray-200">
              <a
                href="/flags"
                class="text-sm text-blue-600 hover:text-blue-800 font-medium"
              >
                View all flags →
              </a>
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>