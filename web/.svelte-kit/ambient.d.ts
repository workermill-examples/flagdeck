
// this file is generated — do not edit it


/// <reference types="@sveltejs/kit" />

/**
 * This module provides access to environment variables that are injected _statically_ into your bundle at build time and are limited to _private_ access.
 * 
 * |         | Runtime                                                                    | Build time                                                               |
 * | ------- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
 * | Private | [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private) | [`$env/static/private`](https://svelte.dev/docs/kit/$env-static-private) |
 * | Public  | [`$env/dynamic/public`](https://svelte.dev/docs/kit/$env-dynamic-public)   | [`$env/static/public`](https://svelte.dev/docs/kit/$env-static-public)   |
 * 
 * Static environment variables are [loaded by Vite](https://vitejs.dev/guide/env-and-mode.html#env-files) from `.env` files and `process.env` at build time and then statically injected into your bundle at build time, enabling optimisations like dead code elimination.
 * 
 * **_Private_ access:**
 * 
 * - This module cannot be imported into client-side code
 * - This module only includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured)
 * 
 * For example, given the following build time environment:
 * 
 * ```env
 * ENVIRONMENT=production
 * PUBLIC_BASE_URL=http://site.com
 * ```
 * 
 * With the default `publicPrefix` and `privatePrefix`:
 * 
 * ```ts
 * import { ENVIRONMENT, PUBLIC_BASE_URL } from '$env/static/private';
 * 
 * console.log(ENVIRONMENT); // => "production"
 * console.log(PUBLIC_BASE_URL); // => throws error during build
 * ```
 * 
 * The above values will be the same _even if_ different values for `ENVIRONMENT` or `PUBLIC_BASE_URL` are set at runtime, as they are statically replaced in your code with their build time values.
 */
declare module '$env/static/private' {
	export const GITHUB_TOKEN: string;
	export const EXECUTION_MODE: string;
	export const WORKER_PERSONA: string;
	export const CLAUDE_CODE_ENTRYPOINT: string;
	export const npm_config_user_agent: string;
	export const TARGET_FILES: string;
	export const TASK_SUMMARY: string;
	export const GITHUB_REPO: string;
	export const GIT_EDITOR: string;
	export const NODE_VERSION: string;
	export const HOSTNAME: string;
	export const GH_TOKEN: string;
	export const YARN_VERSION: string;
	export const PUSH_AFTER_COMMIT: string;
	export const SELF_REVIEW_ENABLED: string;
	export const npm_node_execpath: string;
	export const SHLVL: string;
	export const npm_config_noproxy: string;
	export const HOME: string;
	export const WORKER_PROVIDER: string;
	export const OLDPWD: string;
	export const SCM_TOKEN: string;
	export const SCM_PROVIDER: string;
	export const npm_package_json: string;
	export const TASK_DESCRIPTION: string;
	export const RETRY_NUMBER: string;
	export const MAX_CI_FIX_RETRIES: string;
	export const GITHUB_REVIEWER_TOKEN: string;
	export const BITBUCKET_USERNAME: string;
	export const npm_config_userconfig: string;
	export const npm_config_local_prefix: string;
	export const AUTHOR_EMAIL: string;
	export const COLOR: string;
	export const EXECUTION_MODE_SETTING: string;
	export const BLOCKER_MAX_AUTO_RETRIES: string;
	export const REFERENCE_FILES: string;
	export const TASK_ID: string;
	export const MAX_REVIEW_REVISIONS: string;
	export const ORG_ID: string;
	export const CLAUDE_MODEL: string;
	export const _: string;
	export const npm_config_prefix: string;
	export const npm_config_npm_version: string;
	export const STANDARD_SDK_MODE: string;
	export const GOMODCACHE: string;
	export const TICKET_KEY: string;
	export const OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE: string;
	export const npm_config_cache: string;
	export const MANAGER_MODEL: string;
	export const TARGET_REPO: string;
	export const MAX_PER_STORY_REVISIONS: string;
	export const npm_config_node_gyp: string;
	export const PATH: string;
	export const QUALITY_THRESHOLDS: string;
	export const NODE: string;
	export const npm_package_name: string;
	export const COREPACK_ENABLE_AUTO_PIN: string;
	export const PERSONA: string;
	export const MAIN_BRANCH: string;
	export const PRD_CHILD_TASK: string;
	export const BLOCKER_AUTO_RETRY_ENABLED: string;
	export const REVIEW_ENABLED: string;
	export const NoDefaultCurrentDirectoryInExePath: string;
	export const API_BASE_URL: string;
	export const DEPLOYMENT_ENABLED: string;
	export const JIRA_SUMMARY: string;
	export const QUALITY_GATE_BYPASS: string;
	export const npm_lifecycle_script: string;
	export const EPIC_MODE: string;
	export const SHELL: string;
	export const PIPELINE_VERSION: string;
	export const OLLAMA_CONTEXT_WINDOW: string;
	export const MAX_PARALLEL_EXPERTS: string;
	export const BLOCKER_WAIT_TIMEOUT_MINUTES: string;
	export const GOPATH: string;
	export const npm_package_version: string;
	export const npm_lifecycle_event: string;
	export const REPO_PATH: string;
	export const WORKER_MODEL: string;
	export const JIRA_DESCRIPTION: string;
	export const CLAUDECODE: string;
	export const CODEBASE_INDEXING_ENABLED: string;
	export const ORG_API_KEY: string;
	export const STORY_ID: string;
	export const JIRA_ISSUE_KEY: string;
	export const npm_config_globalconfig: string;
	export const npm_config_init_module: string;
	export const PWD: string;
	export const npm_execpath: string;
	export const IMPROVEMENT_ENABLED: string;
	export const npm_config_global_prefix: string;
	export const GRACEFUL_SHUTDOWN_ENABLED: string;
	export const MANAGER_PROVIDER: string;
	export const npm_command: string;
	export const NVIDIA_VISIBLE_DEVICES: string;
	export const TICKET_SYSTEM: string;
	export const QUALITY_GATE_MAX_RETRIES: string;
	export const PARENT_TASK_ID: string;
	export const INIT_CWD: string;
	export const EDITOR: string;
}

/**
 * This module provides access to environment variables that are injected _statically_ into your bundle at build time and are _publicly_ accessible.
 * 
 * |         | Runtime                                                                    | Build time                                                               |
 * | ------- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
 * | Private | [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private) | [`$env/static/private`](https://svelte.dev/docs/kit/$env-static-private) |
 * | Public  | [`$env/dynamic/public`](https://svelte.dev/docs/kit/$env-dynamic-public)   | [`$env/static/public`](https://svelte.dev/docs/kit/$env-static-public)   |
 * 
 * Static environment variables are [loaded by Vite](https://vitejs.dev/guide/env-and-mode.html#env-files) from `.env` files and `process.env` at build time and then statically injected into your bundle at build time, enabling optimisations like dead code elimination.
 * 
 * **_Public_ access:**
 * 
 * - This module _can_ be imported into client-side code
 * - **Only** variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`) are included
 * 
 * For example, given the following build time environment:
 * 
 * ```env
 * ENVIRONMENT=production
 * PUBLIC_BASE_URL=http://site.com
 * ```
 * 
 * With the default `publicPrefix` and `privatePrefix`:
 * 
 * ```ts
 * import { ENVIRONMENT, PUBLIC_BASE_URL } from '$env/static/public';
 * 
 * console.log(ENVIRONMENT); // => throws error during build
 * console.log(PUBLIC_BASE_URL); // => "http://site.com"
 * ```
 * 
 * The above values will be the same _even if_ different values for `ENVIRONMENT` or `PUBLIC_BASE_URL` are set at runtime, as they are statically replaced in your code with their build time values.
 */
declare module '$env/static/public' {
	
}

/**
 * This module provides access to environment variables set _dynamically_ at runtime and that are limited to _private_ access.
 * 
 * |         | Runtime                                                                    | Build time                                                               |
 * | ------- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
 * | Private | [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private) | [`$env/static/private`](https://svelte.dev/docs/kit/$env-static-private) |
 * | Public  | [`$env/dynamic/public`](https://svelte.dev/docs/kit/$env-dynamic-public)   | [`$env/static/public`](https://svelte.dev/docs/kit/$env-static-public)   |
 * 
 * Dynamic environment variables are defined by the platform you're running on. For example if you're using [`adapter-node`](https://github.com/sveltejs/kit/tree/main/packages/adapter-node) (or running [`vite preview`](https://svelte.dev/docs/kit/cli)), this is equivalent to `process.env`.
 * 
 * **_Private_ access:**
 * 
 * - This module cannot be imported into client-side code
 * - This module includes variables that _do not_ begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) _and do_ start with [`config.kit.env.privatePrefix`](https://svelte.dev/docs/kit/configuration#env) (if configured)
 * 
 * > [!NOTE] In `dev`, `$env/dynamic` includes environment variables from `.env`. In `prod`, this behavior will depend on your adapter.
 * 
 * > [!NOTE] To get correct types, environment variables referenced in your code should be declared (for example in an `.env` file), even if they don't have a value until the app is deployed:
 * >
 * > ```env
 * > MY_FEATURE_FLAG=
 * > ```
 * >
 * > You can override `.env` values from the command line like so:
 * >
 * > ```sh
 * > MY_FEATURE_FLAG="enabled" npm run dev
 * > ```
 * 
 * For example, given the following runtime environment:
 * 
 * ```env
 * ENVIRONMENT=production
 * PUBLIC_BASE_URL=http://site.com
 * ```
 * 
 * With the default `publicPrefix` and `privatePrefix`:
 * 
 * ```ts
 * import { env } from '$env/dynamic/private';
 * 
 * console.log(env.ENVIRONMENT); // => "production"
 * console.log(env.PUBLIC_BASE_URL); // => undefined
 * ```
 */
declare module '$env/dynamic/private' {
	export const env: {
		GITHUB_TOKEN: string;
		EXECUTION_MODE: string;
		WORKER_PERSONA: string;
		CLAUDE_CODE_ENTRYPOINT: string;
		npm_config_user_agent: string;
		TARGET_FILES: string;
		TASK_SUMMARY: string;
		GITHUB_REPO: string;
		GIT_EDITOR: string;
		NODE_VERSION: string;
		HOSTNAME: string;
		GH_TOKEN: string;
		YARN_VERSION: string;
		PUSH_AFTER_COMMIT: string;
		SELF_REVIEW_ENABLED: string;
		npm_node_execpath: string;
		SHLVL: string;
		npm_config_noproxy: string;
		HOME: string;
		WORKER_PROVIDER: string;
		OLDPWD: string;
		SCM_TOKEN: string;
		SCM_PROVIDER: string;
		npm_package_json: string;
		TASK_DESCRIPTION: string;
		RETRY_NUMBER: string;
		MAX_CI_FIX_RETRIES: string;
		GITHUB_REVIEWER_TOKEN: string;
		BITBUCKET_USERNAME: string;
		npm_config_userconfig: string;
		npm_config_local_prefix: string;
		AUTHOR_EMAIL: string;
		COLOR: string;
		EXECUTION_MODE_SETTING: string;
		BLOCKER_MAX_AUTO_RETRIES: string;
		REFERENCE_FILES: string;
		TASK_ID: string;
		MAX_REVIEW_REVISIONS: string;
		ORG_ID: string;
		CLAUDE_MODEL: string;
		_: string;
		npm_config_prefix: string;
		npm_config_npm_version: string;
		STANDARD_SDK_MODE: string;
		GOMODCACHE: string;
		TICKET_KEY: string;
		OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE: string;
		npm_config_cache: string;
		MANAGER_MODEL: string;
		TARGET_REPO: string;
		MAX_PER_STORY_REVISIONS: string;
		npm_config_node_gyp: string;
		PATH: string;
		QUALITY_THRESHOLDS: string;
		NODE: string;
		npm_package_name: string;
		COREPACK_ENABLE_AUTO_PIN: string;
		PERSONA: string;
		MAIN_BRANCH: string;
		PRD_CHILD_TASK: string;
		BLOCKER_AUTO_RETRY_ENABLED: string;
		REVIEW_ENABLED: string;
		NoDefaultCurrentDirectoryInExePath: string;
		API_BASE_URL: string;
		DEPLOYMENT_ENABLED: string;
		JIRA_SUMMARY: string;
		QUALITY_GATE_BYPASS: string;
		npm_lifecycle_script: string;
		EPIC_MODE: string;
		SHELL: string;
		PIPELINE_VERSION: string;
		OLLAMA_CONTEXT_WINDOW: string;
		MAX_PARALLEL_EXPERTS: string;
		BLOCKER_WAIT_TIMEOUT_MINUTES: string;
		GOPATH: string;
		npm_package_version: string;
		npm_lifecycle_event: string;
		REPO_PATH: string;
		WORKER_MODEL: string;
		JIRA_DESCRIPTION: string;
		CLAUDECODE: string;
		CODEBASE_INDEXING_ENABLED: string;
		ORG_API_KEY: string;
		STORY_ID: string;
		JIRA_ISSUE_KEY: string;
		npm_config_globalconfig: string;
		npm_config_init_module: string;
		PWD: string;
		npm_execpath: string;
		IMPROVEMENT_ENABLED: string;
		npm_config_global_prefix: string;
		GRACEFUL_SHUTDOWN_ENABLED: string;
		MANAGER_PROVIDER: string;
		npm_command: string;
		NVIDIA_VISIBLE_DEVICES: string;
		TICKET_SYSTEM: string;
		QUALITY_GATE_MAX_RETRIES: string;
		PARENT_TASK_ID: string;
		INIT_CWD: string;
		EDITOR: string;
		[key: `PUBLIC_${string}`]: undefined;
		[key: `${string}`]: string | undefined;
	}
}

/**
 * This module provides access to environment variables set _dynamically_ at runtime and that are _publicly_ accessible.
 * 
 * |         | Runtime                                                                    | Build time                                                               |
 * | ------- | -------------------------------------------------------------------------- | ------------------------------------------------------------------------ |
 * | Private | [`$env/dynamic/private`](https://svelte.dev/docs/kit/$env-dynamic-private) | [`$env/static/private`](https://svelte.dev/docs/kit/$env-static-private) |
 * | Public  | [`$env/dynamic/public`](https://svelte.dev/docs/kit/$env-dynamic-public)   | [`$env/static/public`](https://svelte.dev/docs/kit/$env-static-public)   |
 * 
 * Dynamic environment variables are defined by the platform you're running on. For example if you're using [`adapter-node`](https://github.com/sveltejs/kit/tree/main/packages/adapter-node) (or running [`vite preview`](https://svelte.dev/docs/kit/cli)), this is equivalent to `process.env`.
 * 
 * **_Public_ access:**
 * 
 * - This module _can_ be imported into client-side code
 * - **Only** variables that begin with [`config.kit.env.publicPrefix`](https://svelte.dev/docs/kit/configuration#env) (which defaults to `PUBLIC_`) are included
 * 
 * > [!NOTE] In `dev`, `$env/dynamic` includes environment variables from `.env`. In `prod`, this behavior will depend on your adapter.
 * 
 * > [!NOTE] To get correct types, environment variables referenced in your code should be declared (for example in an `.env` file), even if they don't have a value until the app is deployed:
 * >
 * > ```env
 * > MY_FEATURE_FLAG=
 * > ```
 * >
 * > You can override `.env` values from the command line like so:
 * >
 * > ```sh
 * > MY_FEATURE_FLAG="enabled" npm run dev
 * > ```
 * 
 * For example, given the following runtime environment:
 * 
 * ```env
 * ENVIRONMENT=production
 * PUBLIC_BASE_URL=http://example.com
 * ```
 * 
 * With the default `publicPrefix` and `privatePrefix`:
 * 
 * ```ts
 * import { env } from '$env/dynamic/public';
 * console.log(env.ENVIRONMENT); // => undefined, not public
 * console.log(env.PUBLIC_BASE_URL); // => "http://example.com"
 * ```
 * 
 * ```
 * 
 * ```
 */
declare module '$env/dynamic/public' {
	export const env: {
		[key: `PUBLIC_${string}`]: string | undefined;
	}
}
