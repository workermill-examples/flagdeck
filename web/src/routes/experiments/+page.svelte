<script lang="ts">
  import { onMount } from "svelte";
  import { api } from "$lib/api.js";
  import type { Experiment, ExperimentVariant } from "$lib/types.js";
  import ExperimentChart from "$lib/components/ExperimentChart.svelte";

  let experiments = $state<Experiment[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let expandedExperiments = $state<Set<string>>(new Set());
  let showCreateForm = $state(false);
  let editingExperiment = $state<Experiment | null>(null);

  // Form state
  let formData = $state({
    key: "",
    name: "",
    description: "",
    flag_key: "",
    status: "draft" as "draft" | "running" | "paused" | "completed",
    variants: [
      { key: "control", name: "Control", weight: 50, value: false },
      { key: "variant", name: "Variant", weight: 50, value: true },
    ] as Omit<ExperimentVariant, "results">[],
  });

  // Delete confirmation
  let deletingExperiment = $state<Experiment | null>(null);

  const statusColors: Record<string, string> = {
    draft: "bg-gray-100 text-gray-800",
    running: "bg-green-100 text-green-800",
    paused: "bg-yellow-100 text-yellow-800",
    completed: "bg-blue-100 text-blue-800",
  };

  const statusIcons: Record<string, string> = {
    draft: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z",
    running: "M14.828 14.828a4 4 0 01-5.656 0M9 10h1.586a1 1 0 01.707.293l2.414 2.414a1 1 0 00.707.293H15M9 10v4a2 2 0 002 2h2a2 2 0 002-2v-4M9 10V9a2 2 0 012-2h2a2 2 0 012 2v1",
    paused: "M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z",
    completed: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z",
  };

  async function loadExperiments() {
    try {
      loading = true;
      const response = await api.getExperiments();
      experiments = response.data;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to load experiments";
    } finally {
      loading = false;
    }
  }

  function toggleExpanded(experimentId: string) {
    if (expandedExperiments.has(experimentId)) {
      expandedExperiments.delete(experimentId);
    } else {
      expandedExperiments.add(experimentId);
    }
    expandedExperiments = new Set(expandedExperiments);
  }

  function resetForm() {
    formData = {
      key: "",
      name: "",
      description: "",
      flag_key: "",
      status: "draft",
      variants: [
        { key: "control", name: "Control", weight: 50, value: false },
        { key: "variant", name: "Variant", weight: 50, value: true },
      ],
    };
    showCreateForm = false;
    editingExperiment = null;
  }

  function startCreate() {
    resetForm();
    showCreateForm = true;
  }

  function startEdit(experiment: Experiment) {
    formData = {
      key: experiment.key,
      name: experiment.name,
      description: experiment.description,
      flag_key: experiment.flag_key,
      status: experiment.status,
      variants: experiment.variants.map(v => ({
        key: v.key,
        name: v.name,
        weight: v.weight,
        value: v.value,
      })),
    };
    editingExperiment = experiment;
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

  function addVariant() {
    const newVariant = {
      key: `variant-${formData.variants.length}`,
      name: `Variant ${formData.variants.length}`,
      weight: 0,
      value: true,
    };
    formData.variants.push(newVariant);
    rebalanceWeights();
  }

  function removeVariant(index: number) {
    formData.variants.splice(index, 1);
    rebalanceWeights();
  }

  function rebalanceWeights() {
    const count = formData.variants.length;
    if (count === 0) return;

    const evenWeight = Math.floor(100 / count);
    const remainder = 100 % count;

    formData.variants.forEach((variant, index) => {
      variant.weight = evenWeight + (index < remainder ? 1 : 0);
    });
    formData.variants = [...formData.variants];
  }

  function updateVariantWeight(index: number, weight: number) {
    formData.variants[index].weight = weight;

    // Adjust other weights to maintain 100% total
    const totalOtherWeights = formData.variants.reduce((sum, v, i) => i === index ? sum : sum + v.weight, 0);
    const remaining = 100 - weight;

    if (totalOtherWeights > 0 && remaining >= 0) {
      const ratio = remaining / totalOtherWeights;
      formData.variants.forEach((variant, i) => {
        if (i !== index) {
          variant.weight = Math.round(variant.weight * ratio);
        }
      });

      // Fix rounding errors
      const actualTotal = formData.variants.reduce((sum, v) => sum + v.weight, 0);
      if (actualTotal !== 100 && formData.variants.length > 1) {
        const diff = 100 - actualTotal;
        const otherIndex = index === 0 ? 1 : 0;
        formData.variants[otherIndex].weight += diff;
      }
    }

    formData.variants = [...formData.variants];
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    // Validate weights sum to 100
    const totalWeight = formData.variants.reduce((sum, v) => sum + v.weight, 0);
    if (Math.abs(totalWeight - 100) > 0.1) {
      error = "Variant weights must sum to 100%";
      return;
    }

    try {
      if (editingExperiment) {
        await api.updateExperiment(editingExperiment.key, formData);
      } else {
        await api.createExperiment(formData);
      }
      await loadExperiments();
      resetForm();
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to save experiment";
    }
  }

  async function handleDelete() {
    if (!deletingExperiment) return;

    try {
      await api.deleteExperiment(deletingExperiment.key);
      await loadExperiments();
      deletingExperiment = null;
    } catch (err) {
      error = err instanceof Error ? err.message : "Failed to delete experiment";
    }
  }

  function formatDate(dateString: string | null): string {
    if (!dateString) return "Not set";
    return new Date(dateString).toLocaleDateString();
  }

  onMount(() => {
    loadExperiments();
  });
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Experiments</h1>
      <p class="mt-1 text-sm text-gray-600">
        A/B test variants and measure performance
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
      Create Experiment
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
          {editingExperiment ? "Edit Experiment" : "Create Experiment"}
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
                placeholder="Checkout Redesign Test"
              />
            </div>
            <div>
              <label for="key" class="block text-sm font-medium text-gray-700">Key</label>
              <input
                type="text"
                id="key"
                bind:value={formData.key}
                required
                readonly={!!editingExperiment}
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm {editingExperiment ? 'bg-gray-50' : ''}"
                placeholder="checkout-redesign-test"
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
              placeholder="Test the new checkout flow design against the current one"
            ></textarea>
          </div>

          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <label for="flag_key" class="block text-sm font-medium text-gray-700">Flag Key</label>
              <input
                type="text"
                id="flag_key"
                bind:value={formData.flag_key}
                required
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                placeholder="new-checkout-flow"
              />
            </div>
            <div>
              <label for="status" class="block text-sm font-medium text-gray-700">Status</label>
              <select
                id="status"
                bind:value={formData.status}
                class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              >
                <option value="draft">Draft</option>
                <option value="running">Running</option>
                <option value="paused">Paused</option>
                <option value="completed">Completed</option>
              </select>
            </div>
          </div>

          <!-- Variants -->
          <div>
            <div class="flex items-center justify-between mb-3">
              <label class="block text-sm font-medium text-gray-700">Experiment Variants</label>
              <button
                type="button"
                onclick={addVariant}
                class="inline-flex items-center px-3 py-1 border border-gray-300 rounded-md text-xs font-medium text-gray-700 bg-white hover:bg-gray-50"
              >
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M12 4v16m8-8H4" />
                </svg>
                Add Variant
              </button>
            </div>
            <div class="space-y-4">
              {#each formData.variants as variant, index (index)}
                <div class="border border-gray-200 rounded-lg p-4">
                  <div class="flex items-center justify-between mb-3">
                    <h4 class="text-sm font-medium text-gray-700">Variant {index + 1}</h4>
                    {#if formData.variants.length > 1}
                      <button
                        type="button"
                        onclick={() => removeVariant(index)}
                        class="text-red-600 hover:text-red-800 text-sm"
                      >
                        Remove
                      </button>
                    {/if}
                  </div>
                  <div class="grid grid-cols-1 gap-4 sm:grid-cols-4">
                    <div>
                      <label class="block text-xs font-medium text-gray-700 mb-1">Key</label>
                      <input
                        type="text"
                        bind:value={variant.key}
                        required
                        class="w-full text-sm border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        placeholder="control"
                      />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-gray-700 mb-1">Name</label>
                      <input
                        type="text"
                        bind:value={variant.name}
                        required
                        class="w-full text-sm border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        placeholder="Control"
                      />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-gray-700 mb-1">Weight (%)</label>
                      <input
                        type="number"
                        bind:value={variant.weight}
                        oninput={(e) => updateVariantWeight(index, parseInt(e.target.value) || 0)}
                        min="0"
                        max="100"
                        required
                        class="w-full text-sm border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                      />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-gray-700 mb-1">Value</label>
                      <input
                        type="text"
                        bind:value={variant.value}
                        required
                        class="w-full text-sm border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
                        placeholder="false"
                      />
                    </div>
                  </div>
                </div>
              {/each}
              <div class="text-center text-sm text-gray-600">
                Total weight: {formData.variants.reduce((sum, v) => sum + v.weight, 0)}%
                {#if Math.abs(formData.variants.reduce((sum, v) => sum + v.weight, 0) - 100) > 0.1}
                  <span class="text-red-600 font-medium">- Must equal 100%</span>
                {/if}
              </div>
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
              {editingExperiment ? "Update" : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  {/if}

  <!-- Experiments list -->
  <div class="bg-white shadow overflow-hidden sm:rounded-md">
    {#if loading}
      <div class="px-4 py-8 text-center">
        <div class="inline-flex items-center text-sm text-gray-500">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Loading experiments...
        </div>
      </div>
    {:else if experiments.length === 0}
      <div class="text-center py-8">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No experiments</h3>
        <p class="mt-1 text-sm text-gray-500">Get started by creating your first experiment.</p>
      </div>
    {:else}
      <ul class="divide-y divide-gray-200">
        {#each experiments as experiment (experiment.id)}
          <li class="px-6 py-4">
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-3">
                  <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d={statusIcons[experiment.status]} />
                  </svg>
                  <div>
                    <div class="flex items-center space-x-2">
                      <h3 class="text-sm font-medium text-gray-900">{experiment.name}</h3>
                      <code class="text-xs bg-gray-100 text-gray-800 px-2 py-1 rounded">{experiment.key}</code>
                      <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium {statusColors[experiment.status]}">
                        {experiment.status}
                      </span>
                      <button
                        type="button"
                        onclick={() => toggleExpanded(experiment.id)}
                        class="text-gray-400 hover:text-gray-600"
                      >
                        <svg class="w-5 h-5 transform {expandedExperiments.has(experiment.id) ? 'rotate-180' : ''}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                    </div>
                    {#if experiment.description}
                      <p class="mt-1 text-sm text-gray-500">{experiment.description}</p>
                    {/if}
                    <div class="mt-1 text-xs text-gray-600">
                      Flag: <code class="bg-gray-100 px-1 rounded">{experiment.flag_key}</code>
                      {#if experiment.start_date}
                        • Started: {formatDate(experiment.start_date)}
                      {/if}
                      {#if experiment.end_date}
                        • Ended: {formatDate(experiment.end_date)}
                      {/if}
                    </div>
                  </div>
                </div>
              </div>
              <div class="flex items-center space-x-2 ml-4">
                <button
                  type="button"
                  class="inline-flex items-center px-3 py-1 border border-gray-300 rounded-md text-xs font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                  onclick={() => startEdit(experiment)}
                >
                  Edit
                </button>
                <button
                  type="button"
                  class="inline-flex items-center px-3 py-1 border border-red-300 rounded-md text-xs font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                  onclick={() => (deletingExperiment = experiment)}
                >
                  Delete
                </button>
              </div>
            </div>

            <!-- Expanded details -->
            {#if expandedExperiments.has(experiment.id)}
              <div class="mt-4 border-t border-gray-200 pt-4">
                <div class="space-y-4">
                  <!-- Variant weights -->
                  <div>
                    <h4 class="text-sm font-medium text-gray-900 mb-2">Variant Configuration</h4>
                    <div class="grid grid-cols-1 gap-2 sm:grid-cols-2 lg:grid-cols-3">
                      {#each experiment.variants as variant}
                        <div class="bg-gray-50 rounded-lg p-3">
                          <div class="flex items-center justify-between mb-1">
                            <span class="text-sm font-medium text-gray-900">{variant.name}</span>
                            <span class="text-xs text-gray-600">{variant.weight}%</span>
                          </div>
                          <div class="text-xs text-gray-600">
                            Key: <code class="bg-white px-1 rounded">{variant.key}</code>
                          </div>
                          <div class="text-xs text-gray-600">
                            Value: <code class="bg-white px-1 rounded">{JSON.stringify(variant.value)}</code>
                          </div>
                        </div>
                      {/each}
                    </div>
                  </div>

                  <!-- Results chart -->
                  <div>
                    <h4 class="text-sm font-medium text-gray-900 mb-2">Performance Results</h4>
                    <ExperimentChart variants={experiment.variants} />
                  </div>
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
{#if deletingExperiment}
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
        <h3 class="text-lg font-medium text-gray-900">Delete Experiment</h3>
        <div class="mt-2">
          <p class="text-sm text-gray-500">
            Are you sure you want to delete the <strong>{deletingExperiment.name}</strong> experiment?
            This action cannot be undone and will remove all associated data.
          </p>
        </div>
      </div>
      <div class="mt-5 flex justify-center space-x-3">
        <button
          type="button"
          class="inline-flex justify-center px-4 py-2 border border-gray-300 shadow-sm bg-white text-sm font-medium text-gray-700 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          onclick={() => (deletingExperiment = null)}
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