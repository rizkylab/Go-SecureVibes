<script lang="ts">
    import { auth, logout } from '$lib/stores/auth';
    import { LogOut, User, Shield } from 'lucide-svelte';
    
    let showUserMenu = false;
    
    function handleLogout() {
        logout();
        window.location.href = '/login';
    }
</script>

<header class="h-16 border-b border-border bg-surface px-6 flex items-center justify-between sticky top-0 z-10">
    <div class="flex items-center gap-3">
        <div class="h-8 w-8 rounded bg-primary flex items-center justify-center text-white">
            <Shield size={20} />
        </div>
        <span class="font-bold text-lg tracking-tight">SecureVibes</span>
    </div>
    
    <div class="flex items-center gap-4">
        {#if $auth.isAuthenticated}
            <div class="relative">
                <button 
                    class="flex items-center gap-2 hover:bg-slate-700 py-1.5 px-3 rounded-md transition-colors"
                    on:click={() => showUserMenu = !showUserMenu}
                >
                    <div class="h-8 w-8 rounded-full bg-slate-700 flex items-center justify-center border border-slate-600">
                        <User size={16} class="text-slate-300" />
                    </div>
                    <span class="text-sm font-medium text-slate-200">{$auth.user?.username}</span>
                </button>
                
                {#if showUserMenu}
                    <div class="absolute right-0 mt-2 w-48 bg-surface border border-border rounded-md shadow-lg py-1 z-50">
                        <button 
                            class="w-full text-left px-4 py-2 text-sm text-slate-300 hover:bg-slate-700 hover:text-white flex items-center gap-2"
                            on:click={handleLogout}
                        >
                            <LogOut size={14} />
                            Sign out
                        </button>
                    </div>
                {/if}
            </div>
        {/if}
    </div>
</header>
