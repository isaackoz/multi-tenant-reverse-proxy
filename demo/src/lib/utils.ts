import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function isValidUrl(str: string) {
	const regex = /^[a-z0-9-]+$/i; // Matches only a-z, 0-9, and the minus sign
	return regex.test(str);
}

export function hostnameToValidUrl(hostname: string) {
	return hostname.toLowerCase().replace(/\./g, '-'); // Replace periods with dashes
}
