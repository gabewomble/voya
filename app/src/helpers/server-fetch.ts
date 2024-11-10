import type { RequestEventBase } from "@builder.io/qwik-city";

export function serverFetch(
  path: string,
  options: RequestInit,
  ctx: RequestEventBase,
) {
  const token = ctx.cookie.get("token")?.value;
  const auth: Record<string, string> = {};

  if (token) {
    auth.Authorization = `Bearer ${token}`;
  }

  const url = new URL(path, ctx.env.get("API_URL"));

  return fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
      ...auth,
    },
  });
}
