// Signup Worker: the only code that runs on the Production Target
// (ADR-0004). It is a port of the Go subscribe handler (subscribe/), which
// still serves the QA Target; the two must keep the same contract:
//
//   POST /api/subscribe  email, source (page path, sent as a Buttondown tag),
//                        website (honeypot)
//   - no-JS form post            -> 303 Location: /subscribed/
//   - Accept: application/json   -> {"ok": bool, "message": "..."}
//   - 400 invalid email, 429 per-IP rate limit, 502 upstream failure or
//     missing key, 405 for any other method; honeypot hits and
//     already-subscribed emails get the normal success answer.
//
// Every other request goes to the static assets binding. On Production the
// Workers config sets run_worker_first to the signup path only, so assets
// are served without invoking this script; the fallthrough below covers
// asset misses and `wrangler dev`.
//
// Environment (Worker secret/var, same names as the Go function):
//   BUTTONDOWN_API_KEY   secret, required. Sent upstream only, never echoed.
//   BUTTONDOWN_BASE_URL  var, optional, default https://api.buttondown.com.
//                        scripts/check-worker.sh points it at a stub.
//
// Verified by scripts/check-worker.sh (wrangler dev + stub Buttondown), not
// by a JavaScript test suite: tests stay Go-only.

const PATH = "/api/subscribe";
const SUBSCRIBED_PATH = "/subscribed/";
const HONEYPOT_FIELD = "website";
const DEFAULT_BASE_URL = "https://api.buttondown.com";

// Rate limit: 5 requests per IP per 10 minutes, as in the Go handler.
// The counters live in this isolate's memory: every Cloudflare location
// and every isolate has its own copy, and an eviction resets it. This is
// best-effort abuse damping, not a quota.
const RATE_LIMIT = 5;
const RATE_WINDOW_MS = 10 * 60 * 1000;
const MAX_TRACKED = 10_000;
const MAX_BODY_BYTES = 16 << 10;
const UPSTREAM_TIMEOUT_MS = 8000;

const MSG_SUBSCRIBED = "Thanks! Check your inbox to confirm your subscription.";
const MSG_INVALID_EMAIL = "Please enter a valid email address.";
const MSG_RATE_LIMITED = "Too many attempts from your network. Please try again in a few minutes.";
const MSG_UPSTREAM = "We couldn't reach our newsletter service. Please try again in a few minutes.";

/** @type {Map<string, {start: number, count: number}>} */
const hits = new Map();

// allow records one request from ip and reports whether it is within the limit.
function allow(ip, now) {
  if (hits.size >= MAX_TRACKED) {
    for (const [k, w] of hits) {
      if (now - w.start >= RATE_WINDOW_MS) hits.delete(k);
    }
  }
  let w = hits.get(ip);
  if (!w || now - w.start >= RATE_WINDOW_MS) {
    w = { start: now, count: 0 };
    hits.set(ip, w);
  }
  w.count++;
  return w.count <= RATE_LIMIT;
}

export default {
  async fetch(request, env) {
    const url = new URL(request.url);
    if (url.pathname === PATH) {
      return subscribe(request, env);
    }
    if (env.ASSETS) {
      return env.ASSETS.fetch(request);
    }
    return new Response("not found\n", { status: 404 });
  },
};

async function subscribe(request, env) {
  if (request.method !== "POST") {
    return new Response("method not allowed\n", {
      status: 405,
      headers: {
        Allow: "POST",
        "Cache-Control": "no-store",
        "Content-Type": "text/plain; charset=utf-8",
        "X-Content-Type-Options": "nosniff",
      },
    });
  }
  const form = await readForm(request);
  const json = wantsJSON(request);
  if ((form.get(HONEYPOT_FIELD) ?? "") !== "") {
    // A bot filled the field humans never see: pretend it worked.
    return succeed(json);
  }
  const email = String(form.get("email") ?? "").trim();
  const source = cleanSource(String(form.get("source") ?? ""));
  const ip = clientIP(request);
  if (!allow(ip, Date.now())) {
    return fail(json, 429, MSG_RATE_LIMITED, source, { "Retry-After": String(RATE_WINDOW_MS / 1000) });
  }
  if (!validEmail(email)) {
    return fail(json, 400, MSG_INVALID_EMAIL, source);
  }
  const apiKey = env.BUTTONDOWN_API_KEY ?? "";
  if (apiKey === "") {
    console.log("subscribe: BUTTONDOWN_API_KEY is not set");
    return fail(json, 502, MSG_UPSTREAM, source);
  }
  const baseURL = (env.BUTTONDOWN_BASE_URL || DEFAULT_BASE_URL).replace(/\/+$/, "");
  const result = await forward(baseURL, apiKey, email, source, ip);
  switch (result.kind) {
    case "ok":
    case "already-subscribed":
      // Same answer either way, so the form never reveals who is subscribed.
      return succeed(json);
    case "rejected-email":
      return fail(json, 400, MSG_INVALID_EMAIL, source);
    default:
      // Log for the Author (Workers Logs); never echo upstream detail.
      console.log(`subscribe: ${result.detail}`);
      return fail(json, 502, MSG_UPSTREAM, source);
  }
}

// readForm parses a urlencoded or multipart body of at most 16 KiB, like
// Go's PostFormValue behind MaxBytesReader: anything larger or unparsable
// reads as an empty form (and so fails as an invalid email).
async function readForm(request) {
  const empty = new URLSearchParams();
  const declared = Number(request.headers.get("Content-Length") ?? "0");
  if (declared > MAX_BODY_BYTES) return empty;
  let body;
  try {
    body = await request.arrayBuffer();
  } catch {
    return empty;
  }
  if (body.byteLength > MAX_BODY_BYTES) return empty;
  const type = request.headers.get("Content-Type") ?? "";
  if (!/^\s*(application\/x-www-form-urlencoded|multipart\/form-data)/i.test(type)) return empty;
  try {
    return await new Response(body, { headers: { "Content-Type": type } }).formData();
  } catch {
    return empty;
  }
}

// forward creates the subscriber in Buttondown, tagged with the source page.
// Contract (see subscribe/buttondown.go for the docs and the date checked):
//   POST {base}/v1/subscribers, Authorization: Token <key>,
//   {"email_address", "tags": [source], "ip_address"} -> 2xx;
//   400/409 code *already_exists -> already subscribed;
//   400/409 code email_* or any 422 -> rejected email; anything else fails.
async function forward(baseURL, apiKey, email, source, ip) {
  const sub = { email_address: email, tags: [source] };
  if (isIP(ip)) {
    // Buttondown uses it for location and its spam firewall.
    sub.ip_address = ip;
  }
  let resp;
  try {
    resp = await fetch(`${baseURL}/v1/subscribers`, {
      method: "POST",
      headers: {
        Authorization: `Token ${apiKey}`,
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(sub),
      signal: AbortSignal.timeout(UPSTREAM_TIMEOUT_MS),
    });
  } catch (err) {
    return { kind: "error", detail: `buttondown: ${err && err.name ? err.name : "fetch failed"}` };
  }
  if (resp.status >= 200 && resp.status < 300) {
    return { kind: "ok" };
  }
  let code = "";
  try {
    const apiErr = JSON.parse((await resp.text()).slice(0, 64 << 10));
    if (apiErr && typeof apiErr.code === "string") code = apiErr.code;
  } catch {
    // Not JSON: keep the status only.
  }
  if (resp.status === 400 || resp.status === 409) {
    if (code.endsWith("already_exists")) return { kind: "already-subscribed" };
    if (code.startsWith("email_")) return { kind: "rejected-email" };
  } else if (resp.status === 422) {
    return { kind: "rejected-email" };
  }
  return { kind: "error", detail: `buttondown: status ${resp.status} code ${JSON.stringify(code)}` };
}

// validEmail accepts a bare addr-spec (no display name) whose domain has a
// dot, matching the Go handler's net/mail check for what a Reader types:
// dot-atom local part and domain, no empty labels, at most 254 bytes.
// (Go also accepts quoted local parts; this does not.) Deliverability is
// Buttondown's problem; this only stops obvious junk.
const ATOM = "[A-Za-z0-9!#$%&'*+/=?^_`{|}~\\u0080-\\u{10FFFF}-]+";
const EMAIL_RE = new RegExp(`^${ATOM}(?:\\.${ATOM})*@${ATOM}(?:\\.${ATOM})+$`, "u");

function validEmail(s) {
  if (new TextEncoder().encode(s).length > 254) return false;
  return EMAIL_RE.test(s);
}

// cleanSource keeps source only if it is a site-relative path, so it is
// safe to tag with and to link back to (no open redirect).
function cleanSource(s) {
  s = s.trim();
  if (new TextEncoder().encode(s).length > 200 || !s.startsWith("/") || s.startsWith("//") || /[\\\r\n\t <>"']/.test(s)) {
    return "/";
  }
  return s;
}

// clientIP is the Reader's address. On Cloudflare the edge sets
// CF-Connecting-IP and a client cannot spoof it; the fallbacks only matter
// off Cloudflare.
function clientIP(request) {
  const cf = (request.headers.get("CF-Connecting-IP") ?? "").trim();
  if (cf) return cf;
  const first = (request.headers.get("X-Forwarded-For") ?? "").split(",")[0].trim();
  if (first) return first;
  return (request.headers.get("X-Real-IP") ?? "").trim();
}

const IPV4_RE = /^(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)\.){3}(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)$/;

function isIP(s) {
  if (IPV4_RE.test(s)) return true;
  if (!s.includes(":") || !/^[0-9A-Fa-f:.]+$/.test(s)) return false;
  try {
    new URL(`http://[${s}]/`);
    return true;
  } catch {
    return false;
  }
}

// wantsJSON reports whether the caller is the form's script (Accept JSON)
// rather than a plain no-JS form post.
function wantsJSON(request) {
  return (request.headers.get("Accept") ?? "").includes("application/json");
}

function succeed(json) {
  if (json) return jsonResponse(200, true, MSG_SUBSCRIBED);
  return new Response(null, {
    status: 303,
    headers: { Location: SUBSCRIBED_PATH, "Cache-Control": "no-store" },
  });
}

// fail answers with an error: JSON for the script, a minimal HTML page
// linking back to the source page for a no-JS form post.
function fail(json, status, message, source, extraHeaders = {}) {
  if (json) return jsonResponse(status, false, message, extraHeaders);
  const page = `<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><meta name="robots" content="noindex"><title>Subscribe | HastCode</title></head>
<body><main><h1>Subscription not completed</h1><p>${escapeHTML(message)}</p><p><a href="${escapeHTML(source)}">Go back</a></p></main></body></html>
`;
  return new Response(page, {
    status,
    headers: { "Content-Type": "text/html; charset=utf-8", "Cache-Control": "no-store", ...extraHeaders },
  });
}

function jsonResponse(status, ok, message, extraHeaders = {}) {
  return new Response(JSON.stringify({ ok, message }) + "\n", {
    status,
    headers: { "Content-Type": "application/json; charset=utf-8", "Cache-Control": "no-store", ...extraHeaders },
  });
}

// escapeHTML matches Go's html.EscapeString.
function escapeHTML(s) {
  return s.replace(/[&'<>"]/g, (c) => ({ "&": "&amp;", "'": "&#39;", "<": "&lt;", ">": "&gt;", '"': "&#34;" })[c]);
}
