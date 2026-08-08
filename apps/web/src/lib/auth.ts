import { cookies } from "next/headers";

const ACCESS_TOKEN_COOKIE = "vakalat_token";
const REFRESH_TOKEN_COOKIE = "vakalat_refresh_token";

// Access tokens are short-lived (backend default: 24h) — the cookie mirrors that
// so the browser doesn't hold on to a cookie the server no longer honors.
const ACCESS_TOKEN_MAX_AGE = 60 * 60 * 24;
const REFRESH_TOKEN_MAX_AGE = 60 * 60 * 24 * 7;

export async function getAuthToken(): Promise<string | undefined> {
  const cookieStore = await cookies();
  return cookieStore.get(ACCESS_TOKEN_COOKIE)?.value;
}

export async function getRefreshToken(): Promise<string | undefined> {
  const cookieStore = await cookies();
  return cookieStore.get(REFRESH_TOKEN_COOKIE)?.value;
}

export async function setAuthCookies(token: string, refreshToken: string): Promise<void> {
  const cookieStore = await cookies();
  const shared = {
    httpOnly: true,
    // NOT tied to NODE_ENV — "production" here just means "built for
    // production," not "served over HTTPS." This deployment currently has no
    // TLS termination (docker-compose serves plain http://localhost:3000).
    // secure:true on plain HTTP means the browser silently drops the cookie,
    // which makes every subsequent request look logged-out. Flip this on
    // (or drive it from an env var) once a real TLS/ingress layer is in front.
    secure: process.env.COOKIE_SECURE === "true",
    sameSite: "lax" as const,
    path: "/",
  };
  cookieStore.set(ACCESS_TOKEN_COOKIE, token, { ...shared, maxAge: ACCESS_TOKEN_MAX_AGE });
  cookieStore.set(REFRESH_TOKEN_COOKIE, refreshToken, { ...shared, maxAge: REFRESH_TOKEN_MAX_AGE });
}

export async function clearAuthCookies(): Promise<void> {
  const cookieStore = await cookies();
  cookieStore.delete(ACCESS_TOKEN_COOKIE);
  cookieStore.delete(REFRESH_TOKEN_COOKIE);
}

export { ACCESS_TOKEN_COOKIE };
