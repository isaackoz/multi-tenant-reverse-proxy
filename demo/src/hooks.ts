import { PUBLIC_DOMAIN } from '$env/static/public';
import { hostnameToValidUrl } from '$lib/utils';
import { type Reroute } from '@sveltejs/kit';

export const reroute: Reroute = ({ url }) => {
	console.log('Incoming URL', url);

	// Host = host + port
	// Hostname = host <--- use this one
	// Port = port

	//
	// If we try to visit the main website, without a subdomain, we return the route unchanged
	if (url.hostname === PUBLIC_DOMAIN) {
		// optionally, we don't allow for viewing subdomain pages from the main site
		if (url.pathname.startsWith('/app/')) {
			return '/not-found';
		}
		return url.pathname;
	}

	// Here we have two choices, if the user is trying to access their page at hostname.localhost.test, then

	// If the hostname is example.com and we receive a request from
	// example.com, the output should be /app/example-com
	// if we receive on from example-com.localhost, the output should be the same, /app/example-com

	// If the url is from an external source (in our case, a custom custom), then we want to add the hostname to the url
	if (!url.hostname.endsWith(PUBLIC_DOMAIN)) {
		// test.com --> test-com        example.test.com ---> example-test-com
		const tenantId = hostnameToValidUrl(url.hostname);

		return `/app/${tenantId}${url.pathname}`; // --> /app/test-com?hello=world
	}

	// If the hostname is not main website, get it's subdomain which should be whatever is to the left of PUBLIC_DOMAIN
	const tenantId = url.hostname.substring(0, url.hostname.lastIndexOf(PUBLIC_DOMAIN) - 1);
	// Optionally ensure that the subdomain doesn't have a sub-subdomain
	if (tenantId.includes('.')) {
		return '/app/error';
	}

	// At this point tenantId should be valid and clean with all edge cases handled

	// Finally we return all subdomain traffic to our app/[tenant] route. Everything under [tenant]/* will have access to params.tenant
	// which is the tenantId we extracted here. You can use this tenant Id in database queries to get the tenant specific info.
	//
	// Don't forget to include url.pathname!
	//
	/* 

  To get a better idea of how this works, check the table out below
  If PUBLIC_DOMAIN = `example.com` in .env

  | -------------------------------- | --------------------------------- |
  | Url the user sees in the browser | pathname handled in svelte-kit    |
  | ---------------------------------| --------------------------------- |
  | example.com                      | /                                 |
  | example.com/dashboard            | /dashboard                        |
  | example.com/app/test             | /not-found                        | <-- we (optionally) handle this above on line 9
  | test.example.com                 | /app/test/                        |
  | hello.example.com                | /app/error/                       | <-- contains an sub-sub-domain which we handle above on line 24
  | user1.example.com                | /app/user1/                       |
  | user1.example.com/blog/1         | /app/user1/blog/1                 |
  | user1.hello.com                  | /not-found                        | <-- handled on line 17 (edge case)
  | -------------------------------- | --------------------------------- |
  */

	return `/app/${tenantId}${url.pathname}`;
};
