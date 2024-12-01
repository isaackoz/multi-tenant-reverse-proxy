<script lang="ts">
	import '../app.css';
	let { children }: { children: any } = $props();
</script>

<div class="min-h-screen bg-gradient-to-b from-stone-900 to-stone-950">
	{@render children()}
</div>
