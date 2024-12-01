<script>
	import UserButton from 'clerk-sveltekit/client/UserButton.svelte';
	import SignedIn from 'clerk-sveltekit/client/SignedIn.svelte';
	import SignedOut from 'clerk-sveltekit/client/SignedOut.svelte';
	let { children } = $props();
</script>

<header class="h-12 bg-black">
	<div class=" mx-auto flex h-full max-w-7xl items-center justify-between">
		<div>Multi-Tenant Demo</div>
		<div class="flex">
			<SignedIn>
				<UserButton afterSignOutUrl="/" />
			</SignedIn>
			<SignedOut>
				<div class="flex space-x-4">
					<a href="/login">Sign in</a> <span> | </span> <a href="/signup">Sign up</a>
				</div>
			</SignedOut>
		</div>
	</div>
</header>
{@render children()}
