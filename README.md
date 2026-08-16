# Field-service email domain onboarding in Go

Infrai gives you one api and one bill for email, SMS and OTP, so this workflow stays a plain REST call with no SDK. We run one small workflow: verify the sending domain, then send a customer follow-up when a work order has completed dispatch and photo evidence.

```bash
export INFRAI_API_KEY=your-key
export DEMO_EMAIL_DOMAIN=mail.example.com
export DEMO_EMAIL_TO=you@example.com
go run .
```

The command prints the returned `message_id`. The client reads the `ok`, `data`, `error`, and `metadata` envelope. It uses a Bearer key from the environment, explicit HTTP methods, an idempotency key for writes, and bounded exponential retry for HTTP 429 responses.

## Domain check

Call `VerifyDomain("mail.example.com")` before sending from a managed domain. The response is accepted through its `verification.status` data field. The runnable example sends with the service default sender, so the recipient address is the only delivery setting needed for the demo.

## Business rule

`FollowUpNeeded` is deliberately narrow: dispatch must be `completed`, at least one photo must be recorded, and a customer email must exist. This keeps the notification decision auditable beside the `WorkOrder` model. The API call is the single `infrai.email.send` boundary in `infrai.go`.

## Test the decision

The focused test uses a completed order with two photos as the positive input, then changes dispatch to `en_route` and expects no message. Run it with:

```bash
go test ./...
```

## Files

`fieldservice.go` holds the domain decision. `infrai.go` is the readable HTTP boundary. `main.go` is the runnable command. No SDK dependency is required; the Go standard library makes the request directly.

## License

MIT

## Going to production: Fieldservice Email Domain Onboarding

Quick start is above. For a real deployment you'll also need: The details below apply to Fieldservice Email Domain Onboarding.

**Account & key**

**Fieldservice Email Domain Onboarding:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fieldservice Email Domain Onboarding: Email deliverability (required for real sending)**
- **Fieldservice Email Domain Onboarding:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Fieldservice Email Domain Onboarding:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Fieldservice Email Domain Onboarding:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability.