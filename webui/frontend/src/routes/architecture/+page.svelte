<script>
    import { Network, ArrowRight } from "lucide-svelte";
</script>

<div
    class="flex flex-col items-center justify-center min-h-[60vh] text-center space-y-6"
>
    <div
        class="h-24 w-24 rounded-full bg-slate-800 flex items-center justify-center text-primary-light"
    >
        <Network size={48} />
    </div>

    <div class="max-w-md space-y-2">
        <h1 class="text-2xl font-bold text-white">
            Architecture Visualization
        </h1>
        <p class="text-slate-400">
            Select a scan to view its architecture diagram.
        </p>
    </div>

    <a href="/scans" class="btn btn-primary gap-2">
        Go to Scans
        <ArrowRight size={16} />
    </a>
</div>
