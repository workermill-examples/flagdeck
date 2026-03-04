<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api.js";
  import type { Segment, SegmentRule, Condition } from "$lib/types.js";

  let segments = $state<Segment[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let expandedSegments = $state<Set<string>>(new Set());
  let showCreateForm = $state(false);
  let editingSegment = $state<Segment | null>(null);

  // Form state
  let formData = $state({
    key: "",
    name: "",
    description: "",
    rules: [] as SegmentRule[],
  });

  // Delete confirmation
  let deletingSegment = $state<Segment | null>(null);

  const operatorLabels: Record<string, string> = {
    equals: "equals",
    not_equals: "does not equal",
    contains: "contains",
    in: "is one of",
    not_in: "is not one of",
    gt: "greater than",
    lt: "less than",
    gte: "greater than or equal",
    lte: "less than or equal",
    regex: "matches regex",
  };

  async function loadSegments() {
    try {
      loading = true;
      const response = await api.getSegments();
      segments = response.data;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load segments";
    } finally {
      loading = false;
    }
  }

  function toggleExpanded(segmentId: string) {
    if (expandedSegments.has(segmentId)) {
      expandedSegments.delete(segmentId);
    } else {
      expandedSegments.add(segmentId);
    }
    expandedSegments = new Set(expandedSegments);
  }

  function resetForm() {
    formData = {
      key: "",
      name: "",
      description: "",
      rules: [{ conditions: [{ property: "", operator: "equals", value: "" }] }],
    };
    showCreateForm = false;
    editingSegment = null;
  }

  function startCreate() {
    resetForm();
    showCreateForm = true;
  }

  function startEdit(segment: Segment) {
    formData = {
      key: segment.key,
      name: segment.name,
      description: segment.description,
      rules: JSON.parse(JSON.stringify(segment.rules)), // Deep copy
    };
    editingSegment = segment;
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

  function addRule() {
    formData.rules.push({ conditions: [{ property: "", operator: "equals", value: "" }] });
    formData.rules = [...formData.rules];
  }

  function removeRule(ruleIndex: number) {
    formData.rules.splice(ruleIndex, 1);
    formData.rules = [...formData.rules];
  }

  function addCondition(ruleIndex: number) {
    formData.rules[ruleIndex].conditions.push({ property: "", operator: "equals", value: "" });
    formData.rules = [...formData.rules];
  }

  function removeCondition(ruleIndex: number, conditionIndex: number) {
    formData.rules[ruleIndex].conditions.splice(conditionIndex, 1);
    formData.rules = [...formData.rules];
  }

  function updateCondition(ruleIndex: number, conditionIndex: number, field: keyof Condition, value: string) {
    formData.rules[ruleIndex].conditions[conditionIndex][field] = value;
    formData.rules = [...formData.rules];
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    // Validate rules
    const validRules = formData.rules.filter(rule =>
      rule.conditions.some(c => c.property.trim() && c.value !== "")
    );

    if (validRules.length === 0) {
      error = "Please add at least one rule with valid conditions";
      return;
    }

    try {
      const payload = { ...formData, rules: validRules };
      if (editingSegment) {
        await api.updateSegment(editingSegment.key, payload);
      } else {
        await api.createSegment(payload);
      }
      await loadSegments();
      resetForm();
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to save segment";
    }
  }

  async function handleDelete() {
    if (!deletingSegment) return;

    try {
      await api.deleteSegment(deletingSegment.key);
      await loadSegments();
      deletingSegment = null;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to delete segment";
    }
  }

  function formatConditions(conditions: Condition[]): string {
    return conditions
      .filter(c => c.property && c.value !== "")
      .map(c => `${c.property} ${operatorLabels[c.operator] || c.operator} ${c.value}`)
      .join(" AND ");
  }

  function formatRules(rules: SegmentRule[]): string {
    return rules
      .map(rule => {
        const conditionText = formatConditions(rule.conditions);
        return conditionText ? `(${conditionText})` : "";
      })
      .filter(text => text)
      .join(" OR ");
  }

  onMount(() => {
    loadSegments();
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Segments</h1>
      <p class="mt-1 text-sm text-gray-600">
        Define user segments based on attributes and conditions
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
      Add Segment
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
          {editingSegment ? "Edit Segment" : "Create Segment"}
        </h3>
        <form onsubmit={handleSubmit} class="space-y-6">
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
                placeholder="Beta Users"
              />
            </div>
            <div>
              <label for="key" class="block text-sm font-medium text-gray-700">Key</label>
              <input
                type="text"
                id="key"
                bind:value={formData.key}
                required
                readonly={!!editingSegment}
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm {editingSegment ? 'bg-gray-50' : ''}"
                placeholder="beta-users"
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
              placeholder="Users participating in beta testing"
            ></textarea>
          </div>

          <!-- Segment Rules -->
          <div>
            <div class="flex items-center justify-between mb-3">
              <label class="block text-sm font-medium text-gray-700">Segment Rules</label>
              <button
                type="button"
                onclick={addRule}
                class="inline-flex items-center px-3 py-1 border border-gray-300 rounded-md text-xs font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M12 4v16m8-8H4" />
                </svg>
                Add Rule
              </button>
            </div>
            <div class="space-y-4">
              {#each formData.rules as rule, ruleIndex (ruleIndex)}
                <div class="border border-gray-200 rounded-lg p-4">
                  <div class="flex items-center justify-between mb-3">
                    <h4 class="text-sm font-medium text-gray-700">Rule {ruleIndex + 1}</h4>
                    {#if formData.rules.length > 1}
                      <button
                        type="button"
                        onclick={() => removeRule(ruleIndex)}
                        class="text-red-600 hover:text-red-800 text-sm"
                      >
                        Remove Rule
                      </button>
                    {/if}
                  </div>
                  <div class="space-y-3">
                    {#each rule.conditions as condition, conditionIndex (conditionIndex)}
                      <div class="flex items-center space-x-2">
                        <input
                          type="text"
                          placeholder="Property"
                          value={condition.property}
                          oninput={(e) => updateCondition(ruleIndex, conditionIndex, "property", e.target.value)}
                          class="flex-1 border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        />
                        <select
                          value={condition.operator}
                          onchange={(e) => updateCondition(ruleIndex, conditionIndex, "operator", e.target.value)}
                          class="border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        >
                          {#each Object.entries(operatorLabels) as [value, label]}
                            <option {value}>{label}</option>
                          {/each}
                        </select>
                        <input
                          type="text"
                          placeholder="Value"
                          value={condition.value}
                          oninput={(e) => updateCondition(ruleIndex, conditionIndex, "value", e.target.value)}
                          class="flex-1 border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        />
                        <button
                          type="button"
                          onclick={() => addCondition(ruleIndex)}
                          class="inline-flex items-center p-1 border border-gray-300 rounded text-gray-600 hover:text-gray-800"
                          title="Add condition"
                        >
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M12 4v16m8-8H4" />
                          </svg>
                        </button>
                        {#if rule.conditions.length > 1}
                          <button
                            type="button"
                            onclick={() => removeCondition(ruleIndex, conditionIndex)}
                            class="inline-flex items-center p-1 border border-gray-300 rounded text-red-600 hover:text-red-800"
                            title="Remove condition"
                          >
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        {/if}
                      </div>
                      {#if conditionIndex < rule.conditions.length - 1}
                        <div class="text-center text-sm text-gray-500 font-medium">AND</div>
                      {/if}
                    {/each}
                  </div>
                </div>
                {#if ruleIndex < formData.rules.length - 1}
                  <div class="text-center text-sm text-gray-500 font-medium py-2">OR</div>
                {/if}
              {/each}
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
              {editingSegment ? "Update" : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}

  <!-- Segments list -->
  <div class="bg-white shadow overflow-hidden sm:rounded-md">
    {#if loading}
      <div class="px-4 py-8 text-center">
        <div class="inline-flex items-center text-sm text-gray-500">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Loading segments...
        </div>
      </div>
    {:else if segments.length === 0}
      <div class="text-center py-8">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No segments</h3>
        <p class="mt-1 text-sm text-gray-500">Get started by creating your first user segment.</p>
      </div>
    {:else}
      <ul class="divide-y divide-gray-200">
        {#each segments as segment (segment.id)}
          <li class="px-6 py-4">
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-2">
                  <h3 class="text-sm font-medium text-gray-900">{segment.name}</h3>
                  <code class="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded">{segment.key}</code>
                  <button
                    type="button"
                    onclick={() => toggleExpanded(segment.id)}
                    class="text-gray-400 hover:text-gray-600"
                  >
                    <svg class="w-5 h-5 transform {expandedSegments.has(segment.id) ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M19 9l-7 7-7-7" />
                    </svg>
                  </button>
                </div>
                {#if segment.description}
                  <p class="mt-1 text-sm text-gray-500">{segment.description}</p>
                {/if}
                <div class="mt-1 text-xs text-gray-600">
                  {formatRules(segment.rules) || "No rules defined"}
                </div>
              </div>
              <div class="flex items-center space-x-2 ml-4">
                <button
                  type="button"
                  class="inline-flex items-center px-3 py-1 border border-gray-300 rounded-md text-xs font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  onclick={() => startEdit(segment)}
                >
                  Edit
                </button>
                <button
                  type="button"
                  class="inline-flex items-center px-3 py-1 border border-red-300 rounded-md text-xs font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                  onclick={() => (deletingSegment = segment)}
                >
                  Delete
                </button>
              </div>
            </div>

            <!-- Expanded detail -->
            {#if expandedSegments.has(segment.id)}
              <div class="mt-4 border-t border-gray-200 pt-4">
                <h4 class="text-sm font-medium text-gray-900 mb-2">Segment Rules</h4>
                <div class="space-y-3">
                  {#each segment.rules as rule, ruleIndex (ruleIndex)}
                    <div class="bg-gray-50 rounded-lg p-3">
                      <h5 class="text-xs font-medium text-gray-700 mb-2">Rule {ruleIndex + 1}</h5>
                      <div class="space-y-1">
                        {#each rule.conditions as condition, conditionIndex (conditionIndex)}
                          <div class="text-sm text-gray-600">
                            <code class="bg-white px-1 rounded">{condition.property}</code>
                            <span class="mx-1">{operatorLabels[condition.operator] || condition.operator}</span>
                            <code class="bg-white px-1 rounded">{condition.value}</code>
                            {#if conditionIndex < rule.conditions.length - 1}
                              <span class="mx-2 text-blue-600 font-medium">AND</span>
                            {/if}
                          </div>
                        {/each}
                      </div>
                    </div>
                    {#if ruleIndex < segment.rules.length - 1}
                      <div class="text-center text-sm text-blue-600 font-medium">OR</div>
                    {/if}
                  {/each}
                </div>
              </div>
            {/if}
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>

<!-- Delete confirmation modal -->
{#if deletingSegment}
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
        <h3 class="text-lg font-medium text-gray-900">Delete Segment</h3>
        <div class="mt-2">
          <p class="text-sm text-gray-500">
            Are you sure you want to delete the <strong>{deletingSegment.name}</strong> segment?
            This action cannot be undone.
          </p>
        </div>
      </div>
      <div class="mt-5 flex justify-center space-x-3">
        <button
          type="button"
          class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm bg-white text-sm font-medium text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          onclick={() => (deletingSegment = null)}
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