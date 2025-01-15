<script lang="ts">
	import type { HTMLButtonAttributes } from 'svelte/elements';

	interface IProps extends HTMLButtonAttributes {}

	const { children, ...props }: IProps = $props();
</script>

<button {...props} class={['rounded border border-white px-3 py-2 text-lg', props.class].join(' ')}
	>{@render children?.()}</button
>
