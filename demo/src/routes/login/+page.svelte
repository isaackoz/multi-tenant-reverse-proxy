<script lang="ts">
	import SignIn from 'clerk-sveltekit/client/SignIn.svelte';
</script>

<div class="mx-auto flex max-w-7xl flex-col items-center justify-center pt-24">
	<h1 class="text-4xl font-bold">Login</h1>
	<div class="mt-12">
		<SignIn redirectUrl="/dashboard" signUpUrl="/signup" />
	</div>
	<div class="space-x-4">
		<a href="/" class="text-blue-500 underline">Back home</a>
		<a href="/signup" class="text-blue-500 underline">Sign up</a>
	</div>
</div>
