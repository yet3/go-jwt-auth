<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '$common/Button.svelte';
	import Input from '$common/Input.svelte';
	import { getUserCtx } from '$contexts/userCtx';
	import { setAccessToken } from '$lib/auth';
	import { signInWithEmail, signUp } from '$services/auth';

	const userCtx = getUserCtx();

	const handleSubmit = (type: 'sing-in' | 'sing-up') => {
		return async (e: SubmitEvent) => {
			e.preventDefault();

			const formEl = e.target as HTMLFormElement;
			if (!(formEl instanceof HTMLFormElement)) return;

			const inputs = formEl.querySelectorAll('input');
			const values: Record<string, string> = {};
			for (const input of inputs) {
				if (input.name) {
					values[input.name] = input.value;
				}
			}

			if (type === 'sing-in') {
				const data = await signInWithEmail(values.email, values.password);
				if (!data) {
          alert("Error signing-in")
          return;
        };
				setAccessToken(data.accessToken);
				userCtx.value = { id: data.userId };
			} else if (type === 'sing-up') {
				if (values.password !== values.repeatPassword) {
					alert('Password must match!');
					return;
				}

				const data = await signUp(values.email, values.password);
				if (!data) {
          alert("Error signing-up")
          return;
        };
				setAccessToken(data.accessToken);
				userCtx.value = { id: data.userId };
			}

			goto('/app', { replaceState: true });
		};
	};
</script>

<main class="flex h-screen w-screen items-center justify-center space-x-8">
	<form class="flex w-80 flex-col" onsubmit={handleSubmit('sing-in')}>
		<label class="mb-2">
			Email
			<Input defaultValue="tak@wp.pl" type="text" name="email" />
		</label>
		<label>
			Password
			<Input defaultValue="I got it" type="password" name="password" />
		</label>
		<Button type="submit" class="mt-4">Sign in</Button>
	</form>

	<form class="flex w-80 flex-col" onsubmit={handleSubmit('sing-up')}>
		<label class="mb-2">
			Email
			<Input type="text" name="email" />
		</label>

		<label class="mb-2">
			Password
			<Input type="password" name="password" />
		</label>

		<label>
			Repeat password
			<Input type="password" name="repeatPassword" />
		</label>
		<Button type="submit" class="mt-4">Sign up</Button>
	</form>
</main>
