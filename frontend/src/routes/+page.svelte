<script lang="ts">
	import { onMount } from 'svelte';
	import { refreshAccessToken, signInWithToken } from '$services/auth';
	import { goto } from '$app/navigation';
	import { getUserCtx } from '$contexts/userCtx';
	import { setAccessToken } from '$lib/auth';

	const userCtx = getUserCtx();

	onMount(async () => {
		await new Promise((resolve) => {
			setTimeout(() => {
				resolve(1);
			}, 250);
		});

		const token = await refreshAccessToken();
		if (!token) {
			goto('/auth', { replaceState: true });
			return;
		}

		const res = await signInWithToken(token);
		if (res == null) {
			goto('/auth', { replaceState: true });
			return;
		}

		setAccessToken(token);
		userCtx.value = { id: res.userId };
		goto('/app', { replaceState: true });
	});
</script>

<main class="grid h-screen w-screen place-items-center text-4xl font-medium text-white/50">
	Loading...
</main>
