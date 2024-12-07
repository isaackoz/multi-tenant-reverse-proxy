# Reverse Dynamic Proxy Proof of Concept
---
## What is it?
A reverse dynamic proxy proof of concept for multi-tenant applications. To explain it better, let's use an example of a web application that hosts user's blogs. Each user is able to create their own blogs through the application's dashboard `example.com/dashboard`. When a user wants to publish their blog, it is accessible on a subdomain with their tenant id, `user1.example.com`. 

The web application is then able to grab the tenant id by parsing anything before it's domain name. In this case it would be `user1` as the tenant id. The web application is then able to use this tenant id to do *anything*. In the Svelte demo, we rewrite the path under the hood in `/demo/src/hooks.ts` to convert `user1.example.com` to `example.com/app/user1`, where `user1` is a dynamic route param that we can use on the server to retrieve data using that tenant id. For example, we might run a SQL query (pseduo) of `SELECT * FROM tenant_blogs WHERE tenant_blogs.tenantid = 'user1'`.

This is great and doesn't require any additional setup to get working. But what if you want users to be able to use their own custom domain for their blog. I.e. `userone.com` -> `user1.example.com`. For this to work, we need to have the following work:
1. The domain is `userone.com`
2. `userone.com` has an SSL cert (https)
3. The application is able to parse a unique tenant id from this domain. i.e. `userone.com` should resolve to the tenant id of `user1`
4. All traffic to `userone.com` is directed to our applications server.


The implementation of this will depend on the application being developed— hence why this is a proof of concept. A library or *one-size fits all* solution to this is not viable. [KISS](https://en.wikipedia.org/wiki/KISS_principle). Below is one way (and what we use) to get the tenant id.

## 1. Custom Domain
In order for the user to use a custom domain, the first thing the application would likely instruct them to do is redirect all traffic to the applications server. This can be done by setting an `A` record. Let's say that our web applications server ip is `123.456.1.1`. We would instruct the user to set the following DNS records for their arbitrary domain:
- Target: `userone.com`, Record Type `A`, Value: `123.456.1.1`, TTL `3600`

This would then direct all traffic from `userone.com` to our server at `123.456.1.1`.

## 2. Generating Dynamic SSL Certficates
The next thing we need to setup is some sort of proxy and generate an SSL certificate for `userone.com`. Manually doing this is trivial. But how do we do it dynamically on-the-fly? We could write our own proxy server... or we could use one with this feature built in. Luckily, it doesn't get any easier than using Caddy. Caddy has a feature called [*on demand TLS*](https://caddyserver.com/on-demand-tls) which will let us automatically obtain a certificate for custom domains.

The way this works is that we will basically say, "Caddy, any traffic that comes from an external domain, check if they are a tenant. If they are, generate a certificate for them on the fly. If they aren't deny their request." This can be translated to a Caddyfile that looks like the following:

```Caddyfile
{
        on_demand_tls {
                ask http://localhost:3000/ask
        }
}

# Handle all traffic to our application
*.example.com, example.com {
        tls {
                dns digitalocean <digital-ocean-api-key>
        }

        handle {
                reverse_proxy 127.0.0.1:3001
        }
}


# For any traffic that isn't directly to our application
https:// {
        tls {
                on_demand
        }
        reverse_proxy 127.0.0.1:3001 {
                header_up Host {host}
        }
}
```

In this Caddyfile, we make it so that any traffic that isn't `*.example.com` has to first check if it's allowed via `on_demand_tls`. The way this works is that Caddy will query `http://localhost:3000/ask?domain=userone.com` and check for a 2xx response. 

This means we will need some sort of microservice that will check our database if the user is allowed. This microservice should ideally respond within 5-10ms, with the exception of the first query. We handle all of this inside `/dynamic-proxy`. In a nutshell, the way it works is like this:
1. Receive a request. Check who the host is. In this case it would be `userone.com`.
2. We check our layered cache for this value.
   1. First we check for a key:value in memory. I.e. `tenant-userone.com:user1`
   2. If that is null, we then check Redis
   3. If redis is null, we then check Postgres
   4. In any of the steps, if a value is found, we pass it back up the chain and save it for any future requests. This means, the first request will likely check Postgres before returning a value. This means the first request will take longer, ~100ms. It will then store the result in Redis and local memory. Since it's in memory, the next request will be near instant, 5-10ms.
3. Upon a value found, we return a 2xx response to Caddy. If it's not found we return a 4xx or 5xx response.

Alas, we succesfully are able to route traffic from an external domain, with SSL, to our application.

## 3. Getting the Tenant Id on the Application
The last step is getting the tenant id on the application. There are a couple of ways to handle this. In our `/demo`, we handle this by converting the hostname into a URL safe string and append that to the domain name as the wildcard. In other words, `userone.com` would be parsed as `userone-com` which is then used to, under the hood, set the path as `example.com/app/userone-com. This likely isn't the most ideal way to identify tenants, but it works at least for this demo. Ideally, in SvelteKit we would add a local param with the tenant id in there. But I handled this slightly different.

Either way, you now have access to the tenant id available via `params` (or `locals` if you configure it that way) in page/layout loads and can use this value in database calls.

Again, this is a proof of concept, so the implementation will vary depending on requirements. 

## Things To Note
- This is a proof of concept. It can be adapted to any application.
- You might consider isolating the dynamic-proxy and make it use its own database and expose an API to manage it.
- You might want to handle setting/getting the tenant ID differently than converting it from `userone.com` -> `userone-com`. You would likely want to make the tenant id static, and the hostname dynamic. Currently, the hostname is the source of truth for determining the tenant id. The tenant id column in the database doesn't do anything, although it should probably be used.
