import { and, eq } from 'drizzle-orm';
import { db } from '../db';
import { tenantsTable } from '../db/schema';
import type { TApiData } from './types';

export const getTenantIdExist = async (tenantId: string): Promise<TApiData<boolean>> => {
	try {
		const tenant = (await db.select().from(tenantsTable).where(eq(tenantsTable.id, tenantId)))[0];

		if (!tenant) {
			return {
				data: false
			};
		} else {
			return {
				data: true
			};
		}
	} catch {
		return {
			error: {
				message: 'Error while fetching tenant'
			}
		};
	}
};

export const getUserTenant = async (
	userId: string
): Promise<TApiData<typeof tenantsTable.$inferSelect | null>> => {
	try {
		const data = (await db.select().from(tenantsTable).where(eq(tenantsTable.userId, userId)))[0];
		if (!data) {
			return {
				data: null
			};
		} else {
			return {
				data: data
			};
		}
	} catch {
		return {
			error: {
				message: 'Error while fetching user tenant'
			}
		};
	}
};

export const createUserTenantId = async (
	tenantId: string,
	userId: string
): Promise<TApiData<undefined>> => {
	try {
		await db.insert(tenantsTable).values({
			id: tenantId,
			userId: userId
		});
		return {
			data: undefined
		};
	} catch {
		return {
			error: {
				message: 'Error while creating user tenant id'
			}
		};
	}
};

export const updatePage = async (
	tenantId: string,
	userId: string,
	title: string,
	message: string
): Promise<TApiData<null>> => {
	try {
		await db
			.update(tenantsTable)
			.set({
				message: message,
				title: title
			})
			.where(and(eq(tenantsTable.id, tenantId), eq(tenantsTable.userId, userId)));
		return {
			data: null
		};
	} catch {
		return {
			error: {
				message: 'Error while updating page'
			}
		};
	}
};

export const getTenantPageData = async (
	tenantId: string
): Promise<
	TApiData<{
		title: string;
		message: string;
	}>
> => {
	try {
		const data = (
			await db
				.select({
					title: tenantsTable.title,
					message: tenantsTable.message
				})
				.from(tenantsTable)
				.where(eq(tenantsTable.id, tenantId))
		)[0];
		if (!data) {
			return {
				error: {
					message: 'Tenant not found'
				}
			};
		}
		return {
			data: data
		};
	} catch {
		return {
			error: {
				message: 'Error while fetching tenant page data'
			}
		};
	}
};

export const setHostname = async (
	tenantId: string,
	userId: string,
	hostname: string
): Promise<TApiData<null>> => {
	try {
		await db
			.update(tenantsTable)
			.set({
				hostname: hostname
			})
			.where(and(eq(tenantsTable.id, tenantId), eq(tenantsTable.userId, userId)));

		return { data: null };
	} catch {
		return {
			error: {
				message: 'Error while updating hostname'
			}
		};
	}
};
