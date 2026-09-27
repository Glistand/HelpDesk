import { NextResponse } from "next/server";
import { AUTH_COOKIE, USER_COOKIE, login } from "@/lib/api";

export async function POST(req: Request) {
  let body: { email?: string; password?: string };
  try {
    body = await req.json();
  } catch {
    return NextResponse.json({ error: "invalid json" }, { status: 400 });
  }
  if (!body.email || !body.password) {
    return NextResponse.json({ error: "email and password required" }, { status: 400 });
  }
  try {
    const result = await login(body.email, body.password);
    const res = NextResponse.json({ user: result.user });
    const cookieOpts = {
      sameSite: "lax" as const,
      path: "/",
      maxAge: 60 * 60 * 12,
    };
    res.cookies.set(AUTH_COOKIE, result.access_token, {
      ...cookieOpts,
      httpOnly: true,
    });
    res.cookies.set(USER_COOKIE, JSON.stringify(result.user), {
      ...cookieOpts,
      httpOnly: true,
    });
    return res;
  } catch {
    return NextResponse.json({ error: "invalid email or password" }, { status: 401 });
  }
}
