import { getTenantPageData } from '$lib/server/data-access';
import type { LayoutServerLoad } from './$types';
import { error } from '@sveltejs/kit';

export const load: LayoutServerLoad = async ({ params }) => {
	// Get this tenants
	const { data, error: dataErr } = await getTenantPageData(params.tenant);
	if (dataErr) {
		error(404, {
			message: 'Not found!'
		});
	}
	return {
		tenantId: params.tenant,
		data
	};
};
