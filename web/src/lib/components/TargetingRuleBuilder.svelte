<script lang="ts">
	import { api } from '../api.js';
	import type { Flag, TargetingRule, Condition } from '../types.js';

	interface Props {
		flag: Flag;
		environment: string;
		onUpdate?: () => void;
	}

	let { flag, environment, onUpdate }: Props = $props();

	let isLoading = $state(false);
	let rules = $state<TargetingRule[]>([]);

	const operatorOptions = [
		{ value: 'equals', label: 'Equals' },
		{ value: 'not_equals', label: 'Not Equals' },
		{ value: 'contains', label: 'Contains' },
		{ value: 'in', label: 'In' },
		{ value: 'not_in', label: 'Not In' },
		{ value: 'gt', label: 'Greater Than' },
		{ value: 'lt', label: 'Less Than' },
		{ value: 'gte', label: 'Greater Than or Equal' },
		{ value: 'lte', label: 'Less Than or Equal' },
		{ value: 'regex', label: 'Regex' }
	];

	// Initialize rules from flag data
	$effect(() => {
		rules = flag.environments[environment]?.targeting_rules ?? [];
	});

	function generateRuleId() {
		return `rule_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
	}

	function addRule() {
		const newRule: TargetingRule = {
			id: generateRuleId(),
			priority: rules.length + 1,
			conditions: [
				{
					property: '',
					operator: 'equals',
					value: ''
				}
			],
			value: flag.default_value
		};
		rules = [...rules, newRule];
	}

	function removeRule(ruleId: string) {
		rules = rules.filter(rule => rule.id !== ruleId);
		// Reorder priorities
		rules = rules.map((rule, index) => ({ ...rule, priority: index + 1 }));
	}

	function addCondition(ruleId: string) {
		rules = rules.map(rule => {
			if (rule.id === ruleId) {
				return {
					...rule,
					conditions: [
						...rule.conditions,
						{
							property: '',
							operator: 'equals',
							value: ''
						}
					]
				};
			}
			return rule;
		});
	}

	function removeCondition(ruleId: string, conditionIndex: number) {
		rules = rules.map(rule => {
			if (rule.id === ruleId) {
				return {
					...rule,
					conditions: rule.conditions.filter((_, index) => index !== conditionIndex)
				};
			}
			return rule;
		});
	}

	function updateCondition(ruleId: string, conditionIndex: number, field: keyof Condition, value: any) {
		rules = rules.map(rule => {
			if (rule.id === ruleId) {
				return {
					...rule,
					conditions: rule.conditions.map((condition, index) => {
						if (index === conditionIndex) {
							return { ...condition, [field]: value };
						}
						return condition;
					})
				};
			}
			return rule;
		});
	}

	function updateRuleValue(ruleId: string, value: any) {
		rules = rules.map(rule => {
			if (rule.id === ruleId) {
				return { ...rule, value };
			}
			return rule;
		});
	}

	function moveRule(ruleId: string, direction: 'up' | 'down') {
		const ruleIndex = rules.findIndex(rule => rule.id === ruleId);
		if (ruleIndex === -1) return;

		const newIndex = direction === 'up' ? ruleIndex - 1 : ruleIndex + 1;
		if (newIndex < 0 || newIndex >= rules.length) return;

		const newRules = [...rules];
		[newRules[ruleIndex], newRules[newIndex]] = [newRules[newIndex], newRules[ruleIndex]];

		// Update priorities
		newRules.forEach((rule, index) => {
			rule.priority = index + 1;
		});

		rules = newRules;
	}

	async function saveRules() {
		if (isLoading) return;

		try {
			isLoading = true;

			const updatedFlag = {
				...flag,
				environments: {
					...flag.environments,
					[environment]: {
						...flag.environments[environment],
						targeting_rules: rules
					}
				}
			};

			await api.updateFlag(flag.key, updatedFlag);
			onUpdate?.();
		} catch (error) {
			console.error('Failed to update targeting rules:', error);
		} finally {
			isLoading = false;
		}
	}

	function formatValueInput(value: any, type: string) {
		if (type === 'json') {
			try {
				return JSON.stringify(value, null, 2);
			} catch {
				return String(value);
			}
		}
		return String(value);
	}

	function parseValueInput(value: string, type: string) {
		if (type === 'json') {
			try {
				return JSON.parse(value);
			} catch {
				return value;
			}
		}
		if (type === 'number') {
			const num = parseFloat(value);
			return isNaN(num) ? value : num;
		}
		if (type === 'boolean') {
			return value.toLowerCase() === 'true';
		}
		return value;
	}
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between">
		<h3 class="text-lg font-medium text-gray-900">Targeting Rules</h3>
		<button
			type="button"
			onclick={addRule}
			class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
		>
			Add Rule
		</button>
	</div>

	{#if rules.length === 0}
		<div class="text-center py-6 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
			<p class="text-sm text-gray-500">
				No targeting rules defined. Flag will use default rollout percentage for all users.
			</p>
		</div>
	{:else}
		<div class="space-y-4">
			{#each rules as rule, ruleIndex (rule.id)}
				<div class="border border-gray-200 rounded-lg p-4 bg-white">
					<div class="flex items-center justify-between mb-4">
						<div class="flex items-center space-x-2">
							<span class="text-sm font-medium text-gray-700">Rule {rule.priority}</span>
							<div class="flex space-x-1">
								<button
									type="button"
									onclick={() => moveRule(rule.id, 'up')}
									disabled={ruleIndex === 0}
									aria-label="Move rule up"
									class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd"></path>
									</svg>
								</button>
								<button
									type="button"
									onclick={() => moveRule(rule.id, 'down')}
									disabled={ruleIndex === rules.length - 1}
									aria-label="Move rule down"
									class="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-50 disabled:cursor-not-allowed"
								>
									<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd"></path>
									</svg>
								</button>
							</div>
						</div>
						<button
							type="button"
							onclick={() => removeRule(rule.id)}
							class="text-red-600 hover:text-red-800 text-sm"
						>
							Remove
						</button>
					</div>

					<div class="space-y-3">
						<div>
							<h4 class="block text-sm font-medium text-gray-700 mb-2">Conditions (ALL must match)</h4>
							<div class="space-y-2">
								{#each rule.conditions as condition, conditionIndex}
									<div class="flex items-center space-x-2">
										<input
											type="text"
											placeholder="Property"
											value={condition.property}
											oninput={(e) => updateCondition(rule.id, conditionIndex, 'property', e.target.value)}
											class="flex-1 min-w-0 block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
										/>
										<select
											value={condition.operator}
											onchange={(e) => updateCondition(rule.id, conditionIndex, 'operator', e.target.value)}
											class="block w-32 px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
										>
											{#each operatorOptions as option}
												<option value={option.value}>{option.label}</option>
											{/each}
										</select>
										<input
											type="text"
											placeholder="Value"
											value={String(condition.value)}
											oninput={(e) => updateCondition(rule.id, conditionIndex, 'value', e.target.value)}
											class="flex-1 min-w-0 block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
										/>
										{#if rule.conditions.length > 1}
											<button
												type="button"
												onclick={() => removeCondition(rule.id, conditionIndex)}
												aria-label="Remove condition"
												class="p-2 text-red-600 hover:text-red-800"
											>
												<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
													<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"></path>
												</svg>
											</button>
										{/if}
									</div>
								{/each}
								<button
									type="button"
									onclick={() => addCondition(rule.id)}
									class="text-sm text-blue-600 hover:text-blue-800"
								>
									+ Add Condition
								</button>
							</div>
						</div>

						<div>
							<h4 class="block text-sm font-medium text-gray-700 mb-2">Return Value</h4>
							{#if flag.type === 'boolean'}
								<select
									value={rule.value}
									onchange={(e) => updateRuleValue(rule.id, e.target.value === 'true')}
									class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
								>
									<option value="true">True</option>
									<option value="false">False</option>
								</select>
							{:else if flag.type === 'number'}
								<input
									type="number"
									value={rule.value}
									oninput={(e) => updateRuleValue(rule.id, parseFloat(e.target.value))}
									class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
								/>
							{:else if flag.type === 'json'}
								<textarea
									value={formatValueInput(rule.value, flag.type)}
									oninput={(e) => updateRuleValue(rule.id, parseValueInput(e.target.value, flag.type))}
									rows="3"
									class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
								></textarea>
							{:else}
								<input
									type="text"
									value={rule.value}
									oninput={(e) => updateRuleValue(rule.id, e.target.value)}
									class="block w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
								/>
							{/if}
						</div>
					</div>
				</div>
			{/each}
		</div>
	{/if}

	<div class="flex justify-end">
		<button
			type="button"
			onclick={saveRules}
			disabled={isLoading}
			class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
		>
			{#if isLoading}
				Saving...
			{:else}
				Save Rules
			{/if}
		</button>
	</div>
</div>