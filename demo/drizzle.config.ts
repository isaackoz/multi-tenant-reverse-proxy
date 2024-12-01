import { defineConfig } from 'drizzle-kit';
// if (!process.env.DATABASE_URL) throw new Error('DATABASE_URL is not set');

export default defineConfig({
	schema: './src/lib/server/db/schema.ts',

	dbCredentials: {
		url: 'postgresql://postgres:password@0.0.0.0:5444/postgres'
	},

	verbose: true,
	strict: true,
	dialect: 'postgresql'
});
