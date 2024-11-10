import type { Cookie } from "@builder.io/qwik-city";

type Context = {
  cookie: Cookie;
};

export function serverFetch(path: string, options: RequestInit, ctx: Context) {
  const token = ctx.cookie.get("token")?.value;
  const auth: Record<string, string> = {};

  if (token) {
    auth.Authorization = `Bearer ${token}`;
  }

  const url = new URL(path, `http://localhost:8080`);

  return fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
      ...auth,
    },
  });
}
