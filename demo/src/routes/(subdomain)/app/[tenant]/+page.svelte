<script lang="ts">
	import ScrollArea from '$lib/components/ui/scroll-area/scroll-area.svelte';
	import { Separator } from '$lib/components/ui/separator';
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<ScrollArea class="mx-auto h-svh max-w-xl px-4 py-12">
	<div class="flex h-full w-full flex-col items-center justify-start">
		<h1 class="text-4xl font-semibold underline underline-offset-4">
			{data.data.title}
		</h1>
		<Separator class="my-4" />
		<p class="mt-8 self-start">
			{data.data.message}
		</p>
	</div>
</ScrollArea>
