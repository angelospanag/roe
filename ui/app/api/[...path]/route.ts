import { type NextRequest, NextResponse } from "next/server";

// The Go API's routes are unprefixed (e.g. /feeds, /posts); this proxy owns
// the /api prefix so the browser only ever talks to relative /api/* paths.
const BACKEND = process.env.API_URL ?? "http://localhost:8080";

async function proxy(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> },
) {
  const { path } = await params;
  const url = `${BACKEND}/${path.join("/")}${request.nextUrl.search}`;

  const init: RequestInit = {
    method: request.method,
    headers: {
      "Content-Type": request.headers.get("content-type") ?? "application/json",
    },
  };
  if (request.method !== "GET" && request.method !== "HEAD") {
    init.body = await request.arrayBuffer();
  }

  try {
    const upstream = await fetch(url, init);
    // 204/304 responses must not carry a body.
    const body = upstream.status === 204 ? null : await upstream.arrayBuffer();
    return new NextResponse(body, {
      status: upstream.status,
      headers: {
        "Content-Type":
          upstream.headers.get("content-type") ?? "application/json",
      },
    });
  } catch {
    return NextResponse.json(
      { error: "upstream unreachable" },
      { status: 502 },
    );
  }
}

export {
  proxy as GET,
  proxy as POST,
  proxy as PUT,
  proxy as PATCH,
  proxy as DELETE,
};
