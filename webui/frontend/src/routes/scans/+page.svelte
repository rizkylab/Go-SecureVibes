<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/api/client";
    import {
        Search,
        Filter,
        Trash2,
        Eye,
        ChevronLeft,
        ChevronRight,
        Calendar,
        Clock,
        GitBranch,
        AlertTriangle,
    } from "lucide-svelte";

    let scans = $state<any[]>([]);
    let isLoading = $state(true);
    let page = $state(1);
    let totalPages = $state(1);
    let totalItems = $state(0);
    let statusFilter = $state("");

    async function loadScans() {
        isLoading = true;
        try {
            const res = await api.get("/scans", {
                params: {
                    page,
                    size: 10,
                    status: statusFilter,
                },
            });

            scans = res.data.data;
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
            loadScans();
        }
    }

    async function handleDelete(id: string) {
        if (!confirm("Are you sure you want to delete this scan?")) return;

        try {
            await api.delete(`/scans/${id}`);
            loadScans();
        } catch (err) {
            console.error(err);
            alert("Failed to delete scan");
        }
    }

    onMount(() => {
        loadScans();
    });
</script>

<div class="space-y-6">
    <div
        class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4"
    >
        <div>
            <h1 class="text-2xl font-bold text-white">Security Scans</h1>
            <p class="text-slate-400">
                History of all security scans performed
            </p>
        </div>

        <div class="flex items-center gap-2">
            <div class="relative">
                <select
                    class="input pl-10 appearance-none cursor-pointer"
                    bind:value={statusFilter}
                    onchange={() => {
                        page = 1;
                        loadScans();
                    }}
                >
                    <option value="">All Status</option>
                    <option value="completed">Completed</option>
                    <option value="failed">Failed</option>
                    <option value="running">Running</option>
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
                        <th class="px-6 py-3">Status</th>
                        <th class="px-6 py-3">Project / Branch</th>
                        <th class="px-6 py-3">Findings</th>
                        <th class="px-6 py-3">Duration</th>
                        <th class="px-6 py-3">Date</th>
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
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-6 w-40 bg-slate-800 rounded"
                                    ></div></td
                                >
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-6 w-16 bg-slate-800 rounded"
                                    ></div></td
                                >
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-6 w-16 bg-slate-800 rounded"
                                    ></div></td
                                >
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-6 w-24 bg-slate-800 rounded"
                                    ></div></td
                                >
                                <td class="px-6 py-4"
                                    ><div
                                        class="h-6 w-16 bg-slate-800 rounded ml-auto"
                                    ></div></td
                                >
                            </tr>
                        {/each}
                    {:else if scans.length > 0}
                        {#each scans as scan}
                            <tr
                                class="border-b border-border hover:bg-slate-800/50 transition-colors"
                            >
                                <td class="px-6 py-4">
                                    <span
                                        class="badge capitalize
                                        {scan.status === 'completed'
                                            ? 'bg-info/20 text-info'
                                            : scan.status === 'failed'
                                              ? 'bg-critical/20 text-critical'
                                              : 'bg-low/20 text-low'}"
                                    >
                                        {scan.status}
                                    </span>
                                </td>
                                <td class="px-6 py-4">
                                    <div class="font-medium text-white">
                                        {scan.project_path}
                                    </div>
                                    {#if scan.branch}
                                        <div
                                            class="flex items-center gap-1 text-xs text-slate-400 mt-1"
                                        >
                                            <GitBranch size={12} />
                                            {scan.branch}
                                            {#if scan.commit_hash}
                                                <span class="mx-1">•</span>
                                                <span class="font-mono"
                                                    >{scan.commit_hash.substring(
                                                        0,
                                                        7,
                                                    )}</span
                                                >
                                            {/if}
                                        </div>
                                    {/if}
                                </td>
                                <td class="px-6 py-4">
                                    {#if scan.summary}
                                        <div class="flex items-center gap-2">
                                            {#if scan.summary.critical > 0}
                                                <span
                                                    class="text-critical font-bold"
                                                    title="Critical"
                                                    >{scan.summary
                                                        .critical}</span
                                                >
                                            {/if}
                                            {#if scan.summary.high > 0}
                                                <span
                                                    class="text-high font-bold"
                                                    title="High"
                                                    >{scan.summary.high}</span
                                                >
                                            {/if}
                                            {#if scan.summary.medium > 0}
                                                <span
                                                    class="text-medium font-bold"
                                                    title="Medium"
                                                    >{scan.summary.medium}</span
                                                >
                                            {/if}
                                            <span class="text-slate-500 text-xs"
                                                >({scan.summary.total_issues} total)</span
                                            >
                                        </div>
                                    {:else}
                                        <span class="text-slate-500">-</span>
                                    {/if}
                                </td>
                                <td class="px-6 py-4 text-slate-400">
                                    <div class="flex items-center gap-1">
                                        <Clock size={14} />
                                        {scan.duration}ms
                                    </div>
                                </td>
                                <td class="px-6 py-4 text-slate-400">
                                    <div class="flex items-center gap-1">
                                        <Calendar size={14} />
                                        {new Date(
                                            scan.timestamp,
                                        ).toLocaleDateString()}
                                    </div>
                                    <div class="text-xs mt-1">
                                        {new Date(
                                            scan.timestamp,
                                        ).toLocaleTimeString()}
                                    </div>
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <div
                                        class="flex items-center justify-end gap-2"
                                    >
                                        <a
                                            href="/scans/{scan.id}"
                                            class="p-2 hover:bg-primary/20 text-primary-light rounded transition-colors"
                                            title="View Details"
                                        >
                                            <Eye size={18} />
                                        </a>
                                        <button
                                            class="p-2 hover:bg-critical/20 text-critical rounded transition-colors"
                                            title="Delete Scan"
                                            onclick={() =>
                                                handleDelete(scan.id)}
                                        >
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    {:else}
                        <tr>
                            <td
                                colspan="6"
                                class="px-6 py-12 text-center text-slate-500"
                            >
                                <div class="flex flex-col items-center gap-3">
                                    <div
                                        class="h-12 w-12 rounded-full bg-slate-800 flex items-center justify-center"
                                    >
                                        <Search size={24} />
                                    </div>
                                    <p>
                                        No scans found matching your criteria.
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
                    Showing {(page - 1) * 10 + 1} to {Math.min(
                        page * 10,
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

                    {#each Array(totalPages) as _, i}
                        <button
                            class="w-8 h-8 rounded text-sm font-medium transition-colors {page ===
                            i + 1
                                ? 'bg-primary text-white'
                                : 'hover:bg-slate-800 text-slate-400'}"
                            onclick={() => handlePageChange(i + 1)}
                        >
                            {i + 1}
                        </button>
                    {/each}

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
