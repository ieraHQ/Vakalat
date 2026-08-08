"use server";

import { redirect } from "next/navigation";
import { ApiError, login } from "@/lib/api";
import { setAuthCookies, clearAuthCookies } from "@/lib/auth";

export interface LoginState {
  error?: string;
}

export async function loginAction(
  _prevState: LoginState,
  formData: FormData
): Promise<LoginState> {
  const email = String(formData.get("email") ?? "").trim();
  const password = String(formData.get("password") ?? "");
  const next = String(formData.get("next") ?? "/");

  if (!email || !password) {
    return { error: "Enter both your email and password." };
  }

  let token: string;
  let refreshToken: string;
  try {
    const result = await login(email, password);
    token = result.token;
    refreshToken = result.refresh_token;
  } catch (err) {
    if (err instanceof ApiError && err.status === 401) {
      return { error: "Incorrect email or password." };
    }
    return { error: "Could not reach the server. Please try again." };
  }

  await setAuthCookies(token, refreshToken);
  redirect(next.startsWith("/") ? next : "/");
}

export async function logoutAction(): Promise<void> {
  await clearAuthCookies();
  redirect("/login");
}
