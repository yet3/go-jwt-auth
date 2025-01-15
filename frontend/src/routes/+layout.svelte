<script lang="ts">
	import '../app.css';
	import { createUserCtx } from '$contexts/userCtx';
	import type { IUserCtx } from '$types/user.types';

	let userCtx = $state<IUserCtx>({ value: null });

	createUserCtx(userCtx);

	let { children } = $props();
</script>

{@render children()}
