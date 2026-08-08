import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

const ACCESS_TOKEN_COOKIE = "vakalat_token";

// Guards authenticated sections of the app. This only checks that a session
// cookie is present — it does not verify the JWT signature (that happens on
// every backend request regardless). Its job is to bounce obviously signed-out
// visitors to /login before they see a page full of failed requests; an
// expired-but-present token still reaches the page and surfaces a normal
// "please sign in again" error from the API layer.
export function proxy(request: NextRequest) {
  const token = request.cookies.get(ACCESS_TOKEN_COOKIE)?.value;

  if (!token) {
    const loginUrl = new URL("/login", request.url);
    loginUrl.searchParams.set("next", request.nextUrl.pathname);
    return NextResponse.redirect(loginUrl);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/matters/:path*", "/clients/:path*"],
};
