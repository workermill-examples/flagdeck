<script lang="ts">
	import { page } from '$app/stores';
	import { onMount } from 'svelte';
	import { api } from '../../../lib/api.js';
	import FlagToggle from '../../../lib/components/FlagToggle.svelte';
	import RolloutSlider from '../../../lib/components/RolloutSlider.svelte';
	import TargetingRuleBuilder from '../../../lib/components/TargetingRuleBuilder.svelte';
	import type { Flag, Environment } from '../../../lib/types.js';

	let flag = $state<Flag | null>(null);
	let environments = $state<Environment[]>([]);
	let isLoading = $state(true);
	let isSaving = $state(false);
	let error = $state<string | null>(null);
	let activeTab = $state('settings');

	// Form state for basic settings
	let editingName = $state('');
	let editingDescription = $state('');
	let editingTags = $state('');

	const flagKey = $derived($page.params.id);

	onMount(async () => {
		await loadData();
	});

	async function loadData() {
		if (!flagKey) return;

		try {
			isLoading = true;
			error = null;

			const [flagResponse, environmentsResponse] = await Promise.all([
				api.getFlag(flagKey),
				api.getEnvironments()
			]);

			flag = flagResponse;
			environments = environmentsResponse.data.sort((a, b) => a.sort_order - b.sort_order);

			// Initialize edit form
			if (flag) {
				editingName = flag.name;
				editingDescription = flag.description;
				editingTags = flag.tags.join(', ');
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load flag';
			console.error('Failed to load flag:', err);
		} finally {
			isLoading = false;
		}
	}

	async function handleToggle() {
		await loadData();
	}

	async function handleUpdate() {
		await loadData();
	}

	async function saveBasicSettings() {
		if (!flag || isSaving) return;

		try {
			isSaving = true;
			error = null;

			const updatedFlag = {
				...flag,
				name: editingName.trim(),
				description: editingDescription.trim(),
				tags: editingTags
					.split(',')
					.map(tag => tag.trim())
					.filter(tag => tag.length > 0)
			};

			await api.updateFlag(flag.key, updatedFlag);
			await loadData(); // Reload to get the updated flag
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to update flag';
			console.error('Failed to update flag:', err);
		} finally {
			isSaving = false;
		}
	}

	async function toggleGlobalStatus() {
		if (!flag || isSaving) return;

		try {
			isSaving = true;
			error = null;

			await api.toggleFlag(flag.key); // No environment = global toggle
			await loadData();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to toggle flag';
			console.error('Failed to toggle flag:', err);
		} finally {
			isSaving = false;
		}
	}

	function formatValue(value: unknown): string {
		if (typeof value === 'string') return value;
		if (typeof value === 'object') {
			try {
				return JSON.stringify(value, null, 2);
			} catch {
				return String(value);
			}
		}
		return String(value);
	}

	function goBack() {
		window.location.href = '/flags';
	}
</script>

<svelte:head>
	<title>{flag?.name ?? 'Flag'} - FlagDeck</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div>
		<div class="flex items-center space-x-2 text-sm text-gray-500 mb-2">
			<button onclick={goBack} class="hover:text-blue-600">Feature Flags</button>
			<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
				<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"></path>
			</svg>
			<span>{flag?.key ?? flagKey}</span>
		</div>
		{#if flag}
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-2xl font-bold text-gray-900">{flag.name}</h1>
					<div class="flex items-center space-x-4 mt-1">
						<span class="text-sm text-gray-500">Key: {flag.key}</span>
						<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
							{flag.type === 'boolean' ? 'bg-blue-100 text-blue-800' :
							flag.type === 'string' ? 'bg-green-100 text-green-800' :
							flag.type === 'number' ? 'bg-purple-100 text-purple-800' :
							'bg-yellow-100 text-yellow-800'}">
							{flag.type}
						</span>
						<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
							{flag.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
							{flag.is_active ? 'Active' : 'Inactive'}
						</span>
					</div>
				</div>
				<div class="flex items-center space-x-3">
					<button
						onclick={toggleGlobalStatus}
						disabled={isSaving}
						class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md
							{flag.is_active
								? 'text-red-700 bg-red-100 hover:bg-red-200 focus:ring-red-500'
								: 'text-green-700 bg-green-100 hover:bg-green-200 focus:ring-green-500'}
							focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
					>
						{#if isSaving}
							...
						{:else}
							{flag.is_active ? 'Disable Flag' : 'Enable Flag'}
						{/if}
					</button>
				</div>
			</div>
		{/if}
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
		</div>
	{:else if error && !flag}
		<div class="bg-red-50 border border-red-200 rounded-md p-4">
			<div class="flex">
				<svg class="flex-shrink-0 h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
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
	{:else if flag}
		<!-- Tab Navigation -->
		<div class="border-b border-gray-200">
			<nav class="-mb-px flex space-x-8" aria-label="Tabs">
				<button
					onclick={() => activeTab = 'settings'}
					class="py-2 px-1 border-b-2 font-medium text-sm
						{activeTab === 'settings'
							? 'border-blue-500 text-blue-600'
							: 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
				>
					Settings
				</button>
				<button
					onclick={() => activeTab = 'environments'}
					class="py-2 px-1 border-b-2 font-medium text-sm
						{activeTab === 'environments'
							? 'border-blue-500 text-blue-600'
							: 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'}"
				>
					Environments
				</button>
			</nav>
		</div>

		<!-- Error Display -->
		{#if error}
			<div class="bg-red-50 border border-red-200 rounded-md p-4">
				<div class="flex">
					<svg class="flex-shrink-0 h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
					</svg>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Error</h3>
						<p class="text-sm text-red-700 mt-1">{error}</p>
					</div>
				</div>
			</div>
		{/if}

		<!-- Tab Content -->
		{#if activeTab === 'settings'}
			<!-- Basic Settings Tab -->
			<div class="space-y-6">
				<!-- Basic Information -->
				<div class="bg-white shadow rounded-lg p-6">
					<h2 class="text-lg font-medium text-gray-900 mb-4">Basic Information</h2>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
						<div>
							<label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name</label>
							<input
								id="name"
								type="text"
								bind:value={editingName}
								class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
							/>
						</div>
						<div>
							<label for="key" class="block text-sm font-medium text-gray-700 mb-1">Key (read-only)</label>
							<input
								id="key"
								type="text"
								value={flag.key}
								disabled
								class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm bg-gray-50 text-gray-500"
							/>
						</div>
					</div>
					<div class="mt-4">
						<label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
						<textarea
							id="description"
							bind:value={editingDescription}
							rows="3"
							class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>
					<div class="mt-4">
						<label for="tags" class="block text-sm font-medium text-gray-700 mb-1">Tags</label>
						<input
							id="tags"
							type="text"
							bind:value={editingTags}
							placeholder="comma, separated, tags"
							class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>
					<div class="mt-6 flex justify-end">
						<button
							onclick={saveBasicSettings}
							disabled={isSaving}
							class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
						>
							{#if isSaving}
								Saving...
							{:else}
								Save Changes
							{/if}
						</button>
					</div>
				</div>

				<!-- Default Value -->
				<div class="bg-white shadow rounded-lg p-6">
					<h2 class="text-lg font-medium text-gray-900 mb-4">Default Value</h2>
					<div class="space-y-4">
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Type</label>
							<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
								{flag.type === 'boolean' ? 'bg-blue-100 text-blue-800' :
								flag.type === 'string' ? 'bg-green-100 text-green-800' :
								flag.type === 'number' ? 'bg-purple-100 text-purple-800' :
								'bg-yellow-100 text-yellow-800'}">
								{flag.type}
							</span>
						</div>
						<div>
							<label class="block text-sm font-medium text-gray-700 mb-1">Value</label>
							<div class="mt-1 p-3 border border-gray-300 rounded-md bg-gray-50 text-sm text-gray-900 font-mono">
								{formatValue(flag.default_value)}
							</div>
							<p class="text-xs text-gray-500 mt-1">
								Returned when the flag is disabled or no targeting rules match
							</p>
						</div>
					</div>
				</div>
			</div>
		{:else if activeTab === 'environments'}
			<!-- Environments Tab -->
			<div class="space-y-6">
				{#each environments as env}
					<div class="bg-white shadow rounded-lg p-6">
						<div class="flex items-center space-x-3 mb-6">
							<div class="w-4 h-4 rounded-full" style="background-color: {env.color}"></div>
							<h2 class="text-lg font-medium text-gray-900">{env.name}</h2>
							<span class="text-sm text-gray-500">({env.key})</span>
						</div>

						<div class="space-y-6">
							<!-- Environment Toggle -->
							<div class="flex items-center justify-between py-4 border-b border-gray-200">
								<div>
									<h3 class="text-sm font-medium text-gray-900">Environment Status</h3>
									<p class="text-sm text-gray-500">Enable or disable this flag in {env.name}</p>
								</div>
								<FlagToggle {flag} environment={env.key} onToggle={handleToggle} />
							</div>

							{#if flag.environments[env.key]?.enabled}
								<!-- Current Value -->
								<div>
									<h3 class="text-sm font-medium text-gray-900 mb-2">Current Value</h3>
									<div class="p-3 border border-gray-300 rounded-md bg-gray-50 text-sm text-gray-900 font-mono">
										{formatValue(flag.environments[env.key]?.value ?? flag.default_value)}
									</div>
								</div>

								<!-- Rollout Percentage -->
								<RolloutSlider {flag} environment={env.key} onUpdate={handleUpdate} />

								<!-- Targeting Rules -->
								<div>
									<TargetingRuleBuilder {flag} environment={env.key} onUpdate={handleUpdate} />
								</div>
							{:else}
								<div class="text-center py-8 text-gray-500">
									<svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
										<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L5.636 5.636" />
									</svg>
									<p class="mt-2 text-sm">
										Flag is disabled in this environment. Enable it to configure rollout and targeting.
									</p>
								</div>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{/if}
	{/if}
</div>