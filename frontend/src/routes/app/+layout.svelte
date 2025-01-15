<script lang="ts">
	import { goto } from '$app/navigation';
	import { getUserCtx } from '$contexts/userCtx';
	import { getAccessToken } from '$lib/auth';
	import { onMount } from 'svelte';

	const userCtx = getUserCtx().value;

	onMount(() => {
		if (getAccessToken() == null || !userCtx) {
			goto('/', { replaceState: true });
		}
	});

	let { children } = $props();
</script>

{#if userCtx}
	{@render children()}
{/if}
