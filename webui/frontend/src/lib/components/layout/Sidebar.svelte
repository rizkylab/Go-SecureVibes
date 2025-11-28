<script lang="ts">
    import { page } from '$app/stores';
    import { 
        LayoutDashboard, 
        Scan, 
        Bug, 
        Network, 
        ShieldAlert, 
        Settings,
        GitCompare
    } from 'lucide-svelte';
    
    const menuItems = [
        { href: '/', label: 'Dashboard', icon: LayoutDashboard },
        { href: '/scans', label: 'Scans', icon: Scan },
        { href: '/findings', label: 'Findings', icon: Bug },
        { href: '/architecture', label: 'Architecture', icon: Network },
        { href: '/threats', label: 'Threat Model', icon: ShieldAlert },
        { href: '/compare', label: 'Compare', icon: GitCompare },
        { href: '/settings', label: 'Settings', icon: Settings },
    ];
    
    $: activeUrl = $page.url.pathname;
</script>

<aside class="w-64 bg-surface border-r border-border h-[calc(100vh-4rem)] sticky top-16 overflow-y-auto">
    <nav class="p-4 space-y-1">
        {#each menuItems as item}
            <a 
                href={item.href} 
                class="flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors {activeUrl === item.href || (item.href !== '/' && activeUrl.startsWith(item.href)) ? 'bg-primary/20 text-primary-light' : 'text-slate-400 hover:bg-slate-800 hover:text-slate-100'}"
            >
                <svelte:component this={item.icon} size={18} />
                {item.label}
            </a>
        {/each}
    </nav>
</aside>
