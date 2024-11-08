<script lang="ts">
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<h1>Tenant ID is {data.tenantId}</h1>
