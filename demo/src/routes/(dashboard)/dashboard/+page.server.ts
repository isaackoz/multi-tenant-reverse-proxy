import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import {
	createUserTenantId,
	getTenantIdExist,
	getUserTenant,
	setHostname,
	updatePage
} from '$lib/server/data-access';
import { invalidateTenantProxyCache } from '$lib/server/data-access/proxy';

export const load: PageServerLoad = async ({ locals }) => {
	const userId = locals?.session?.claims?.sub;
	if (!userId) {
		throw redirect(301, '/login');
	}

	// Get the current data
	const { data } = await getUserTenant(userId);

	return {
		tenantData: data
	};
};

export const actions = {
	create: async (event) => {
		const formdata = await event.request.formData();
		const tenantId = formdata.get('tenantid');

		if (!tenantId || typeof tenantId !== 'string') {
			return fail(401, { message: 'Tenant Id required' });
		}

		const userId = event.locals?.session?.claims?.sub;
		if (!userId) {
			return fail(401, { message: 'unauthorized' });
		}

		// Check if user already has a tenant id
		const { data: userData, error: userError } = await getUserTenant(userId);
		if (userError) {
			return fail(500, { message: userError.message });
		}

		if (userData) {
			return fail(401, { message: 'User already has a tenant Id' });
		}

		// Check if already exists
		const { data, error } = await getTenantIdExist(tenantId);
		if (error) {
			return fail(500, { message: error.message });
		}

		if (data === true) {
			return fail(401, { message: 'Tenant Id already taken' });
		}

		// If user does not have a tenant id AND the requested one doesn't exist, create it for them
		const { error: createErr } = await createUserTenantId(tenantId, userId);
		if (createErr) {
			return fail(500, { message: createErr.message });
		}

		return { success: true };
	},
	isTenantIdTaken: async (event) => {
		const userId = event.locals?.session?.claims?.sub;
		if (!userId) {
			return fail(401, { message: 'unauthorized' });
		}
		const { data, error } = await getTenantIdExist(userId);
		if (error) {
			return fail(500, { message: error.message });
		}
		return { success: true, data: data };
	},
	update: async (event) => {
		const userId = event.locals?.session?.claims?.sub;
		if (!userId) {
			return fail(401, { message: 'unauthorized' });
		}

		const formData = await event.request.formData();
		const title = formData.get('title');
		const message = formData.get('message');

		if (!title || !message || typeof title !== 'string' || typeof message !== 'string') {
			return fail(401, { message: 'Title and message required' });
		}

		const { data: userData, error: userError } = await getUserTenant(userId);

		if (userError || !userData) {
			return fail(500, { message: userError?.message ?? 'Unknown error' });
		}

		const { error } = await updatePage(userData.id, userId, title, message);

		if (error) {
			return fail(500, { message: error.message });
		}

		return {
			success: true
		};
	},
	updateHostname: async (event) => {
		const userId = event.locals?.session?.claims?.sub;
		if (!userId) {
			return fail(401, { hostmessage: 'unauthorized' });
		}

		const formData = await event.request.formData();
		const hostname = formData.get('hostname');
		if (!hostname || typeof hostname !== 'string' || hostname.length === 0) {
			return fail(401, { hostmessage: 'Hostname required' });
		}

		if (hostname.startsWith('https://') || hostname.startsWith('http://')) {
			return fail(401, { hostmessage: 'Hostname must not contain protocol (http:// or https://)' });
		}

		const { data: userData, error: userError } = await getUserTenant(userId);

		if (userError || !userData) {
			return fail(500, { message: userError?.message ?? 'Unknown error' });
		}

		// Invalidate data in the proxy
		if (userData.hostname) {
			const { error: invalError } = await invalidateTenantProxyCache(userData.hostname);
			if (invalError) {
				// We could handle this, we are just going to ignore it for now
				console.warn(
					'There was an error invalidating the tenant cache in the proxy: ',
					invalError.message
				);
			}
		}

		const { error: hostErr } = await setHostname(userData.id, userId, hostname);

		if (hostErr) {
			return fail(500, { hostmessage: hostErr.message });
		}

		return {
			hostsuccess: true
		};
		//
	}
} satisfies Actions;
