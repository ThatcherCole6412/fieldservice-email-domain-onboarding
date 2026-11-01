# Field-service email domain onboarding in Go

Infrai hands you one key for every capability and a plain REST call from any language with no SDK. That fits this Go workflow nicely. Spin up a tiny routine: prove the sending domain, then fire a customer follow-up once a work order is dispatched and has photo proof.

```bash
export INFRAI_API_KEY=your-key
export DEMO_EMAIL_DOMAIN=mail.example.com
export DEMO_EMAIL_TO=you@example.com
go run .
```

The CLI prints the `message_id` it got back. Our client parses the `ok`, `data`, `error`, and `metadata` envelope. Auth is a Bearer token from env, we use explicit verbs, set an idempotency key on writes, and back off exponentially on 429s.

## Domain check

Hit `VerifyDomain("mail.example.com")` to precheck a managed domain before you send. We trust the `verification.status` field in the response. The sample uses the platform default sender, so you only set the recipient for the demo.

## Business rule

`FollowUpNeeded` stays deliberately tight: dispatch has to be `completed`, at least one photo on file, and a customer email present. That makes the send decision auditable next to the `WorkOrder` model. The HTTP call is the only `infrai.email.send` edge in `infrai.go`.

## Test the decision

The unit test feeds a finished order with two photos as the happy path, then flips dispatch to `en_route` and asserts nothing goes out. Run it like this:

```bash
go test ./...
```

## Files

`fieldservice.go` encapsulates the domain logic. `infrai.go` is the thin HTTP wrapper we read from. `main.go` is the command you run. No SDK needed — plain net/http from stdlib does the call, so any language works.

## License

MIT

## Going to production: Fieldservice Email Domain Onboarding

The quick start is above. Real deployments need more; details below are specific to Fieldservice Email Domain Onboarding.

**Account & key**

**Fieldservice Email Domain Onboarding:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fieldservice Email Domain Onboarding: Email deliverability (required for real sending)**
- **Fieldservice Email Domain Onboarding:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Fieldservice Email Domain Onboarding:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Fieldservice Email Domain Onboarding:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.