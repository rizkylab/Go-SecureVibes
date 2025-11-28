<script lang="ts">
    import '../app.css';
    import Navbar from '$lib/components/layout/Navbar.svelte';
    import Sidebar from '$lib/components/layout/Sidebar.svelte';
    import { page } from '$app/stores';
    import { auth } from '$lib/stores/auth';
    import { onMount } from 'svelte';
    
    let { children } = $props();
    
    // Check if we are on the login page
    let isLoginPage = $derived($page.url.pathname === '/login');
    
    onMount(() => {
        // Redirect to login if not authenticated and not already on login page
        if (!$auth.isAuthenticated && !isLoginPage) {
            window.location.href = '/login';
        }
    });
</script>

<div class="min-h-screen bg-background text-slate-100 font-sans">
    {#if !isLoginPage}
        <Navbar />
        <div class="flex">
            <Sidebar />
            <main class="flex-1 p-6 overflow-x-hidden">
                {@render children()}
            </main>
        </div>
    {:else}
        {@render children()}
    {/if}
</div>
