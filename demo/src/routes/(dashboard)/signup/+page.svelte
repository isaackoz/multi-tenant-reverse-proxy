<script lang="ts">
	import SignUp from 'clerk-sveltekit/client/SignUp.svelte';
</script>

<div class="mx-auto flex max-w-7xl flex-col items-center justify-center pt-24">
	<h1 class="text-4xl font-bold">Signup</h1>
	<div class="mt-12">
		<SignUp redirectUrl="/dashboard" />
	</div>
	<div class="space-x-4">
		<a href="/" class="text-blue-500 underline">Back home</a>
		<a href="/login" class="text-blue-500 underline">Sign in</a>
	</div>
</div>
