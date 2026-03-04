<script lang="ts">
	import { api } from '../api.js';
	import type { Flag } from '../types.js';

	interface Props {
		flag: Flag;
		environment: string;
		onToggle?: () => void;
	}

	let { flag, environment, onToggle }: Props = $props();

	let isLoading = $state(false);

	async function handleToggle() {
		if (isLoading) return;

		try {
			isLoading = true;
			await api.toggleFlag(flag.key, environment);
			onToggle?.();
		} catch (error) {
			console.error('Failed to toggle flag:', error);
		} finally {
			isLoading = false;
		}
	}

	const isEnabled = $derived(flag.environments[environment]?.enabled ?? false);
</script>

<button
	type="button"
	class="relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 {isEnabled
		? 'bg-blue-600'
		: 'bg-gray-200'} {isLoading ? 'opacity-50 cursor-not-allowed' : ''}"
	onclick={handleToggle}
	disabled={isLoading}
	aria-label="Toggle flag for {environment}"
>
	<span
		class="pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow transition duration-200 ease-in-out {isEnabled
			? 'translate-x-5'
			: 'translate-x-0'}"
	></span>
</button>