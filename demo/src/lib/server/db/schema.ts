import { pgTable, text } from 'drizzle-orm/pg-core';

// Normally we'd have a userTable and use foreign keys, but for simplicity we will de-normalize the userId
// to avoid having to sync data with endpoints and clerk, etc.
// userid is whatever the auth id is (clerk id in this case)
export const tenantsTable = pgTable('tenants', {
	id: text('id').primaryKey(),
	hostname: text('hostname').unique(),
	userId: text('user_id').unique().notNull(),
	title: text('title').notNull().default('Hello world'),
	message: text('message').notNull().default('lorem ipsum')
});
