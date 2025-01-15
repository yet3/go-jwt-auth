<script lang="ts">
	import { signOut } from '$services/auth';
	import { useUserCtx } from '$contexts/userCtx';
	import { goto } from '$app/navigation';
	import { getAccessToken } from '$lib/auth';
	import Button from '$common/Button.svelte';

	const handleSignOut = async () => {
		const signedOut = await signOut(getAccessToken()!);
		if (signedOut) {
			goto('/', { replaceState: true });
		}
	};

	const user = useUserCtx();
</script>

<main class="w-screen h-screen flex flex-col justify-center items-center">
	<h1 class="text-4xl mb-4">App</h1>
	<h2 class="text-2xl mb-4">User ID: {user.id}</h2>
  <Button onclick={handleSignOut}>Sign out</Button>
</main>
