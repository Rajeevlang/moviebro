export default {
  async fetch(request, env, ctx) {
    const url = new URL(request.url);

    // 1. Forward API and healthcheck requests directly to your EC2 backend
    if (url.pathname.startsWith("/api/") || url.pathname === "/healthcheck") {
      // Set your EC2 Public IP and API Gateway port here
      const backendHost = env.BACKEND_HOST || "http://YOUR_EC2_PUBLIC_IP:8080";
      const backendUrl = `${backendHost}${url.pathname}${url.search}`;

      const backendRequest = new Request(backendUrl, {
        method: request.method,
        headers: request.headers,
        body: request.body,
        redirect: "manual"
      });

      try {
        const response = await fetch(backendRequest);
        return response;
      } catch (err) {
        return new Response(
          JSON.stringify({ error: "Backend service unreachable", details: err.message }),
          {
            status: 502,
            headers: { "Content-Type": "application/json" }
          }
        );
      }
    }

    // 2. Serve static frontend assets (SPA) from the public/ directory
    let response = await env.ASSETS.fetch(request);

    // SPA client-side routing fallback (for paths like /login, /admin, etc.)
    if (response.status === 404 && !url.pathname.includes(".")) {
      const indexRequest = new Request(new URL("/index.html", request.url), request);
      response = await env.ASSETS.fetch(indexRequest);
    }

    return response;
  }
};
