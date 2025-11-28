<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/api/client";
    import {
        Search,
        Filter,
        Bug,
        ChevronLeft,
        ChevronRight,
        FileCode,
        GitBranch,
        ExternalLink,
    } from "lucide-svelte";

    let findings = $state<any[]>([]);
    let isLoading = $state(true);
    let page = $state(1);
    let totalPages = $state(1);
    let totalItems = $state(0);
    let severityFilter = $state("");
    let categoryFilter = $state("");

    async function loadFindings() {
        isLoading = true;
        try {
            const res = await api.get("/findings", {
                params: {
                    page,
                    size: 20,
                    severity: severityFilter,
                    category: categoryFilter,
                },
            });

            findings = res.data.data;
            totalPages = res.data.total_pages;
            totalItems = res.data.total_items;
        } catch (err) {
            console.error(err);
        } finally {
            isLoading = false;
        }
    }

    function handlePageChange(newPage: number) {
        if (newPage >= 1 && newPage <= totalPages) {
            page = newPage;
            loadFindings();
        }
    }

    onMount(() => {
        loadFindings();
    });
</script>

<div class="space-y-6">
    <div
        class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4"
    >
        <div>
            <h1 class="text-2xl font-bold text-white">Findings Explorer</h1>
            <p class="text-slate-400">
                Search and analyze security findings across all projects
            </p>
        </div>

        <div class="flex items-center gap-2">
            <div class="relative">
                <select
                    class="input pl-10 appearance-none cursor-pointer"
                    bind:value={severityFilter}
                    onchange={() => {
                        page = 1;
                        loadFindings();
                    }}
                >
                    <option value="">All Severities</option>
                    <option value="critical">Critical</option>
                    <option value="high">High</option>
                    <option value="medium">Medium</option>
                    <option value="low">Low</option>
                    <option value="info">Info</option>
                </select>
                <Filter
                    class="absolute left-3 top-2.5 text-slate-400"
                    size={16}
                />
            </div>

            <div class="relative">
                <select
                    class="input pl-10 appearance-none cursor-pointer"
                    bind:value={categoryFilter}
                    onchange={() => {
                        page = 1;
                        loadFindings();
                    }}
                >
                    <option value="">All Categories</option>
                    <option value="security">Security</option>
                    <option value="sast">SAST</option>
                    <option value="secret">Secret</option>
                    <option value="dependency">Dependency</option>
                </select>
                <Filter
                    class="absolute left-3 top-2.5 text-slate-400"
                    size={16}
                />
            </div>
        </div>
    </div>

    <div class="card overflow-hidden">
        <div class="overflow-x-auto">
            <table class="w-full text-sm text-left">
                <thead
                    class="text-xs text-slate-400 uppercase bg-slate-900/50 border-b border-border"
                >
                    <tr>
                        <th class="px-6 py-3">Severity</th>
                        <th class="px-6 py-3">Finding</th>
                        <th class="px-6 py-3">Location</th>
                        <th class="px-6 py-3">Project</th>
                        <th class="px-6 py-3 text-right">Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#if isLoading}
                        {#each Array(5) as _}
                            <tr class="border-b border-border animate-pulse">
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-6 w-20 bg-slate-800 rounded"
                                    ></div></td
                                >
                                <td class="px-6 py-4">
                                    <div
                                        class="h-4 w-48 bg-slate-800 rounded mb-2"
                                    ></div>
                                    <div
                                        class="h-3 w-32 bg-slate-800 rounded"
                                    ></div>
                                </td>
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-4 w-32 bg-slate-800 rounded"
                                    ></div></td
                                >
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-4 w-24 bg-slate-800 rounded"
                                    ></div></td
                                >
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-6 w-8 bg-slate-800 rounded ml-auto"
                                    ></div></td
                                >
                            </tr>
                        {/each}
                    {:else if findings.length > 0}
                        {#each findings as finding}
                            <tr
                                class="border-b border-border hover:bg-slate-800/50 transition-colors group"
                            >
                                <td class="px-6 py-4 align-top">
                                    <span
                                        class="badge capitalize
                                        {finding.severity === 'critical'
                                            ? 'bg-critical/20 text-critical'
                                            : finding.severity === 'high'
                                              ? 'bg-high/20 text-high'
                                              : finding.severity === 'medium'
                                                ? 'bg-medium/20 text-medium'
                                                : finding.severity === 'low'
                                                  ? 'bg-low/20 text-low'
                                                  : 'bg-info/20 text-info'}"
                                    >
                                        {finding.severity}
                                    </span>
                                </td>
                                <td class="px-6 py-4 align-top">
                                    <div
                                        class="font-medium text-white group-hover:text-primary-light transition-colors"
                                    >
                                        {finding.title}
                                    </div>
                                    <div
                                        class="text-xs text-slate-400 mt-1 line-clamp-2"
                                    >
                                        {finding.description}
                                    </div>
                                    <div class="flex items-center gap-2 mt-2">
                                        {#if finding.cwe}
                                            <span
                                                class="text-xs bg-slate-800 px-1.5 py-0.5 rounded text-slate-300"
                                                >CWE-{finding.cwe}</span
                                            >
                                        {/if}
                                        <span
                                            class="text-xs bg-slate-800 px-1.5 py-0.5 rounded text-slate-300 capitalize"
                                            >{finding.category}</span
                                        >
                                    </div>
                                </td>
                                <td class="px-6 py-4 align-top text-slate-400">
                                    <div
                                        class="flex items-center gap-1.5 text-xs font-mono break-all"
                                    >
                                        <FileCode size={12} class="shrink-0" />
                                        {finding.file_path}:{finding.line_number}
                                    </div>
                                </td>
                                <td class="px-6 py-4 align-top text-slate-400">
                                    <div
                                        class="flex items-center gap-1.5 text-xs"
                                    >
                                        <div class="font-medium text-slate-300">
                                            {finding.project_path}
                                        </div>
                                    </div>
                                    {#if finding.branch}
                                        <div
                                            class="flex items-center gap-1 text-xs mt-1"
                                        >
                                            <GitBranch size={12} />
                                            {finding.branch}
                                        </div>
                                    {/if}
                                </td>
                                <td class="px-6 py-4 align-top text-right">
                                    <a
                                        href="/scans/{finding.scan_id}?tab=findings&finding={finding.id}"
                                        class="p-2 hover:bg-primary/20 text-primary-light rounded transition-colors inline-block"
                                        title="View in Context"
                                    >
                                        <ExternalLink size={16} />
                                    </a>
                                </td>
                            </tr>
                        {/each}
                    {:else}
                        <tr>
                            <td
                                colspan="5"
                                class="px-6 py-12 text-center text-slate-500"
                            >
                                <div class="flex flex-col items-center gap-3">
                                    <div
                                        class="h-12 w-12 rounded-full bg-slate-800 flex items-center justify-center"
                                    >
                                        <Bug size={24} />
                                    </div>
                                    <p>
                                        No findings found matching your
                                        criteria.
                                    </p>
                                </div>
                            </td>
                        </tr>
                    {/if}
                </tbody>
            </table>
        </div>

        <!-- Pagination -->
        {#if totalPages > 1}
            <div
                class="px-6 py-4 border-t border-border flex items-center justify-between"
            >
                <div class="text-sm text-slate-400">
                    Showing {(page - 1) * 20 + 1} to {Math.min(
                        page * 20,
                        totalItems,
                    )} of {totalItems} results
                </div>

                <div class="flex items-center gap-2">
                    <button
                        class="p-2 rounded hover:bg-slate-800 disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={page === 1}
                        onclick={() => handlePageChange(page - 1)}
                    >
                        <ChevronLeft size={18} />
                    </button>

                    <span class="text-sm text-slate-400 px-2">
                        Page {page} of {totalPages}
                    </span>

                    <button
                        class="p-2 rounded hover:bg-slate-800 disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={page === totalPages}
                        onclick={() => handlePageChange(page + 1)}
                    >
                        <ChevronRight size={18} />
                    </button>
                </div>
            </div>
        {/if}
    </div>
</div>
