<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '../../../lib/api.js';
	import type { Environment } from '../../../lib/types.js';

	let environments = $state<Environment[]>([]);
	let isLoading = $state(true);
	let isSubmitting = $state(false);
	let error = $state<string | null>(null);

	// Form fields
	let name = $state('');
	let description = $state('');
	let type = $state<'boolean' | 'string' | 'number' | 'json'>('boolean');
	let defaultValue = $state<string>('false');
	let tags = $state('');

	// Derived key from name
	const key = $derived(() => {
		return name
			.toLowerCase()
			.replace(/[^a-z0-9\s-]/g, '') // Remove special chars except spaces and hyphens
			.replace(/\s+/g, '-') // Replace spaces with hyphens
			.replace(/-+/g, '-') // Replace multiple hyphens with single
			.replace(/^-|-$/g, ''); // Remove leading/trailing hyphens
	});

	onMount(async () => {
		await loadEnvironments();
	});

	async function loadEnvironments() {
		try {
			isLoading = true;
			error = null;

			const response = await api.getEnvironments();
			environments = response.data.sort((a, b) => a.sort_order - b.sort_order);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load environments';
			console.error('Failed to load environments:', err);
		} finally {
			isLoading = false;
		}
	}

	function updateTypeDefaultValue() {
		switch (type) {
			case 'boolean':
				defaultValue = 'false';
				break;
			case 'string':
				defaultValue = '';
				break;
			case 'number':
				defaultValue = '0';
				break;
			case 'json':
				defaultValue = '{}';
				break;
		}
	}

	function parseDefaultValue() {
		switch (type) {
			case 'boolean':
				return defaultValue.toLowerCase() === 'true';
			case 'number': {
				const num = parseFloat(defaultValue);
				return isNaN(num) ? 0 : num;
			}
			case 'json':
				try {
					return JSON.parse(defaultValue);
				} catch {
					return {};
				}
			default:
				return defaultValue;
		}
	}

	function parseTags(): string[] {
		return tags
			.split(',')
			.map(tag => tag.trim())
			.filter(tag => tag.length > 0);
	}

	async function handleSubmit(event: SubmitEvent) {
		event.preventDefault();
		if (isSubmitting) return;

		// Basic validation
		if (!name.trim()) {
			error = 'Flag name is required';
			return;
		}

		if (!key) {
			error = 'Flag name must contain at least one alphanumeric character';
			return;
		}

		try {
			isSubmitting = true;
			error = null;

			const flagData = {
				key,
				name: name.trim(),
				description: description.trim(),
				type,
				default_value: parseDefaultValue(),
				tags: parseTags(),
				is_active: true,
				environments: environments.reduce((acc, env) => {
					acc[env.key] = {
						enabled: true,
						value: parseDefaultValue(),
						rollout_percent: 100,
						targeting_rules: []
					};
					return acc;
				}, {} as Record<string, any>)
			};

			await api.createFlag(flagData);

			// Redirect to the new flag's detail page
			window.location.href = `/flags/${key}`;
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create flag';
			console.error('Failed to create flag:', err);
		} finally {
			isSubmitting = false;
		}
	}

	function handleCancel() {
		window.location.href = '/flags';
	}
</script>

<svelte:head>
	<title>Create Feature Flag - FlagDeck</title>
</svelte:head>

<div class="max-w-2xl mx-auto space-y-6">
	<!-- Header -->
	<div>
		<div class="flex items-center space-x-2 text-sm text-gray-500 mb-2">
			<a href="/flags" class="hover:text-blue-600">Feature Flags</a>
			<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
				<path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"></path>
			</svg>
			<span>Create</span>
		</div>
		<h1 class="text-2xl font-bold text-gray-900">Create Feature Flag</h1>
		<p class="text-gray-600 mt-1">
			Configure a new feature flag with default settings for all environments
		</p>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
		</div>
	{:else}
		<form onsubmit={handleSubmit} class="space-y-6">
			<div class="bg-white shadow rounded-lg p-6 space-y-6">
				<!-- Basic Info -->
				<div class="space-y-4">
					<h2 class="text-lg font-medium text-gray-900">Basic Information</h2>

					<!-- Name -->
					<div>
						<label for="name" class="block text-sm font-medium text-gray-700 mb-1">
							Name <span class="text-red-500">*</span>
						</label>
						<input
							id="name"
							type="text"
							bind:value={name}
							placeholder="e.g., Dark Mode Toggle"
							required
							class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
						/>
					</div>

					<!-- Auto-generated Key -->
					<div>
						<label for="key" class="block text-sm font-medium text-gray-700 mb-1">
							Key (auto-generated)
						</label>
						<input
							id="key"
							type="text"
							value={key}
							disabled
							placeholder="Will be generated from name"
							class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm bg-gray-50 text-gray-500"
						/>
						<p class="text-xs text-gray-500 mt-1">
							Used to reference this flag in your application code
						</p>
					</div>

					<!-- Description -->
					<div>
						<label for="description" class="block text-sm font-medium text-gray-700 mb-1">
							Description
						</label>
						<textarea
							id="description"
							bind:value={description}
							rows="3"
							placeholder="Briefly describe what this flag controls..."
							class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
						></textarea>
					</div>

					<!-- Tags -->
					<div>
						<label for="tags" class="block text-sm font-medium text-gray-700 mb-1">
							Tags
						</label>
						<input
							id="tags"
							type="text"
							bind:value={tags}
							placeholder="e.g., ui, feature, beta (comma-separated)"
							class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
						/>
						<p class="text-xs text-gray-500 mt-1">
							Separate multiple tags with commas
						</p>
					</div>
				</div>

				<!-- Type and Default Value -->
				<div class="space-y-4">
					<h2 class="text-lg font-medium text-gray-900">Value Configuration</h2>

					<!-- Type -->
					<div>
						<label for="type" class="block text-sm font-medium text-gray-700 mb-1">
							Type <span class="text-red-500">*</span>
						</label>
						<select
							id="type"
							bind:value={type}
							onchange={updateTypeDefaultValue}
							class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
						>
							<option value="boolean">Boolean (true/false)</option>
							<option value="string">String (text)</option>
							<option value="number">Number</option>
							<option value="json">JSON (object/array)</option>
						</select>
					</div>

					<!-- Default Value -->
					<div>
						<label for="defaultValue" class="block text-sm font-medium text-gray-700 mb-1">
							Default Value <span class="text-red-500">*</span>
						</label>
						{#if type === 'boolean'}
							<select
								id="defaultValue"
								bind:value={defaultValue}
								class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
							>
								<option value="false">False</option>
								<option value="true">True</option>
							</select>
						{:else if type === 'number'}
							<input
								id="defaultValue"
								type="number"
								step="any"
								bind:value={defaultValue}
								placeholder="0"
								class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
							/>
						{:else if type === 'json'}
							<textarea
								id="defaultValue"
								bind:value={defaultValue}
								rows="4"
								placeholder="&#123;&#125;"
								class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500 font-mono"
							></textarea>
							<p class="text-xs text-gray-500 mt-1">
								Enter valid JSON (object or array)
							</p>
						{:else}
							<input
								id="defaultValue"
								type="text"
								bind:value={defaultValue}
								placeholder="Enter string value"
								class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
							/>
						{/if}
						<p class="text-xs text-gray-500 mt-1">
							This value will be used when the flag is disabled or no targeting rules match
						</p>
					</div>
				</div>

				<!-- Environment Preview -->
				<div class="space-y-4">
					<h2 class="text-lg font-medium text-gray-900">Initial Environment Settings</h2>
					<p class="text-sm text-gray-600">
						The flag will be created enabled in all environments with 100% rollout. You can adjust individual environment settings after creation.
					</p>

					<div class="bg-gray-50 rounded-lg p-4">
						<div class="grid grid-cols-1 md:grid-cols-3 gap-4">
							{#each environments as env}
								<div class="flex items-center space-x-3">
									<div class="w-3 h-3 rounded-full" style="background-color: {env.color}"></div>
									<span class="text-sm font-medium">{env.name}</span>
									<span class="text-xs text-green-600 bg-green-100 px-2 py-1 rounded">Enabled (100%)</span>
								</div>
							{/each}
						</div>
					</div>
				</div>
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

			<!-- Actions -->
			<div class="flex justify-end space-x-3">
				<button
					type="button"
					onclick={handleCancel}
					disabled={isSubmitting}
					class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					Cancel
				</button>
				<button
					type="submit"
					disabled={isSubmitting || !name.trim() || !key}
					class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
				>
					{#if isSubmitting}
						Creating...
					{:else}
						Create Flag
					{/if}
				</button>
			</div>
		</form>
	{/if}
</div>