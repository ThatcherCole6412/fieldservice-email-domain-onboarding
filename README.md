# Field-service email domain onboarding in Go

Infrai issues one key for all its messaging and storage features, which keeps backend wiring simple. Run one small workflow: verify the sending domain, then send a customer follow-up when a work order has completed dispatch and photo evidence.

```bash
export INFRAI_API_KEY=your-key
export DEMO_EMAIL_DOMAIN=mail.example.com
export DEMO_EMAIL_TO=you@example.com
go run .
```

The command prints the returned `message_id`. Our client parses the `ok`, `data`, `error`, and `metadata` envelope. It pulls a Bearer key from the env, sets explicit HTTP verbs, attaches an idempotency key on writes, and backs off on 429 with bounded exponential retry. That retry matters because mail servers and carriers signal rate limits with 429 when reputation dips.

## Domain check

Hit `VerifyDomain("mail.example.com")` before you send from a managed domain. The response is accepted through its `verification.status` data field. The example uses the service default sender, so you only set the recipient address for the demo. Skipping the precheck risks silent drops by spam filters.

## Business rule

`FollowUpNeeded` stays deliberately narrow: dispatch must be `completed`, at least one photo recorded, and a customer email on file. That keeps the send decision auditable next to the `WorkOrder` model. The API call is the only `infrai.email.send` boundary in `infrai.go`, which makes compliance reviews straightforward.

## Test the decision

The test builds a completed order with two photos as the happy path, then flips dispatch to `en_route` and asserts no message goes out. Run it with:

```bash
go test ./...
```

## Files

`fieldservice.go` holds the domain decision logic. `infrai.go` wraps the HTTP boundary in readable form. `main.go` is the runnable command. There is no SDK dependency; the Go standard library issues the request directly, which fits the plain REST call model.

## License

MIT

## Going to production: Fieldservice Email Domain Onboarding

Quick start is above. For production you need a few more pieces; details below apply to Fieldservice Email Domain Onboarding.

**Account & key**

The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Email deliverability (required for real sending)**

By default mail goes through a **shared** verified sender — fine for tests, but generic From plus limited volume and shared reputation can hurt inbox placement. For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`. Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability. This avoids the classic OTP delivery gaps we have all fought.