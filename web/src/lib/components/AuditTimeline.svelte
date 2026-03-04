<script lang="ts">
  import type { AuditLogEntry } from "$lib/types.js";

  interface Props {
    entries: AuditLogEntry[];
    className?: string;
  }

  let { entries, className = "" }: Props = $props();

  const actionIcons: Record<string, string> = {
    "flag.created":
      "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z",
    "flag.updated":
      "M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h14a2 2 0 002-2V7a2 2 0 00-2-2h-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z",
    "flag.toggled":
      "M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z",
    "flag.deleted":
      "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16",
    "environment.created":
      "M19 11H5m14-7H3a2 2 0 00-2 2v12a2 2 0 002 2h16a2 2 0 002-2V6a2 2 0 00-2-2z",
    "environment.updated":
      "M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h14a2 2 0 002-2V7a2 2 0 00-2-2h-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z",
    "environment.deleted":
      "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16",
    "segment.created":
      "M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z",
    "segment.updated":
      "M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h14a2 2 0 002-2V7a2 2 0 00-2-2h-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z",
    "segment.deleted":
      "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16",
    "experiment.created":
      "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z",
    "experiment.updated":
      "M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h14a2 2 0 002-2V7a2 2 0 00-2-2h-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z",
    "experiment.deleted":
      "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16",
    "api_key.created":
      "M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z",
    "api_key.deleted":
      "M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16",
    default: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z",
  };

  const actionColors: Record<string, string> = {
    "flag.created": "bg-green-100 text-green-600",
    "flag.updated": "bg-blue-100 text-blue-600",
    "flag.toggled": "bg-yellow-100 text-yellow-600",
    "flag.deleted": "bg-red-100 text-red-600",
    "environment.created": "bg-green-100 text-green-600",
    "environment.updated": "bg-blue-100 text-blue-600",
    "environment.deleted": "bg-red-100 text-red-600",
    "segment.created": "bg-green-100 text-green-600",
    "segment.updated": "bg-blue-100 text-blue-600",
    "segment.deleted": "bg-red-100 text-red-600",
    "experiment.created": "bg-green-100 text-green-600",
    "experiment.updated": "bg-blue-100 text-blue-600",
    "experiment.deleted": "bg-red-100 text-red-600",
    "api_key.created": "bg-green-100 text-green-600",
    "api_key.deleted": "bg-red-100 text-red-600",
    default: "bg-gray-100 text-gray-600",
  };

  let expandedEntries = $state<Set<string>>(new Set());

  function toggleExpanded(entryId: string) {
    if (expandedEntries.has(entryId)) {
      expandedEntries.delete(entryId);
    } else {
      expandedEntries.add(entryId);
    }
    expandedEntries = new Set(expandedEntries);
  }

  function formatAction(action: string): string {
    return action
      .split(".")
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(" ");
  }

  function formatTimestamp(timestamp: string): string {
    const date = new Date(timestamp);
    const now = new Date();
    const diff = now.getTime() - date.getTime();

    // Less than a minute
    if (diff < 60000) {
      return "Just now";
    }

    // Less than an hour
    if (diff < 3600000) {
      const minutes = Math.floor(diff / 60000);
      return `${minutes} minute${minutes === 1 ? "" : "s"} ago`;
    }

    // Less than a day
    if (diff < 86400000) {
      const hours = Math.floor(diff / 3600000);
      return `${hours} hour${hours === 1 ? "" : "s"} ago`;
    }

    // Less than a week
    if (diff < 604800000) {
      const days = Math.floor(diff / 86400000);
      return `${days} day${days === 1 ? "" : "s"} ago`;
    }

    // More than a week, show the date
    return date.toLocaleDateString();
  }

  function getAbsoluteTimestamp(timestamp: string): string {
    return new Date(timestamp).toLocaleString();
  }

  function getUserInitials(email: string): string {
    return email.split("@")[0].slice(0, 2).toUpperCase();
  }

  function hasChanges(entry: AuditLogEntry): boolean {
    return (
      entry.changes != null &&
      typeof entry.changes === "object" &&
      Object.keys(entry.changes).length > 0
    );
  }

  function formatChanges(
    changes: unknown,
  ): { field: string; before: string; after: string }[] {
    if (!changes || typeof changes !== "object") return [];

    const result: { field: string; before: string; after: string }[] = [];

    for (const [field, change] of Object.entries(changes)) {
      if (
        change &&
        typeof change === "object" &&
        "before" in change &&
        "after" in change
      ) {
        result.push({
          field,
          before: JSON.stringify(change.before),
          after: JSON.stringify(change.after),
        });
      }
    }

    return result;
  }
</script>

<div class="audit-timeline {className}">
  {#if entries.length === 0}
    <div class="text-center py-8">
      <svg
        class="mx-auto h-12 w-12 text-gray-400"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width={2}
          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">
        No audit log entries
      </h3>
      <p class="mt-1 text-sm text-gray-500">
        Activity will appear here as changes are made.
      </p>
    </div>
  {:else}
    <div class="flow-root">
      <ul class="-mb-8">
        {#each entries as entry, index (entry.id)}
          <li>
            <div class="relative pb-8">
              {#if index !== entries.length - 1}
                <span
                  class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                  aria-hidden="true"
                ></span>
              {/if}

              <div class="relative flex space-x-3">
                <!-- Icon -->
                <div class="relative">
                  <div
                    class="flex items-center justify-center h-8 w-8 rounded-full {actionColors[
                      entry.action
                    ] || actionColors.default}"
                  >
                    <svg
                      class="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width={2}
                        d={actionIcons[entry.action] || actionIcons.default}
                      />
                    </svg>
                  </div>
                </div>

                <!-- Content -->
                <div class="flex-1 min-w-0">
                  <div>
                    <div class="text-sm">
                      <span class="font-medium text-gray-900"
                        >{formatAction(entry.action)}</span
                      >
                    </div>
                    <div class="mt-1 text-sm text-gray-500">
                      <div class="flex items-center space-x-2">
                        <!-- User avatar -->
                        <div class="flex items-center space-x-1">
                          <div
                            class="flex items-center justify-center h-5 w-5 rounded-full bg-gray-300 text-xs font-medium text-gray-600"
                          >
                            {getUserInitials(entry.user_email)}
                          </div>
                          <span class="text-xs">{entry.user_email}</span>
                        </div>

                        <!-- Resource info -->
                        <span class="text-gray-400">•</span>
                        <span class="text-xs">
                          {entry.resource}
                          {#if entry.resource_id}
                            <code class="bg-gray-100 px-1 rounded text-xs"
                              >{entry.resource_id}</code
                            >
                          {/if}
                        </span>
                      </div>
                    </div>

                    <!-- Timestamp -->
                    <div class="mt-1">
                      <span
                        class="text-xs text-gray-500 cursor-help"
                        title={getAbsoluteTimestamp(entry.timestamp)}
                      >
                        {formatTimestamp(entry.timestamp)}
                      </span>

                      <!-- Expand button -->
                      {#if hasChanges(entry)}
                        <button
                          type="button"
                          onclick={() => toggleExpanded(entry.id)}
                          class="ml-2 text-xs text-blue-600 hover:text-blue-800"
                        >
                          {expandedEntries.has(entry.id) ? "Hide" : "Show"} changes
                        </button>
                      {/if}
                    </div>

                    <!-- Expanded changes -->
                    {#if expandedEntries.has(entry.id) && hasChanges(entry)}
                      <div class="mt-3 p-3 bg-gray-50 rounded-md">
                        <h4 class="text-xs font-medium text-gray-700 mb-2">
                          Changes
                        </h4>
                        <div class="space-y-2">
                          {#each formatChanges(entry.changes) as change}
                            <div class="text-xs">
                              <div class="font-medium text-gray-700">
                                {change.field}
                              </div>
                              <div class="mt-1 grid grid-cols-2 gap-2">
                                <div>
                                  <span class="text-gray-500">Before:</span>
                                  <code
                                    class="block bg-red-50 text-red-800 px-1 py-0.5 rounded text-xs mt-0.5 break-all"
                                    >{change.before}</code
                                  >
                                </div>
                                <div>
                                  <span class="text-gray-500">After:</span>
                                  <code
                                    class="block bg-green-50 text-green-800 px-1 py-0.5 rounded text-xs mt-0.5 break-all"
                                    >{change.after}</code
                                  >
                                </div>
                              </div>
                            </div>
                          {/each}
                        </div>
                      </div>
                    {/if}
                  </div>
                </div>
              </div>
            </div>
          </li>
        {/each}
      </ul>
    </div>
  {/if}
</div>

<style>
  .audit-timeline {
    @apply w-full;
  }
</style>
