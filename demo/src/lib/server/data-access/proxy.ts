import { PROXY_API_URL, PROXY_AUTH_TOKEN } from '$env/static/private';
import type { TApiData } from './types';

export const invalidateTenantProxyCache = async (hostname: string): Promise<TApiData<null>> => {
	try {
		const res = await fetch(`${PROXY_API_URL}/invalidate?hostname=${hostname}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${PROXY_AUTH_TOKEN}`
			}
		});
		if (!res.ok) {
			return {
				error: {
					message: `Error while invalidating cache: ${await res.text()}`
				}
			};
		}

		return {
			data: null
		};
	} catch (e) {
		console.warn(e);
		return {
			error: {
				message: 'Unknown error while invalidating cache'
			}
		};
	}
};
