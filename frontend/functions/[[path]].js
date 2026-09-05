export async function onRequest(context) {
  // Extract the incoming request from the frontend
  const request = context.request;
  
  // Extract the URL path segments after /api/
  const pathString = context.params.path ? context.params.path.join('/') : '';
  
  // Preserve any query parameters (like ?user_id=123)
  const searchParams = new URL(request.url).search; 
  
  // IMPORTANT: Replace this with your actual EC2 IP and port
  const backendUrl = `http://98.85.29.80:8080/${pathString}${searchParams}`;
  
  // Clone the original request to preserve headers, methods (POST, GET), and the body payload
  const backendRequest = new Request(backendUrl, {
    method: request.method,
    headers: request.headers,
    body: request.body,
    redirect: "manual"
  });
  
  // Fetch the data directly from the EC2 instance and return it to the frontend
  try {
    const response = await fetch(backendRequest);
    return response;
  } catch (error) {
    return new Response("Error communicating with backend microservice", { status: 502 });
  }
}