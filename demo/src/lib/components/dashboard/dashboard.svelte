<script lang="ts">
	import { dev } from '$app/environment';
	import { enhance } from '$app/forms';
	import { PUBLIC_DOMAIN } from '$env/static/public';
	import type { ActionData, PageData } from '../../../routes/(dashboard)/dashboard/$types';
	import { Button } from '../ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardFooter,
		CardHeader,
		CardTitle
	} from '../ui/card';
	import { ExternalLinkIcon } from 'lucide-svelte';
	import { Textarea } from '../ui/textarea';
	import { Label } from '../ui/label';
	import { Input } from '../ui/input';
	import { hostnameToValidUrl } from '$lib/utils';

	let { tenantData, form }: { tenantData: PageData['tenantData']; form: ActionData } = $props();
	let title = $state(tenantData?.title ?? '');
	let message = $state(tenantData?.message ?? '');
</script>

<div class="mt-4 flex w-full justify-center">
	<Card class="w-full max-w-lg">
		<CardHeader>
			<CardTitle class="text-center font-mono">
				<a
					href={`${dev ? 'http://' : 'https://'}${tenantData?.hostname ? hostnameToValidUrl(tenantData.hostname) : ''}.${PUBLIC_DOMAIN}${dev ? ':5173' : ''}`}
					class="flex justify-center gap-2 text-blue-500 underline underline-offset-4"
					target="_blank"
				>
					<span>
						{tenantData?.id}
					</span>
					<ExternalLinkIcon class="size-5" />
				</a>
			</CardTitle>
			<CardDescription>Modify your page</CardDescription>
		</CardHeader>
		<CardContent class="space-y-4">
			<form
				use:enhance={() => {
					return ({ update }) => update({ reset: false });
				}}
				method="POST"
				action="?/updateHostname"
			>
				<div>
					<Label for="hostname">Hostname</Label>
					<Input id="hostname" name="hostname" value={tenantData?.hostname} />
					<p class="text-muted-foreground mb-1 text-xs">
						If the hostname has a port, make sure to include it
					</p>
				</div>
				{#if form?.hostsuccess}
					<p class="text-sm text-green-500">Data updated successfully</p>
				{:else if form?.hostmessage}
					<p class="text-sm text-red-500">{form.hostmessage}</p>
				{/if}
				<Button type="submit" class="mt-4 w-full">Update Hostname</Button>
			</form>
			<form
				use:enhance={() => {
					return ({ update }) => update({ reset: false });
				}}
				method="POST"
				action="?/update"
				data-sveltekit-reload
			>
				<div>
					<Label for="title">Title</Label>
					<Input id="title" name="title" bind:value={title} />
				</div>
				<div>
					<Label for="message">Message</Label>
					<Textarea name="message" id="message" bind:value={message} />
				</div>
				<div>
					{#if form?.success}
						<p class="text-sm text-green-500">Data updated successfully</p>
					{:else if form?.message}
						<p class="text-sm text-red-500">{form.message}</p>
					{/if}
				</div>
				<Button type="submit" class="mt-4 w-full">Update Page Data</Button>
			</form>
		</CardContent>
	</Card>
</div>
