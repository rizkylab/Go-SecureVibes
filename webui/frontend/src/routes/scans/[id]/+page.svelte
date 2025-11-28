<script lang="ts">
    import { page } from "$app/stores";
    import { onMount } from "svelte";
    import api from "$lib/api/client";
    import ArchitectureGraph from "$lib/components/viz/ArchitectureGraph.svelte";
    import {
        ArrowLeft,
        Calendar,
        Clock,
        GitBranch,
        Download,
        Shield,
        Bug,
        Network,
        ShieldAlert,
        Activity,
        FileCode,
        ExternalLink,
        ChevronDown,
        ChevronRight,
    } from "lucide-svelte";

    let scanId = $derived($page.params.id);
    let scan = $state<any>(null);
    let isLoading = $state(true);
    let activeTab = $state("findings");

    // Tab data
    let findings = $state<any[]>([]);
    let architecture = $state<any>(null);
    let threats = $state<any[]>([]);
    let isLoadingTab = $state(false);

    async function loadScan() {
        isLoading = true;
        try {
            const res = await api.get(`/scans/${scanId}`);
            scan = res.data.data;
        } catch (err) {
            console.error(err);
        } finally {
            isLoading = false;
        }
    }

    async function loadFindings() {
        if (findings.length > 0) return;
        isLoadingTab = true;
        try {
            const res = await api.get(`/scans/${scanId}/findings?size=100`);
            findings = res.data.data;
        } catch (err) {
            console.error(err);
        } finally {
            isLoadingTab = false;
        }
    }

    async function loadArchitecture() {
        if (architecture) return;
        isLoadingTab = true;
        try {
            const res = await api.get(`/scans/${scanId}/architecture`);
            // Transform data for Cytoscape if needed, assuming backend returns compatible JSON
            architecture = res.data.data;
        } catch (err) {
            console.error(err);
        } finally {
            isLoadingTab = false;
        }
    }

    async function loadThreats() {
        if (threats.length > 0) return;
        isLoadingTab = true;
        try {
            const res = await api.get(`/scans/${scanId}/threat-model`);
            threats = res.data.data;
        } catch (err) {
            console.error(err);
        } finally {
            isLoadingTab = false;
        }
    }

    $effect(() => {
        if (activeTab === "findings") loadFindings();
        if (activeTab === "architecture") loadArchitecture();
        if (activeTab === "threats") loadThreats();
    });

    onMount(() => {
        loadScan();
    });
</script>

<div class="space-y-6">
    <div class="flex items-center gap-4">
        <a
            href="/scans"
            class="p-2 hover:bg-slate-800 rounded-full transition-colors text-slate-400 hover:text-white"
        >
            <ArrowLeft size={20} />
        </a>

        {#if scan}
            <div>
                <h1
                    class="text-2xl font-bold text-white flex items-center gap-3"
                >
                    Scan Details
                    <span
                        class="badge capitalize text-sm font-normal
                        {scan.status === 'completed'
                            ? 'bg-info/20 text-info'
                            : scan.status === 'failed'
                              ? 'bg-critical/20 text-critical'
                              : 'bg-low/20 text-low'}"
                    >
                        {scan.status}
                    </span>
                </h1>
                <div
                    class="flex items-center gap-4 text-sm text-slate-400 mt-1"
                >
                    <span class="flex items-center gap-1">
                        <Calendar size={14} />
                        {new Date(scan.timestamp).toLocaleString()}
                    </span>
                    <span class="flex items-center gap-1">
                        <Clock size={14} />
                        {scan.duration}ms
                    </span>
                    {#if scan.branch}
                        <span class="flex items-center gap-1">
                            <GitBranch size={14} />
                            {scan.branch}
                        </span>
                    {/if}
                </div>
            </div>

            <div class="ml-auto flex gap-2">
                <button class="btn btn-secondary gap-2">
                    <Download size={16} />
                    Export Report
                </button>
            </div>
        {/if}
    </div>

    {#if isLoading}
        <div class="h-64 bg-surface rounded-lg animate-pulse"></div>
    {:else if scan}
        <!-- Summary Cards -->
        <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
            <div class="card p-4 bg-critical/10 border-critical/20">
                <p class="text-xs font-medium text-critical uppercase">
                    Critical
                </p>
                <p class="text-2xl font-bold text-white mt-1">
                    {scan.summary?.critical || 0}
                </p>
            </div>
            <div class="card p-4 bg-high/10 border-high/20">
                <p class="text-xs font-medium text-high uppercase">High</p>
                <p class="text-2xl font-bold text-white mt-1">
                    {scan.summary?.high || 0}
                </p>
            </div>
            <div class="card p-4 bg-medium/10 border-medium/20">
                <p class="text-xs font-medium text-medium uppercase">Medium</p>
                <p class="text-2xl font-bold text-white mt-1">
                    {scan.summary?.medium || 0}
                </p>
            </div>
            <div class="card p-4 bg-low/10 border-low/20">
                <p class="text-xs font-medium text-low uppercase">Low</p>
                <p class="text-2xl font-bold text-white mt-1">
                    {scan.summary?.low || 0}
                </p>
            </div>
            <div class="card p-4 bg-info/10 border-info/20">
                <p class="text-xs font-medium text-info uppercase">Info</p>
                <p class="text-2xl font-bold text-white mt-1">
                    {scan.summary?.info || 0}
                </p>
            </div>
        </div>

        <!-- Tabs -->
        <div class="border-b border-border">
            <nav class="flex gap-6">
                <button
                    class="pb-4 text-sm font-medium border-b-2 transition-colors flex items-center gap-2
                    {activeTab === 'findings'
                        ? 'border-primary text-primary-light'
                        : 'border-transparent text-slate-400 hover:text-slate-200'}"
                    onclick={() => (activeTab = "findings")}
                >
                    <Bug size={16} />
                    Findings
                </button>
                <button
                    class="pb-4 text-sm font-medium border-b-2 transition-colors flex items-center gap-2
                    {activeTab === 'architecture'
                        ? 'border-primary text-primary-light'
                        : 'border-transparent text-slate-400 hover:text-slate-200'}"
                    onclick={() => (activeTab = "architecture")}
                >
                    <Network size={16} />
                    Architecture
                </button>
                <button
                    class="pb-4 text-sm font-medium border-b-2 transition-colors flex items-center gap-2
                    {activeTab === 'threats'
                        ? 'border-primary text-primary-light'
                        : 'border-transparent text-slate-400 hover:text-slate-200'}"
                    onclick={() => (activeTab = "threats")}
                >
                    <ShieldAlert size={16} />
                    Threat Model
                </button>
                {#if scan.dast_enabled}
                    <button
                        class="pb-4 text-sm font-medium border-b-2 transition-colors flex items-center gap-2
                        {activeTab === 'dast'
                            ? 'border-primary text-primary-light'
                            : 'border-transparent text-slate-400 hover:text-slate-200'}"
                        onclick={() => (activeTab = "dast")}
                    >
                        <Activity size={16} />
                        DAST Results
                    </button>
                {/if}
            </nav>
        </div>

        <!-- Tab Content -->
        <div class="min-h-[400px]">
            {#if isLoadingTab}
                <div class="flex items-center justify-center h-64">
                    <div
                        class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"
                    ></div>
                </div>
            {:else if activeTab === "findings"}
                <div class="space-y-4">
                    {#each findings as finding}
                        <div
                            class="card p-4 hover:border-primary/50 transition-colors"
                        >
                            <div class="flex items-start gap-4">
                                <span
                                    class="badge capitalize mt-1
                                    {finding.severity === 'critical'
                                        ? 'bg-critical/20 text-critical'
                                        : finding.severity === 'high'
                                          ? 'bg-high/20 text-high'
                                          : finding.severity === 'medium'
                                            ? 'bg-medium/20 text-medium'
                                            : 'bg-low/20 text-low'}"
                                >
                                    {finding.severity}
                                </span>
                                <div class="flex-1">
                                    <h3 class="font-medium text-white">
                                        {finding.title}
                                    </h3>
                                    <p class="text-sm text-slate-400 mt-1">
                                        {finding.description}
                                    </p>

                                    <div
                                        class="flex items-center gap-4 mt-3 text-xs text-slate-500"
                                    >
                                        <div
                                            class="flex items-center gap-1 font-mono bg-slate-900 px-2 py-1 rounded"
                                        >
                                            <FileCode size={12} />
                                            {finding.file_path}:{finding.line_number}
                                        </div>
                                        <span
                                            class="capitalize bg-slate-800 px-2 py-1 rounded"
                                            >{finding.category}</span
                                        >
                                        {#if finding.cwe}
                                            <span
                                                class="bg-slate-800 px-2 py-1 rounded"
                                                >CWE-{finding.cwe}</span
                                            >
                                        {/if}
                                    </div>

                                    {#if finding.line_content}
                                        <div
                                            class="mt-3 bg-slate-950 p-3 rounded border border-slate-800 font-mono text-xs text-slate-300 overflow-x-auto"
                                        >
                                            <code
                                                >{finding.line_content.trim()}</code
                                            >
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        </div>
                    {:else}
                        <div class="text-center py-12 text-slate-500">
                            <Bug size={48} class="mx-auto mb-4 opacity-50" />
                            <p>No findings detected in this scan.</p>
                        </div>
                    {/each}
                </div>
            {:else if activeTab === "architecture"}
                {#if architecture}
                    <ArchitectureGraph elements={architecture.elements || []} />
                {:else}
                    <div class="text-center py-12 text-slate-500">
                        <Network size={48} class="mx-auto mb-4 opacity-50" />
                        <p>No architecture data available.</p>
                    </div>
                {/if}
            {:else if activeTab === "threats"}
                <div class="space-y-4">
                    {#each threats as threat}
                        <div class="card p-4 border-l-4 border-l-critical">
                            <div class="flex justify-between items-start">
                                <h3 class="font-bold text-white">
                                    {threat.title}
                                </h3>
                                <span class="badge bg-slate-800 text-slate-300"
                                    >{threat.stride_category}</span
                                >
                            </div>
                            <p class="text-sm text-slate-400 mt-2">
                                {threat.description}
                            </p>
                            <div class="mt-3 p-3 bg-slate-900/50 rounded">
                                <p class="text-xs font-semibold text-info mb-1">
                                    Mitigation:
                                </p>
                                <p class="text-sm text-slate-300">
                                    {threat.mitigation}
                                </p>
                            </div>
                        </div>
                    {:else}
                        <div class="text-center py-12 text-slate-500">
                            <ShieldAlert
                                size={48}
                                class="mx-auto mb-4 opacity-50"
                            />
                            <p>No threat model generated.</p>
                        </div>
                    {/each}
                </div>
            {:else if activeTab === "dast"}
                <div
                    class="text-center py-12 text-slate-500 bg-surface rounded-lg border border-border border-dashed"
                >
                    <Activity size={48} class="mx-auto mb-4 opacity-50" />
                    <p>DAST Results Component (Coming Soon)</p>
                </div>
            {/if}
        </div>
    {/if}
</div>
