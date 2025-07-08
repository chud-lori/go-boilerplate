import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  // Adjust these options based on your load testing goals
  vus: 10,       // 10 virtual users
  duration: '30s', // for 30 seconds

  // Define thresholds for success/failure criteria
  thresholds: {
    'http_req_duration': ['p(95)<700'], // 95% of requests must complete within 700ms (adjusted slightly for multiple requests)
    'http_req_failed': ['rate<0.01'],    // less than 1% of requests should fail
  },
};

export default function () {
  const baseUrl = 'http://localhost:1234/api/post';

  // --- 1. POST Request (Create Post) ---
  const postHeaders = {
    'accept': 'application/json',
    'X-API-KEY': 'secret-api-key',
    'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTIwNDU2NDQsInVzZXJfaWQiOiIxYWRlMGU5NS1kODZjLTQ2ZTQtYmQ4Ny04Njc3ZTBmZGI0Y2YifQ.hKWQdURIIMnC58Do-Oyup8m78XJyI_eDu9NXKrjgIUQ',
    'Content-Type': 'application/json',
  };

  const postPayload = JSON.stringify({
    "author_id": "1ade0e95-d86c-46e4-bd87-8677e0fdb4cf",
    "body": "this is body " + __VU + "-" + __ITER, // Dynamically change body for unique posts
    "title": "this is title " + __VU + "-" + __ITER
  });

  // console.log(`[VU:${__VU} ITER:${__ITER}] Sending POST request...`); // Keep this for debugging if needed
  const postRes = http.post(baseUrl, postPayload, { headers: postHeaders });

  check(postRes, {
    'POST status is 201': (r) => r.status === 201, // Changed to 201
    // Corrected check: Access 'id' inside the 'data' object
    'POST response has id': (r) => r.json() && r.json().data && r.json().data.id !== undefined,
  });

  // Add a sleep after creating the post before fetching
  sleep(1);

  // --- 2. First GET Request (Get Posts) ---
  const getHeaders = {
    'accept': 'application/json',
    'X-API-KEY': 'secret-api-key',
  };

  // console.log(`[VU:${__VU} ITER:${__ITER}] Sending first GET request...`); // Keep for debugging
  const getRes1 = http.get(baseUrl, { headers: getHeaders });

  check(getRes1, {
    'GET 1 status is 200': (r) => r.status === 200,
    // Corrected check: Access the array inside the 'data' object
    'GET 1 response is array': (r) => r.json() && Array.isArray(r.json().data),
  });

  // Optional: Add a small sleep between the two GET requests
  // sleep(0.5);

  // --- 3. Second GET Request (Get Posts - Identical) ---
  // console.log(`[VU:${__VU} ITER:${__ITER}] Sending second GET request...`); // Keep for debugging
  const getRes2 = http.get(baseUrl, { headers: getHeaders });

  check(getRes2, {
    'GET 2 status is 200': (r) => r.status === 200,
    // Corrected check: Access the array inside the 'data' object
    'GET 2 response is array': (r) => r.json() && Array.isArray(r.json().data),
  });

  // Add a general sleep at the end of the iteration to simulate user think time
  sleep(2);
}
