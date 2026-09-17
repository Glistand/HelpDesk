import { NextResponse } from "next/server";
import { AUTH_COOKIE, login } from "@/lib/api";

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
    res.cookies.set(AUTH_COOKIE, result.access_token, {
      httpOnly: true,
      sameSite: "lax",
      path: "/",
      maxAge: 60 * 60 * 12,
    });
    return res;
  } catch {
    return NextResponse.json({ error: "invalid email or password" }, { status: 401 });
  }
}
