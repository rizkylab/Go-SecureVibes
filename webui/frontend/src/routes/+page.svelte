<script lang="ts">
    import { onMount } from "svelte";
    import api from "$lib/api/client";
    import { chart } from "$lib/actions/chart";
    import {
        Shield,
        AlertTriangle,
        CheckCircle,
        Activity,
    } from "lucide-svelte";

    let summary = $state<any>(null);
    let isLoading = $state(true);

    // Chart options
    let severityChartOptions = $state<any>({
        chart: { type: "donut", height: 350, background: "transparent" },
        labels: ["Critical", "High", "Medium", "Low", "Info"],
        colors: ["#ef4444", "#f97316", "#eab308", "#3b82f6", "#22c55e"],
        theme: { mode: "dark" },
        plotOptions: { pie: { donut: { size: "70%" } } },
        dataLabels: { enabled: false },
        stroke: { show: false },
        legend: {
            position: "bottom",
            fontFamily: "Inter",
            labels: { colors: "#94a3b8" },
        },
    });

    let trendChartOptions = $state<any>({
        chart: {
            type: "area",
            height: 350,
            toolbar: { show: false },
            background: "transparent",
        },
        colors: ["#1f3b73"],
        stroke: { curve: "smooth", width: 2 },
        fill: {
            type: "gradient",
            gradient: { shadeIntensity: 1, opacityFrom: 0.7, opacityTo: 0.3 },
        },
        theme: { mode: "dark" },
        xaxis: {
            type: "datetime",
            tooltip: { enabled: false },
            axisBorder: { show: false },
            axisTicks: { show: false },
            labels: { style: { colors: "#94a3b8" } },
        },
        yaxis: { show: false },
        grid: { show: false, padding: { left: 0, right: 0 } },
        dataLabels: { enabled: false },
    });

    onMount(async () => {
        try {
            const res = await api.get("/dashboard/summary");
            summary = res.data.data;

            // Update charts
            if (summary) {
                severityChartOptions = {
                    ...severityChartOptions,
                    series: [
                        summary.severity_breakdown.critical,
                        summary.severity_breakdown.high,
                        summary.severity_breakdown.medium,
                        summary.severity_breakdown.low,
                        summary.severity_breakdown.info,
                    ],
                };

                trendChartOptions = {
                    ...trendChartOptions,
                    series: [
                        {
                            name: "Scans",
                            data: summary.trend_data.map((d: any) => ({
                                x: d.date,
                                y: d.count,
                            })),
                        },
                    ],
                };
            }
        } catch (err) {
            console.error(err);
        } finally {
            isLoading = false;
        }
    });
</script>

<div class="space-y-6">
    <div>
        <h1 class="text-2xl font-bold text-white">Dashboard</h1>
        <p class="text-slate-400">Overview of your security posture</p>
    </div>

    {#if isLoading}
        <div
            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 animate-pulse"
        >
            {#each Array(4) as _}
                <div class="h-32 bg-surface rounded-lg"></div>
            {/each}
        </div>
    {:else if summary}
        <!-- Stats Cards -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="card p-6 flex items-center justify-between">
                <div>
                    <p class="text-sm font-medium text-slate-400">
                        Total Scans
                    </p>
                    <h3 class="text-2xl font-bold mt-1">
                        {summary.total_scans}
                    </h3>
                </div>
                <div
                    class="h-12 w-12 rounded-full bg-primary/20 flex items-center justify-center text-primary-light"
                >
                    <Activity size={24} />
                </div>
            </div>

            <div class="card p-6 flex items-center justify-between">
                <div>
                    <p class="text-sm font-medium text-slate-400">
                        Total Findings
                    </p>
                    <h3 class="text-2xl font-bold mt-1">
                        {summary.total_findings}
                    </h3>
                </div>
                <div
                    class="h-12 w-12 rounded-full bg-red-500/20 flex items-center justify-center text-red-500"
                >
                    <AlertTriangle size={24} />
                </div>
            </div>

            <div class="card p-6 flex items-center justify-between">
                <div>
                    <p class="text-sm font-medium text-slate-400">
                        Critical Issues
                    </p>
                    <h3 class="text-2xl font-bold mt-1">
                        {summary.severity_breakdown.critical}
                    </h3>
                </div>
                <div
                    class="h-12 w-12 rounded-full bg-critical/20 flex items-center justify-center text-critical"
                >
                    <Shield size={24} />
                </div>
            </div>

            <div class="card p-6 flex items-center justify-between">
                <div>
                    <p class="text-sm font-medium text-slate-400">
                        Recent Scans
                    </p>
                    <h3 class="text-2xl font-bold mt-1">
                        {summary.recent_scans}
                    </h3>
                </div>
                <div
                    class="h-12 w-12 rounded-full bg-info/20 flex items-center justify-center text-info"
                >
                    <CheckCircle size={24} />
                </div>
            </div>
        </div>

        <!-- Charts -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div class="card p-6">
                <h3 class="text-lg font-semibold mb-4">Severity Breakdown</h3>
                <div use:chart={severityChartOptions} class="w-full"></div>
            </div>

            <div class="card p-6">
                <h3 class="text-lg font-semibold mb-4">
                    Scan Activity (30 Days)
                </h3>
                <div use:chart={trendChartOptions} class="w-full"></div>
            </div>
        </div>

        <!-- Recent Findings -->
        <div class="card overflow-hidden">
            <div class="p-6 border-b border-border">
                <h3 class="text-lg font-semibold">Recent Findings</h3>
            </div>
            <div class="overflow-x-auto">
                <table class="w-full text-sm text-left">
                    <thead
                        class="text-xs text-slate-400 uppercase bg-slate-900/50"
                    >
                        <tr>
                            <th class="px-6 py-3">Severity</th>
                            <th class="px-6 py-3">Title</th>
                            <th class="px-6 py-3">File</th>
                            <th class="px-6 py-3">Time</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#if summary.recent_findings && summary.recent_findings.length > 0}
                            {#each summary.recent_findings as finding}
                                <tr
                                    class="border-b border-border hover:bg-slate-800/50"
                                >
                                    <td class="px-6 py-4">
                                        <span
                                            class="badge
                                            {finding.severity === 'critical'
                                                ? 'bg-critical/20 text-critical'
                                                : finding.severity === 'high'
                                                  ? 'bg-high/20 text-high'
                                                  : finding.severity ===
                                                      'medium'
                                                    ? 'bg-medium/20 text-medium'
                                                    : 'bg-low/20 text-low'} capitalize"
                                        >
                                            {finding.severity}
                                        </span>
                                    </td>
                                    <td class="px-6 py-4 font-medium"
                                        >{finding.title}</td
                                    >
                                    <td
                                        class="px-6 py-4 text-slate-400 font-mono text-xs"
                                        >{finding.file_path}</td
                                    >
                                    <td class="px-6 py-4 text-slate-400">
                                        {new Date(
                                            finding.timestamp,
                                        ).toLocaleDateString()}
                                    </td>
                                </tr>
                            {/each}
                        {:else}
                            <tr>
                                <td
                                    colspan="4"
                                    class="px-6 py-8 text-center text-slate-500"
                                >
                                    No findings yet. Run a scan to see results.
                                </td>
                            </tr>
                        {/if}
                    </tbody>
                </table>
            </div>
        </div>
    {/if}
</div>
