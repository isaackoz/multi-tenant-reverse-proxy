<script lang="ts">
	import {
		Card,
		CardHeader,
		CardTitle,
		CardDescription,
		CardContent,
		CardFooter
	} from '../ui/card';
	import { Button } from '../ui/form';
	import { Input } from '../ui/input';
	import { Label } from '../ui/label';
	import { enhance } from '$app/forms';
	import type { ActionData } from '../../../routes/(dashboard)/dashboard/$types';

	let { createForm }: { createForm: ActionData } = $props();
</script>

<form class="mt-4 flex items-center font-bold" action="?/create" method="POST" use:enhance>
	<Card class="max-w-sm">
		<CardHeader>
			<CardTitle>Create your subdomain</CardTitle>
			<CardDescription>
				To get started, create your tenant id. You must use a unique identifier. Once created, you
				will be able to access it.
			</CardDescription>
		</CardHeader>
		<CardContent>
			<Label for="tenantid">Tenant Id</Label>
			<Input placeholder="tenant id" class="my-2" name="tenantid" id="tenantid" />
			<p class="text-muted-foreground text-xs font-medium">
				Must contain only alphanumeric characters and hyphens
			</p>
			{#if createForm?.message}
				<p class="mt-4 text-xs text-red-500">
					{createForm.message}
				</p>
			{/if}
		</CardContent>
		<CardFooter>
			<Button class="w-full" type="submit">Create</Button>
		</CardFooter>
	</Card>
</form>
