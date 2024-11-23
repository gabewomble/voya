import type { RequestEventBase } from "@builder.io/qwik-city";

export function setCookie(key: string, value: string, ctx: RequestEventBase) {
  ctx.cookie.set(key, value, {
    path: "/",
    httpOnly: true,
    sameSite: true,
    secure: false,
  });
}
