<script lang="ts">
	import CreateTenantView from '$lib/components/dashboard/create-tenant-view.svelte';
	import Dashboard from '$lib/components/dashboard/dashboard.svelte';
	import type { ActionData, PageData } from './$types';
	let { data, form }: { data: PageData; form: ActionData } = $props();
</script>

<div class="mx-auto flex max-w-7xl flex-col items-center justify-center pt-24">
	<h1 class="text-4xl font-bold">Dashboard</h1>
	{#if !data.tenantData}
		<CreateTenantView createForm={form} />
	{:else}
		<Dashboard tenantData={data.tenantData} {form} />
	{/if}
</div>
