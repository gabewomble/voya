import type { RequestEventBase } from "@builder.io/qwik-city";
import { setCookie } from "./set-cookie";
import { ErrorResponseSchema } from "~/types/server-errors";
import { SERVER_ERROR_MESSAGES } from "~/constants/server-errors";

export async function serverFetch(
  path: string,
  options: RequestInit,
  ctx: RequestEventBase,
) {
  let token = ctx.cookie.get("token")?.value;
  const auth: Record<string, string> = {};

  if (token) {
    auth.Authorization = `Bearer ${token}`;
  }

  const url = new URL(path, ctx.env.get("API_URL"));

  let res = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
      ...auth,
    },
  });

  if (res.status === 401 && token) {
    const errorResponse = ErrorResponseSchema.safeParse(await res.json());
    if (!errorResponse.success) return res;

    const isAuthError = errorResponse.data.errors.some(
      (e) =>
        e.message === SERVER_ERROR_MESSAGES.INVALID_CREDENTIALS ||
        e.message === SERVER_ERROR_MESSAGES.MISSING_CREDENTIALS,
    );

    if (!isAuthError) return res;

    // Try to refresh the token
    let refreshToken = ctx.cookie.get("refreshToken")?.value;

    if (!refreshToken) return res;

    const refreshRes = await fetch(`${ctx.env.get("API_URL")}/refresh`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ refreshToken }),
    });

    if (!refreshRes.ok) return res;

    const json = await refreshRes.json();
    token = json.token;
    refreshToken = json.refresh_token;

    if (token) setCookie("token", token, ctx);
    if (refreshToken) setCookie("refreshToken", refreshToken, ctx);

    // Retry the original request
    res = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
        Authorization: `Bearer ${token}`,
      },
    });
  }

  return res;
}
