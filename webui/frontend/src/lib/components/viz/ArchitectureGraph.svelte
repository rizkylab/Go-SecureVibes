<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import cytoscape from "cytoscape";

    let { elements } = $props();
    let container: HTMLElement;
    let cy: any;

    onMount(() => {
        if (!container) return;

        cy = cytoscape({
            container: container,
            elements: elements || [],
            style: [
                {
                    selector: "node",
                    style: {
                        "background-color": "#1f3b73",
                        label: "data(label)",
                        color: "#fff",
                        "text-valign": "bottom",
                        "text-halign": "center",
                        "text-margin-y": 8,
                        "font-size": "12px",
                        width: 40,
                        height: 40,
                        "border-width": 2,
                        "border-color": "#3b82f6",
                    },
                },
                {
                    selector: 'node[type="database"]',
                    style: {
                        shape: "barrel",
                        "background-color": "#22c55e",
                        "border-color": "#16a34a",
                    },
                },
                {
                    selector: 'node[type="service"]',
                    style: {
                        shape: "round-rectangle",
                        "background-color": "#3b82f6",
                        "border-color": "#2563eb",
                    },
                },
                {
                    selector: 'node[type="external"]',
                    style: {
                        shape: "cloud",
                        "background-color": "#64748b",
                        "border-color": "#475569",
                    },
                },
                {
                    selector: "edge",
                    style: {
                        width: 2,
                        "line-color": "#475569",
                        "target-arrow-color": "#475569",
                        "target-arrow-shape": "triangle",
                        "curve-style": "bezier",
                        label: "data(label)",
                        "font-size": "10px",
                        color: "#94a3b8",
                        "text-background-color": "#0f172a",
                        "text-background-opacity": 1,
                        "text-background-padding": 2,
                    },
                },
            ],
            layout: {
                name: "grid",
            },
        });
    });

    onDestroy(() => {
        if (cy) cy.destroy();
    });

    $effect(() => {
        if (cy && elements) {
            cy.json({ elements });
            // Use cose layout for better organic arrangement if available, otherwise grid or circle
            // Note: cose requires an extension usually, but basic layouts are built-in
            cy.layout({
                name: "grid",
                fit: true,
                padding: 50,
            }).run();
        }
    });
</script>

<div
    bind:this={container}
    class="w-full h-[600px] bg-slate-900 rounded-lg border border-border"
></div>
