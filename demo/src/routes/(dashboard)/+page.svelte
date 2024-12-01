<script lang="ts">
	import { Button } from '$lib/components/ui/button';
</script>

<div class="mx-auto flex max-w-7xl flex-col items-center justify-center pt-24">
	<h1 class="text-4xl font-bold">Multi-Tenant Demo</h1>
	<p class="text-muted-foreground mt-2">
		This is a demo application for the dynamic reverse proxy. Login below to try it out.
	</p>
	<div class="mt-10">
		<Button href="/login">Login</Button>
		<Button href="/signup">Sign up</Button>
	</div>
</div>
