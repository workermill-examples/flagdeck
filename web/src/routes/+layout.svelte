<script lang="ts">
import { onMount } from "svelte";
import { goto } from "$app/navigation";
import { page } from "$app/stores";
import { auth } from "$lib/auth.js";
import Sidebar from "$lib/components/Sidebar.svelte";
import "../../app.css";

interface Props {
  children: import('svelte').Snippet;
}

let { children }: Props = $props();

// Check if we're on the login page
let isLoginPage = $derived($page.url.pathname === "/login");

// Auth guard - redirect to login if not authenticated
onMount(async () => {
  await auth.checkAuth();

  if (!auth.isAuthenticated && !isLoginPage) {
    goto("/login", { replaceState: true });
  }
});

// Reactive redirect when auth state changes
$effect(() => {
  if (!auth.isLoading && !auth.isAuthenticated && !isLoginPage) {
    goto("/login", { replaceState: true });
  }
});
</script>

{#if auth.isLoading}
  <!-- Loading state while checking authentication -->
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <div class="w-8 h-8 mx-auto mb-4 border-4 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
      <p class="text-gray-600">Loading...</p>
    </div>
  </div>
{:else if isLoginPage}
  <!-- Login page - no sidebar -->
  <main class="min-h-screen bg-gray-50">
    {@render children()}
  </main>
{:else if auth.isAuthenticated}
  <!-- Authenticated layout with sidebar -->
  <div class="min-h-screen bg-gray-50 flex">
    <Sidebar />
    <main class="flex-1 overflow-hidden">
      {@render children()}
    </main>
  </div>
{:else}
  <!-- Fallback - should not normally be seen -->
  <div class="min-h-screen flex items-center justify-center bg-gray-50">
    <div class="text-center">
      <p class="text-gray-600">Redirecting to login...</p>
    </div>
  </div>
{/if}