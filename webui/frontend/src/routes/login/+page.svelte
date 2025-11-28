<script lang="ts">
    import { auth, login } from "$lib/stores/auth";
    import api from "$lib/api/client";
    import { Shield, Loader2 } from "lucide-svelte";
    import { goto } from "$app/navigation";

    let username = $state("");
    let password = $state("");
    let isLoading = $state(false);
    let error = $state("");

    async function handleSubmit(e: Event) {
        e.preventDefault();
        isLoading = true;
        error = "";

        try {
            const response = await api.post("/auth/login", {
                username,
                password,
            });
            const { token, user } = response.data.data;

            login(token, user);
            goto("/");
        } catch (err: any) {
            console.error(err);
            error =
                err.response?.data?.error ||
                "Login failed. Please check your credentials.";
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="flex min-h-screen items-center justify-center bg-background p-4">
    <div
        class="w-full max-w-md space-y-8 rounded-xl bg-surface p-8 shadow-lg border border-border"
    >
        <div class="text-center">
            <div
                class="mx-auto flex h-12 w-12 items-center justify-center rounded-lg bg-primary text-white"
            >
                <Shield size={32} />
            </div>
            <h2 class="mt-6 text-3xl font-bold tracking-tight text-white">
                SecureVibes
            </h2>
            <p class="mt-2 text-sm text-slate-400">
                Sign in to access your security dashboard
            </p>
        </div>

        <form class="mt-8 space-y-6" onsubmit={handleSubmit}>
            {#if error}
                <div
                    class="rounded-md bg-red-500/10 p-4 text-sm text-red-500 border border-red-500/20"
                >
                    {error}
                </div>
            {/if}

            <div class="space-y-4">
                <div>
                    <label
                        for="username"
                        class="block text-sm font-medium text-slate-300"
                        >Username</label
                    >
                    <input
                        id="username"
                        name="username"
                        type="text"
                        required
                        class="input mt-1"
                        placeholder="admin"
                        bind:value={username}
                    />
                </div>

                <div>
                    <label
                        for="password"
                        class="block text-sm font-medium text-slate-300"
                        >Password</label
                    >
                    <input
                        id="password"
                        name="password"
                        type="password"
                        required
                        class="input mt-1"
                        placeholder="••••••••"
                        bind:value={password}
                    />
                </div>
            </div>

            <button
                type="submit"
                class="btn btn-primary w-full flex justify-center py-2.5"
                disabled={isLoading}
            >
                {#if isLoading}
                    <Loader2 class="mr-2 h-4 w-4 animate-spin" />
                    Signing in...
                {:else}
                    Sign in
                {/if}
            </button>
        </form>

        <div class="text-center text-xs text-slate-500">
            <p>Default credentials: admin / admin123</p>
        </div>
    </div>
</div>
