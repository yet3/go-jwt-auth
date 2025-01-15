<script lang="ts">
	import type { HTMLInputAttributes } from 'svelte/elements';

	interface IProps extends HTMLInputAttributes {}

	const props: IProps = $props();
</script>

<input
	{...props}
	class={['w-full rounded border border-white p-2 text-lg text-gray-900', props.class].join(' ')}
/>
