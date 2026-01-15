import type { NextRequest } from "next/server"
import { NextResponse } from "next/server"

export const middleware = (req: NextRequest) => {
  const start = Date.now()
  const response = NextResponse.next()

  console.log(
    `[${new Date().toISOString()}] ${req.method} ${req.nextUrl.pathname} - ${Date.now() - start}ms`
  )

  return response
}

export const config = {
  matcher: ["/((?!_next/static|_next/image|favicon.ico).*)"],
}
