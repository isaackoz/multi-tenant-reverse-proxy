<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import {
		Card,
		CardContent,
		CardDescription,
		CardFooter,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
</script>

<div class="mx-auto flex max-w-7xl flex-col items-center justify-center pt-24">
	<h1 class="text-4xl font-bold">Dashboard</h1>
	<div class="mt-4 flex items-center font-bold">
		<Card class="max-w-sm">
			<CardHeader>
				<CardTitle>Create your subdomain</CardTitle>
				<CardDescription>
					To get started, create your own subdomain. You must use a unique identifier. Once created,
					you will be able to access it.
				</CardDescription>
			</CardHeader>
			<CardContent>
				<Label>Tenant Id</Label>
				<Input placeholder="tenant id" class="my-2" />
				<p class="text-muted-foreground text-xs font-medium">
					Must contain only alphanumeric characters and hyphens
				</p>
			</CardContent>
			<CardFooter>
				<Button class="w-full">Create</Button>
			</CardFooter>
		</Card>
	</div>
</div>
